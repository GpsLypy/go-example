package moverconfig

import (
	"bytes"
	"strings"
)

const (
	CT_JSON = iota
	CT_BINARY
)

const (
	SeedLog = `<seelog>
    <outputs formatid="main">
        <filter levels="debug,info,warn,critical,error">
            <console />
        </filter>
        <filter levels="debug,info,warn,critical,error">
            <file path="debuglog/debug.log"/>
        </filter>
    </outputs>
 
    <formats>
        <format id="main" format="%Date %Time [%LEV] %Msg (%File::%Func %Line)%n"/>
    </formats>
</seelog>`
)

const (
	DBType = iota
	DB_MYSQL
	DB_MSSQL
	// append here, same with protobuffer define

)

const (
	TargetDatabase = iota
	TargetTurboMQ
)

// mysql data type
const (
	MYSQL_TYPE_DECIMAL byte = iota
	MYSQL_TYPE_TINY
	MYSQL_TYPE_SHORT
	MYSQL_TYPE_LONG
	MYSQL_TYPE_FLOAT
	MYSQL_TYPE_DOUBLE
	MYSQL_TYPE_NULL
	MYSQL_TYPE_TIMESTAMP
	MYSQL_TYPE_LONGLONG
	MYSQL_TYPE_INT24
	MYSQL_TYPE_DATE
	MYSQL_TYPE_TIME
	MYSQL_TYPE_DATETIME
	MYSQL_TYPE_YEAR
	MYSQL_TYPE_NEWDATE
	MYSQL_TYPE_VARCHAR
	MYSQL_TYPE_BIT

	//mysql 5.6
	MYSQL_TYPE_TIMESTAMP2
	MYSQL_TYPE_DATETIME2
	MYSQL_TYPE_TIME2
)

const (
	MYSQL_TYPE_NEWDECIMAL byte = iota + 0xf6
	MYSQL_TYPE_ENUM
	MYSQL_TYPE_SET
	MYSQL_TYPE_TINY_BLOB
	MYSQL_TYPE_MEDIUM_BLOB
	MYSQL_TYPE_LONG_BLOB
	MYSQL_TYPE_BLOB
	MYSQL_TYPE_VAR_STRING
	MYSQL_TYPE_STRING
	MYSQL_TYPE_GEOMETRY
)

const (
	FieldDataTypeUnknown = iota
	FDT_INT
	FDT_FLOAT
	FDT_BOOL
	FDT_STRING
	FDT_BINARY
	FDT_DATE
	// just list special type, like date, because 0000-00-00 00:00:00 is NULL in mysql + python, will cause exception
)

type Task struct {
	TaskID int64 `json:"TaskID"`

	// source
	FromDBType     int        `json:"FromDBType"`
	FromIsShard    bool       `json:"FromIsShard"`
	FromEndpoint   Endpoint   `json:"FromEndpoint"`
	FromEndpoints  []Endpoint `json:"FromEndpoints"`
	FromDBName     string     `json:"FromDBName"`
	FromTable      string     `json:"FromTable"`
	OrgDBName      string     `json:"OrgDBName"`
	OrgTbName      string     `json:"OrgTbName"`
	FromField      string     `json:"FromField"`
	CurrWriteId    int64      `json:"CurrWriteId"`
	StrCurrWriteId string     `json:"StrCurrWriteId"`
	FromWhere      string     `json:"FromWhere"`

	// task data
	ClosedInterval int    `json:"ClosedInterval"`
	StrStartId     string `json:"strStartId"`
	StrEndId       string `json:"strEndId"`
	StrCurrId      string `json:"strCurrId"`

	StartId        int64  `json:"StartId"`
	EndId          int64  `json:"EndId"`
	PrimaryKey     string `json:"PrimaryKey"`
	PrimaryKeyType int    `json:"PrimaryKeyType"`
	CurrId         int64  `json:"CurrtId"`
	StartSkip      bool   `json:"startSkip"`
	Skip           int    `json:"skip"`
	WinSize        int    `json:"WinSize"`
	// destination
	DestTo           int      `json:"DestTo"`
	DestDBType       int      `json:"DestDBType"`
	DestEndpoint     Endpoint `json:"DestEndpoint"`
	DestDBName       string   `json:"DestDBName"`
	DestTable        string   `json:"DestTable"`
	DestField        []string `json:"DestField"`
	DestDBIsSharding bool     `json:"DestDBIsSharding"`
	DestMQInf        DestMQInf
	IsSelfPri        bool `json:"isSelfPri"`
	RowCountLimit    int  `json:"rowCountLimit"`
}

func (task *Task) IsEmpty() bool {
	if (task.PrimaryKeyType == FDT_INT && task.ClosedInterval == 0 && task.StartId == 0 && task.EndId == 0) ||
		(task.PrimaryKeyType == FDT_STRING && task.ClosedInterval == 0 && task.StrStartId == "" && task.StrEndId == "") {
		return true
	}

	return false
}

//strings.Compare() => 相等为0，不相等为1
func (task *Task) Completed() bool {
	if (FDT_INT == task.PrimaryKeyType && task.CurrId >= task.EndId) ||
		(FDT_STRING == task.PrimaryKeyType && strings.Compare(task.StrCurrId, task.StrEndId) >= 0) {
		return true
	}

	return false
}

type KeyRange struct {
	keyType int    `json:"keyType"`
	StartId string `json:"StartId"`
	EndId   string `json:"EndId"`
	CurrId  string `json:"CurrId"`
}

type RowData struct {
	Row []interface{}
}

type RowInfo struct {
	Data           *RowData
	Ds             *DataSource
	Task           *Task
	Fields         string
	CurrWriteId    int64
	StrCurrWriteId string

	MQInf DestMQInf
}

type FieldInfo struct {
	FieldName     string
	FieldTypeEx   byte
	FieldType     byte
	FieldTypeSelf byte
	IsUnsigned    bool
	IsNullAble    bool
	IsPriamryKey  bool
	IsAutoIncr    bool
	Idx           int
}

type TableStat struct {
	HostPort  string `json:"HostPort"`
	DbName    string `json:"DbName"`
	TableName string `json:"TableName"`

	SumRows  int64  `json:"SumRows"`
	TranRows int64  `json:"TranRows"`
	Sql      string `json:"sql"`
}

type CacheData struct {
	ExpireTime int64
	Data       interface{}
	DataS      interface{}
}

type ApiResult struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type column struct {
	idx            int
	name           string
	unsigned       bool
	auto_increment bool
	primary_key    bool
	null           bool
}

type table struct {
	schema string
	name   string

	columns []*column
}

func (t *table) ToString() string {
	tb := bytes.NewBuffer(nil)
	tb.WriteString("DATABASE[")
	tb.WriteString(t.schema)
	tb.WriteString("] TABLE[")
	tb.WriteString(t.name)
	tb.WriteString("] COLUMNS: | ")

	for _, v := range t.columns {
		tb.WriteString(v.name)
		tb.WriteString(" | ")
	}

	return tb.String()
}
