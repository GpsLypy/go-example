package mysqlib

import (
	"context"
	"data-transmission-service/application/components/addrwrapper"
	"data-transmission-service/application/components/inbound"
	"database/sql"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/cihub/seelog"
	"github.com/eapache/channels"
	"github.com/juju/errors"
	"golang.org/x/sync/errgroup"
)

type source struct {
	// From addr wrapper
	addrWrapper *addrwrapper.AddrWrapper
	// From fixed host
	host string
	port int
	// Authenticate
	username, password string
	queryDBs           map[string]*sql.DB
}

//固定信息
func NewHostSource(host string, port int, username, password string) *source {
	return &source{
		host:     host,
		port:     port,
		username: username,
		password: password,
		queryDBs: make(map[string]*sql.DB),
	}
}

func NewAddrWrapperSource(addrWrapper *addrwrapper.AddrWrapper, username, password string) *source {
	return &source{
		addrWrapper: addrWrapper,
		username:    username,
		password:    password,
		queryDBs:    make(map[string]*sql.DB),
	}
}

func (s *source) create() error {
	if (s.host == "" || s.port == 0) &&
		nil == s.addrWrapper {
		return errors.New("Invalid datasource configuration")
	}

	heartbeatInterval := time.Second * 30
	if s.parent.getHeartbeatIntervalSec() != 0 {
		heartbeatInterval = time.Second * time.Duration(s.parent.getHeartbeatIntervalSec())
	}
	readTimeout := heartbeatInterval * 3 / 2

	cfg := replication.BinlogSyncerConfigEx2{
		ServerID:        s.parent.getServerID(),
		Flavor:          mysql.MySQLFlavor,
		User:            s.username,
		Password:        s.password,
		HeartbeatPeriod: heartbeatInterval,
		ReadTimeout:     readTimeout,
		UseDecimal:      s.parent.getUseDecimal(),
	}

	var addrFunc func() (string, error)
	if nil != s.addrWrapper {
		// Create using addr wrapper
		addrFunc = func() (string, error) {
			return s.addrWrapper.Next()
		}
	} else {
		cfg.Host = s.host
		cfg.Port = uint16(s.port)
	}

	var err error
	if err = s.initQueryDB(); nil != err {
		return errors.Trace(err)
	}

	// Validate
	if err = s.validate(); nil != err {
		return errors.Trace(err)
	}
	// Apply flavor
	cfg.Flavor = s.sourceType

	s.syncer = replication.NewBinlogSyncerEx2(cfg, addrFunc)
	return nil
}

//记录工作线程内部的状态信息
type workState struct {
	ErrMsg        string
	IsHealthy     bool
	DoingTasks    int32
	Prepared      bool
	ProgressRatio int
	SumTasks      int32
	Errresult     []error
}

type MysqlInbound struct {
	ctx        context.Context
	cancel     context.CancelFunc
	taskChan   *channels.InfiniteChannel
	dataChan   chan *RowInfo
	workerSize int
	sources    []*source
	//stopWorkerWg sync.WaitGroup
	eventCh   chan *inbound.InboundEvent
	currState *workState
	config    MoverConfig
	eg        *errgroup.Group
}

func NewMysqlInbound(ctx context.Context, config MoverConfig, taskChan *channels.InfiniteChannel, dataChan chan *RowInfo, workerSize int) (*MysqlInbound, error) {
	State := workState{ErrMsg: "", IsHealthy: true, DoingTasks: 0, SumTasks: 0, Prepared: true, ProgressRatio: 0}
	g, _ := errgroup.WithContext(context.Background())
	ib := &MysqlInbound{
		taskChan:   taskChan,
		dataChan:   dataChan,
		ctx:        ctx,
		workerSize: workerSize,
		currState:  &State,
		config:     config,
		eg:         g,
	}
	return ib, nil
}

func (m *MysqlInbound) Start(positions []inbound.Position) error {
	if len(positions) != 0 && len(positions) != len(m.sources) {
		return errors.Errorf("can't start inbound, datasource count = %d, position count = %d",
			len(m.sources), len(positions))
	}
	m.ctx, m.cancel = context.WithCancel(context.Background())
	for i := 0; i < m.workerSize; i++ {
		m.eg.Go(m.worker)
	}
	return nil
}

func (m *MysqlInbound) Close() error {
	m.cancel()
	if err := m.eg.Wait(); err == nil {
		return nil
	} else {
		return err
	}
}

//暂停
func (m *MysqlInbound) Pause() error {

	return nil
}

//恢复
func (m *MysqlInbound) Resume([]inbound.Position) error {

	return nil
}

//Rds
func (m *MysqlInbound) AddAddrWrapper(envir, addr, api string, username, password string, jobID int64) error {
	// Create to pull all nodes
	addrw := addrwrapper.NewAddrWrapper(envir, addr, api, jobID, 0, "")
	dsCount, err := addrw.GetDatasourceCount()
	if nil != err {
		return errors.Trace(err)
	}
	if dsCount == 0 {
		return errors.New("Get empty nodes from rds source")
	}
	for i := 0; i < dsCount; i++ {
		shardingIndex := i + 1
		if dsCount == 1 {
			// Non-sharding db
			shardingIndex = 0
		}
		src := NewAddrWrapperSource(addrwrapper.NewAddrWrapper(envir, addr, api, jobID, shardingIndex, ""),
			username, password)
		src.tag = strconv.Itoa(len(m.sources))
		m.sources = append(m.sources, src)
	}
	return nil
}

func (m *MysqlInbound) AddFixedHostSource(host string, port int, username string, password string) error {
	src := NewHostSource(host, port, username, password)
	src.tag = strconv.Itoa(len(m.sources))
	m.sources = append(m.sources, src)
	src.fixedIndex = len(m.sources) - 1
	return nil
}

func (m *MysqlInbound) GetEventChannel() <-chan *inbound.InboundEvent {
	return m.eventCh
}

func (m *MysqlInbound) GetDataSourceCount() int {
	return len(m.sources)
}

func (m *MysqlInbound) GetStartPoint(int) inbound.Position {
	return inbound.Position{}
}

func (m *MysqlInbound) GetQueryDB(int) *sql.DB {
	return nil
}

func (m *MysqlInbound) GetAddress(int) string {
	return ""
}

func (m *MysqlInbound) PushEvent(event *inbound.InboundEvent) {
	m.eventCh <- event
}

func (m *MysqlInbound) worker() error {
	defer func() {
		atomic.AddInt32(&NumWorkers, -1)
	}()

	//Ticker是一个周期触发定时的计时器，它会按照一个时间间隔往channel发送系统当前时间，而channel的接收者可以以固定的时间间隔从channel中读取事件
	ticker := time.NewTicker(time.Second * 15)

	LogInfo("Worker start, task left %d", m.taskChan.Len())
	for {
		select {
		case data := <-m.taskChan.Out():
			atomic.AddInt32(&m.currState.DoingTasks, 1)
			task := data.(Task)
			srcPolicy, err := getPolicy(task.FromDBType)
			if nil != err {
				LogError(EC_Runtime_ERR, "Worker exit, fail to get policy")
				m.currState.IsHealthy = false
				m.currState.ErrMsg = fmt.Sprintf("Worker exit, fail to get policy, err=%v", err)
				return err
			}
			seelog.Infof("task: %+v", task)

			// do task
			for task.ClosedInterval == 1 || !task.Completed() || !task.IsEmpty() {
				// check stop flag
				if StopFlag {
					LogInfo("Worker stop")
					return nil
				}

				// read data
				datas, err := srcPolicy.ReadData(&task, m.config)
				if nil != err {
					endpoint := task.FromEndpoint
					Endpoint := &Endpoint{endpoint.Host, endpoint.Port, endpoint.User, endpoint.Password}
					mysqlAddr := getConnstr(*Endpoint, "", task.DestDBIsSharding)
					LogError(EC_Runtime_ERR, "Worker ReadData error, task=%v,mysqlAddr=%s, err=%v", task, mysqlAddr, err)

					m.currState.IsHealthy = false
					m.currState.ErrMsg = fmt.Sprintf("Worker ReadData error, task=%v,mysqlAddr=%s, err=%v", task, mysqlAddr, err)

					LogInfo("Worker  exit")
					return err
				}

				if nil == datas || len(datas) == 0 {
					break
				}

				// write data
				err = WriteData4CacheChan(task, datas, m.dataChan)
				if nil != err {
					LogError(EC_Runtime_ERR, "Worker write error, task=%v, err=%s", task, err.Error())
					m.currState.IsHealthy = false
					m.currState.ErrMsg = fmt.Sprintf("Worker write error, task=%v, err=%s", task, err.Error())
					return err
				}
				TableMutex.Lock()
				hostTableKey := fmt.Sprintf("%s:%s:%s", task.FromEndpoint.Host+":"+strconv.Itoa(task.FromEndpoint.Port), task.FromDBName, task.FromTable)
				atomic.AddInt64(&StatsArr[StatsMap[hostTableKey]].TranRows, int64(len(datas)))
				TableMutex.Unlock()
			}

			// set state
			newdoingTask := atomic.AddInt32(&m.currState.DoingTasks, -1)

			// delay to the dispatch phase
			//标记此datachan数据结尾
			WriteData4CacheChan(task, []RowData{}, m.dataChan)

			leftTask := m.taskChan.Len() + int(newdoingTask)
			if leftTask == 0 && m.currState.Prepared {
				// for sync
				time.Sleep(time.Second * 15)
				/* push eof event */
				m.dataChan <- nil
				LogInfo("The job's gather completed.")

				offlineWorkerChOnce.Do(func() {
					close(offlineWorkerCh)
				})
			} else {
				progress := (int(m.currState.SumTasks) - leftTask) * 10000 / int(m.currState.SumTasks)
				if progress >= 10000 {
					progress = 9999
				}
				m.currState.ProgressRatio = progress
			}
		case <-StopWorkerCh:
			LogInfo("Worker stop")
			return nil
		case <-ticker.C:
			/* finish dump worker */
			if m.currState.Prepared &&
				m.currState.ProgressRatio < 10000 &&
				m.taskChan.Len()+int(atomic.AddInt32(&m.currState.DoingTasks, 0)) == 0 {
				m.currState.ProgressRatio = 9999
				m.dataChan <- nil
				seelog.Infof("Worker gather completed")
				return nil
			}
		case <-stopJobChan:
			LogInfo("Worker was stop")
			return nil
		case <-m.ctx.Done():
			return m.ctx.Err()
		}
	}
}

func (m *MysqlInbound) GetConfig() *MoverConfig {
	return &m.config
}

func (m *MysqlInbound) GetContext() context.Context {
	return m.ctx
}
