功能：
	全量订阅功能
config说明：
    1. TaskList      []DataSourcePair
        数据迁移的任务队列，详见：DataSourcePair说明
    2. WorkerSize    int
        工作协程数量
    3. Target      int
		if 0, mover to db; 1, mover to MQ
	4. FloatToString int
		if 1 ==> 投递过程float 转换为string 投递 
	5. Dispatcher 	Dispatcher
		Mq 配置信息
	6. DataDir  string 断点传输的临时文件存放位置

DataSourcePair说明：
    1. From    DataSource
        数据源头，详见：DataSource说明
        
        如果是普通的，根据PK拆分任务
        如果是分片的，根据实例列表，并根据PK拆分任务
    2. FromField  string
        空：  不需要清洗字段，且与目的地字段名相同
        非空：对数据源进行字段清洗，
            如果和Dest的字段名不同，或者需要转换数据内容，需要使用alias
            例如：
                field1, srcField2 as destField2, substring(srcDatetimeField3, 0, 10) as destDate

                field1：和dest字段类型、名称相同，不需要alias
                srcField2：和dest类型相同，但名称不同，需要alias
                srcDatetimeField3：和dest类型、名称都不同，进行数据转换并alias
        Note：当From.TableName='*'，此处必须为空
    3. FromWhere  string
        空：  不需要过滤数据
        非空：需要对原始数据进行过滤
            例如：
                createTime>'2017-04-01' and state=0
        Note：当From.TableName='*'，此处必须为空
    4. Dest    DataSource
        数据目的地，详见：DataSource说明
        Note：当From.TableName='*'，Dest.TableName必须为空

DataSource说明：
    1. DBType     int
        db类型：1 mysql、2 mssql
        sharding必须是mysql
    2. IsSharding bool
        是否分片
        如果分片，使用dbName获取分片配置，并得到真实的DBName、tableName、endpoint
    3. DBName     string
        db名称，分片表为逻辑库名
    4. TableName  string
        表名称，分片表为逻辑表名
        *：  src库下的所有表，此时src与dest的表名称必须一致（如果dest.TableName不为空，则src的所有表的数据导入一张dest表）
        非*：指定的表名称，当src与dest的表名不一致时，需要明确指定出
    5. Endpoints  []Endpoint
        db的账号、地址，普通DB只有一个
        如果src为sharding，此项为多个后端实例地址
        如果dest为sharding，此项可为多个，利于插入速度的提高（sharding时不能batch insert，只能多个dbRouter加速插入）
        


配置举例： // 从分表表 同步到mysql

{
  "TaskList": [
    {
      "From": {
        "DBType": 1,
        "IsSharding": true,
        "DBName": "db_ymj_0524", // 分片表逻辑库名
        "TableName": "tbl_ymj_0525", // 分片表逻辑表名
        "Endpoints": [  			// 后端真实实例地址
          {
            "Host": "192.168.1.2",
            "Port": 4406,
            "User": "root",
            "Password": "P@ss1234"
          },
          {
            "Host": "192.168.1.2",
            "Port": 4407,
            "User": "root",
            "Password": "P@ss1234"
          }
        ]
      },
	  "Target": 0, //  if 0, mover to db; 1, mover to MQ
      "FromField": "id,age,numInt",  // 清洗字段 配的字段是只需要的字段
      "FromWhere": "numInt>=4 and create_time>='2019-02-12'", // 根据条件过滤
      "Dest": {
        "DBType": 1,
        "IsSharding": false,
        "DBName": "db_ymj_0524",
        "TableName": "tbl_ymj_0525",
        "Endpoints": [
          {
            "Host": "192.168.1.2",
            "Port": 3306,
            "User": "root",
            "Password": "password#dbr"
          }
        ]
      }
    }
  ],
  "DataDir": "/data/xxx",    // 记录文件临时存放位置 传输完成会删除
  "WorkerSize": 1,
  "ShardingConfigApi": "http://192.168.1.2:9021"
}


//  非分片表 同步到mysql 
{
  "TaskList": [
    {
      "From": {
        "DBType": 1,
        "IsSharding": false,
        "DBName": "test_mover",
        "TableName": "mover_02",  // 如果是* dest的填表名就是同步到一张表 src的表结构必须相同
        "Endpoints": [
          {
            "Host": "192.168.1.2",
            "Port": 3306,
            "User": "root",
            "Password": "password#dbr"
          }
        ]
      },
      "FromField": "",
      "FromWhere": "",
      "Dest": {
        "DBType": 1,
        "IsSharding": false,
        "DBName": "test_mover_02",
        "TableName": "",   // 表名不写就同src表名
        "Endpoints": [
          {
            "Host": "192.168.1.2",
            "Port": 3306,
            "User": "root",
            "Password": "password#dbr"
          }
        ]
      }
    }
  ],
  "DataDir": "/data/xxx",
  "WorkerSize": 1,
  "ShardingConfigApi": "http://192.168.1.2:9021"
}


// 分表表同步到mq
{
  "TaskList": [
    {
      "From": {
        "DBType": 1,
        "IsSharding": true,
        "DBName": "db_ymj_0524",
        "TableName": "tbl_ymj_0525",
        "Endpoints": [    // 后端真实实例地址
          {
            "Host": "192.168.1.2",
            "Port": 4406,
            "User": "root",
            "Password": "P@ss1234"
          },
          {
            "Host": "192.168.1.2",
            "Port": 4407,
            "User": "root",
            "Password": "P@ss1234"
          }
        ]
      },
      "FromField": "",
      "FromWhere": "",
      "DestMQ": { // 投递到哪个topic  为空则投递到 DBName_TableName_full 
        "Topic": "test_tc001",
		"primaryKeys":["id"]   // 根据 
      }
    }
  ],
  "DataDir": "/data/xxx",
  "Target": 1, // to mq
  "WorkerSize": 6,
  "Dispatcher": { // mq 配置
    "TurboMQConf": {
      "producer_group": "DataSync",
      "namesrv_addr": "192.168.1.2:9876;192.168.1.2:9876",
      "brokers": [
        "broker-c"
      ],
      "queue_number": 1
    }
  },
  "ShardingConfigApi": "http://192.168.1.2:9021"
}


// 非分片表 投递mq
{
  "TaskList": [
    {
      "From": {
        "DBType": 1,
        "IsSharding": false,
        "DBName": "SyncMover",
        "TableName": "SyncMover_tw",
        "Endpoints": [
          {
            "Host": "192.168.1.2",
            "Port": 3306,
            "User": "root",
            "Password": "password#dbr"
          }
        ]
      },
      "FromField": "",
      "FromWhere": "",
      "DestMQ": {
        "Topic": "SyncMover_tt03"
      }
    }
  ],
  "FloatToString": 1,  // float 投递为string
  "RowsPerTask": 1000,
  "Target": 1,
  "DataDir": "/data/xxx",
  "WorkerSize": 6,
  "Dispatcher": {
    "TurboMQConf": {
      "producer_group": "DataSync",
      "namesrv_addr": "192.168.1.2:9876;192.168.1.2:9876",
      "brokers": [
        "broker-c",
        "broker-d"
      ],
      "queue_number": 8
    }
  },
  "ShardingConfigApi": "http://192.168.1.2:9021"
}
