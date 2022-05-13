package main

import (
	"app/tools/excavator/localpprof"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"

	"app/tools/mover/help"
	"app/tools/mover/http"
	"app/tools/mover/info"
	"app/tools/mover/moverconfig"
	"encoding/json"
	"github.com/ngaut/log"
	"golang.org/x/net/context"
	"runtime"

	"app/tools/mover/breakpoint"
	"app/tools/mover/logging"
	"app/tools/mover/service"

	"github.com/cihub/seelog"
	"github.com/pkg/errors"
)

var (
	flagJobID  string
	flagEnv    string
	flagCfg    string
	flagNoQuit bool
	stopGenWg  sync.WaitGroup
	once       = sync.Once{}
)

func main() {
	defer func() {
		e := recover()
		if nil != e {
			logging.LogError(logging.EC_Success, "", "App quit with panic : ", e)
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("==> %s\n", string(buf[:n]))
		} else {
			logging.LogInfo("App quit normally")
		}
		seelog.Flush()
	}()

	// set logger first 先设置日志记录器
	logger, err := seelog.LoggerFromConfigAsFile("./conf/seelog.xml")
	if nil != err {
		logger, err = seelog.LoggerFromConfigAsBytes([]byte(moverconfig.SeedLog))
		if nil != err {
			panic(err)
		}
	}

	// set caller stack frame depth
	logger.SetAdditionalStackDepth(1)
	seelog.ReplaceLogger(logger)

	//解析命令行参数
	flag.StringVar(&flagJobID, "job", "", fmt.Sprintf("default: EV(%s)", help.ENV_SPECIFY_JOB_ID))
	flag.StringVar(&flagEnv, "env", "", fmt.Sprintf("default: EV(%s)", help.ENV_SPECIFY_JOB_ID))
	flag.StringVar(&flagCfg, "config", "", "default: {}")
	flag.BoolVar(&flagNoQuit, "noquit", false, "Do not quit if stopped")
	flag.Parse()

	//初始化环境信息
	if err := info.InitEnvInfo(flagEnv, flagJobID); nil != err {
		seelog.Error(err)
		return
	}

	var host string
	var env string
	var jobID int64

	rand.Seed(time.Now().UTC().UnixNano())
	host = info.GetEnvInfo().Host
	env = info.GetEnvInfo().EnvType
	jobID = info.GetEnvInfo().JobID
	service.NumWorkers = -1

	logging.LogInfo("Service container: %s, env: %s, jobID: %d, config: %s", host, env, jobID, flagCfg)
	logging.LogInfo("Go version: %s", runtime.Version())
	logging.LogInfo("%s version: %d, v%s", info.JobName, info.VersionCode, info.Version)

	// create etcd
	//etcd 是一个高可用强一致性的键值仓库[服务发现],存储了系统中服务的配置信息
	err = help.CreateEtcd(info.GetEnvInfo().EtcdEndpoints)
	if nil != err {
		logging.LogError(logging.EC_Lost_Etcd, "Can't create etcd api, envInfo: %v, err: %v", *info.GetEnvInfo(), err)
		return
	}

	// fetch config
	configData, err := fetchConfig(flagCfg, jobID)
	if err != nil {
		return
	}
	//设置配置信息
	err = setConfig(configData)
	if err != nil {
		return
	}
	// breakpoint module
	breakpoint.InitBreakPoint(jobID)

	// Enable local pprof
	if moverconfig.GetConfig().LocalPprofPath != "" {
		if err = localpprof.InitLocalPprof(context.TODO(), moverconfig.GetConfig().LocalPprofPath, "mover"); nil != err {
			seelog.Errorf("Can't start local pprof on %s: %s",
				moverconfig.GetConfig().LocalPprofPath, err.Error())
		}
	}

	// start mover
	go service.StartMover(jobID)

	//监听信号
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			// Print stack
			syscall.Signal(10))

		for {
			sig := <-sc

			if sig == syscall.Signal(10) {
				log.Infof("User signal1 to print stack information")
				//收集当前正在使用的所有goroutine的堆栈跟踪信息。
				pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
				continue
			}
			log.Infof("Got signal [%d] to exit.", sig)

			break
		}
	}()

	// wait quit & upload status
	//我们得使用带缓冲 channel,否则，发送信号时我们还没有准备好接收，就有丢失信号的风险
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	timer := time.NewTimer(time.Second * time.Duration(service.GetUploaderTime()))
	breakFlag := false
	tickerC := 0
	for {
		select {
		case <-timer.C:
			tickerC++
			timer.Reset(time.Second * time.Duration(service.GetUploaderTime()))

			go service.UploadStatusInfo("")

			if 0 == tickerC%10 /*%10*/ {
				if err := breakpoint.FlushBreakPoint(); err != nil {
					logging.LogError(logging.EC_Runtime_ERR, "flushBreakPoint file %s, err: %v", breakpoint.BpFile, err)
				}
			}

			// break, if all workers exit.
			if moverconfig.GetConfig().Target == moverconfig.TargetTurboMQ {
				if service.SlaveMgr.Finished {
					logging.LogInfo("Mq Job's all workers already exit.")
					breakFlag = true
				}
			} else {
				if service.NumWorkers == 0 {
					logging.LogInfo("Job's all workers already exit.")
					breakFlag = true
				}
			}

			if moverconfig.GetState().IsHealthy == false {
				logging.LogInfo("Job's status is unhealthy.")
				breakFlag = true
			}
		case <-c:
			logging.LogInfo("Job recv stopped command, ready to exit")
			once.Do(func() {
				close(service.StopJobChan)
			})
			breakFlag = true
			service.RecvStopped = true
		case <-service.OfflineWorkerCh:

			service.Stopworkers()
			if moverconfig.GetConfig().Target == moverconfig.TargetDatabase {
				breakFlag = true
			}

		case <-service.StopJobChan:
			logging.LogInfo("Job ready to exit")
			breakFlag = true
		}

		if breakFlag {
			break
		}
	}

	err = service.StopMover()
	if nil != err {
		logging.LogError(logging.EC_Job_Run_ERR, "Fail to stop job, envInfo: %v, err: %v", *info.GetEnvInfo(), err)
		return
	}

	service.ShowSummary()
	service.UploadSummary("")

	if flagNoQuit {
		for {
			// Do not quit if specified flag is set
			time.Sleep(time.Second)
		}
	}

	logging.LogInfo("Job exit, envInfo: %v", *info.GetEnvInfo())
}

// fetch config
func fetchConfig(flagCfg string, jobID int64) (string, error) {
	var configData string

	if flagCfg != "" {
		//读取配置文件
		data, err := ioutil.ReadFile(flagCfg)
		if nil != err {
			service.UploadStatusInfo(err.Error())
			logging.LogError(logging.EC_Config_ERR, "Fail to parse config, envInfo: %v, config: %s, err: %v", *info.GetEnvInfo(), flagCfg, err)
			return configData, err
		}
		configData = string(data)
	} else {
		configKey := help.GetEtcdJobConfigKey(jobID)
		confdata, err := help.GetEtcdKey(configKey)
		if nil != err {
			service.UploadStatusInfo(err.Error())
			logging.LogError(logging.EC_Config_ERR, "Fail to get config, envInfo: %v, key: %s, err: %v", *info.GetEnvInfo(), configKey, err)
			return "", err
		}
		configData = confdata
	}
	return configData, nil
}

func prepareRdsConfig() error {
	etcd_gateway := info.GetEnvInfo().EtcdEndpoints[0]
	jobid := info.GetEnvInfo().JobID
	rdsconf := moverconfig.GetConfig().RdsConf
	env := moverconfig.GetConfig().RdsConf.Env

	var url string

	if len(rdsconf.RdsAddr) > 0 && len(rdsconf.RdsApi) > 0 && "qatest" == env {
		url = fmt.Sprintf("%s%s", rdsconf.RdsAddr, rdsconf.RdsApi)
	} else {
		url = fmt.Sprintf("%s/datasource?jobid=%d&env=%s", etcd_gateway, jobid, env)
	}

	//url = "http://10.101.20.50:15338/mock/datasource"
	seelog.Infof("Get rds config url: %s", url)
	body, err := http.HttpGetBody(url)
	if nil != err {
		return err
	}

	res := &moverconfig.ApiResult{}
	err = json.Unmarshal(body, res)
	if nil != err {
		seelog.Errorf("Json unmarshal rds http Response, err=%v", err)
		return err
	}

	if res.Code != 0 {
		return errors.Errorf("Can't get rds datasource: %v", res.Message)
	}
	(*moverconfig.GetCurrState()).RdsMessage = res.Message

	return moverconfig.SetRdsConfig([]byte(res.Message))
}

func setConfig(configData string) error {
	err := moverconfig.SetConfig(moverconfig.CT_JSON, []byte(configData))
	if nil != err {
		service.UploadStatusInfo(err.Error())
		logging.LogError(logging.EC_Config_ERR, "Fail to set config, envInfo: %v, err: %v", *info.GetEnvInfo(), err)
		return err
	}

	if moverconfig.GetConfig().LogLevel == "" {
		moverconfig.GetConfig().LogLevel = "info"
	}
	if moverconfig.GetConfig().DataSizeLimit == 0 {
		moverconfig.GetConfig().DataSizeLimit = 100
	}
	log.SetLevelByString(strings.ToLower(moverconfig.GetConfig().LogLevel))

	// pull rds config, if needed
	if moverconfig.GetConfig().IsNeedRdsSources() {
		err = prepareRdsConfig()
		if nil != err {
			seelog.Errorf("Fail to pull rds datasource, envInfo: %v, err: %v", *info.GetEnvInfo(), err)
			return err
		}
	}
	return nil
}
