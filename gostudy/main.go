package main

import (
	"fmt"
	"strings"
)

//golang字符串操作
func main() {
	s := "Hello world hello world"
	str := "Hello"
	//var s = []string{"11","22","33"}

	//比较字符串，区分大小写，比”==”速度快。相等为0，不相等为-1。
	ret := strings.Compare(s, str)
	fmt.Println(ret) //  1
}

func getDatasourceJobs(datasourceID int, needSync, needPub bool, needRedis bool) (*distributerSubscribeDetail, error) {
	distributers, err := model.GetDistributersByDatasource(datasourceID)
	if err != nil {
		return nil, errors.Trace(err)
	}
	var rsp distributerSubscribeDetail
	for _, distribute := range distributers {
		if err = fetchDistributerJobs(datasourceID, distribute, &rsp, nil); nil != err {
			return nil, errors.Trace(err)
		}
	}
	details := make([]*JobDetail, 0, len(rsp.Detail))
	for _, d := range rsp.Detail {
		if (d.TaskType == TaskTypeMq && needPub) ||
			(d.TaskType == TaskTypeDb && needSync) ||
			(d.TaskType == TaskTypeRedis && needRedis) {
			details = append(details, d)
		}
	}
	rsp.Detail = details
	return &rsp, nil
}

func fetchDistributerJobs(datasourceID int,
	distribute *model.DistributerModel,
	rsp *distributerSubscribeDetail,
	tools []int) error {
	haGet, err := hagroup.GetHaGroupInfo(int64(distribute.ID))
	if err != nil {
		return errors.Annotatef(err, "getHaGroupJob,dtid=%v,err=%v", distribute.ID, err)
	}
	if len(tools) == 0 {
		tools = []int{tooltype.ToolMover, tooltype.ToolLabour}
	}
	jobIds, err := model.GetJobsByTaskTools(strconv.Itoa(distribute.ID),
		model.TaskTypeDistributer,
		tools)
	if err != nil {
		return errors.Trace(err)
	}
	if len(jobIds) == 0 {
		return nil
	}
	jobPairMap := make(map[int64]int64)
	jobs, err := model.GetJobNonSDKInfoBulk(jobIds)
	if err != nil {
		return errors.Trace(err)
	}
	for _, v := range jobs {
		jobPairMap[v.ID] = v.UpstreamJobId
		jobPairMap[v.UpstreamJobId] = v.ID
	}
	jobCfgs, err := model.GetJobConfigBatch(jobIds)
	if err != nil {
		return errors.Trace(err)
	}
	jobConfigMap := make(map[int64]*model.JobConfigModel)
	for _, cfg := range jobCfgs {
		jobConfigMap[cfg.JobID] = cfg
	}

	for i := range jobs {
		job := &jobs[i]
		jobCfg := jobConfigMap[job.ID]
		if jobCfg == nil {
			// Config not available
			fillJobDetail(rsp,
				job,
				errConfigNotAvailable.Error(),
				[]string{errConfigNotAvailable.Error()},
				distribute.ID,
				datasourceID,
				[]string{errConfigNotAvailable.Error()},
				errConfigNotAvailable.Error(),
				errConfigNotAvailable.Error(),
				0,
				errConfigNotAvailable.Error())
			continue
		}

		if job.Tool == tooltype.ToolMover {
			var cfg moverconfig.MoverConfig
			if err = json.Unmarshal([]byte(jobCfg.Content), &cfg); nil != err {
				return errors.Trace(err)
			}

			tables := make([]string, 0, 32)
			topics := make([]string, 0, 32)
			destDB := ""

			for _, task := range cfg.TaskList {
				tables = append(tables, task.From.DBName+"."+task.From.TableName)

				if task.DestTo == 0 {
					if destDB != "" {
						destDB += ","
					}
					for _, v := range task.Dest.Endpoints {
						destDB += fmt.Sprintf("%s:%d", v.Host, v.Port)
					}
				} else if task.DestTo == 1 {
					if task.DestMQ.Topic != "" {
						topics = append(topics, task.DestMQ.Topic)
					} else {
						topics = append(topics, task.From.DBName+"_"+task.From.TableName+"_full")
					}
				}
			}

			if len(topics) != 0 {
				fillJobDetail(rsp,
					job,
					"",
					tables,
					distribute.ID,
					datasourceID,
					topics,
					cfg.Dispatcher.TurboMQConf.NamesrvAddr,
					strings.Join(cfg.Dispatcher.TurboMQConf.Brokers, ","),
					cfg.Dispatcher.TurboMQConf.QueueNumber,
					TaskTypeMq)
			}
			if len(destDB) != 0 {
				fillJobDetail(rsp,
					job,
					destDB,
					tables,
					distribute.ID,
					datasourceID,
					nil,
					"",
					"",
					0,
					TaskTypeDb)
			}
		} else if job.Tool == tooltype.ToolLabour {
			if haGet != nil {
				// Only shows the active group job ?
				if job.HaGroup != haGet.HaGroup {
					continue
				}
				if job.UpstreamJobId == 0 {
					job.UpstreamJobId = jobPairMap[job.ID]
				}
			}
			var cfg labourconfig.AppConfig
			if err = json.Unmarshal([]byte(jobCfg.Content), &cfg); nil != err {
				return errors.Trace(err)
			}
			if nil == cfg.SyncRule {
				continue
			}
			if nil == cfg.TurboMQOb && nil == cfg.MysqlOb {
				continue
			}
			cfg.SyncRule.Init()

			tables := make([]string, 0, 32)
			topics := make([]string, 0, 32)
			for dbName, dv := range cfg.SyncRule.DBRule {
				finalDBName := dbName
				if dv.RewriteName != "" {
					finalDBName = dv.RewriteName
				}
				if len(dv.TableRule) == 0 || dv.Sharding {
					tables = append(tables, finalDBName+".*")
					break
				}
				for tableName, tv := range dv.TableRule {
					finalTableName := tableName
					if tv.RewriteName != "" {
						finalTableName = tv.RewriteName
					}
					tables = append(tables, finalDBName+"."+finalTableName)

					if cfg.TurboMQOb != nil {
						if tv.Attributes["topic"] != "" {
							topics = append(topics, tv.Attributes["topic"])
						} else {
							topics = append(topics, dbName+"_"+tableName)
						}
					}
				}
			}
			sort.Strings(tables)

			if cfg.TurboMQOb != nil {
				fillJobDetail(rsp,
					job,
					"",
					tables,
					distribute.ID,
					datasourceID,
					topics,
					cfg.TurboMQOb.To.NamesvrAddr,
					strings.Join(cfg.TurboMQOb.To.Brokers, ","),
					cfg.TurboMQOb.To.QueueNumber,
					TaskTypeMq)
			} else if cfg.MysqlOb != nil {
				var destDB string
				for _, v := range cfg.MysqlOb.Tos {
					if destDB != "" {
						destDB += ","
					}
					destDB += fmt.Sprintf("%s:%d", v.Host, v.Port)
				}

				fillJobDetail(rsp,
					job,
					destDB,
					tables,
					distribute.ID,
					datasourceID,
					nil,
					"",
					"",
					0,
					TaskTypeDb)
			}
		} else if job.Tool == tooltype.ToolRedisShake {
			var cfg toolconfig.RedisShakeConfig
			if _, err = toml.Decode(jobCfg.Content, &cfg); nil != err {
				return errors.Trace(err)
			}

			fillJobDetail(rsp,
				job,
				cfg.TargetAddress,
				[]string{cfg.TargetName},
				distribute.ID,
				datasourceID,
				nil,
				"",
				"",
				0,
				TaskTypeRedis)
		}
	}
	return nil
}



1、根据数据源id所有任务
[
	{
		"jobid": 1,
		"isdel": 0,
		"haGroup": "",
		"upstreamJobId": 0,
		"slaveId": 0,
		"updateversion": 0,
		"target": 0,
		"status": 1,
		"owner": 123,
		"alive": 0,
		"tag": "",
		"jobtype": 1,
		"tool": 8,
		"disablealarm": 0,
		"createtime": 1652061932,
		"updatetime": null,
		"starttime": null,
		"lastchecktime": null,
		"stoptime": null
	}
]

2、根据数据源id 获取升级双中心的任务
[
	{
		"groupId": 119,
		"dsId": 48,
		"activatedJobIds": [1710],
		"activatedHaGroup":"激活的中心名称",
		"jobList": [{
			{
				"jobid": 1710,
				"isdel": 0,
				"haGroup": "",
				"upstreamJobId": 0,
				"slaveId": 1914,
				"updateversion": 0,
				"target": 0,
				"status": 1,
				"owner": 167,
				"alive": 0,
				"tag": "DataSub:=TETrmsReport:48[0]->10.100.159.200:9876;10.100.157.34:9876",
				"jobtype": 1,
				"tool": 11,
				"disablealarm": 0,
				"createtime": 1631674560,
				"updatetime": 1636701049
			},
			{
				"jobid": 1914,
				"isdel": 0,
				"haGroup": "",
				"upstreamJobId": 1710,
				"slaveId": 1710,
				"updateversion": 0,
				"target": 0,
				"status": 1,
				"owner": 167,
				"alive": 0,
				"tag": "DataSub:=TETrmsReport:48[0]->10.100.159.200:9876;10.100.157.34:9876",
				"jobtype": 1,
				"tool": 11,
				"disablealarm": 0,
				"createtime": 1631674560,
				"updatetime": 1636701049
			}
		}]
	},
	{
		
		"groupId": 11111,
		"dsId": 48,
		"activatedJobIds": null,
		"activatedHaGroup": null,
		"jobList":[
		]
	}
]





