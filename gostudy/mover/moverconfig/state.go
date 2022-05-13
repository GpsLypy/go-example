package moverconfig

import "sync"

const (
	AdditionDataKeyMin = iota
	ADK_ROW_COUNT
	// append here
	AdditionDataKeyMax
)

//ProgressRatio进度条
type MoverState struct {
	// base
	IsHealthy     bool
	Started       bool
	ProgressRatio int
	ErrMsg        string
	Prepared      bool
	SumTasks      int32
	DoingTasks    int32
	DoneTasks     int32

	// dispatcher
	DataChanOutRows uint64

	RdsMessage string
	FromHost   map[string]struct{}

	// addition
	AdditionData [AdditionDataKeyMax]int32
	Tasks        []string
}

var (
	currState      MoverState
	currStateMutex sync.Mutex
)

func GetCurrState() *MoverState {
	return &currState
}

func init() {
	currState.IsHealthy = true
	currState.ProgressRatio = 0
	currState.FromHost = make(map[string]struct{})
}

func GetState() MoverState {
	return currState
}

func AppendStateFromHost(host string) {
	if currState.FromHost == nil {
		currState.FromHost = make(map[string]struct{})
	}
	currStateMutex.Lock()
	defer currStateMutex.Unlock()

	currState.FromHost[host] = struct{}{}
	return
}

func AppendTask(task string) {
	currStateMutex.Lock()
	defer currStateMutex.Unlock()

	currState.Tasks = append(currState.Tasks, task)
	return
}
