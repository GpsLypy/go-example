package service

import (
	"app/tools/excavator/hauler/msgtracer"
	"app/tools/excavator/replay/writer/jsonevent"
	"app/tools/excavator/replay/writer/turbo"
	"app/tools/mover/moverconfig"
	"app/tools/mover/utils"
	"bytes"
	"github.com/juju/errors"
	"github.com/siddontang/go/sync2"
	"sync"
)

const (
	// Dispatch event by primary key (default)
	dispatchMethodPK = iota
	// Dispatch event by table name
	dispatchMethodTable
)

var (
	WorkerPerQ int
	SlaveMgr   *SlaveManager
)

type msgLogWriter struct {
	file *msgtracer.MsgLogFile
	mu   sync.Mutex
}

type SlaveManager struct {
	Finished bool
	workers  []*ExecuteWorker
	cfg      *moverconfig.MoverConfig
	workerWg sync.WaitGroup
	jobWg    sync.WaitGroup
	// Debug
	jobAdd  sync2.AtomicInt64
	jobDone sync2.AtomicInt64
	// Dispatch method
	dispatchMethod int
	dispatchKeyBuf bytes.Buffer
	// Worker count per queue of broker
	workerPerQ int
	// Message log files
	msgLogWriters []*msgLogWriter
	topicsMtx     sync.Mutex
	topics        map[string]string
}

func AppendOneTopic(k, topic string) {
	SlaveMgr.topicsMtx.Lock()
	defer SlaveMgr.topicsMtx.Unlock()

	if _, ok := SlaveMgr.topics[k]; !ok {
		SlaveMgr.topics[k] = topic
	}
}

func initSlaveManager(cfg *moverconfig.MoverConfig) error {
	m, err := NewSlaveManager(cfg)
	if nil != err {
		return err
	}

	SlaveMgr = m
	return nil
}

func NewSlaveManager(cfg *moverconfig.MoverConfig) (*SlaveManager, error) {

	m := &SlaveManager{}
	m.cfg = cfg
	m.workerPerQ = 0
	//感觉有bug，要修复
	if 0 == m.workerPerQ {
		m.workerPerQ = 1
	}
	//拿到MQ的配置信息
	MQConf := cfg.Dispatcher.TurboMQConf
	if len(MQConf.Brokers) == 0 {
		return nil, errors.New("Configure lack of brokers.")
	}

	if MQConf.QueueNumber <= 0 {
		MQConf.QueueNumber = 1
	}

	//统计共需要多少worker
	totalWorkerCount := len(MQConf.Brokers) * MQConf.QueueNumber * m.workerPerQ
	m.workers = make([]*ExecuteWorker, totalWorkerCount)
	if nil == m.workers {
		return nil, errors.New("lack of workers in the slave-manager")
	}

	m.topics = make(map[string]string, 0)
	if nil == m.topics {
		return nil, errors.New("no topics in the slave-manager")
	}

	var i int
	for bi, v := range MQConf.Brokers {
		for q := 0; q < MQConf.QueueNumber; q++ {
			// Create worker-count workers per queue of broker
			for wi := 0; wi < m.workerPerQ; wi++ {
				var workerCfg WorkerConfig
				workerCfg.FromAppConfig(&cfg.Dispatcher)
				workerCfg.BrokerName = v
				workerCfg.BrokerIndex = bi
				workerCfg.QueueID = q
				workerCfg.ProducerGroup = MQConf.ProducerGroup

				w := NewExecuteWorker(&workerCfg, &m.jobWg)
				w.mgr = m
				m.workers[i] = w

				w.SetWorkerID(i)
				i++
			}
		}
	}

	return m, nil
}

func (m *SlaveManager) Start() error {
	var err error
	for _, v := range m.workers {
		if err = v.Start(&m.workerWg); nil != err {
			return err
		}
	}

	return nil
}

func (m *SlaveManager) Close() {
	for _, v := range m.workers {
		v.Close()
	}
	m.workerWg.Wait()
}

func (m *SlaveManager) DispatchJob(job *SQLContext) error {
	widx := -1

	// Get JSON format with field filter
	currentJob := job

	// dispatch eof job for worker(0), to avoid repeating
	if currentJob.action == actionEof {
		for _, worker := range SlaveMgr.workers {
			// One queue-brocker has several workers
			//if worker.ID == 0 {
			worker.PushJob(currentJob)
			//}
		}

		return nil
	}

	// Determine which worker to handle  确定哪个工人处理
	m.dispatchKeyBuf.Reset()
	for _, v := range currentJob.indexColumns {
		if m.dispatchKeyBuf.Len() != 0 {
			m.dispatchKeyBuf.WriteByte('_')
		}

		m.dispatchKeyBuf.WriteString(utils.FormatFieldValue(utils.CastUnsigned(job.datas[v.Idx], v.IsUnsigned)))
	}
	hashcode := jsonevent.StringHashCode(m.dispatchKeyBuf.String())
	if hashcode < 0 {
		hashcode = -hashcode
	}
	// Hash to which queue of broker
	widx = hashcode % (len(m.workers) / m.workerPerQ)
	if widx < 0 {
		return errors.New("Can't find destinate dispatch worker")
	}
	// Get which queue should push
	widx *= m.workerPerQ
	if m.workerPerQ != 1 {
		hashcode = jsonevent.StringHashCode(job.topic)
		if hashcode < 0 {
			hashcode = 0
		}
		widx += hashcode % m.workerPerQ
	}

	m.workers[widx].PushJob(currentJob)
	return nil
}

func (m *SlaveManager) WaitJobsDone() {
	m.jobWg.Wait()
}

func (m *SlaveManager) GetTotalExecuted() int64 {
	total := int64(0)
	for _, v := range m.workers {
		total += v.executed.Get()
	}
	return total
}

func (m *SlaveManager) GetExecuteCostMS(output []int64) {
	for i, v := range m.workers {
		if i >= len(output) {
			return
		}
		output[i] = v.executeCostMS.Get()
	}
}

func (m *SlaveManager) PersisteMessageIDs(worker *ExecuteWorker, messages []*turbo.Message) error {
	if nil == m.msgLogWriters {
		return nil
	}

	writer := m.msgLogWriters[worker.cfg.BrokerIndex]
	writer.mu.Lock()
	defer writer.mu.Unlock()

	var err error
	for _, m := range messages {
		if werr := writer.file.WriteMessagePosition(m.GetMessageID(),
			m.Userdata.(*SQLContext).pos,
			m.Userdata.(*SQLContext).gtidSeq); nil != werr {
			err = werr
		}
	}
	return errors.Trace(err)
}
