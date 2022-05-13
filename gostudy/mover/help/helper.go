package help

import (
	"app/tools/mover/logging"
	"app/tools/mover/utils"
	_ "fmt"
	"os"
	"strconv"
)

const ENV_SPECIFY_HOST = "MACHINEIP"
const ENV_SPECIFY_SYS_ENV = "DAOKEENVTYPE"
const ENV_SPECIFY_JOB_ID = "JOBID"

func GetHost() string {
	// env
	host := os.Getenv(ENV_SPECIFY_HOST)
	if "" != host {
		return host
	}

	// read hardware
	ifaMap, err := utils.GetAllInterface()
	if nil != err {
		logging.LogError(logging.EC_Net_PKG_ERR, "Get all interface error, err: %v", err)

		return "127.0.0.1"
	}

	var ifaName string
	ifaName = "bond0"
	if ip, ok := ifaMap[ifaName]; ok {
		return ip
	}
	ifaName = "em1"
	if ip, ok := ifaMap[ifaName]; ok {
		return ip
	}
	ifaName = "eth0"
	if ip, ok := ifaMap[ifaName]; ok {
		return ip
	}
	ifaName = "本地连接" // windows
	if ip, ok := ifaMap[ifaName]; ok {
		return ip
	}
	return "127.0.0.1"
}

func GetEnvType(flagEnv string) string {
	if "" == flagEnv {
		//检索并返回名为ENV_SPECIFY_SYS_ENV的环境变量的值。如果不存在该环境变量会返回空字符串。
		envStr := os.Getenv(ENV_SPECIFY_SYS_ENV)
		flagEnv = envStr
	}

	if "" == flagEnv {
		flagEnv = "test"
	}

	return flagEnv
}

func GetJobID(flagJobID string) int64 {
	var jobID int64 = 0

	if flagJobID != "" {
		jobIDInt, err := strconv.Atoi(flagJobID)
		if nil != err {
			jobIDInt = 0
		}
		jobID = int64(jobIDInt)
	}

	if 0 == jobID {
		jobIDStr := os.Getenv(ENV_SPECIFY_JOB_ID)
		//返回字符串表示的整数值
		jobIDEnv, err := strconv.ParseInt(jobIDStr, 10, 64)
		if nil != err {
			return 0
		}
		jobID = jobIDEnv
	}

	return jobID
}
