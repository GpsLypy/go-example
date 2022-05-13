package service

import (
	"app/platform/jobmgr/logger"
	"app/tools/mover/breakpoint"
	"app/tools/mover/info"
	"app/tools/mover/logging"
	"app/tools/mover/moverconfig"
	"app/tools/mover/utils"
	"database/sql"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"strings"

	"strconv"

	"github.com/eapache/channels"
	_ "github.com/go-sql-driver/mysql"
	"unicode"
)

var (
	mysqlDbMutex sync.Mutex
	mysqlDbMap   map[string]*sql.DB

	mysqlFieldInfoMutex    sync.Mutex
	mysqlFieldInfoCacheMap map[string]moverconfig.CacheData
)

type MysqlPolicy struct {
}

func getMysqlDB4Ds(ds moverconfig.DataSource, DBName string) (*sql.DB, string, error) {
	var dbConn *sql.DB = nil
	var err error = nil
	var mysqlAddr string

	for i := 0; i < len(ds.Endpoints); i++ {
		endpoint := ds.Endpoints[i]

		mysqlAddr = GetConnstr(endpoint, DBName, ds.IsSharding)
		dbConn, err = getMysqlDB(mysqlAddr)
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "connect ep:%v, err: %v.", endpoint, err)
			continue
		}

		err = dbConn.Ping()
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "connect ep:%v, err: %v.", endpoint, err)
			continue
		}
		moverconfig.AppendStateFromHost(endpoint.Host + ":" + strconv.Itoa(endpoint.Port))
		WorkerSize := moverconfig.GetConfig().WorkerSize
		dbConn.SetMaxIdleConns(WorkerSize)

		break
	}

	return dbConn, mysqlAddr, errors.Trace(err)
}

func getMysqlDB4DsFilterDs(ds *moverconfig.DataSource, DBName string) (*sql.DB, string, error) {
	var dbConn *sql.DB = nil
	var err error = nil
	var mysqlAddr string

	for i := 0; i < len(ds.Endpoints); i++ {
		Endpoint := ds.Endpoints[i]

		mysqlAddr = GetConnstr(Endpoint, DBName, ds.IsSharding)
		dbConn, err = getMysqlDB(mysqlAddr)
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "Fail to get mysql db,mysqlAddr=%s err: %v", mysqlAddr, err)
			continue
		}

		err = dbConn.Ping()
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "Fail to ping mysql db, ds: %v, err: %v", ds, err)
			continue
		}
		moverconfig.AppendStateFromHost(Endpoint.Host + ":" + strconv.Itoa(Endpoint.Port))
		ds.Endpoints = []moverconfig.Endpoint{Endpoint}
		break
	}

	return dbConn, mysqlAddr, errors.Trace(err)
}

func (p *MysqlPolicy) CheckDataSource(ds moverconfig.DataSource) error {
	var dbConn *sql.DB = nil
	var err error = nil
	for i := 0; i < len(ds.Endpoints); i++ {
		Endpoint := ds.Endpoints[i]

		mysqlAddr := GetConnstr(Endpoint, ds.DBName, ds.IsSharding)
		dbConn, err = getMysqlDB(mysqlAddr)
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "Fail to get mysql db,mysqlAddr=%s err: %v", mysqlAddr, err)
			continue
		}

		err = dbConn.Ping()
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "Fail to ping mysql db, ds: %v, err: %v", ds, err)
			continue
		}
		moverconfig.AppendStateFromHost(Endpoint.Host + ":" + strconv.Itoa(Endpoint.Port))
		break
	}

	return err
}

func (p *MysqlPolicy) GetShardingDataSource(ds moverconfig.DataSource, baseDbName string, baseTableName string) ([]moverconfig.DataSource, error) {
	oriDbName := ds.DBName
	if "" != baseDbName {
		oriDbName = baseDbName
	}

	oriTableName := ds.TableName
	if "" != baseTableName {
		oriTableName = baseTableName
	}

	if len(ds.Endpoints) == 0 {
		return nil, errors.Trace(errors.New("Get sharding datasource error, endpoint empty"))
	}

	dbConn, _, err := getMysqlDB4DsFilterDs(&ds, "")
	if nil != err {
		return nil, errors.Trace(err)
	}

	// read all db
	dbNameList, err := showMysqlDB(dbConn)
	if nil != err {
		return nil, errors.Trace(err)
	}

	// read every table
	dsList := make([]moverconfig.DataSource, 0)
	for _, dbName := range dbNameList {
		if !isShardingDB(oriDbName, dbName) {
			continue
		}

		tableNameList, err := showMysqlTable(dbConn, dbName)
		if nil != err {
			return nil, errors.Trace(err)
		}
		for _, tableName := range tableNameList {
			if !IsShardingTable(oriTableName, tableName) {
				continue
			}

			newDs := ds
			newDs.DBName = dbName
			newDs.TableName = tableName
			dsList = append(dsList, newDs)
		}
	}

	// get datasource
	return dsList, nil
}

func (p *MysqlPolicy) GetTables(ds moverconfig.DataSource) ([]string, error) {
	if len(ds.Endpoints) == 0 {
		return nil, errors.Trace(errors.New("Ds no endpoint"))
	}

	dbConn, _, err := getMysqlDB4Ds(ds, "")
	if nil != err {
		logging.LogError(logging.EC_Datasource_ERR, "Get mysql db error, %v", err)
		return nil, errors.Trace(err)
	}

	return showMysqlTable(dbConn, ds.DBName)
}

func (p *MysqlPolicy) GetFields(ds moverconfig.DataSource, fields string, isUserDef bool) (map[string]moverconfig.FieldInfo, []moverconfig.FieldInfo, error) {
	mysqlFieldInfoMutex.Lock()
	defer mysqlFieldInfoMutex.Unlock()

	if nil == mysqlFieldInfoCacheMap {
		mysqlFieldInfoCacheMap = make(map[string]moverconfig.CacheData)
	}

	// check cache
	key := fmt.Sprintf("%s:%s:%s", ds.DBName, ds.TableName, fields)
	if v, ok := mysqlFieldInfoCacheMap[key]; ok {
		//检查缓存是否过期
		if v.ExpireTime >= time.Now().Unix() {
			return v.Data.(map[string]moverconfig.FieldInfo), v.DataS.([]moverconfig.FieldInfo), nil
		}
		delete(mysqlFieldInfoCacheMap, key)
	}

	// query
	logging.LogInfo("Get fields, key=%s", key)
	if len(ds.Endpoints) == 0 {
		return nil, nil, errors.Trace(errors.New("Ds no endpoint"))
	}
	//拿到操作数据库的dbConn等信息
	dbConn, mysqlAddr, err := getMysqlDB4Ds(ds, "")
	if nil != err {
		logging.LogError(logging.EC_Datasource_ERR, "Fail to get mysql conn, addr=%s, err: %v", mysqlAddr, err)
		return nil, nil, errors.Trace(err)
	}

	index := 0
	fieldMap := make(map[string]moverconfig.FieldInfo)
	fieldArr := make([]moverconfig.FieldInfo, 0)
	if "" == fields {
		sql := fmt.Sprintf("SELECT COLUMN_NAME,DATA_TYPE,COLUMN_KEY,IS_NULLABLE,COLUMN_TYPE,EXTRA FROM information_schema.`COLUMNS` WHERE TABLE_SCHEMA = '%s' and TABLE_NAME='%s'",
			ds.DBName, ds.TableName)
		logging.LogDebug("Sql: %s", sql)

		rows, err := dbConn.Query(sql)
		if nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "Fail to query sql, addr=%s, sql=%s, err: %v", mysqlAddr, sql, err)
			return nil, nil, errors.Trace(err)
		}
		defer rows.Close()

		for rows.Next() {
			var columnName, dataType, columnKey, isNullAble, columnType, extra string
			err = rows.Scan(&columnName, &dataType, &columnKey, &isNullAble, &columnType, &extra)
			if nil != err {
				logging.LogError(logging.EC_Runtime_ERR, "Fail to scan row, addr=%s, sql=%s, err: %v", mysqlAddr, sql, err)
				return nil, nil, errors.Trace(err)
			}

			var fieldInfo moverconfig.FieldInfo
			fieldInfo.FieldName = columnName
			fieldInfo.FieldType = getMysqlDataType(dataType, isUserDef)
			fieldInfo.FieldTypeSelf = getMysqlDataType4SelfToMq(dataType)

			if strings.ToUpper(isNullAble) == "YES" {
				fieldInfo.IsNullAble = true
			}

			if strings.Contains(columnType, "unsigned") || strings.Contains(columnType, "UNSIGNED") {
				fieldInfo.IsUnsigned = true
			}

			if strings.Index(strings.ToUpper(columnKey), "PRI") >= 0 {
				fieldInfo.IsPriamryKey = true
			}

			if strings.Index(strings.ToUpper(extra), "AUTO_INCREMENT") >= 0 {
				fieldInfo.IsAutoIncr = true
			}

			fieldInfo.Idx = index
			index++

			fieldMap[columnName] = fieldInfo
			fieldArr = append(fieldArr, fieldInfo)
		}
	} else {
		// can't get field type
		sql := fmt.Sprintf("SELECT %s FROM `%s`.`%s` LIMIT 1", fields, ds.DBName, ds.TableName)
		logging.LogDebug("Sql: %s", sql)

		rows, err := dbConn.Query(sql)
		if nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "Fail to query sql, addr=%s, sql=%s, err: %v", mysqlAddr, sql, err)
			return nil, nil, errors.Trace(err)
		}
		defer rows.Close()

		columnNames, err := rows.Columns()
		if nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "Fail to get row column, addr=%s, sql=%s, err: %v", mysqlAddr, sql, err)
			return nil, nil, errors.Trace(err)
		}

		var fieldInfo moverconfig.FieldInfo
		fieldInfo.FieldType = moverconfig.FieldDataTypeUnknown
		for _, columnName := range columnNames {
			fieldInfo.FieldName = columnName
			fieldMap[columnName] = fieldInfo
			fieldArr = append(fieldArr, fieldInfo)
		}

		sql = fmt.Sprintf("SELECT COLUMN_NAME,DATA_TYPE,COLUMN_KEY,IS_NULLABLE,COLUMN_TYPE,EXTRA FROM information_schema.`COLUMNS` WHERE TABLE_SCHEMA = '%s' and TABLE_NAME='%s' and COLUMN_NAME in('%s')",
			ds.DBName, ds.TableName, strings.Join(columnNames, "','"))
		seelog.Debugf("Sql: %s", sql)

		rows, err = dbConn.Query(sql)
		if nil != err {
			return nil, nil, errors.Trace(err)
		}
		defer rows.Close()

		for rows.Next() {
			var columnName, dataType, columnKey, isNullAble, columnType, extra string
			err = rows.Scan(&columnName, &dataType, &columnKey, &isNullAble, &columnType, &extra)
			if nil != err {
				return nil, nil, errors.Trace(err)
			}

			var fieldInfo *moverconfig.FieldInfo
			// find FieldInfo item
			for i, fi := range fieldArr {
				if fi.FieldName == columnName {
					fieldInfo = &fieldArr[i]
				}
			}

			if fieldInfo == nil {
				errorMsg := fmt.Sprintf("Can't fine fields %s", columnName)
				return nil, nil, errors.Trace(errors.New(errorMsg))
			}

			// var fieldInfo FieldInfo
			fieldInfo.FieldName = columnName
			fieldInfo.FieldType = getMysqlDataType(dataType, IsSelfDef)
			fieldInfo.FieldTypeSelf = getMysqlDataType4SelfToMq(dataType)
			if strings.ToUpper(isNullAble) == "YES" {
				fieldInfo.IsNullAble = true
			}
			if strings.Index(strings.ToUpper(columnType), " UNSIGNED") > 0 {
				fieldInfo.IsUnsigned = true
			}
			if strings.Index(strings.ToUpper(columnKey), "PRI") >= 0 {
				fieldInfo.IsPriamryKey = true
			}
			if strings.Index(strings.ToUpper(extra), "AUTO_INCREMENT") >= 0 {
				fieldInfo.IsAutoIncr = true
			}

			fieldInfo.Idx = index
			index++

			fieldMap[columnName] = *fieldInfo
		}
	}

	// add cache
	var cacheData moverconfig.CacheData
	cacheData.DataS = fieldArr
	cacheData.Data = fieldMap
	cacheData.ExpireTime = time.Now().Unix() + int64(info.CacheTime)
	key = fmt.Sprintf("%s:%s:%s", ds.DBName, ds.TableName, fields)
	mysqlFieldInfoCacheMap[key] = cacheData

	return fieldMap, fieldArr, nil
}

func (p *MysqlPolicy) GetTableStat(ds moverconfig.DataSource, whereExpr string) (moverconfig.TableStat, error) {
	dbConn, mysqlAddr, err := getMysqlDB4Ds(ds, "")
	if nil != err {
		logging.LogError(logging.EC_Datasource_ERR, "Fail to get mysql conn, addr=%s, err: %v", mysqlAddr, err)
		return moverconfig.TableStat{}, err
	}

	stats := moverconfig.TableStat{
		TableName: ds.TableName,
		DbName:    ds.DBName,
		SumRows:   0,
	}

	var rowSum sql.NullInt64
	var sql string
	if "" == whereExpr {
		sql = fmt.Sprintf("SELECT count(*) FROM `%s`.`%s`", ds.DBName, ds.TableName)
	} else {
		sql = fmt.Sprintf("SELECT count(*) FROM `%s`.`%s` WHERE %s", ds.DBName, ds.TableName, whereExpr)
	}

	logging.LogInfo("Sql: %s", sql)

	rows, err := dbConn.Query(sql)
	if nil != err {
		logging.LogError(logging.EC_Runtime_ERR, "Fail to query sql, addr=%s, sql=%s, err: %v", mysqlAddr, sql, err)
		return moverconfig.TableStat{}, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&rowSum)
		if nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "Fail to scan row, addr=%s, sql=%s, err: %v", mysqlAddr, sql, err)
			return moverconfig.TableStat{}, err
		}

		stats.SumRows = rowSum.Int64
	}

	return stats, nil
}

func GetNameByHostAndPort(host string, port int, dbname string) string {
	return fmt.Sprintf("%s:%d:%s", host, port, dbname)
}

func ReserveBuffer(buf []byte, appendSize int) []byte {
	newSize := len(buf) + appendSize
	if cap(buf) < newSize {
		// Grow buffer exponentially
		newBuf := make([]byte, len(buf)*2+appendSize)
		copy(newBuf, buf)
		buf = newBuf
	}
	return buf[:newSize]
}

//单一主键任务派发
func mysqlDispatchTask(dsp moverconfig.DataSourcePair, taskChan *channels.InfiniteChannel, fieldList []string, priKey string, priKeyTyp int) error {
	// get db conn
	dbConn, _, err := getMysqlDB4Ds(dsp.From, "")
	if nil != err {
		return errors.Trace(err)
	}
	var conditions []string
	if "" != dsp.FromWhere {
		conditions = append(conditions, dsp.FromWhere)
	}

	bp := breakpoint.FindBreakPoint(GetNameByHostAndPort(dsp.From.Endpoints[0].Host, dsp.From.Endpoints[0].Port, dsp.From.DBName), dsp.From.TableName)
	//如果主键类型是int型
	if moverconfig.FDT_INT == priKeyTyp && nil != bp && bp.PrimaryKeyType == priKeyTyp {
		bpCond := fmt.Sprintf("`%s` > %d", bp.PrimaryKey, bp.EndId)
		conditions = append(conditions, bpCond)
	} else if nil != bp && bp.PrimaryKeyType == priKeyTyp {
		bpCond := fmt.Sprintf("`%s` > '%s'", bp.PrimaryKey, bp.StrEndId)
		conditions = append(conditions, bpCond)
	}
	var rowSum, rowMin, rowMax sql.NullInt64
	var strRowMin, strRowMax, strRowCur sql.NullString
	var sum, min, max int64
	var strmin, strmax string
	var Xsql string
	// get table sum & maxid
	if moverconfig.FDT_INT == priKeyTyp {
		// var rowSum, rowMin, rowMax sql.NullInt64
		// var sql string
		if 0 == len(conditions) {
			Xsql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s`", priKey, priKey, dsp.From.DBName, dsp.From.TableName)
		} else {
			Xsql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s` WHERE %s", priKey, priKey, dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "))
		}

		logging.LogDebug("Sql: %s", Xsql)

		row := dbConn.QueryRow(Xsql)
		err = row.Scan(&rowSum, &rowMin, &rowMax)
		if nil != err {
			logger.Errorf("mysqlDispatchTaskByInt-sql=%s-err=%v", Xsql, err)
			return errors.Trace(err)
		}
		logging.LogInfo("Sum=%v, min=%v, max=%v, dsp: %v", rowSum, rowMin, rowMax, dsp)

		//var sum, min, max int64
		sum = rowSum.Int64
		if rowMin.Valid {
			min = rowMin.Int64
		}
		if rowMax.Valid {
			max = rowMax.Int64
		}
	} else {
		if 0 == len(conditions) {
			Xsql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s`", priKey, priKey, dsp.From.DBName, dsp.From.TableName)
		} else {
			Xsql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s` WHERE %s", priKey, priKey, dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "))
		}

		logging.LogDebug("Sql: %s", Xsql)

		row := dbConn.QueryRow(Xsql)
		err = row.Scan(&rowSum, &strRowMin, &strRowMax)
		if nil != err {
			logger.Errorf("mysqlDispatchTaskByInt-sql=%s-err=%v", Xsql, err)
			return errors.Trace(err)
		}

		logging.LogInfo("Sum=%v, min=%v, max=%v, dsp: %v", rowSum, strRowMin, strRowMax, dsp)
		// var sum int64
		// var min, max string
		sum = rowSum.Int64
		if strRowMin.Valid {
			strmin = strRowMin.String
		}
		if strRowMax.Valid {
			strmax = strRowMax.String
		}
	}

	ResetAppendTableStat(moverconfig.TableStat{
		HostPort:  dsp.From.Endpoints[0].Host + ":" + strconv.Itoa(dsp.From.Endpoints[0].Port),
		DbName:    dsp.From.DBName,
		TableName: dsp.From.TableName,
		SumRows:   rowSum.Int64,
		Sql:       Xsql,
	})
	// split job
	var MyDataSizePerTask int64 = info.DataSizePerTask
	if dsp.Dest.IsSharding {
		MyDataSizePerTask = info.DataSizePerTask / 10
		if MyDataSizePerTask <= 0 {
			MyDataSizePerTask = 1
		}
	}

	// for qatest
	{
		if moverconfig.GetConfig().RowsPerTask > 0 {
			MyDataSizePerTask = int64(moverconfig.GetConfig().RowsPerTask)
		}
	}
	//根据表的数据总量和任务的最大处理能力得到具体的任务数量
	var taskSize int
	if sum%MyDataSizePerTask == 0 {
		taskSize = int(sum / MyDataSizePerTask)
	} else {
		taskSize = int(sum/MyDataSizePerTask + 1)
	}
	var step int64
	if moverconfig.FDT_INT == priKeyTyp {
		// var step int64
		idSize := max - min + 1
		if idSize%int64(taskSize) == 0 {
			step = idSize / int64(taskSize)
		} else {
			step = idSize/int64(taskSize) + 1
		}
	}

	if 0 == taskSize {

		var fields string
		if "" == dsp.FromField {
			fields = "`" + strings.Join(fieldList, "`,`") + "`"
		} else {
			fields = priKey + "," + dsp.FromField
		}

		empTask := moverconfig.Task{
			FromDBType:    dsp.From.DBType,
			FromEndpoints: dsp.From.Endpoints,
			FromEndpoint:  dsp.From.Endpoints[0],

			FromDBName:     dsp.From.DBName,
			FromTable:      dsp.From.TableName,
			FromField:      fields,
			FromWhere:      dsp.FromWhere,
			OrgDBName:      dsp.From.OrgDBName,
			OrgTbName:      dsp.From.OrgTbName,
			PrimaryKey:     priKey,
			PrimaryKeyType: priKeyTyp,
			DestTo:         dsp.DestTo,
			DestMQInf:      dsp.DestMQ,
			ClosedInterval: 0,
			StartId:        0,
			EndId:          0,
		}

		seelog.Debugf("Empty task: %v", empTask)
		taskChan.In() <- empTask
		atomic.AddInt32(&(*moverconfig.GetCurrState()).SumTasks, 1)

		return nil
	}

	//准备定义任务
	var task moverconfig.Task
	task.FromDBType = dsp.From.DBType
	task.FromIsShard = dsp.From.IsSharding
	task.FromEndpoint = dsp.From.Endpoints[0]
	moverconfig.AppendStateFromHost(task.FromEndpoint.Host + ":" + strconv.Itoa(task.FromEndpoint.Port))
	task.FromEndpoints = dsp.From.Endpoints
	task.FromDBName = dsp.From.DBName
	task.FromTable = dsp.From.TableName
	task.OrgDBName = dsp.From.OrgDBName
	task.OrgTbName = dsp.From.OrgTbName
	if "" == dsp.FromField {
		task.FromField = "`" + strings.Join(fieldList, "`,`") + "`"
	} else {
		if moverconfig.GetConfig().Target == moverconfig.TargetDatabase {
			task.FromField = dsp.FromField
		} else {
			task.FromField = priKey + "," + dsp.FromField
		}
	}
	task.FromWhere = dsp.FromWhere
	task.PrimaryKey = priKey
	task.PrimaryKeyType = priKeyTyp
	task.DestTo = dsp.DestTo

	if dsp.DestTo == moverconfig.TargetDatabase {
		task.DestDBType = dsp.Dest.DBType
		task.DestEndpoint = dsp.Dest.Endpoints[0]
		task.DestDBName = dsp.Dest.DBName
		task.DestTable = dsp.Dest.TableName
		task.DestField = fieldList
		task.DestDBIsSharding = dsp.Dest.IsSharding
	} else {
		task.DestMQInf = dsp.DestMQ
	}

	destEndpointSize := len(dsp.Dest.Endpoints)

	if moverconfig.FDT_INT == priKeyTyp {
		curId := min - 1
		for i := 0; i < taskSize; i++ {
			if curId >= max {
				logging.LogWarn("Premature termination, curId=%d, maxId=%d, dsp: %v", curId, max, dsp)
				break
			}

			task.TaskID = int64(i)
			task.StartId = curId
			task.CurrId = curId
			curId += step
			if curId > max {
				curId = max
			}
			task.EndId = curId
			if dsp.DestTo == 0 && destEndpointSize > 1 {
				// round robin 轮询
				task.DestEndpoint = dsp.Dest.Endpoints[i%destEndpointSize]
			}
			//限流策略，对查询的数据量作出限制，防止用户配置过大导致内存爆掉
			dataSizeLimit := moverconfig.GetConfig().DataSizeLimit

			if (dataSizeLimit*len(fieldList))*2 < 65535 {
				task.RowCountLimit = dataSizeLimit
			} else {
				logging.LogError(logging.EC_Config_ERR, "Invalid config,DataSizelimit(%d) too big,It should be under(%d) ", dataSizeLimit, 65535/(len(fieldList)*2))
				return errors.New("Invalid config")
			}

			logging.LogDebug("Task: %+v", task)
			taskChan.In() <- task
			atomic.AddInt32(&(*moverconfig.GetCurrState()).SumTasks, 1)
		}
	} else {
		curId := strmin
		for i := 0; i < taskSize; i++ {
			if strings.Compare(curId, strmax) > 0 {
				logging.LogWarn("Premature termination, curId=%s, maxId=%s, dsp: %v", curId, max, dsp)
				break
			}

			task.ClosedInterval = 0
			if 0 == i {
				task.ClosedInterval = 1
			} else {
				task.ClosedInterval = 0
			}

			task.StrStartId = curId
			task.StrCurrId = curId
			if strings.Compare(curId, strmax) >= 0 {
				curId = strmax
			} else {
				if "" == dsp.FromWhere {
					Xsql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` > '%s' LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, priKey, curId, MyDataSizePerTask)
				} else {
					Xsql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE (%s) and `%s` > '%s' LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, dsp.FromWhere, priKey, curId, MyDataSizePerTask)
				}

				logging.LogDebug("Sql: %s", Xsql)

				row := dbConn.QueryRow(Xsql)
				err = row.Scan(&strRowCur)
				if nil != err && sql.ErrNoRows != err {
					return errors.Trace(err)
				}
				if sql.ErrNoRows == err || strings.Compare(strRowCur.String, strmax) >= 0 {
					curId = strmax
				} else {
					curId = strRowCur.String
				}
			}

			task.StrEndId = curId

			if dsp.DestTo == 0 && destEndpointSize > 1 {
				// round robin
				task.DestEndpoint = dsp.Dest.Endpoints[i%destEndpointSize]
			}

			logging.LogDebug("Task: %v", task)
			taskChan.In() <- task
			atomic.AddInt32(&(*moverconfig.GetCurrState()).SumTasks, 1)
		}
	}
	return nil
}

//联合主键任务派发
func mysqlDispatchTaskByPri(dsp moverconfig.DataSourcePair, taskChan *channels.InfiniteChannel, fieldList []string, priKey string, priKeyTyp int) error {
	// get db conn
	dbConn, _, err := getMysqlDB4Ds(dsp.From, "")
	if nil != err {
		return errors.Trace(err)
	}
	var conditions []string
	if "" != dsp.FromWhere {
		conditions = append(conditions, dsp.FromWhere)
	}

	bp := breakpoint.FindBreakPoint(GetNameByHostAndPort(dsp.From.Endpoints[0].Host, dsp.From.Endpoints[0].Port, dsp.From.DBName), dsp.From.TableName)
	//如果主键类型是int型
	if moverconfig.FDT_INT == priKeyTyp && nil != bp && bp.PrimaryKeyType == priKeyTyp {
		bpCond := fmt.Sprintf("`%s` > %d", bp.PrimaryKey, bp.EndId)
		conditions = append(conditions, bpCond)
	} else if nil != bp && bp.PrimaryKeyType == priKeyTyp {
		bpCond := fmt.Sprintf("`%s` > '%s'", bp.PrimaryKey, bp.StrEndId)
		conditions = append(conditions, bpCond)
	}

	// get table sum & maxid
	var rowSum sql.NullInt64
	var rowMin, rowMax, rowCount sql.NullInt64
	var strRowMin, strRowMax, rowCur sql.NullString
	var xSql, countSql string
	if moverconfig.FDT_INT == priKeyTyp {
		if len(conditions) == 0 {
			xSql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s`", priKey, priKey, dsp.From.DBName, dsp.From.TableName)
		} else {
			xSql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s` WHERE %s", priKey, priKey, dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "))
		}
		if len(conditions) == 0 {
			countSql = fmt.Sprintf("SELECT count(*) as counts FROM `%s`.`%s` group by `%s` order by counts desc limit 1", dsp.From.DBName, dsp.From.TableName, priKey)
		} else {
			countSql = fmt.Sprintf("SELECT count(*) as counts FROM `%s`.`%s` WHERE %s group by `%s` order by counts desc limit 1", dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "), priKey)
		}

		logging.LogDebug("Sql: %s CountSql :%s", xSql, countSql)

		row := dbConn.QueryRow(xSql)
		err = row.Scan(&rowSum, &rowMin, &rowMax)

		if nil != err {
			logger.Errorf("mysqlDispatchTaskBy-sql=%s-err=%v", xSql, err)
			return errors.Trace(err)
		}

		logging.LogInfo("Sum=%v, min=%v, max=%v, dsp: %v", rowSum, rowMin, rowMax, dsp)

	} else {
		if len(conditions) == 0 {
			xSql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s`", priKey, priKey, dsp.From.DBName, dsp.From.TableName)
		} else {
			xSql = fmt.Sprintf("SELECT count(*), min(`%s`), max(`%s`) FROM `%s`.`%s` WHERE %s", priKey, priKey, dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "))
		}
		if len(conditions) == 0 {
			countSql = fmt.Sprintf("SELECT count(*) as counts FROM `%s`.`%s` group by `%s` order by counts desc limit 1", dsp.From.DBName, dsp.From.TableName, priKey)
		} else {
			countSql = fmt.Sprintf("SELECT count(*) as counts FROM `%s`.`%s` WHERE %s group by `%s` order by counts desc limit 1", dsp.From.DBName, dsp.From.TableName, strings.Join(conditions, " AND "), priKey)
		}

		logging.LogDebug("Sql: %s CountSql :%s", xSql, countSql)

		row := dbConn.QueryRow(xSql)
		err = row.Scan(&rowSum, &strRowMin, &strRowMax)

		if nil != err {
			logger.Errorf("mysqlDispatchTaskBy-sql=%s-err=%v", xSql, err)
			return errors.Trace(err)
		}

		logging.LogInfo("Sum=%v, min=%v, max=%v, dsp: %v", rowSum, strRowMin, strRowMax, dsp)
	}
	ResetAppendTableStat(moverconfig.TableStat{
		HostPort:  dsp.From.Endpoints[0].Host + ":" + strconv.Itoa(dsp.From.Endpoints[0].Port),
		DbName:    dsp.From.DBName,
		TableName: dsp.From.TableName,
		SumRows:   rowSum.Int64,
		Sql:       xSql,
	})

	if rowSum.Int64 == 0 {
		return nil
	}

	rowC := dbConn.QueryRow(countSql)
	err = rowC.Scan(&rowCount)
	if nil != err {
		logger.Errorf("mysqlDispatchTaskBy-countSql=%s-err=%v", countSql, err)
		return errors.Trace(err)
	}

	var sum int64
	var min, max int64
	var strmin, strmax string
	if moverconfig.FDT_INT == priKeyTyp {
		sum = rowSum.Int64
		if rowMin.Valid {
			min = rowMin.Int64
		}
		if rowMax.Valid {
			max = rowMax.Int64
		}
	} else {
		sum = rowSum.Int64
		if strRowMin.Valid {
			strmin = strRowMin.String
		}
		if strRowMax.Valid {
			strmax = strRowMax.String
		}
	}
	// split job
	var MyDataSizePerTask int64 = info.DataSizePerTask
	if dsp.Dest.IsSharding {
		MyDataSizePerTask = info.DataSizePerTask / 10
		if MyDataSizePerTask <= 0 {
			MyDataSizePerTask = 1
		}
	}

	// for qatest
	{
		if moverconfig.GetConfig().RowsPerTask > 0 {
			MyDataSizePerTask = int64(moverconfig.GetConfig().RowsPerTask)
		}
	}
	if MyDataSizePerTask < rowCount.Int64+1 {
		MyDataSizePerTask = rowCount.Int64 + 1
	}

	var taskSize int
	if sum%MyDataSizePerTask == 0 {
		taskSize = int(sum / MyDataSizePerTask)
	} else {
		taskSize = int(sum/MyDataSizePerTask + 1)
	}

	if 0 == taskSize {
		var fields string
		if "" == dsp.FromField {
			fields = "`" + strings.Join(fieldList, "`,`") + "`"
		} else {
			fields = priKey + "," + dsp.FromField
		}

		empTask := moverconfig.Task{
			FromDBType:     dsp.From.DBType,
			FromEndpoints:  dsp.From.Endpoints,
			FromEndpoint:   dsp.From.Endpoints[0],
			FromDBName:     dsp.From.DBName,
			FromTable:      dsp.From.TableName,
			FromField:      fields,
			FromWhere:      dsp.FromWhere,
			OrgDBName:      dsp.From.OrgDBName,
			OrgTbName:      dsp.From.OrgTbName,
			PrimaryKey:     priKey,
			PrimaryKeyType: priKeyTyp,
			DestTo:         dsp.DestTo,
			DestMQInf:      dsp.DestMQ,
			ClosedInterval: 0,
			StartId:        0,
			EndId:          0,
		}

		seelog.Debugf("Empty task: %v", empTask)
		taskChan.In() <- empTask
		atomic.AddInt32(&(*moverconfig.GetCurrState()).SumTasks, 1)

		return nil
	}

	var task moverconfig.Task
	task.IsSelfPri = true
	task.FromDBType = dsp.From.DBType
	task.FromIsShard = dsp.From.IsSharding
	task.FromEndpoint = dsp.From.Endpoints[0]
	moverconfig.AppendStateFromHost(task.FromEndpoint.Host + ":" + strconv.Itoa(task.FromEndpoint.Port))
	task.FromEndpoints = dsp.From.Endpoints
	task.FromDBName = dsp.From.DBName
	task.FromTable = dsp.From.TableName
	task.OrgDBName = dsp.From.OrgDBName
	task.OrgTbName = dsp.From.OrgTbName
	if "" == dsp.FromField {
		task.FromField = "`" + strings.Join(fieldList, "`,`") + "`"
	} else {
		if moverconfig.GetConfig().Target == moverconfig.TargetDatabase {
			task.FromField = dsp.FromField
		} else {
			task.FromField = priKey + "," + dsp.FromField
		}
	}
	task.FromWhere = dsp.FromWhere
	task.PrimaryKey = priKey
	task.PrimaryKeyType = priKeyTyp
	task.DestTo = dsp.DestTo

	if dsp.DestTo == 0 {
		task.DestDBType = dsp.Dest.DBType
		task.DestEndpoint = dsp.Dest.Endpoints[0]
		task.DestDBName = dsp.Dest.DBName
		task.DestTable = dsp.Dest.TableName
		task.DestField = fieldList
		task.DestDBIsSharding = dsp.Dest.IsSharding
	} else {
		task.DestMQInf = dsp.DestMQ
	}

	var curId int64
	var strCurId string
	var destEndpointSize int
	if moverconfig.FDT_INT == priKeyTyp {
		//限流，联合主键与rowCount.Int64有关
		if ((int(rowCount.Int64)+1)*len(fieldList))*2 < 65535 {
			task.RowCountLimit = int(rowCount.Int64) + 1
			if task.RowCountLimit < 100 {
				task.RowCountLimit = 100
			}
		} else {
			task.RowCountLimit = 65535 / (len(fieldList) * 2)
		}
		destEndpointSize = len(dsp.Dest.Endpoints)
		curId = min
	} else {
		task.RowCountLimit = int(rowCount.Int64) + 1
		destEndpointSize = len(dsp.Dest.Endpoints)
		strCurId = strmin
	}

	if moverconfig.FDT_INT == priKeyTyp {
		for i := 0; i < taskSize; i++ {
			var rowCur sql.NullInt64
			if curId >= max {
				logging.LogWarn("Premature termination, curId=%d, maxId=%d, dsp: %v", curId, max, dsp)
				break
			}
			task.TaskID = int64(i)
			task.StartId = curId
			task.CurrId = curId

			if curId >= max {
				curId = max
			} else {
				if "" == dsp.FromWhere {
					xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` >= %d order by `%s` LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, priKey, curId, priKey, MyDataSizePerTask)
				} else {
					xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE (%s) and `%s` >= %d  order by `%s` LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, dsp.FromWhere, priKey, curId, priKey, MyDataSizePerTask)
				}

				logging.LogDebug("Sql: %s", xSql)

				row := dbConn.QueryRow(xSql)
				err = row.Scan(&rowCur)
				if nil != err && sql.ErrNoRows != err {
					return errors.Trace(err)
				}
				if sql.ErrNoRows == err || rowCur.Int64 >= max {
					curId = max
				} else {
					curId = rowCur.Int64
				}
			}

			task.EndId = curId
			moverconfig.AppendTask(fmt.Sprintf("taskId=%d,startId=%v,taskEndId=%v,currId=%v,endId=%v,rowCur.Int64=%v,xSql=%v", task.TaskID, task.StartId, task.CurrId, task.EndId, curId, rowCur.Int64, xSql))

			if dsp.DestTo == 0 && destEndpointSize > 1 {
				// round robin
				task.DestEndpoint = dsp.Dest.Endpoints[i%destEndpointSize]
			}

			logging.LogDebug("Task: %v", task)
			taskChan.In() <- task
			atomic.AddInt32(&(*moverconfig.GetCurrState()).SumTasks, 1)
		}
	} else {
		for i := 0; i < taskSize; i++ {
			if strings.Compare(strCurId, strmax) > 0 {
				logging.LogWarn("Premature termination, curId=%s, maxId=%s, dsp: %v", curId, max, dsp)
				break
			}

			task.ClosedInterval = 0
			if 0 == i {
				task.ClosedInterval = 1
			} else {
				task.ClosedInterval = 0
			}

			task.StrStartId = strCurId
			task.StrCurrId = strCurId
			if strings.Compare(strCurId, strmax) >= 0 {
				strCurId = strmax
			} else {
				if "" == dsp.FromWhere {

					xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` >= '%s' order by `%s` LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, priKey, strCurId, priKey, MyDataSizePerTask)
					//xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` >= %s order by `%s` LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, priKey, curId, priKey, MyDataSizePerTask)
				} else {

					xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE (%s) and `%s` >= '%s' order by `%s`  LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, dsp.FromWhere, priKey, strCurId, priKey, MyDataSizePerTask)
					//xSql = fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE (%s) and `%s` >= %s order by `%s`  LIMIT %d,1;", priKey, dsp.From.DBName, dsp.From.TableName, dsp.FromWhere, priKey, curId, priKey, MyDataSizePerTask)
				}

				logging.LogDebug("Sql: %s", xSql)

				row := dbConn.QueryRow(xSql)
				err = row.Scan(&rowCur)
				if nil != err && sql.ErrNoRows != err {
					return errors.Trace(err)
				}
				if sql.ErrNoRows == err || strings.Compare(rowCur.String, strmax) >= 0 {
					strCurId = strmax
				} else {
					strCurId = rowCur.String
				}
			}

			task.StrEndId = strCurId

			if dsp.DestTo == 0 && destEndpointSize > 1 {
				// round robin
				task.DestEndpoint = dsp.Dest.Endpoints[i%destEndpointSize]
			}

			logging.LogDebug("Task: %v", task)
			taskChan.In() <- task
			atomic.AddInt32(&(*moverconfig.GetCurrState()).SumTasks, 1)
		}
	}
	return nil
}

//任务生成
func (p *MysqlPolicy) GenerateTask(dsp moverconfig.DataSourcePair, taskChan *channels.InfiniteChannel) error {
	if len(dsp.From.Endpoints) == 0 ||
		(dsp.DestTo == moverconfig.TargetDatabase && len(dsp.Dest.Endpoints) == 0) {
		return errors.Trace(errors.New("Ds no endpoint"))
	}

	breakpoint.FinishTask(GetNameByHostAndPort(dsp.From.Endpoints[0].Host, dsp.From.Endpoints[0].Port, dsp.From.DBName), dsp.From.TableName, nil)

	// get field & primary key
	fieldList, priKey, priKeyTyp, priKeyInfos, err := getFromFieldListMysql(dsp)
	if nil != err {
		return errors.Trace(err)
	}
	if "" == priKey {
		errMsg := fmt.Sprintf("Lost primary key in datasourcePair: %v", dsp)
		return errors.Trace(errors.New(errMsg))
	}

	if len(priKeyInfos) > 1 {
		selfPriKey := moverconfig.GetConfig().PriKey
		for _, v := range priKeyInfos {
			if v.PriKey == selfPriKey {
				priKey = v.PriKey
				priKeyTyp = v.PriKeyType
			}
		}
		//联合主键表
		mysqlDispatchTaskByPri(dsp, taskChan, fieldList, priKey, priKeyTyp)
	}
	//单一主键表
	mysqlDispatchTask(dsp, taskChan, fieldList, priKey, priKeyTyp)
	return nil
}

func mysqlReadDataByRange(task *moverconfig.Task) ([]moverconfig.RowData, error) {
	logging.LogDebug("Do task: %v", task)

	dataSizeLimit := moverconfig.GetConfig().DataSizeLimit

	if dataSizeLimit <= 100 {
		dataSizeLimit = 100
	} else if dataSizeLimit > task.RowCountLimit {
		//增加错误处理逻辑
		logging.LogError(logging.EC_Config_ERR, "Invalid config,DataSizelimit(%d) too big,It should be under(%d) ", dataSizeLimit, task.RowCountLimit)
		return nil, errors.New("Invalid config")
	}
	// get conn
	ep := task.FromEndpoint
	mysqlAddr := GetConnstr(ep, "", false)
	dbConn, err := getMysqlDB(mysqlAddr)
	if nil != err {
		return nil, err
	}

	// create sql
	if moverconfig.FDT_INT == task.PrimaryKeyType {
		// create sql
		query := fmt.Sprintf("SELECT `%s`,%s FROM `%s`.`%s` WHERE `%s`>%d and `%s`<=%d",
			task.PrimaryKey, task.FromField, task.FromDBName, task.FromTable, task.PrimaryKey, task.CurrId, task.PrimaryKey, task.EndId)
		if task.IsSelfPri {
			query = fmt.Sprintf("SELECT `%s`,%s FROM `%s`.`%s` WHERE `%s`>= %d and `%s`<=%d",
				task.PrimaryKey, task.FromField, task.FromDBName, task.FromTable, task.PrimaryKey, task.CurrId, task.PrimaryKey, task.EndId)
		}
		if "" != task.FromWhere {
			query += fmt.Sprintf(" AND (%s)", task.FromWhere)
		}
		if task.IsSelfPri && task.CurrId == task.EndId {
			task.StartSkip = true
			log.Debugf("==task.taskId=%v======task.currId=%v===task.EndId=%v", task.TaskID, task.CurrId, task.EndId)
			query += fmt.Sprintf(" ORDER BY `%s` LIMIT %d,%d", task.PrimaryKey, task.Skip, dataSizeLimit)
		} else {
			query += fmt.Sprintf(" ORDER BY `%s` LIMIT %d,%d", task.PrimaryKey, task.WinSize, dataSizeLimit)
		}

		// query
		logging.LogDebug("Query sql: %s", query)

		var rows *sql.Rows
		for i := 0; true; i++ {
			rows, err = dbConn.Query(query)
			if nil != err {
				sleep_time := rand.Intn(3000)
				logging.LogError(logging.EC_Ignore_ERR, "Retry, i=%d, sleep=%d,sql=%s, err: %v", i, sleep_time, query, err)
				time.Sleep(time.Duration(sleep_time) * time.Millisecond)
			} else {
				break
			}
		}
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		colNames, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		curId := task.CurrId //记录当前任务开始的ID值，用来和任务读取数据结束后的ID值比较，从而确定偏移量WInsize
		var priKeyValue int64
		cols := make([]interface{}, len(colNames)-1)
		colPtrs := make([]interface{}, len(colNames))
		colPtrs[0] = &priKeyValue
		for i := 1; i < len(colNames); i++ {
			colPtrs[i] = &cols[i-1]
		}

		// get data
		results := make([]moverconfig.RowData, 0)
		for rows.Next() {
			err = rows.Scan(colPtrs...)
			if err != nil {
				return nil, err
			}

			var rd moverconfig.RowData
			rd.Row = make([]interface{}, len(cols))
			for i, v := range cols {
				if nil != v {
					// convert every thing to string
					rd.Row[i] = utils.ByteToString(v.([]uint8))
				} else {
					rd.Row[i] = v
				}
			}
			results = append(results, rd)
			task.CurrId = priKeyValue
		}

		//利用类滑动窗口原理记录偏移量
		if task.CurrId == curId {
			task.WinSize += len(results)
		} else {
			task.WinSize = 0
		}

		if task.IsSelfPri && task.StartSkip {
			task.Skip += len(results)
		}

		return results, nil
	} else {

		operator := ">"
		// the 1st's task is closed interval "[]", the others is semi-closed interval "(]"
		if 1 == task.ClosedInterval {
			operator = ">="
			task.ClosedInterval = 0
		}
		if task.IsSelfPri {
			operator = ">="
		}
		// create sql
		//SELECT `id` FROM `t1` WHERE `id` >= '"1"' and `id` <= '"9"'ORDER BY `id` LIMIT 5;
		query := fmt.Sprintf("SELECT `%s`,%s FROM `%s`.`%s` WHERE `%s` %s '%s' and `%s` <= '%s'",
			task.PrimaryKey, task.FromField, task.FromDBName, task.FromTable, task.PrimaryKey, operator, task.StrCurrId, task.PrimaryKey, task.StrEndId)
		if "" != task.FromWhere {
			query += fmt.Sprintf(" AND (%s)", task.FromWhere)
		}
		if task.IsSelfPri && task.StrCurrId == task.StrEndId {
			task.StartSkip = true
			log.Debugf("==task.taskId=%v======task.currId=%v===task.EndId=%v", task.TaskID, task.StrCurrId, task.StrEndId)
			query += fmt.Sprintf(" ORDER BY `%s` LIMIT %d,%d", task.PrimaryKey, task.Skip, dataSizeLimit)
		} else {
			query += fmt.Sprintf(" ORDER BY `%s` LIMIT %d, %d", task.PrimaryKey, task.WinSize, dataSizeLimit)
		}

		// query
		logging.LogDebug("Query sql: %s", query)

		var rows *sql.Rows
		for i := 0; true; i++ {
			rows, err = dbConn.Query(query)
			if nil != err {
				sleep_time := rand.Intn(3000)
				logging.LogError(logging.EC_Ignore_ERR, "Retry, i=%d, sleep=%d,sql=%s, err: %v", i, sleep_time, query, err)
				time.Sleep(time.Duration(sleep_time) * time.Millisecond)
			} else {
				break
			}
		}
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		colNames, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		//局部偏移ID
		Localoffset := task.StrCurrId
		batch := true
		var currentBatchId string
		var priKeyValue string
		cols := make([]interface{}, len(colNames)-1)
		colPtrs := make([]interface{}, len(colNames))
		colPtrs[0] = &priKeyValue
		for i := 1; i < len(colNames); i++ {
			colPtrs[i] = &cols[i-1]
		}

		// get data
		results := make([]moverconfig.RowData, 0)
		for rows.Next() {
			err = rows.Scan(colPtrs...)
			if err != nil {
				return nil, err
			}
			var rd moverconfig.RowData
			rd.Row = make([]interface{}, len(cols))
			for i, v := range cols {
				if nil != v {
					// convert every thing to string
					rd.Row[i] = utils.ByteToString(v.([]uint8))
				} else {
					rd.Row[i] = v
				}
			}
			results = append(results, rd)
			task.StrCurrId = priKeyValue
			if batch == true {
				//本读取批次的开始ID
				currentBatchId = priKeyValue
				batch = false
			}
			//task.strCurrId = Stringescape(priKeyValue)
		}
		if strings.Compare(task.StrCurrId, Localoffset) == 0 && strings.Compare(task.StrCurrId, currentBatchId) == 0 {
			task.WinSize += len(results)
		} else if strings.Compare(task.StrCurrId, currentBatchId) == 0 {
			task.WinSize = len(results)
		} else {
			task.WinSize = 0
		}
		if task.IsSelfPri && task.StartSkip {
			task.Skip += len(results)
		}

		return results, nil
	}
}

//8、读数据派发逻辑
func (p *MysqlPolicy) ReadData(task *moverconfig.Task) ([]moverconfig.RowData, error) {
	return mysqlReadDataByRange(task)
}

func (p *MysqlPolicy) WriteData(task *moverconfig.Task, datas []moverconfig.RowData) error {
	if nil == datas || len(datas) == 0 {
		return nil
	}

	writeMode := "INSERT"
	if 1 == moverconfig.GetConfig().WriteMode {
		writeMode = "REPLACE"
	}

	// get field
	var ds moverconfig.DataSource
	ds.DBType = task.DestDBType
	ds.IsSharding = task.DestDBIsSharding
	ds.DBName = task.DestDBName
	ds.TableName = task.DestTable
	ds.Endpoints = make([]moverconfig.Endpoint, 1)
	ds.Endpoints[0] = task.DestEndpoint
	fieldInfo, _, err := p.GetFields(ds, "", IsSelfDef)
	if nil != err {
		return err
	}

	// get date field not null
	checkNotNullDateType := false
	notNullDateTypeField := make([]int, 0)
	for i, fieldName := range task.DestField {
		info, ok := fieldInfo[fieldName] //拿出描述该字段的信息
		if !ok {
			logging.LogError(logging.EC_Runtime_ERR, "Lost field info in dest, ignore, db: %s, table: %s, field: %s", task.DestDBName, task.DestTable, fieldName)
			continue
		}
		if moverconfig.FDT_DATE == info.FieldType && !info.IsNullAble {
			checkNotNullDateType = true
			notNullDateTypeField = append(notNullDateTypeField, i)
		}
	}

	// create sql
	if task.DestDBIsSharding {
		// raw string
		if len(datas) != 1 {
			return errors.New("Dest is sharding, can't batch insert")
		}

		values := datas[0]
		fieldString := "`" + strings.Join(task.DestField, "`,`") + "`"
		var sql string
		sql = fmt.Sprintf("%s INTO `%s`.`%s`(%s) VALUES(", writeMode, task.DestDBName, task.DestTable, fieldString)
		for i, v := range values.Row {
			if 0 != i {
				sql += ", "
			}
			switch t := v.(type) {
			case int64:
				/*
					mssql: tinyint, smallint, int, bigint
				*/
				sql += fmt.Sprintf("%d", v)
			case float64:
				/*
					mssql: float, real
				*/
				sql += fmt.Sprintf("%f", v)
			case bool:
				/*
					mssql: bit
				*/
				if true == v.(bool) {
					sql += "\"1\""
				} else {
					sql += "\"0\""
				}
			case string:
				/*
					mssql:
						char, varchar, text, nchar, nvarchar, ntext
				*/
				sql += fmt.Sprintf("%s", strconv.Quote(v.(string)))
			case []uint8:
				// binary
				/*
					mssql:
						decimal, numeric, smallmoney, money
						binary, varbinary, images, timestamp (can't use those type when mssql -> mysql/sharding)
				*/
				sql += fmt.Sprintf("%s", strconv.Quote(utils.ByteToString(v.([]uint8))))
			case time.Time:
				/*
					mssql:
						date, datetime, datetime2, datetimeoffset, smalldatetime, time
				*/
				sql += fmt.Sprintf("%s", strconv.Quote(v.(time.Time).Format("2006-01-02 15:04:05")))
			case nil:
				sql += "NULL"
			default:
				errMsg := fmt.Sprintf("Error value type(%v), task: %v, data: %v", t, *task, datas)
				return errors.New(errMsg)
			}
		}
		sql += ")"

		// get conn
		ep := task.DestEndpoint
		//mysqlAddr := getMysqlConnstr(ep, "")
		mysqlAddr := GetConnstr(ep, "", true)
		dbConn, err := getMysqlDB(mysqlAddr)
		if nil != err {
			return err
		}

		// execute
		logging.LogDebug("Sql: %s", sql)

		for i := 0; true; i++ {
			_, err = dbConn.Exec(sql)
			if nil != err {
				mysqlAddr := GetConnstr(task.FromEndpoint, "", task.DestDBIsSharding)
				sleep_time := rand.Intn(500) + 1
				logging.LogError(logging.EC_Ignore_ERR, "Retry, i=%d, sleep=%d, sql=%s, fromMysqlAddr=%s,err: %v", i, sleep_time, sql, mysqlAddr, err)

				time.Sleep(time.Duration(sleep_time) * time.Millisecond)

				sql = strings.Replace(sql, "INSERT", "REPLACE", 1)
			} else {
				break
			}
		}
		if nil != err {
			return err
		}
	} else {
		//将datas分批次插入，首先将datas分批次
		size := 65535 / (len(task.DestField) * 3)
		var end int
		for start := 0; start < len(datas); start += size {
			end += size
			if end > len(datas) {
				end = len(datas)
			}

			var placeHolder string
			for i := 0; i < len(datas[0].Row); i++ {
				if "" == placeHolder {
					placeHolder = "?"
				} else {
					placeHolder += ", ?"
				}
			}
			values := make([]interface{}, 0, len(task.DestField)*len(datas[start:end]))
			fieldString := "`" + strings.Join(task.DestField, "`,`") + "`"

			sql := fmt.Sprintf("%s INTO `%s`.`%s`(%s) VALUES", writeMode, task.DestDBName, task.DestTable, fieldString)
			for i := 0; i < len(datas[start:end]); i++ {
				if 0 == i {
					sql += fmt.Sprintf("(%s)", placeHolder)
				} else {
					sql += fmt.Sprintf(", (%s)", placeHolder)
				}
				if checkNotNullDateType {
					for _, ofst := range notNullDateTypeField {
						if nil == datas[start:end][i].Row[ofst] {
							logging.LogError(logging.EC_Runtime_ERR, "Source has NULL datetime, but dest can't be NULL, db=%s, table=%s, field=%d, row=%v", task.DestDBName, task.DestTable, ofst, datas[i])
							datas[start:end][i].Row[ofst] = 0
						}
					}
				}
				values = append(values, datas[start:end][i].Row...)
			}
			// get conn
			ep := task.DestEndpoint
			//mysqlAddr := getMysqlConnstr(ep, "")
			mysqlAddr := GetConnstr(ep, "", false)
			dbConn, err := getMysqlDB(mysqlAddr)
			if nil != err {
				return err
			}

			// execute
			logging.LogDebug("Sql: %s", sql)
			for i := 0; true; i++ {
				_, err = dbConn.Exec(sql, values...)
				if nil != err {
					mysqlAddr := GetConnstr(task.FromEndpoint, "", task.DestDBIsSharding)
					sleep_time := rand.Intn(500)
					logging.LogError(logging.EC_Ignore_ERR, "Retry, i=%d, sleep=%d, sql=%s,fromMysqlAddr=%s,values=%v,err: %v", i, sleep_time, sql, mysqlAddr, values, err)
					time.Sleep(time.Duration(sleep_time) * time.Millisecond)
					sql = strings.Replace(sql, "INSERT", "REPLACE", 1)
				} else {
					break
				}
			}
			if nil != err {
				return err
			}

		}

	}

	return nil
}

func showMysqlDB(dbConn *sql.DB) ([]string, error) {
	sql := "show databases"
	logging.LogDebug("Sql: %s", sql)

	rows, err := dbConn.Query(sql)
	if nil != err {
		return nil, errors.Trace(err)
	}
	defer rows.Close()

	dbNameList := make([]string, 0)
	for rows.Next() {
		var dbName string
		err = rows.Scan(&dbName)
		if nil != err {
			return nil, errors.Trace(err)
		}
		dbNameList = append(dbNameList, dbName)
	}

	return dbNameList, nil
}

func showMysqlTable(dbConn *sql.DB, dbName string) ([]string, error) {
	sql := fmt.Sprintf("use %s", dbName)
	dbConn.Exec(sql)

	sql = fmt.Sprintf("show tables from %s", dbName)
	logging.LogDebug("Sql: %s", sql)

	rows, err := dbConn.Query(sql)
	if nil != err {
		logging.LogError(logging.EC_Runtime_ERR, "Fail to query sql, sql=%s, err: %v", sql, err)
		return nil, errors.Trace(err)
	}
	defer rows.Close()

	tableNameList := make([]string, 0)
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "Fail to scan row, sql=%s, err: %v", sql, err)
			return nil, errors.Trace(err)
		}
		tableNameList = append(tableNameList, tableName)
	}

	return tableNameList, nil
}

func getMysqlDataType(dataType string, isSelf bool) byte {
	if isSelf {
		return getMysqlDataType4Self(dataType)
	} else {
		return getMysqlDataType4Real(dataType)
	}
}

func getMysqlDataType4Real(dataType string) byte {
	// type: int/float/string/time
	dataType = strings.ToUpper(dataType)
	switch dataType {
	// integer
	case "TINYINT":
		{
			return moverconfig.MYSQL_TYPE_TINY
		}
	case "SMALLINT":
		{
			return moverconfig.MYSQL_TYPE_SHORT
		}
	case "MEDIUMINT":
		{
			return moverconfig.MYSQL_TYPE_INT24
		}
	case "INT":
		{
			return moverconfig.MYSQL_TYPE_LONG
		}
	case "BIGINT":
		{
			return moverconfig.MYSQL_TYPE_LONGLONG
		}
	// float
	case "FLOAT":
		{
			return moverconfig.MYSQL_TYPE_FLOAT
		}
	case "DOUBLE":
		{
			return moverconfig.MYSQL_TYPE_DOUBLE
		}
	case "DECIMAL":
		{
			return moverconfig.MYSQL_TYPE_NEWDECIMAL
		}
	// string
	case "CHAR":
		{
			return moverconfig.MYSQL_TYPE_VARCHAR
		}
	case "VARCHAR":
		{
			return moverconfig.MYSQL_TYPE_VARCHAR
		}
	case "TINYTEXT":
		{
			return moverconfig.MYSQL_TYPE_VARCHAR
		}
	case "TEXT":
		{
			return moverconfig.MYSQL_TYPE_VARCHAR
		}
	case "MEDIUMTEXT":
		{
			return moverconfig.MYSQL_TYPE_VARCHAR
		}
	case "LONGTEXT":
		{
			return moverconfig.MYSQL_TYPE_VARCHAR
		}
	// binary
	case "TINYBLOB":
		{
			return moverconfig.MYSQL_TYPE_BLOB
		}
	case "MEDIUMBLOB":
		{
			return moverconfig.MYSQL_TYPE_BLOB
		}
	case "BLOB":
		{
			return moverconfig.MYSQL_TYPE_BLOB
		}
	case "LONGBLOB":
		{
			return moverconfig.MYSQL_TYPE_BLOB
		}
	case "VARBINARY":
		{
			return moverconfig.MYSQL_TYPE_BLOB
		}
	case "BINARY":
		{
			return moverconfig.MYSQL_TYPE_BLOB
		}
	case "BIT":
		{
			return moverconfig.MYSQL_TYPE_BLOB
		}
	// time
	case "DATE":
		{
			return moverconfig.MYSQL_TYPE_DATE
		}
	case "TIME":
		{
			return moverconfig.MYSQL_TYPE_TIME
		}
	case "YEAR":
		{
			return moverconfig.MYSQL_TYPE_YEAR
		}
	case "TIMESTAMP":
		{
			return moverconfig.MYSQL_TYPE_DATETIME
		}
	case "DATETIME":
		{
			return moverconfig.MYSQL_TYPE_DATETIME
		}
	// other
	case "ENUM":
		{
			return moverconfig.MYSQL_TYPE_VARCHAR
		}
	case "SET":
		{
			return moverconfig.MYSQL_TYPE_VARCHAR
		}
	default:
		{
			return 0
		}
	}
}

func getMysqlDataType4Self(dataType string) byte {
	// type: int/float/string/time
	dataType = strings.ToUpper(dataType)
	switch dataType {
	// integer
	case "TINYINT":
		{
			return moverconfig.FDT_INT
		}
	case "SMALLINT":
		{
			return moverconfig.FDT_INT
		}
	case "MEDIUMINT":
		{
			return moverconfig.FDT_INT
		}
	case "INT":
		{
			return moverconfig.FDT_INT
		}
	case "BIGINT":
		{
			return moverconfig.FDT_INT
		}
	// float
	case "FLOAT":
		{
			return moverconfig.FDT_FLOAT
		}
	case "DOUBLE":
		{
			return moverconfig.FDT_FLOAT
		}
	case "DECIMAL":
		{
			return moverconfig.FDT_FLOAT
		}
	// string
	case "CHAR":
		{
			return moverconfig.FDT_STRING
		}
	case "VARCHAR":
		{
			return moverconfig.FDT_STRING
		}
	case "TINYTEXT":
		{
			return moverconfig.FDT_STRING
		}
	case "TEXT":
		{
			return moverconfig.FDT_STRING
		}
	case "MEDIUMTEXT":
		{
			return moverconfig.FDT_STRING
		}
	case "LONGTEXT":
		{
			return moverconfig.FDT_STRING
		}
	// binary
	case "TINYBLOB":
		{
			return moverconfig.FDT_BINARY
		}
	case "MEDIUMBLOB":
		{
			return moverconfig.FDT_BINARY
		}
	case "BLOB":
		{
			return moverconfig.FDT_BINARY
		}
	case "LONGBLOB":
		{
			return moverconfig.FDT_BINARY
		}
	case "VARBINARY":
		{
			return moverconfig.FDT_BINARY
		}
	case "BINARY":
		{
			return moverconfig.FDT_BINARY
		}
	// time
	case "DATE":
		{
			return moverconfig.FDT_DATE
		}
	case "TIME":
		{
			return moverconfig.FDT_DATE
		}
	case "YEAR":
		{
			return moverconfig.FDT_DATE
		}
	case "TIMESTAMP":
		{
			return moverconfig.FDT_DATE
		}
	case "DATETIME":
		{
			return moverconfig.FDT_DATE
		}
	// other
	case "ENUM":
		{
			return moverconfig.FDT_STRING
		}
	case "SET":
		{
			return moverconfig.FDT_STRING
		}
	default:
		{
			return moverconfig.FieldDataTypeUnknown
		}
	}
}

func getMysqlDataType4SelfToMq(dataType string) byte {
	// type: int/float/string/time
	dataType = strings.ToUpper(dataType)
	switch dataType {
	// integer
	case "TINYINT":
		{
			return moverconfig.FDT_INT
		}
	case "SMALLINT":
		{
			return moverconfig.FDT_INT
		}
	case "MEDIUMINT":
		{
			return moverconfig.FDT_INT
		}
	case "INT":
		{
			return moverconfig.FDT_INT
		}
	case "BIGINT":
		{
			return moverconfig.FDT_INT
		}
	// float
	case "FLOAT":
		{
			return moverconfig.FDT_FLOAT
		}
	case "DOUBLE":
		{
			return moverconfig.FDT_FLOAT
		}
	case "DECIMAL":
		{
			return moverconfig.FDT_FLOAT
		}
	// string
	case "CHAR":
		{
			return moverconfig.FDT_STRING
		}
	case "VARCHAR":
		{
			return moverconfig.FDT_STRING
		}
	case "TINYTEXT":
		{
			return moverconfig.FDT_STRING
		}
	case "TEXT":
		{
			return moverconfig.FDT_STRING
		}
	case "MEDIUMTEXT":
		{
			return moverconfig.FDT_STRING
		}
	case "LONGTEXT":
		{
			return moverconfig.FDT_STRING
		}
	// binary
	case "BIT":
		{
			return moverconfig.FDT_BINARY
		}
	case "TINYBLOB":
		{
			return moverconfig.FDT_BINARY
		}
	case "MEDIUMBLOB":
		{
			return moverconfig.FDT_BINARY
		}
	case "BLOB":
		{
			return moverconfig.FDT_BINARY
		}
	case "LONGBLOB":
		{
			return moverconfig.FDT_BINARY
		}
	case "VARBINARY":
		{
			return moverconfig.FDT_BINARY
		}
	case "BINARY":
		{
			return moverconfig.FDT_BINARY
		}
	// time
	case "DATE":
		{
			return moverconfig.FDT_DATE
		}
	case "TIME":
		{
			return moverconfig.FDT_DATE
		}
	case "YEAR":
		{
			return moverconfig.FDT_INT
		}
	case "TIMESTAMP":
		{
			return moverconfig.FDT_DATE
		}
	case "DATETIME":
		{
			return moverconfig.FDT_DATE
		}
	// other
	case "ENUM":
		{
			return moverconfig.FDT_STRING
		}
	case "SET":
		{
			return moverconfig.FDT_STRING
		}
	default:
		{
			return moverconfig.FieldDataTypeUnknown
		}
	}
}

type PriKeyInfo struct {
	PriKey     string
	PriKeyType int
}

//改动1：由默认返回string类型的主键改为默认优先返回整形主键
func getFromFieldListMysql(dsp moverconfig.DataSourcePair) ([]string, string, int, []PriKeyInfo, error) {
	//ep := dsp.From.Endpoints[0]
	//mysqlAddr := getMysqlConnstr(ep, "")
	//dbConn, err := getMysqlDB(mysqlAddr)
	dbConn, _, err := getMysqlDB4Ds(dsp.From, "")
	if nil != err {
		return nil, "", moverconfig.FieldDataTypeUnknown, nil, errors.Trace(err)
	}

	var srcPriKey string
	var srcPriTyp int

	//新增
	var retSrcPriKey string
	var retSrcPriTyp int

	fieldList := make([]string, 0)
	priKeyInfos := []PriKeyInfo{}
	sql := fmt.Sprintf("SELECT COLUMN_NAME,DATA_TYPE,COLUMN_KEY,IS_NULLABLE FROM information_schema.`COLUMNS` WHERE TABLE_SCHEMA = '%s' and TABLE_NAME='%s'",
		dsp.From.DBName, dsp.From.TableName)
	logging.LogDebug("Sql: %s", sql)

	rows1, err := dbConn.Query(sql)
	if nil != err {
		return nil, "", moverconfig.FieldDataTypeUnknown, nil, errors.Trace(err)
	}
	defer rows1.Close()
	//新增
	var isFirst bool

	for rows1.Next() {
		var columnName, dataType, columnKey, isNullAble string
		err = rows1.Scan(&columnName, &dataType, &columnKey, &isNullAble)
		if nil != err {
			return nil, "", moverconfig.FieldDataTypeUnknown, nil, errors.Trace(err)
		}
		fieldList = append(fieldList, columnName)
		if "PRI" == columnKey {
			srcPriKey = columnName
			srcPriTyp = int(getMysqlDataType(dataType, true))
			//增加这部分逻辑：优先返回整型主键字段
			if moverconfig.FDT_INT == srcPriTyp && !isFirst {
				retSrcPriKey = columnName
				retSrcPriTyp = int(getMysqlDataType(dataType, true))
				isFirst = true
			}
			priKeyInfos = append(priKeyInfos, PriKeyInfo{
				srcPriKey,
				srcPriTyp,
			})
		}
	}

	if len(fieldList) == 0 {
		errMsg := fmt.Sprintf("table `%s`.`%s` is not exist", dsp.From.DBName, dsp.From.TableName)
		return nil, "", moverconfig.FieldDataTypeUnknown, nil, errors.Trace(errors.New(errMsg))
	}

	if "" != dsp.FromField {
		sql := fmt.Sprintf("SELECT %s FROM `%s`.`%s` LIMIT 1", dsp.FromField, dsp.From.DBName, dsp.From.TableName)
		logging.LogDebug("Sql: %s", sql)

		rows2, err := dbConn.Query(sql)
		if nil != err {
			return nil, "", moverconfig.FieldDataTypeUnknown, nil, errors.Trace(err)
		}
		defer rows2.Close()

		fieldList, err = rows2.Columns()
		if nil != err {
			return nil, "", moverconfig.FieldDataTypeUnknown, nil, errors.Trace(err)
		}
	}

	//return fieldList, srcPriKey, srcPriTyp, priKeyInfos, nil
	if retSrcPriKey == "" {
		return fieldList, srcPriKey, srcPriTyp, priKeyInfos, nil
	} else {
		return fieldList, retSrcPriKey, retSrcPriTyp, priKeyInfos, nil
	}
}

func GetConnstr(ep moverconfig.Endpoint, dbName string, isSharding bool) string {
	if isSharding {
		return getDbrouterConnstr(ep, dbName)
	} else {
		return getMysqlConnstr4Utf8mb4(ep, dbName)
	}
	//unreachable code
	return getMysqlConnstr(ep, dbName)
}

func getMysqlConnstr(ep moverconfig.Endpoint, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", ep.User, ep.Password, ep.Host, ep.Port, dbName)
}

func getMysqlConnstr4Utf8mb4(ep moverconfig.Endpoint, dbName string) string {
	// root:password@/name?charset=utf8mb4&collation=utf8mb4_unicode_ci
	// resolve emoji character
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci", ep.User, ep.Password, ep.Host, ep.Port, dbName)
}

func getDbrouterConnstr(ep moverconfig.Endpoint, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?interpolateParams=True", ep.User, ep.Password, ep.Host, ep.Port, dbName)
}

func getMysqlDB(dbMetaAddr string) (*sql.DB, error) {
	mysqlDbMutex.Lock()
	defer mysqlDbMutex.Unlock()
	//mysqlDbMap保存 【数据库信息+操作此数据库的对象】
	if nil == mysqlDbMap {
		mysqlDbMap = make(map[string]*sql.DB)
	}

	var mysqlDB *sql.DB
	mysqlDB, ok := mysqlDbMap[dbMetaAddr]
	if ok {
		return mysqlDB, nil
	}

	//	create connection
	var err error
	mysqlDB, err = sql.Open("mysql", dbMetaAddr)
	if nil != err {
		errMsg := fmt.Sprintf("Fail to open mysql(%s) err: %v", dbMetaAddr, err)
		logging.LogError(logging.EC_Datasource_ERR, "Get mysql db error, addr=%s, %s, err: %v", dbMetaAddr, errMsg, err)
		err = errors.New(errMsg)
		return nil, err
	}

	err = mysqlDB.Ping()
	if nil != err {
		errMsg := fmt.Sprintf("Fail to ping mysql(%s) err: %v", dbMetaAddr, err)
		logging.LogError(logging.EC_Datasource_ERR, "Ping mysql db error, addr=%s, %s, err: %v", dbMetaAddr, errMsg, err)
		err = errors.New(errMsg)
		return nil, err
	}

	workerSize := 1
	//SetMaxIdleConns设置空闲连接池中的最大连接数。
	mysqlDB.SetMaxIdleConns(workerSize)

	// save
	mysqlDbMap[dbMetaAddr] = mysqlDB

	return mysqlDB, nil
}

func isShardingDB(baseDBName string, dbName string) bool {
	baseDBNameLen := len(baseDBName)
	dbNameLen := len(dbName)
	if moverconfig.GetConfig().ShardingTable == 1 {
		if baseDBName == dbName {
			return true
		}
	}

	if dbNameLen < baseDBNameLen {
		return false
	}

	if dbNameLen == baseDBNameLen {
		if strings.ToUpper(baseDBName) != strings.ToUpper(dbName) {
			return false
		}

		// not match the same database name
		return false
	} else {
		if strings.ToUpper(dbName[:baseDBNameLen]) != strings.ToUpper(baseDBName) {
			return false
		}

		if baseDBNameLen+1 >= dbNameLen {
			return false
		}

		if dbName[baseDBNameLen] != '_' {
			return false
		}

		if !unicode.IsDigit(rune(dbName[baseDBNameLen+1])) {
			return false
		}
	}

	return true
}

func IsShardingTable(baseTableName string, tableName string) bool {
	baseTableNameLen := len(baseTableName)
	tableNameLen := len(tableName)

	//if moverconfig.GetConfig().ShardingTable == 1 {
	if baseTableName == tableName {
		return true
	}
	//}

	if tableNameLen < baseTableNameLen {
		return false
	}

	if tableNameLen == baseTableNameLen {
		if strings.ToUpper(baseTableName) != strings.ToUpper(tableName) {
			return false
		}

		// not match the same table name
		return false
	} else {
		if strings.ToUpper(tableName[:baseTableNameLen]) != strings.ToUpper(baseTableName) {
			return false
		}

		if tableNameLen < baseTableNameLen+3 {
			return false
		}

		// OrderPay match OrderPay_Detail sharding table
		// if tableName[baseTableNameLen] != '_' {
		if strings.ToUpper(tableName[baseTableNameLen:baseTableNameLen+3]) != "_S_" {
			return false
		}
	}

	return true
}
