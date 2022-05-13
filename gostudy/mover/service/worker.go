package service

import (
	"app/platform/jobmgr/logger"
	"app/tools/mover/breakpoint"
	"app/tools/mover/info"
	"app/tools/mover/logging"
	"app/tools/mover/moverconfig"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cihub/seelog"
	"github.com/eapache/channels"
	"github.com/juju/errors"
)

var (
	jobId          int64
	RecvStopped    bool
	stopFlag       bool
	stopWorkerFlag bool
	stopWorkerCh   chan struct{}
	stopWorkerWg   sync.WaitGroup
	starting       bool
	workerMutex    sync.Mutex
	doingTask      int32
	sumTask        int
	NumWorkers     int32
	tableMutex     sync.Mutex
	statMutex      sync.Mutex
	StatsMap       map[string]int
	TableMap       map[string]int
	StatsArr       []moverconfig.TableStat
	IsSelfDef      bool
)
var OfflineWorkerChOnce = sync.Once{}
var OfflineWorkerCh chan struct{}
var StopJobChan chan struct{}
var StopGenWg sync.WaitGroup
var Once = sync.Once{}

func init() {
	StopJobChan = make(chan struct{})
	OfflineWorkerCh = make(chan struct{})
}

// can't re-entry
func StartMover(jid int64) (errRet error) {
	if starting {
		// job has been started
		logging.LogWarn("Job is in progress")
		return nil
	}

	workerMutex.Lock()
	defer func() {
		if nil != errRet {
			starting = false
			(*moverconfig.GetCurrState()).IsHealthy = false
			(*moverconfig.GetCurrState()).ErrMsg = "[" + errRet.Error() + " ! ]"
			(*moverconfig.GetCurrState()).Started = false
		}
		workerMutex.Unlock()
	}()

	if starting {
		// job has been started
		// can't return error, will cause to start job twice
		logging.LogWarn("Job is in progress")

		return nil
	}

	starting = true
	jobId = jid

	if nil == stopWorkerCh {
		stopWorkerCh = make(chan struct{})
	}

	config := moverconfig.GetConfig()
	if config.Target == moverconfig.TargetDatabase {
		IsSelfDef = true
	} else {
		IsSelfDef = false
	}

	//如果需要订阅至MQ
	if config.Target == moverconfig.TargetTurboMQ {
		if err := initSlaveManager(config); nil != err {
			seelog.Errorf("init slave manager, envInfo: %v, dispatcher: %v, err: %v", *info.GetEnvInfo(), config.Dispatcher, err)
			return err
		}
	}

	//遍历存放任务的切片
	for i, task := range config.TaskList {
		//检查数据源中服务节点的信息是否有效
		if !task.From.Validate() {
			//获取rds数据源中的服务节点信息
			endpoints, err := moverconfig.GetRdsDataSourceEndpoints()
			if nil != err {
				seelog.Errorf("Get rds datasource error, envInfo: %v, task: %v, err: %v", *info.GetEnvInfo(), config.TaskList[0], err)
				return err
			}
			srcDs := &config.TaskList[i].From
			srcDs.SetEndpoints(endpoints)
		}
	}

	/*if config.IsNeedRdsSources() {
		endpoints, err := getRdsDataSourceEndpoints()
		if nil != err {
			seelog.Errorf("Get rds datasource error, envInfo: %v, task: %v, err: %v", envInfo, config.TaskList[0], err)
			return
		}

		// Tip: value reference lead to fatal error
		for i, _ := range config.TaskList {
			srcDs := &config.TaskList[i].From
			srcDs.SetEndpoints(endpoints)
		}
	}*/

	// bring forward 10 seconds
	//MQ相关
	dumptime := uint32(0)
	if config.FixedExeTime == 1 {
		dumptime = uint32(time.Now().Unix() - 10)
	}

	// get policy
	policy, err := GetPolicy(config.TaskList[0].From.DBType)
	if nil != err {
		seelog.Errorf("Get policy err, envInfo: %v, task: %v, err: %v", *info.GetEnvInfo(), config.TaskList[0], err)
		return err
	}

	logging.LogInfo("start check job...")
	if err := checkJob(); nil != err {
		logging.LogError(logging.EC_Job_Run_ERR, "Check job fail, err: %v", err)
		return err
	}

	// start slave, if needed.
	if config.Target == moverconfig.TargetTurboMQ {
		if err := SlaveMgr.Start(); nil != err {
			logging.LogError(logging.EC_Job_Run_ERR, "Start dispatcher fail, err: %v", err)
			return err
		}
	}

	// generate task
	logging.LogInfo("start generate task...")
	taskChan := channels.NewInfiniteChannel()
	//generateTask(taskChan)
	go generateTask(taskChan)
	// 1. wait 30 seconds for generate tasks
	// 2. wait first task etc.
	for ttl := 30; ttl > 0; ttl-- {
		if (*moverconfig.GetCurrState()).ErrMsg != "" {
			seelog.Errorf("Datasource error, envInfo: %v, err: %v", *info.GetEnvInfo(), (*moverconfig.GetCurrState()).ErrMsg)
			return err
		}
		if !(*moverconfig.GetCurrState()).Prepared &&
			(*moverconfig.GetCurrState()).IsHealthy &&
			(*moverconfig.GetCurrState()).SumTasks == 0 {
			ttl = 30
		}

		time.Sleep(time.Second)
	}

	//判断是否有任务并且任务是否准备好
	if (*moverconfig.GetCurrState()).SumTasks > 0 && 0 == config.Prepare {
		// start worker routine
		logging.LogInfo("start worker routine...")

		dataChan := make(chan *moverconfig.RowInfo, 10000)
		for i := 0; i < config.WorkerSize; i++ {
			if NumWorkers == -1 {
				NumWorkers += 2
			} else {
				NumWorkers = NumWorkers + 1
			}
			stopWorkerWg.Add(1)
			go worker(i, taskChan, dataChan)
			//worker(i, taskChan, dataChan)
		}

		//上面的worker去处理任务，如果是数任务是数据库表同步，datachan里面不放数据，如果发现任务是订阅任务，则datachan会有数据
		//执行到这里后，判断是否需要订阅消息，是的话取出datachan里面的数据丢到mq里
		if config.Target == moverconfig.TargetTurboMQ {
			err = DispatchJob(policy, dataChan, dumptime)
			if nil != err {
				// exit successfully
				errRet = err
				return
			}

			// wait workers to quit successfully
			time.Sleep(time.Second * 5)
			(*moverconfig.GetCurrState()).ProgressRatio = 10000
			SlaveMgr.Finished = true
			UploadStatusInfo("")
		}
	}

	logging.LogInfo("mover finished.")

	(*moverconfig.GetCurrState()).IsHealthy = true
	(*moverconfig.GetCurrState()).ErrMsg = ""
	(*moverconfig.GetCurrState()).Started = true

	if (*moverconfig.GetCurrState()).SumTasks == 0 {
		SlaveMgr.Finished = true
		(*moverconfig.GetCurrState()).ProgressRatio = 10000
		UploadStatusInfo("No task")
	}

	return nil
}

func Stopworkers() error {
	if stopWorkerFlag {
		return nil
	}

	logging.LogInfo("Workers are going to stop.")
	close(stopWorkerCh)
	stopWorkerWg.Wait()
	stopWorkerFlag = true

	return nil
}

func StopMover() error {
	logging.LogInfo("Job mover ready to stop")

	stopFlag = true
	if !RecvStopped {
		for times := 10; (*moverconfig.GetCurrState()).ProgressRatio == 10000 && times > 0; times-- {
			if err := os.Remove(breakpoint.BpFile); nil != err {
				logging.LogError(logging.EC_Runtime_ERR, "remove breakpoint.json, err: %v", err)
				(*moverconfig.GetCurrState()).IsHealthy = false
				(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("remove breakpoint.json, err: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			logging.LogInfo("remove breakpoint.json=%s", breakpoint.BpFile)
			break
		}
	}
	// purge breakpoint.json file, if finish

	logging.LogInfo("Job mover stopped")

	return nil
}

/* append tableInfo to tableArr and set it's index */
func appendTableStat(stat moverconfig.TableStat) error {
	tableMutex.Lock()
	defer tableMutex.Unlock()

	key := fmt.Sprintf("%s:%s", stat.DbName, stat.TableName)
	_, ok := TableMap[key]
	if ok {
		return nil
	}

	idx := len(TableMap)
	// StatsArr = append(StatsArr, stat)
	TableMap[key] = idx

	return nil
}

func ResetAppendTableStat(stat moverconfig.TableStat) error {
	tableMutex.Lock()
	defer tableMutex.Unlock()

	key := fmt.Sprintf("%s:%s:%s", stat.HostPort, stat.DbName, stat.TableName)
	_, ok := StatsMap[key]
	if ok {
		return nil
	}

	idx := len(StatsMap)
	StatsArr = append(StatsArr, stat)
	StatsMap[key] = idx

	return nil
}

/*
  1. check datasource alive
  2. check fields
  3. collect tablestats
*/
func checkDataSource(dsp moverconfig.DataSourcePair) error {
	logging.LogInfo("Check datasourcePair: %v", dsp)

	var srcPolicy, dstPolicy MoverPolicy
	var err error
	srcPolicy, err = GetPolicy(dsp.From.DBType)
	if nil != err {
		logging.LogError(logging.EC_Datasource_ERR, "Fail to get policy, err: %v", err)
		return err
	}

	// check datasource
	for {
		err = srcPolicy.CheckDataSource(dsp.From)
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "Fail to check data source, err: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

	// get src fields
	srcFieldMap := make(map[string]moverconfig.FieldInfo)
	//从数据库中获取所有字段保存在临时map中
	srcCheckMap, _, err := srcPolicy.GetFields(dsp.From, "", true)
	if nil != err {
		logging.LogError(logging.EC_Runtime_ERR, "Fail to get fields, err:%v", err)
		return err
	}
	//若有需要过滤的字段
	if dsp.FromField != "" {
		fromFiels := strings.Split(dsp.FromField, ",")
		//检查过滤字段是否有效
		for _, v := range fromFiels {
			if _, ok := srcCheckMap[v]; !ok {
				logging.LogError(logging.EC_Runtime_ERR, "Fail to get fromFields ! err:%v  or  Table %v  doesn't exist ", v, dsp.From.TableName)
				//return errors.New("Fail to get fromFields=" + v)
				return errors.New((*moverconfig.GetCurrState()).ErrMsg)
			}
		}
	}
	//srcFieldMap保存最终的字段信息
	srcFieldMap, _, err = srcPolicy.GetFields(dsp.From, dsp.FromField, true)
	if nil != err {
		logging.LogError(logging.EC_Runtime_ERR, "Fail to get fields, err:%v", err)
		return err
	}

	if "" == dsp.FromField && dsp.DestTo == moverconfig.TargetDatabase && dsp.Dest.IsSharding {
		err = checkFieldType(srcFieldMap)
		if nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "Fail to check field, err:%v", err)
			return err
		}
	}
	//目标端是数据库
	if dsp.DestTo == moverconfig.TargetDatabase {
		// dest
		dstPolicy, err = GetPolicy(dsp.Dest.DBType)
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "Fail to get policy, err:%v", err)
			return err
		}

		// check datasource
		err = dstPolicy.CheckDataSource(dsp.Dest)
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "Fail to check data source, err:%v", err)
			return err
		}

		// get fields
		destFieldMap := make(map[string]moverconfig.FieldInfo)
		destFieldMap, _, err = dstPolicy.GetFields(dsp.Dest, "", true)
		if nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "Fail to get fields, err:%v", err)
			return err
		}

		// check different field
		if !IsSameFields(srcFieldMap, destFieldMap) {
			errMsg := fmt.Sprintf("Different fields in task When checkDataSource, src.dbName=%s, src.tableName=%s, dest.dbName=%s, dest.tableName=%s",
				dsp.From.DBName, dsp.From.TableName, dsp.Dest.DBName, dsp.Dest.TableName)
			logging.LogError(logging.EC_Runtime_ERR, "Has diff fields, err: %s", errMsg)
			return errors.New(errMsg)
		}
	}

	// get tablestat & append tablestat
	logging.LogInfo("Get table information %v", dsp.From)
	//dsp.FromWhere为过滤条件
	//stat包含了表中有多少条数据
	stat, err := srcPolicy.GetTableStat(dsp.From, dsp.FromWhere)
	if nil != err {
		logging.LogError(logging.EC_Datasource_ERR, "Fail to get table info, err:%v", err)
		return err
	}
	//StatsArr = append(StatsArr, stat)
	//StatsMap[key] = idx
	appendTableStat(stat)

	return nil
}

func checkJob() error {
	var err error
	config := moverconfig.GetConfig()

	if nil == config {
		return errors.New("Config nil")
	}

	StatsMap = make(map[string]int)
	TableMap = make(map[string]int)

	allListedTable := make(map[string]int)
	for _, task := range config.TaskList {
		if "*" == task.From.TableName {
			continue
		}
		//配置中明确指定了表名，记录下来
		tableKey := fmt.Sprintf("%s:%s", task.From.DBName, task.From.TableName)
		allListedTable[tableKey] = 1
	}

	checkedTable := make(map[string]int)
	for _, task := range config.TaskList {
		logging.LogInfo("Check task: %v", task)
		taskList := make([]moverconfig.DataSourcePair, 0)
		if task.From.IsSharding {
			dataSource, err := GetShardingEndpoints(task)
			if err != nil {
				logging.LogError(logging.EC_Datasource_ERR, "Fail to GetShardingEndpoints, %v", err)
				(*moverconfig.GetCurrState()).IsHealthy = false
				(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Fail to GetShardingEndpoints, %v", err)
				return errors.Trace(err)
			}
			for _, from := range dataSource {
				newTasks, err := getTask(task, from, allListedTable, checkedTable)
				if err != nil {
					logging.LogError(logging.EC_Datasource_ERR, "Fail to getTask, %v", err)
					(*moverconfig.GetCurrState()).IsHealthy = false
					(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Fail to getTask, %v", err)
					continue
				}
				if len(newTasks) > 0 {
					taskList = append(taskList, newTasks...)
				}
			}
		} else {
			newTasks, err := getTask(task, task.From, allListedTable, checkedTable)
			if err != nil {
				logging.LogError(logging.EC_Datasource_ERR, "Fail to get policy, %v", err)
				return err
			}
			if len(newTasks) > 0 {
				taskList = append(taskList, newTasks...)
			}
		}

		for _, dsp := range taskList {
			err = checkDataSource(dsp)
			if nil != err {
				logging.LogError(logging.EC_Runtime_ERR, "Fail to check datasource, dsp: %v, err: %v", dsp, err)
				return err
			}
		}
	}

	if len(TableMap) <= 0 {
		return errors.New("no table")
	}

	return nil
}

func getTask(task moverconfig.DataSourcePair, taskFrom moverconfig.DataSource, allListedTable, checkedTable map[string]int) ([]moverconfig.DataSourcePair, error) {
	taskList := []moverconfig.DataSourcePair{}
	//无效
	task.From = taskFrom
	if "*" == task.From.TableName {
		policy, err := GetPolicy(task.From.DBType)
		if nil != err {
			logging.LogError(logging.EC_Datasource_ERR, "Fail to get policy, %v", err)
			return nil, errors.Trace(err)
		}

		tables, err := policy.GetTables(task.From)
		if nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "Fail to get tables, %v", err)
			return nil, errors.Trace(err)
		}

		for _, tableName := range tables {
			tableKey := fmt.Sprintf("%s:%s", task.From.DBName, tableName)
			_, ok := allListedTable[tableKey]
			if ok {
				continue
			}
			_, ok = checkedTable[tableKey]
			if ok {
				continue
			}
			checkedTable[tableKey] = 1

			newTask := task
			newTask.From.TableName = tableName
			if "" == newTask.Dest.DBName {
				newTask.Dest.DBName = newTask.From.DBName
			}
			if "" == newTask.Dest.TableName {
				newTask.Dest.TableName = newTask.From.TableName
			}
			taskList = append(taskList, newTask)
		}
	} else {
		tableKey := fmt.Sprintf("%s:%s", task.From.DBName, task.From.TableName)
		_, ok := checkedTable[tableKey]
		if ok {
			return nil, nil
		}
		checkedTable[tableKey] = 1

		if "" == task.Dest.DBName {
			task.Dest.DBName = task.From.DBName
		}
		if "" == task.Dest.TableName {
			task.Dest.TableName = task.From.TableName
		}
		taskList = append(taskList, task)
	}
	return taskList, nil
}

func generateNormalTask(dsPair moverconfig.DataSourcePair, tableName string, policy MoverPolicy, ch *channels.InfiniteChannel) error {
	var err error

	StopGenWg.Add(1)
	defer StopGenWg.Done()

	newDsPair := dsPair
	newDsPair.From.TableName = tableName
	newDsPair.From.OrgDBName = newDsPair.From.DBName
	newDsPair.From.OrgTbName = newDsPair.From.TableName
	if "" == newDsPair.Dest.DBName {
		newDsPair.Dest.DBName = newDsPair.From.DBName
	}
	if "" == newDsPair.Dest.TableName {
		newDsPair.Dest.TableName = newDsPair.From.TableName
	}

	if dsPair.From.IsSharding {
		dataSources, err := GetShardingEndpoints(newDsPair)

		if nil != err {
			(*moverconfig.GetCurrState()).ErrMsg = errors.Details(err) + "<<generateNormalTask()::GetShardingEndpoints() filed : Please check whether the source is shard or ShardingEndpoints IsEnabled>>"
			return err
		}

		for _, ds := range dataSources {
			logger.Infof("=====sharding-last-datasources==%v,db=%s,table=%s,oriDB=%s,oriTable=%s", ds.Endpoints, ds.DBName, ds.TableName, ds.OrgDBName, ds.OrgTbName)
			newDsPair.From = ds

			err = policy.GenerateTask(newDsPair, ch)
			if nil != err {
				(*moverconfig.GetCurrState()).ErrMsg = errors.Details(err) + "<<In the case of From.IsSharding==true,generateNormalTask()::GenerateTask() filed>>"
				return err
			}
		}
	} else {
		err = policy.GenerateTask(newDsPair, ch)
		if nil != err {
			(*moverconfig.GetCurrState()).ErrMsg = errors.Details(err) + "<<In the case of From.IsSharding==false,generateNormalTask()::GenerateTask() filed>>"
			return err
		}
	}

	return err
}

func generateRdsTask(dsPair moverconfig.DataSourcePair, tableName string, policy MoverPolicy, ch *channels.InfiniteChannel) error {
	var err error

	StopGenWg.Add(1)
	defer StopGenWg.Done()

	newDsPair := dsPair
	newDsPair.From.TableName = tableName
	if "" == newDsPair.Dest.DBName {
		newDsPair.Dest.DBName = newDsPair.From.DBName
	}
	if "" == newDsPair.Dest.TableName {
		newDsPair.Dest.TableName = newDsPair.From.TableName
	}

	if dsPair.From.IsSharding {
		dataSources, err := GetRdsShardingEndpoints(newDsPair)
		if nil != err {
			(*moverconfig.GetCurrState()).ErrMsg = errors.Details(err) + "<<generateRdsTask()::GetRdsShardingEndpoints() filed : Please check whether the source is shard or ShardingEndpoints IsEnabled>>"
			return err
		}

		for _, ds := range dataSources {
			newDsPair.From = ds

			err = policy.GenerateTask(newDsPair, ch)
			if nil != err {
				(*moverconfig.GetCurrState()).ErrMsg = errors.Details(err) + "<<In the case of From.IsSharding==true,generateRdsTask::GenerateTask() filed>>"
				return err
			}
		}
	} else {
		err = policy.GenerateTask(newDsPair, ch)
		if nil != err {
			(*moverconfig.GetCurrState()).ErrMsg = errors.Details(err) + "<<In the case of From.IsSharding==false,generateRdsTask::GenerateTask() filed>>"
			return err
		}
	}

	return err
}

func generateTask(taskChan *channels.InfiniteChannel) error {
	config := moverconfig.GetConfig()

	if nil == config {
		(*moverconfig.GetCurrState()).ErrMsg = "Config(MoverConfig) nil"
		return errors.New("Config(MoverConfig) nil")
	}
	//统计所有taskList中源端的表
	allListedTable := make(map[string]int)
	for _, taskInfo := range config.TaskList {
		if "*" == taskInfo.From.TableName {
			continue
		}
		tableKey := fmt.Sprintf("%s:%s", taskInfo.From.DBName, taskInfo.From.TableName)
		allListedTable[tableKey] = 1
	}

	/* multi-Task channel slice */
	processedTable := make(map[string]int)
	for _, taskInfo := range config.TaskList {
		policy, err := GetPolicy(taskInfo.From.DBType)
		if nil != err {
			(*moverconfig.GetCurrState()).ErrMsg = errors.Details(err) + "<<generateTask()::getPolicy()>>"
			return err
		}

		tables := make([]string, 0)
		if "*" == taskInfo.From.TableName {
			allTableNames, err := policy.GetTables(taskInfo.From)
			if nil != err {
				(*moverconfig.GetCurrState()).ErrMsg = errors.Details(err) + "<<generateTask()::GetTables()>>"
				return err
			}

			for _, tableName := range allTableNames {
				tableKey := fmt.Sprintf("%s:%s", taskInfo.From.DBName, tableName)
				_, ok := allListedTable[tableKey]
				if ok {
					continue
				}
				_, ok = processedTable[tableKey]
				if ok {
					continue
				}
				processedTable[tableKey] = 1
				tables = append(tables, tableName)
			}
		} else {
			tableKey := fmt.Sprintf("%s:%s", taskInfo.From.DBName, taskInfo.From.TableName)
			_, ok := processedTable[tableKey]
			if ok {
				continue
			}
			processedTable[tableKey] = 1
			tables = append(tables, taskInfo.From.TableName)
		}

		/* multi-Table channel slice */
		for _, tableName := range tables {
			go generateNormalTask(taskInfo, tableName, policy, taskChan) ///go
			//generateNormalTask(taskInfo, tableName, policy, taskChan)
		}
	}

	time.Sleep(time.Second * 15)
	StopGenWg.Wait()
	(*moverconfig.GetCurrState()).Prepared = true

	if taskChan.Len() == 0 {
		logging.LogInfo("No tasks.")
		return nil
	}

	return nil
}

func worker(wid int, taskChan *channels.InfiniteChannel, dataChan chan *moverconfig.RowInfo) {
	defer func() {
		stopWorkerWg.Done()
		atomic.AddInt32(&NumWorkers, -1)
	}()

	config := moverconfig.GetConfig()
	if nil == config {
		logging.LogError(logging.EC_Config_ERR, "Worker(%d) exit, nil config", wid)
		(*moverconfig.GetCurrState()).IsHealthy = false
		(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Worker(%d) exit, nil config", wid)
		return
	}
	//Ticker是一个周期触发定时的计时器，它会按照一个时间间隔往channel发送系统当前时间，而channel的接收者可以以固定的时间间隔从channel中读取事件
	ticker := time.NewTicker(time.Second * 15)

	logging.LogInfo("Worker(%d) start, task left %d", wid, taskChan.Len())
	for {
		select {
		case data := <-taskChan.Out():
			atomic.AddInt32(&(*moverconfig.GetCurrState()).DoingTasks, 1)
			task := data.(moverconfig.Task)
			srcPolicy, err := GetPolicy(task.FromDBType)
			if nil != err {
				logging.LogError(logging.EC_Runtime_ERR, "Worker(%d) exit, fail to get policy", wid)

				(*moverconfig.GetCurrState()).IsHealthy = false
				(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Worker(%d) exit, fail to get policy, err=%v", wid, err)
				return
			}
			seelog.Infof("task: %+v", task)

			var destPolicy MoverPolicy
			if task.DestTo == moverconfig.TargetDatabase && !task.IsEmpty() {
				destPolicy, err = GetPolicy(task.DestDBType)
				if nil != err {
					logging.LogError(logging.EC_Runtime_ERR, "Worker(%d) exit, fail to get policy", wid)

					(*moverconfig.GetCurrState()).IsHealthy = false
					(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Worker(%d) exit, fail to get policy, err=%v", wid, err)
					return
				}
			}

			// do task
			for 1 == task.ClosedInterval || !task.Completed() || !task.IsEmpty() {
				// check stop flag
				//ctx
				if stopFlag {
					logging.LogInfo("Worker(%d) stop", wid)
					return
				}

				// read data
				datas, err := srcPolicy.ReadData(&task)
				if nil != err {
					mysqlAddr := GetConnstr(task.FromEndpoint, "", task.DestDBIsSharding)
					logging.LogError(logging.EC_Runtime_ERR, "Worker(%d) ReadData error, task=%v,mysqlAddr=%s, err=%v", wid, task, mysqlAddr, err)

					(*moverconfig.GetCurrState()).IsHealthy = false
					(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Worker(%d) ReadData error, task=%v,mysqlAddr=%s, err=%v", wid, task, mysqlAddr, err)

					logging.LogInfo("Worker(%d) exit", wid)
					return
				}

				if nil == datas || len(datas) == 0 {
					break
				}

				// write data
				if task.DestTo == moverconfig.TargetDatabase {
					err = writeData4Database(task, datas, destPolicy, wid)
					if nil != err {
						logging.LogError(logging.EC_Runtime_ERR, "Worker(%d) write error, task=%v, err=%s", wid, task, err.Error())
						(*moverconfig.GetCurrState()).IsHealthy = false
						(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Worker(%d) write error, task=%v, err=%s", wid, task, err.Error())
						return
					}
				} else {
					err = writeData4CacheChan(task, datas, dataChan)
					if nil != err {
						logging.LogError(logging.EC_Runtime_ERR, "Worker(%d) write error, task=%v, err=%s", wid, task, err.Error())
						(*moverconfig.GetCurrState()).IsHealthy = false
						(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Worker(%d) write error, task=%v, err=%s", wid, task, err.Error())
						return
					}
				}

				// tableKey := fmt.Sprintf("%s:%s", task.FromDBName, task.FromTable)
				tableMutex.Lock()
				// atomic.AddInt64(&StatsArr[StatsMap[tableKey]].TranRows, int64(len(datas)))
				hostTableKey := fmt.Sprintf("%s:%s:%s", task.FromEndpoint.Host+":"+strconv.Itoa(task.FromEndpoint.Port), task.FromDBName, task.FromTable)
				atomic.AddInt64(&StatsArr[StatsMap[hostTableKey]].TranRows, int64(len(datas)))
				tableMutex.Unlock()
			}

			// set state
			newdoingTask := atomic.AddInt32(&(*moverconfig.GetCurrState()).DoingTasks, -1)
			if task.DestTo == moverconfig.TargetDatabase {
				atomic.AddInt32(&(*moverconfig.GetCurrState()).DoneTasks, 1)
				breakpoint.FinishTask(GetNameByHostAndPort(task.FromEndpoint.Host, task.FromEndpoint.Port, task.DestDBName), task.FromTable, &task)
			} else {
				// delay to the dispatch phase
				//标记此datachan数据结尾
				writeData4CacheChan(task, []moverconfig.RowData{}, dataChan)
			}

			leftTask := taskChan.Len() + int(newdoingTask)
			if 0 == leftTask && (*moverconfig.GetCurrState()).Prepared {
				// for sync
				time.Sleep(time.Second * 15)

				if task.DestTo == moverconfig.TargetDatabase {
					(*moverconfig.GetCurrState()).ProgressRatio = 10000
					logging.LogInfo("The job's mover completed.")
				} else { /* push eof event */
					dataChan <- nil
					logging.LogInfo("The job's gather completed.")
				}
				OfflineWorkerChOnce.Do(func() {
					close(OfflineWorkerCh)
				})
			} else {
				progress := (int((*moverconfig.GetCurrState()).SumTasks) - leftTask) * 10000 / int((*moverconfig.GetCurrState()).SumTasks)
				if progress >= 10000 {
					progress = 9999
				}
				(*moverconfig.GetCurrState()).ProgressRatio = progress
			}
		case <-stopWorkerCh:
			logging.LogInfo("Worker(%d) stop", wid)
			return
		case <-ticker.C:
			/* finish dump worker */
			if moverconfig.GetConfig().Target == moverconfig.TargetTurboMQ &&
				(*moverconfig.GetCurrState()).Prepared &&
				(*moverconfig.GetCurrState()).ProgressRatio < 10000 &&
				0 == taskChan.Len()+int(atomic.AddInt32(&(*moverconfig.GetCurrState()).DoingTasks, 0)) {
				(*moverconfig.GetCurrState()).ProgressRatio = 9999
				dataChan <- nil
				seelog.Infof("Worker(%d) gather completed", wid)
				return
			}
		case <-StopJobChan:
			logging.LogInfo("Worker(%d) was stop", wid)
			return
		}
	}
}

func writeData4CacheChan(task moverconfig.Task, datas []moverconfig.RowData, dataChan chan *moverconfig.RowInfo) error {
	logging.LogInfo("write data to cache, chan len: %d.", len(dataChan))

	if len(datas) == 0 {
		rowInfo := &moverconfig.RowInfo{
			Data:  nil,
			Task:  &task,
			MQInf: task.DestMQInf,
		}

		dataChan <- rowInfo
	}

	for idx, _ := range datas {
		var rowInfo *moverconfig.RowInfo = new(moverconfig.RowInfo)
		if nil == rowInfo {
			return errors.New("Memory is exhaust.")
		}

		rowInfo.Data = &datas[idx]
		rowInfo.Ds = &moverconfig.DataSource{
			DBType:     task.FromDBType,
			IsSharding: task.FromIsShard,
			DBName:     task.FromDBName,
			TableName:  task.FromTable,
			OrgDBName:  task.OrgDBName,
			OrgTbName:  task.OrgTbName,
			//Endpoint:   task.FromEndpoint,
			Endpoints: task.FromEndpoints,
		}
		rowInfo.Task = &task
		/*if len(rowInfo.Data.Row) > 0 {
			if rowInfo.Task.PrimaryKeyType == FDT_INT {
				i, err := strconv.ParseInt(rowInfo.Data.Row[0].(string), 10, 64)
				if err != nil {
					logging.LogError(logging.EC_Runtime_ERR, "strconv.ParseInt(rowInfo.Data.Row[0].(string): err=%v.", err)
				}
				rowInfo.CurrWriteId = i
			} else {
				rowInfo.StrCurrWriteId = rowInfo.Data.Row[0].(string)
			}
		}*/
		rowInfo.Fields = task.FromField
		rowInfo.MQInf = task.DestMQInf
		dataChan <- rowInfo
	}

	return nil
}

func writeData4Database(task moverconfig.Task, datas []moverconfig.RowData, destPolicy MoverPolicy, wid int) error {
	var err error

	if task.DestDBIsSharding {
		for _, d := range datas {
			var newDatas []moverconfig.RowData
			newDatas = append(newDatas, d)
			err = destPolicy.WriteData(&task, newDatas)
			if nil != err {
				logging.LogError(logging.EC_Runtime_ERR, "Worker(%d) WriteData error, task=%v, err=%v", wid, task, err)

				(*moverconfig.GetCurrState()).IsHealthy = false
				(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Worker(%d) WriteData error, task=%v, err=%v", wid, task, err)

				// find error row
				found, errorRow := FindErrorData(destPolicy, &task, newDatas)
				if found {
					logging.LogError(logging.EC_Runtime_ERR, "Found failure record, db=%s, table=%s, data: %v", task.DestDBName, task.DestTable, errorRow)
					(*moverconfig.GetCurrState()).IsHealthy = false
					(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Found failure record, db=%s, table=%s, data: %v", task.DestDBName, task.DestTable, errorRow)
				} else {
					logging.LogError(logging.EC_Runtime_ERR, "Can't find failure record, db=%s, table=%s", task.DestDBName, task.DestTable)
					(*moverconfig.GetCurrState()).IsHealthy = false
					(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Can't find failure record, db=%s, table=%s", task.DestDBName, task.DestTable)
				}

				// worker exit
				logging.LogInfo("Worker(%d) exit", wid)
				return err
			}
		}
	} else {
		err = destPolicy.WriteData(&task, datas)
		if nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "Worker(%d) WriteData error, task=%v, err=%v", wid, task, err)

			(*moverconfig.GetCurrState()).IsHealthy = false
			(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Worker(%d) WriteData error, task=%v, err=%v", wid, task, err)

			// find error row
			found, errorRow := FindErrorData(destPolicy, &task, datas)
			if found {
				logging.LogError(logging.EC_Runtime_ERR, "Found failure record, db=%s, table=%s, data: %v", task.DestDBName, task.DestTable, errorRow)
				(*moverconfig.GetCurrState()).IsHealthy = false
				(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Found failure record, db=%s, table=%s, data: %v", task.DestDBName, task.DestTable, errorRow)
			} else {
				logging.LogError(logging.EC_Runtime_ERR, "Can't find failure record, db=%s, table=%s", task.DestDBName, task.DestTable)
				(*moverconfig.GetCurrState()).IsHealthy = false
				(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Can't find failure record, db=%s, table=%s", task.DestDBName, task.DestTable)
			}

			// worker exit
			logging.LogInfo("Worker(%d) exit", wid)
			return err
		}
	}

	return nil
}

func GetFullTopicName(db, table string) string {
	return fmt.Sprintf("%s_%s_full", db, table)
}

func DispatchEofJob() error {

	SlaveMgr.topicsMtx.Lock()
	defer SlaveMgr.topicsMtx.Unlock()
	flag := make(map[string]struct{})

	for _, topic := range SlaveMgr.topics {
		var job SQLContext
		job.action = actionEof
		job.topic = topic
		if _, ok := flag[topic]; ok {
			continue
		}
		flag[topic] = struct{}{}

		if err := SlaveMgr.DispatchJob(&job); nil != err {
			logging.LogError(logging.EC_Runtime_ERR, "dispatcher error, job:%v, err: %v", job, err)
			(*moverconfig.GetCurrState()).IsHealthy = false
			(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("dispatcher error, job:%v, err: %v", job, err)
			return errors.New("dispatcher failed")
		}
		logging.LogInfo("dispatch eof job topic: %s", topic)
	}

	// close slave manager
	SlaveMgr.WaitJobsDone()

	return nil
}

func DispatchJob(policy MoverPolicy, dataChan chan *moverconfig.RowInfo, dumptime uint32) error {
	for {
		select {
		case data := <-dataChan:
			{
				if nil == data {
					// close slave manager
					SlaveMgr.WaitJobsDone()

					// dispatch eof job
					DispatchEofJob()

					seelog.Infof("All data write over")
					(*moverconfig.GetCurrState()).ProgressRatio = 10000

					return nil
				}

				// sync push MQ operation
				if nil == data.Data {
					SlaveMgr.WaitJobsDone()
					seelog.Infof("Sync one task, wait jobs done.")
					if nil != data.Task {
						atomic.AddInt32(&(*moverconfig.GetCurrState()).DoneTasks, 1)
						breakpoint.FinishTask(GetNameByHostAndPort(data.Task.FromEndpoint.Host, data.Task.FromEndpoint.Port, data.Task.DestDBName), data.Task.FromTable, data.Task)
					}

					// gather each topic
					dbname, tbname := data.Task.OrgDBName, data.Task.OrgTbName
					topic := data.MQInf.Topic
					if topic == "" {
						topic = GetFullTopicName(dbname, tbname)
					}

					k := fmt.Sprintf("%s.%s", dbname, tbname)
					AppendOneTopic(k, topic)

					continue
				}

				var job SQLContext
				job.action = actionInsert
				job.database = data.Ds.OrgDBName
				job.table = data.Ds.OrgTbName
				job.realDatabase = data.Ds.DBName
				job.realTable = data.Ds.TableName
				job.timestamp = dumptime

				job.datas = data.Data.Row

				_, FieldInfo, err := policy.GetFields(*data.Ds, data.Fields, false)
				if nil != err {
					seelog.Errorf("Fail to get table field info, ds: %v, err: %v", *data.Ds, err)
					return err
				}

				columns := make([]*moverconfig.FieldInfo, 0)
				indColumns := make([]*moverconfig.FieldInfo, 0)
				for i := range FieldInfo {
					v := FieldInfo[i]
					columns = append(columns, &v)

					if len(data.MQInf.PrimaryKeys) > 0 {
						for _, k := range data.MQInf.PrimaryKeys {
							if k == v.FieldName {
								indColumns = append(indColumns, &v)
							}
						}
					} else if v.IsPriamryKey {
						indColumns = append(indColumns, &v)
					}
				}

				job.columns = columns
				job.indexColumns = indColumns

				job.topic = data.MQInf.Topic
				if job.topic == "" {
					job.topic = GetFullTopicName(job.database, job.table)
				}
				//finishmoverconfig.RowData(data)
				if err := SlaveMgr.DispatchJob(&job); nil != err {
					logging.LogInfo("dispatcher error, err: %v", err)
					return errors.New("dispatcher failed")
				}
			}
		case <-StopJobChan:
			logging.LogInfo("dispatcher stop.")
			return errors.New("dispatcher exit(-1)")
		}
	}
}

func ShowSummary() error {
	logging.LogInfo("Summary:")

	for _, tableInfo := range StatsArr {
		logging.LogInfo("host=%s, Table %s.%s SumRows:%d, TransRows:%d", tableInfo.HostPort, tableInfo.DbName, tableInfo.TableName, tableInfo.SumRows, tableInfo.TranRows)
	}

	return nil
}

func GetTableInfos() []moverconfig.TableStat {
	return StatsArr
}
