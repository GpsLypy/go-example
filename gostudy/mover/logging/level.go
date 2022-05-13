package logging

import (
	"github.com/cihub/seelog"
	"sync"
)

type ErrorCode int32

const (
	EC_Success        ErrorCode = 10000
	EC_Warn           ErrorCode = 10001
	EC_Info           ErrorCode = 10002
	EC_Debug          ErrorCode = 10003
	EC_Trace          ErrorCode = 10004 // reserve
	EC_Lost_Etcd      ErrorCode = 10005
	EC_Config_ERR     ErrorCode = 10006
	EC_Job_Run_ERR    ErrorCode = 10007
	EC_Net_PKG_ERR    ErrorCode = 10008
	EC_Datasource_ERR ErrorCode = 10009
	EC_Runtime_ERR    ErrorCode = 10010
	EC_Ignore_ERR     ErrorCode = 10011
	EC_Encode_ERR     ErrorCode = 10012
	EC_Decode_ERR     ErrorCode = 10013

	EC_MAX ErrorCode = 10100
)

const (
	RedisLogMaxNum = 10000
	RedisLogPopNum = 100
)

var RedisLogCurNum = 0
var ErrorCnt = [EC_MAX - EC_Success]int{}
var ErrorLock sync.Mutex

func (p ErrorCode) String() string {
	switch p {
	case EC_Success:
		return "Success"
	case EC_Warn:
		return "Warning"
	case EC_Info:
		return "Info"
	case EC_Trace:
		return "Trace"
	case EC_Debug:
		return "Debug"
	case EC_Lost_Etcd:
		return "Lost etcd"
	case EC_Config_ERR:
		return "Config error"
	case EC_Job_Run_ERR:
		return "Job running error"
	case EC_Net_PKG_ERR:
		return "Host network error"
	case EC_Datasource_ERR:
		return "Lost data source"
	case EC_Runtime_ERR:
		return "Runtime error"
	case EC_Ignore_ERR:
		return "Ignore error"
	case EC_Encode_ERR:
		return "Encode error"
	case EC_Decode_ERR:
		return "Decode error"
	default:
		return "Unknown"
	}
}

func LogDebug(formated string, params ...interface{}) {
	if "" == formated {
		seelog.Debug(params...)
	} else {
		seelog.Debugf(formated, params...)
	}
}

func LogWarn(formated string, params ...interface{}) {
	if "" == formated {
		seelog.Warn(params...)
	} else {
		seelog.Warnf(formated, params...)
	}
}

func LogInfo(formated string, params ...interface{}) {
	if "" == formated {
		seelog.Info(params...)
	} else {
		seelog.Infof(formated, params...)
	}
}

func LogError(code ErrorCode, formated string, params ...interface{}) {
	if "" == formated {
		seelog.Error(params...)
	} else {
		seelog.Errorf(formated, params...)
	}
}
