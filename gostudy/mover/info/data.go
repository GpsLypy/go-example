package info

import (
	"lib/tooltype"
	"os"
)

const VersionCode = 10
const Version = "0.1.0.7"
const JobName = "mover"

const CacheTime = 30
const DataSizePerTask = 100000
const DataSizePerSql = 100

const MssqlZeroDate = "1900-01-01 00:00:00"
const ParamMaxSizePerSqlInMssql = 2100

const ToolType = tooltype.ToolMover

const (
	DefaultUploadTime          = 5
	DefaultPoolMaxActive       = 10
	DefaultPoolMaxIdle         = 1
	DefaultPoolIdleTimeoutSecs = 180
)

const EtcdStatusTtl = 259200 // 3 day
const EtcdMonitorTtl = 30

var (
	// Separator /
	Separator = string(os.PathSeparator)
)
