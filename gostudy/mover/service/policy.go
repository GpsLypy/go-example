package service

import (
	"app/tools/mover/logging"
	"app/tools/mover/moverconfig"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/eapache/channels"
	"github.com/juju/errors"
)

type MoverPolicy interface {
	// check db alive
	CheckDataSource(ds moverconfig.DataSource) error
	// if table sharding, get all matched table belonged to this db instance
	// every table is a datasource
	GetShardingDataSource(ds moverconfig.DataSource, baseDbName string, baseTableName string) ([]moverconfig.DataSource, error)
	// get all tableName in DB
	GetTables(ds moverconfig.DataSource) ([]string, error)
	// get table fields
	// if user configure src fields, using user specified
	GetFields(ds moverconfig.DataSource, fields string, isUserDef bool) (map[string]moverconfig.FieldInfo, []moverconfig.FieldInfo, error)
	// get table stat
	GetTableStat(ds moverconfig.DataSource, where string) (moverconfig.TableStat, error)

	// get source primary key & fields
	// get source size/max id/min id
	// split job into small tasks
	// range: (StartId, EndId]
	GenerateTask(dsp moverconfig.DataSourcePair, taskChan *channels.InfiniteChannel) error
	ReadData(task *moverconfig.Task) ([]moverconfig.RowData, error)
	WriteData(task *moverconfig.Task, datas []moverconfig.RowData) error
}

func GetPolicy(dbType int) (MoverPolicy, error) {
	var policy MoverPolicy
	switch dbType {
	case moverconfig.DB_MYSQL:
		{
			policy = &MysqlPolicy{}
		}
	case moverconfig.DB_MSSQL:
		{
			policy = &MysqlPolicy{}
		}
	default:
		{
			errMsg := fmt.Sprintf("Unkown dbType(%d)", dbType)
			logging.LogError(logging.EC_Decode_ERR, "Fail to get policy, %s", errMsg)
			return nil, errors.New(errMsg)
		}
	}

	return policy, nil
}

func stringMatchInArr(str string, strPatArr []string) bool {
	for _, strPat := range strPatArr {
		if str == strPat || "*" == strPat {
			return true
		}
	}
	return false
}

func GetShardingEndpoints(dsPair moverconfig.DataSourcePair) ([]moverconfig.DataSource, error) {
	/* Get datasources from RDS source, if non-validate */
	if !dsPair.From.Validate() {
		return GetRdsShardingEndpoints(dsPair)
	}

	BackendEPs := dsPair.From.Endpoints

	// load FieldInfo into cache
	policy, err := GetPolicy(dsPair.From.DBType)
	if nil != err {
		return nil, errors.Trace(err)
	}

	if !dsPair.From.IsSharding {
		policy.GetFields(dsPair.From, dsPair.FromField, IsSelfDef)
	}

	newDatasources := make([]moverconfig.DataSource, 0)
	for _, endpoint := range BackendEPs {

		newDs := dsPair.From
		newDs.SetEndpoints([]moverconfig.Endpoint{endpoint})

		newDsList, err := GetShardignDataSourceUntilSuccess(policy, newDs, dsPair.From.DBName, dsPair.From.TableName)
		if nil != err {
			return nil, errors.Trace(err)
		}

		newDatasources = append(newDatasources, newDsList...)
	}

	return newDatasources, nil
}

func GetShardignDataSourceUntilSuccess(policy MoverPolicy, ds moverconfig.DataSource, baseDbName string, baseTableName string) ([]moverconfig.DataSource, error) {
	for {
		res, err := policy.GetShardingDataSource(ds, baseDbName, baseTableName)
		if nil != err {
			/* TODO: Check global context to quit */
			logging.LogError(logging.EC_Datasource_ERR, "Can't get sharding datasource: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		return res, nil
	}
}

func GetRdsShardingEndpoints(dsPair moverconfig.DataSourcePair) ([]moverconfig.DataSource, error) {
	// get endpoints from rds config
	endpointMatrix, err := moverconfig.GetRdsShardingBackendEndpoints()
	if nil != err {
		return nil, errors.Trace(err)
	}

	if nil == endpointMatrix || len(endpointMatrix) == 0 {
		return nil, errors.Trace(errors.Errorf("Wrong rds config, dnName=%s, no endpoints", dsPair.From.DBName))
	}

	// load FieldInfo into cache
	policy, err := GetPolicy(dsPair.From.DBType)
	if nil != err {
		return nil, errors.Trace(err)
	}
	if !dsPair.From.IsSharding {
		policy.GetFields(dsPair.From, dsPair.FromField, IsSelfDef)
	}

	newDatasources := make([]moverconfig.DataSource, 0)
	for _, endpointRow := range endpointMatrix {
		if len(endpointRow) == 0 {
			return nil, errors.New("Maybe exist empty instanceList element.")
		}

		newDs := dsPair.From
		newDs.SetEndpoints(endpointRow)

		newDsList, err := policy.GetShardingDataSource(newDs, dsPair.From.DBName, dsPair.From.TableName)
		if nil != err {
			return nil, errors.Trace(err)
		}

		for _, newDs := range newDsList {
			newDatasources = append(newDatasources, newDs)
		}
	}

	return newDatasources, nil
}

func FindErrorData(policy MoverPolicy, task *moverconfig.Task, datas []moverconfig.RowData) (bool, *moverconfig.RowData) {
	if nil == datas || len(datas) == 0 {
		return false, nil
	}

	size := len(datas)
	if 1 == size {
		err := policy.WriteData(task, datas)
		if nil != err {
			return true, &datas[0]
		}
		return false, nil
	} else {
		part1 := datas[:size/2]
		err := policy.WriteData(task, part1)
		if nil != err {
			return FindErrorData(policy, task, part1)
		}

		part2 := datas[size/2:]
		return FindErrorData(policy, task, part2)
	}
}

func IsSameFields(srcFiledMap map[string]moverconfig.FieldInfo, destFieldMap map[string]moverconfig.FieldInfo) bool {
	srcFileds := make([]string, 0)
	destFields := make([]string, 0)
	for k, _ := range srcFiledMap {
		srcFileds = append(srcFileds, strings.ToUpper(k))
	}
	for k, _ := range destFieldMap {
		destFields = append(destFields, strings.ToUpper(k))
	}

	sort.Strings(srcFileds)
	sort.Strings(destFields)
	srcLen := len(srcFileds)
	destLen := len(destFields)
	if srcLen > destLen {
		logging.LogError(logging.EC_Runtime_ERR, "Diff field failure, dest field less than src")
		return false
	}

	for k, v := range srcFiledMap {
		srcFiledMap[strings.ToUpper(k)] = v
	}

	for k, v := range destFieldMap {
		destFieldMap[strings.ToUpper(k)] = v
	}

	for i := 0; i < srcLen; i++ {
		field := srcFileds[i]
		_, ok := destFieldMap[field]
		if !ok {
			logging.LogError(logging.EC_Runtime_ERR, "Lost field in dest, field=%s", field)

			return false
		}

		srcFieldInfo := srcFiledMap[field]
		destFieldInfo := destFieldMap[field]
		//如果是未知字段类型为0
		if srcFieldInfo.FieldType == 0 || destFieldInfo.FieldType == 0 {
			logging.LogError(logging.EC_Runtime_ERR, "Diff field failure, field: %s, srcType: %d, destType: %d", field, srcFiledMap[field], destFieldMap[field])

			return false
		}

		if moverconfig.FieldDataTypeUnknown != srcFieldInfo.FieldType &&
			moverconfig.FieldDataTypeUnknown != destFieldInfo.FieldType &&
			srcFieldInfo.FieldType != destFieldInfo.FieldType {
			logging.LogError(logging.EC_Runtime_ERR, "Diff field failure, field: %s, srcType: %d, destType: %d", field, srcFieldInfo.FieldType, destFieldInfo.FieldType)

			return false
		}
	}

	return true
}

func isZeroDate(dbType int, date string) bool {
	if "0000-00-00 00:00:00" == date || "0000-00-00" == date || "00:00:00" == date || "0" == date {
		return true
	}
	return false
}

func checkFieldType(fieldMap map[string]moverconfig.FieldInfo) error {
	for k, v := range fieldMap {
		if moverconfig.FDT_BINARY == v.FieldType {
			return errors.New(fmt.Sprintf("Field '%s' is binary", k))
		}
	}
	return nil
}
