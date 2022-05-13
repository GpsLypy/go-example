package service

import (
	"app/tools/excavator/alarm"
	"app/tools/excavator/bformat"
	"app/tools/excavator/replay/writer/jsonevent"
	"app/tools/excavator/replay/writer/turbo"
	"app/tools/mover/logging"
	"app/tools/mover/moverconfig"
	"app/tools/mover/utils"
	"fmt"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go/sync2"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBatchSize = 128
	defaultCommitMS  = 100
)

const (
	_ = iota
	actionInsert
	actionUpdate
	actionDelete
	actionEof
)

type SQLContext struct {
	text   bool
	action int

	// After rewrite
	database string
	table    string

	// Original
	realDatabase string
	realTable    string
	timestamp    uint32

	//
	pos mysql.Position
	// The gtid sequence
	gtidSeq uint64

	// Columns
	columns      []*moverconfig.FieldInfo
	indexColumns []*moverconfig.FieldInfo

	// Binlog meta
	columnTypes []byte
	columnMetas []uint16

	// Datas
	datas []interface{}

	// For mq publish
	topic string
	msg   *turbo.Message

	// Filter strategy
	strategy      *moverconfig.StrategyConfig
	subscribeType string
}

const (
	topicMessageBatchCount = 128
	topicMessageBatchSize  = 1000 * 1024
)

type TopicMessages struct {
	topic            string
	Messages         []*turbo.Message
	PendingMessages  []*turbo.Message
	bodySizeEstimate int64
	seed             int64
	lastEmit         int64
}

type MQTaskQueue struct {
	Topics map[string]*TopicMessages
}

func NewTaskQueue() *MQTaskQueue {
	return &MQTaskQueue{
		Topics: make(map[string]*TopicMessages),
	}
}

// Return should flush immediately
func (t *MQTaskQueue) Push(ctx *SQLContext) bool {
	tms, ok := t.Topics[ctx.topic]
	if !ok {
		tms = &TopicMessages{
			topic:           ctx.topic,
			Messages:        make([]*turbo.Message, 0, 128),
			PendingMessages: make([]*turbo.Message, 0, 32),
		}
		t.Topics[ctx.topic] = tms
	}
	// Convert to message
	columnNames := make([]string, 0, 64)
	for _, v := range ctx.columns {
		columnNames = append(columnNames, v.FieldName)
	}
	dataType := 0
	if ctx.action == actionInsert {
		dataType = bformat.DataTypeInsert
	} else if ctx.action == actionUpdate {
		dataType = bformat.DataTypeUpdate
	} else if ctx.action == actionDelete {
		dataType = bformat.DataTypeDelete
	} else if ctx.action == actionEof {
		dataType = bformat.DataTypeEof
	}

	var jsonData []byte
	if dataType == bformat.DataTypeEof {
		jsonData = []byte(ctx.topic + "_eof")
	} else {
		je, err := jsonevent.GetJSONRowEventRawSubscribeType(ctx.database, ctx.table, columnNames,
			ctx.datas, nil, "", dataType, ctx.timestamp, ctx.subscribeType, 0)
		if nil != err {
			panic(err)
		}
		jsonData, err = je.GetJSON()
		if nil != err {
			panic(err)
		}
	}

	// Get turboMQ message
	var mqMsg turbo.Message
	mqMsg.SetID(atomic.AddInt64(&tms.seed, 1))
	mqMsg.Body = jsonData
	mqMsg.Userdata = ctx

	// Should into messages or pending messages
	emitNow := false
	messageSize := mqMsg.EncodeSize(ctx.topic)
	if messageSize+tms.bodySizeEstimate >= topicMessageBatchSize ||
		len(tms.Messages) >= topicMessageBatchCount {
		// Into pending messages
		emitNow = true
		tms.PendingMessages = append(tms.PendingMessages, &mqMsg)
	} else {
		tms.bodySizeEstimate += messageSize
		tms.Messages = append(tms.Messages, &mqMsg)
	}

	return emitNow
}

func (t *MQTaskQueue) GetTopicMessage(topic string) *TopicMessages {
	tms, _ := t.Topics[topic]
	return tms
}

func (tms *TopicMessages) PendingPrepare() {
	if len(tms.PendingMessages) == 0 {
		return
	}
	if len(tms.Messages) > 0 {
		panic("Can't prepare pending messages while messages are unsent")
	}
	for _, m := range tms.PendingMessages {
		tms.Messages = append(tms.Messages, m)
	}
	tms.PendingMessages = tms.PendingMessages[0:0]

	// Update estimate size
	tms.bodySizeEstimate = 0
	for _, m := range tms.Messages {
		tms.bodySizeEstimate += int64(m.EncodeSize(tms.topic))
	}
}

func (tms *TopicMessages) Emitable(interval int64) bool {
	tn := time.Now().UnixNano() / 1e6
	if tn-tms.lastEmit > interval {
		return true
	}
	return false
}

// Write message id for search
func (tms *TopicMessages) PersistMessageID(w *ExecuteWorker, msgIDList []string) error {
	for mi, q := range tms.Messages {
		q.SetMessageID(msgIDList[mi])
	}

	return errors.Trace(w.mgr.PersisteMessageIDs(w, tms.Messages))
}

// Emit the specified topic messages
func (tms *TopicMessages) Emit(w *ExecuteWorker) error {
	if nil == tms.Messages || 0 == len(tms.Messages) {
		// Check pending messages
		tms.PendingPrepare()
		return nil
	}

	var err error
	retryTimes := 0
	retryInterval := time.Second * 2
	maxRetryTimes := 100000000
	executeStartTm := time.Now().UnixNano() / 1e6
	errorReport := false

	for {
		if nil != err {
			// Max retry
			if retryTimes >= maxRetryTimes {
				break
			}

			// Error report
			if !errorReport {
				errorReport = true
			}

			// TODO: Cancel context

			log.Warnf("Retry commit job queue, count %v, retry times %v, lastError %v", len(tms.Messages), retryTimes, err)
			// We meet error last loop, so we enter sleep and try again
			time.Sleep(retryInterval)
			// If receive signal to quit
			if qv := atomic.LoadInt64(&w.closed); 0 != qv {
				return err
			}
			retryTimes++
		}

		// Force to update the broker addr list, avoid address change
		if retryTimes != 0 && retryTimes%15 == 0 {
			if err = w.client.UpdateBrokerAddress(); nil != err {
				log.Warnf("Worker %d update broker address failed, error: %v", w.ID, err)
			} else {
				log.Infof("Worker %d update broker address after %d retry", w.ID, retryTimes)
				// Once broker address updated, we should release the connection to make sure
				// new connection will use the updated address from nameserver
				w.client.FreeConnection()
			}
		}

		var rsp *turbo.SendResult
		if len(tms.Messages) == 1 {
			rsp, err = w.client.PublishWithTopic(tms.Messages[0], tms.topic)
		} else {
			rsp, err = w.client.PublishBatchWithTopic(tms.Messages, tms.topic)
		}
		if nil != err {
			continue
		}

		// Read response
		// Not error, but response is not SUCCESS
		if rsp.SendStatus != turbo.Success {
			err = errors.Errorf("Send response not success(%d)", rsp.SendStatus)
			continue
		}
		log.Infof("Worker %d publish %d message for topic %s success", w.ID, len(tms.Messages), tms.topic)
		// Record the message id if id is not zero
		msgIDList := strings.Split(strings.TrimRight(rsp.MsgId, ","), ",")
		if len(msgIDList) != len(tms.Messages) {
			err = errors.Errorf("Message id from MQ response count mismatch, want %d, got %d",
				len(tms.Messages), len(msgIDList))
			continue
		}

		// Below is not require success
		if persistErr := tms.PersistMessageID(w, msgIDList); nil != persistErr {
			log.Warnf("Error on persisting message id: %v", err)
		}

		// Publish done
		binlogTimestamp := uint32(0)
		for _, m := range tms.Messages {
			job, ok := m.Userdata.(*SQLContext)
			if !ok {
				panic("Wrong type")
			}
			if job.timestamp != 0 {
				binlogTimestamp = job.timestamp
			}
		}

		if 0 != binlogTimestamp {
			w.binlogDelayTime.Set(time.Now().Unix() - int64(binlogTimestamp))
		}

		// Set execute statistic
		w.executeCostMS.Set((time.Now().UnixNano()/1e6 - executeStartTm))

		// Emit over
		emitCount := len(tms.Messages)
		tms.Messages = tms.Messages[0:0]
		tms.bodySizeEstimate = 0
		tms.PendingPrepare()
		tms.lastEmit = time.Now().UnixNano() / 1e6

		// Update worker
		w.jobWg.Add(-1 * emitCount)
		w.executed.Add(int64(emitCount))
		w.mgr.jobDone.Add(int64(emitCount))

		return nil
	}

	return err
}

type WorkerConfig struct {
	NameSvrAddr   string `toml:"namesvr-addr" json:"namesvr-addr"`
	QueueID       int
	BrokerName    string
	BrokerIndex   int
	ProducerGroup string
	BatchSize     int   `toml:"batch-size" json:"batch-size"`
	CommitMS      int64 `toml:"commit-ms" json:"commit-ms"`
}

func (c *WorkerConfig) FromAppConfig(dispatcher *moverconfig.Dispatcher) error {
	c.NameSvrAddr = dispatcher.TurboMQConf.NamesrvAddr

	if 0 == dispatcher.TurboMQConf.BatchSize {
		c.BatchSize = defaultBatchSize
	} else {
		c.BatchSize = dispatcher.TurboMQConf.BatchSize
	}

	if 0 == dispatcher.BatchCommitMS {
		c.CommitMS = defaultCommitMS
	} else {
		c.CommitMS = int64(dispatcher.BatchCommitMS)
	}

	return nil
}

type ExecuteWorker struct {
	// The dest MQ
	client *turbo.MQClient

	ich             chan *SQLContext
	closed          int64
	cfg             *WorkerConfig
	cfgLock         sync.RWMutex
	jobWg           *sync.WaitGroup
	mgr             *SlaveManager
	executed        sync2.AtomicInt64
	binlogDelayTime sync2.AtomicInt64
	seed            int64
	// Worker statistics
	executeCostMS sync2.AtomicInt64
	// Export fields
	ID int
	// Task queue
	taskQueue *MQTaskQueue
}

func NewExecuteWorker(cfg *WorkerConfig, jobWg *sync.WaitGroup) *ExecuteWorker {
	w := &ExecuteWorker{}
	w.ich = make(chan *SQLContext, 1000)
	w.cfg = cfg
	if cfg.BatchSize == 0 {
		cfg.BatchSize = defaultBatchSize
	}
	if cfg.BatchSize > defaultBatchSize {
		log.Warnf("MQ batch size is out of limit(%d), auto reset to the default value",
			defaultBatchSize)
		cfg.BatchSize = defaultBatchSize
	}
	if cfg.CommitMS == 0 {
		cfg.CommitMS = defaultCommitMS
	}

	w.jobWg = jobWg
	w.taskQueue = NewTaskQueue()
	return w
}

func (w *ExecuteWorker) SetWorkerID(id int) {
	w.ID = id
}

func (w *ExecuteWorker) Close() {
	if !atomic.CompareAndSwapInt64(&w.closed, 0, 1) {
		return
	}
	close(w.ich)
	w.jobWg.Wait()
}

func (w *ExecuteWorker) PushJob(job *SQLContext) {
	w.jobWg.Add(1)
	w.mgr.jobAdd.Add(1)
	if job.timestamp == 0 {
		job.timestamp = uint32(time.Now().Unix())
	}
	job.subscribeType = "FULL"
	w.ich <- job
}

func (w *ExecuteWorker) SetConfig(cfg *WorkerConfig) {
	w.cfgLock.Lock()
	w.cfg = cfg
	w.cfgLock.Unlock()
}

func (w *ExecuteWorker) GetConfig() *WorkerConfig {
	var cfg *WorkerConfig
	w.cfgLock.RLock()
	cfg = w.cfg
	w.cfgLock.RUnlock()

	return cfg
}

func (w *ExecuteWorker) Start(wg *sync.WaitGroup) error {
	var err error
	w.client = turbo.NewTurboMQClient(&turbo.MQConfig{
		ProducerGroup:  w.cfg.ProducerGroup,
		NameSvrAddr:    w.cfg.NameSvrAddr,
		Topic:          "",
		SendTimeoutSec: 5,
		ReadTimeoutSec: 5,
	}, "Replayer", w.cfg.BrokerName, int32(w.cfg.QueueID))
	if err = w.client.StartNoTopic(); nil != err {
		return errors.Trace(err)
	}

	wg.Add(1)
	go w.workv2(wg)

	return nil
}

func (w *ExecuteWorker) workv2(wg *sync.WaitGroup) {
	var err error

	cfg := w.GetConfig()
	ticker := time.NewTicker(time.Millisecond * time.Duration(cfg.CommitMS))

	defer func() {
		ticker.Stop()

		if err == nil {
			log.Infof("Worker %d quit ok", w.ID)
		} else {
			alarm.SendAlarm(fmt.Sprintf("Worker %d quit with error %v", w.ID, err))
		}
		pi := recover()
		if nil != pi {
			log.Errorf("Panic: %v", pi)
			buf := make([]byte, 4<<10) // 4 KB should be enough
			runtime.Stack(buf, false)
			log.Errorf("Stack: %v", string(buf))
		}
		wg.Done()
	}()

	idleTimeContinue := int64(0)

	for {
		select {
		case evt, ok := <-w.ich:
			{
				if !ok {
					// Closed
					return
				}
				idleTimeContinue = 0
				parseToData(evt)
				// Push to task queue
				if w.taskQueue.Push(evt) {
					// Emit immediately
					if err = w.taskQueue.GetTopicMessage(evt.topic).Emit(w); nil != err {
						return
					}
				}
			}
		case <-ticker.C:
			{
				tn := time.Now().UnixNano() / 1e6

				// Loop for every topic
				purgeTopics := make(map[string]struct{})
				emitCount := 0
				notEmitCount := 0

				// Check timeout topics
				for _, tms := range w.taskQueue.Topics {
					if len(tms.Messages) > 0 {
						if tn-tms.lastEmit > cfg.CommitMS {
							if err = tms.Emit(w); nil != err {
								return
							}
							idleTimeContinue = 0
							emitCount++
						} else {
							notEmitCount++
						}
					}

					if tn-tms.lastEmit > 10*60*1000 &&
						len(tms.Messages) == 0 &&
						len(tms.PendingMessages) == 0 {
						purgeTopics[tms.topic] = struct{}{}
					}
				}

				// Remove purge topic
				for topic := range purgeTopics {
					delete(w.taskQueue.Topics, topic)
					log.Infof("Topic %s was removed due too idle too long", topic)
				}

				// No tasks in all topics, it is in idle空闲的 state
				if 0 == emitCount+notEmitCount {
					idleTimeContinue += cfg.CommitMS
					if idleTimeContinue > 300*1000 /* 5 minutes */ {
						// Idle time continues for 5 minutes, it normally happens when
						// upstream mysql has no binlog datas. Ignore the case when the binlog connections
						// is corrupt
						w.binlogDelayTime.Set(0)
					}
					// Handle idle
					if 0 != idleTimeContinue && idleTimeContinue%(30*1000) == 0 {
						w.idle()
					}
				}
			}
		}
	}
}

func parseToData(ctx *SQLContext) {
	columnIdx := make(map[int]*moverconfig.FieldInfo)
	for k, column := range ctx.columns {
		columnIdx[k] = column
	}
	for k, data := range ctx.datas {
		if column, ok := columnIdx[k]; ok {
			if value, ok := data.(string); ok {
				switch column.FieldTypeSelf {
				case moverconfig.FDT_INT:
					var d interface{}
					var err error
					if column.IsUnsigned {
						d, err = utils.StringToUInt64(value)
					} else {
						d, err = utils.StringToInt64(value)
					}
					if err != nil {
						logging.LogError(logging.EC_Warn, "stringToInt64-err= %v,err=%v", data, err)
						continue
					}
					ctx.datas[k] = d
				case moverconfig.FDT_FLOAT:
					if moverconfig.GetConfig().FloatToString != 1 {
						d, err := utils.StringToFloat64(value)
						if err != nil {
							logging.LogError(logging.EC_Warn, "stringToFloat64-err= %v,err=%v", data, err)
							continue
						}
						ctx.datas[k] = d
					}
				case moverconfig.FDT_BINARY:
					ctx.datas[k] = []byte(value)
				}
			}
		}
	}
}

func (w *ExecuteWorker) idle() {
	// Send heartbeat and update broker address
	if err := w.client.HeartBeat(); nil != err {
		log.Warnf("Worker %d send heartbeat error: %v", w.ID, err)
	} else {
		log.Debugf("Worker %d heartbeat sent", w.ID)
	}
}
