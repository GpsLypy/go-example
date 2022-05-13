package moverconfig

import (
	"app/platform/etcd-gateway/rdsdef"
	"app/tools/excavator/replay/writer"
	"github.com/pkg/errors"
	"sync"
)

var (
	currRdsDS   *rdsdef.DsGetRsp
	currConfig  *MoverConfig
	configMutex sync.Mutex
)

type UserInf struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

type RdsConfig struct {
	RdsAddr           string     `json:"RdsAddr"` // qa test so ignore myself
	RdsApi            string     `json:"RdsApi"`  // qa test so ignore myself
	Env               string     `json:"Env"`
	User              UserInf    `json:"User"`              // Data source user information in tidiness mode
	BackendUser       UserInf    `json:"BackendUser"`       // If sharding, backend instance user information in tidiness mode
	DetailedEndpoints []Endpoint `json:"DetailedEndpoints"` // Enumerate each endpoint in detailed mode
}

type Endpoint struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
}

type DataSource struct {
	DBType     int        `json:"DBType"`
	IsSharding bool       `json:"IsSharding"` // if sharding, get db config from configCenter
	BackendEPs []Endpoint `json:"BackendEPs"` // if has ids, get all backend endpoints
	DBName     string     `json:"DBName"`
	TableName  string     `json:"TableName"` // *: represent all table, except has been listed
	OrgDBName  string     `json:"-"`
	OrgTbName  string     `json:"-"`
	Endpoints  []Endpoint `json:"Endpoints"` // maybe more than one, it's useful when destination is sharding
}

type DataSourcePair struct {
	From      DataSource `json:"From"`
	FromField string     `json:"FromField"`
	FromWhere string     `json:"FromWhere"`
	DestTo    int        `json:"-"` // dest type, 0: to database, 1: to MQ
	Dest      DataSource `json:"Dest"`
	DestMQ    DestMQInf  `json:"DestMQ"`
}

// TableInfo 数据库表配置
type DestMQInf struct {
	Strategy    *StrategyConfig `json:"strategy"`
	Topic       string          `json:"Topic"`
	SplitKeys   []string        `json:"SplitKeys"`
	PrimaryKeys []string        `json:"primaryKeys"`
}

// TurboMQConf MQ配置信息
type TurboMQConf struct {
	ProducerGroup string   `json:"producer_group"`
	NamesrvAddr   string   `json:"namesrv_addr"`
	Brokers       []string `json:"brokers"`
	QueueNumber   int      `json:"queue_number"`
	BatchSize     int      `json:"batch_size"`
}

type MoverConfig struct {
	JobId          int              `json:"JobId"`
	LogLevel       string           `json:"log-level"`
	DataSizeLimit  int              `json:"dataSizeLimit"`
	Prepare        int              `json:"Prepare"`
	TaskList       []DataSourcePair `json:"TaskList"`
	PriKey         string           `json:"priKey"`
	ShardingTable  int              `json:"ShardingTable"`
	WorkerSize     int              `json:"WorkerSize"`
	RowsPerTask    int              `json:"RowsPerTask"`
	FixedExeTime   int              `json:"FixedExeTime"`
	Target         int              `json:"Target"`        // if 0, mover to db; 1, mover to MQ
	FloatToString  int              `json:"FloatToString"` // if 1 float to string
	WriteMode      int              `json:"WriteMode"`
	RdsConf        RdsConfig        `json:"RdsConf"`
	Dispatcher     Dispatcher       `json:"Dispatcher"`
	DataDir        string           `json:"DataDir"`
	IsMaster       int              `json:"isMaster"`
	LocalPprofPath string           `json:"localPprofPath"` // Enable local pprof if path is not empty
}

type Dispatcher struct {
	TurboMQConf   TurboMQConf         `json:"TurboMQConf"`
	Kafka         *writer.KafkaConfig `json:"kafka"`
	BatchCommitMS int                 `json:"BatchCommitMS"`
}

func (ds *DataSource) SetEndpoints(endpoints []Endpoint) error {
	if len(endpoints) == 0 {
		return errors.New("invalid endpoints")
	}

	ds.Endpoints = endpoints

	return nil
}

//查看ip和port
func (ep Endpoint) Validate() bool {
	if ep.Host != "" && ep.Port > 0 {
		return true
	}

	return false
}

//检查数据源是否有效
func (ds DataSource) Validate() bool {
	//是否存在实例信息
	if len(ds.Endpoints) == 0 {
		return false
	}
	//检查实例信息的有效性
	for _, ep := range ds.Endpoints {
		if !ep.Validate() {
			return false
		}
	}
	//检查是否需要分片
	if ds.IsSharding {
		if nil == ds.Endpoints || len(ds.Endpoints) == 0 {
			return false
		}

		for _, ep := range ds.Endpoints {
			if !ep.Validate() {
				return false
			}
		}
	}

	return true
}

func (tc MoverConfig) IsNeedRdsSources() bool {
	for _, task := range tc.TaskList {
		if !task.From.Validate() {
			return true
		}
	}

	return false
}
