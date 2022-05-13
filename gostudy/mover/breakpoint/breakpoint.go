package breakpoint

import (
	"app/tools/mover/logging"
	"app/tools/mover/moverconfig"
	"app/tools/mover/utils"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"sync"
)

const DefaultDataDir = "/data/data_distribute"

type BreakPoint struct {
	// source
	FromDBType int    `json:"FromDBType"`
	FromDBName string `json:"FromDBName"`
	FromTable  string `json:"FromTable"`
	OrgDBName  string `json:"OrgDBName"`
	OrgTbName  string `json:"OrgTbName"`

	// task data
	PrimaryKey     string `json:"PrimaryKey"`
	PrimaryKeyType int    `json:"PrimaryKeyType"`
	ClosedInterval int    `json:"ClosedInterval"`

	StrStartId string                       `json:"StrStartId"`
	StrEndId   string                       `json:"StrEndId"`
	StartId    int64                        `json:"StartId"`
	EndId      int64                        `json:"EndId"`
	TaskInd    int64                        `json:"TaskInd"`
	Tasks      []moverconfig.Task           `json:"tasks"`
	TaskMap    map[string]*moverconfig.Task `json:"task_map"`
}

var (
	BpFile     string
	mapBPoints map[string]*BreakPoint
	mtxBPoints sync.RWMutex
)

func GetBreakPointFile(datadir string, jobId int64) string {
	// dir := filepath.Dir(datadir)
	if datadir == "" {
		datadir = DefaultDataDir
	}
	datadir += "/" + strconv.Itoa(int(jobId))
	if _, err := os.Stat(datadir); os.IsNotExist(err) {
		os.MkdirAll(datadir, 0777)
		seelog.Infof("create datadir %s", datadir)
	}

	return fmt.Sprintf("%s/breakpoints.json", datadir)
}

func (bp *BreakPoint) TakeTask(taskId int64) *moverconfig.Task {

	if len(bp.Tasks) == 0 {
		return nil
	}

	for i, task := range bp.Tasks {
		if task.TaskID == taskId {

			if i == 0 {
				bp.Tasks = bp.Tasks[i+1:]
			} else {
				bp.Tasks = append(bp.Tasks[:i-1], bp.Tasks[i+1:]...)
			}
			return &task
		}
	}

	return nil
}

func (bp *BreakPoint) getKey() string {
	return fmt.Sprintf("%s.%s", bp.FromDBName, bp.FromTable)
}

func (bp *BreakPoint) freshBreakPoint(task *moverconfig.Task) {
	if bp.TaskInd != task.TaskID {
		bp.Tasks = append(bp.Tasks, *task)
	} else {
		curTask := task
		for {
			nextTask := bp.TakeTask(curTask.TaskID + 1)
			if nil == nextTask {
				break
			}

			curTask = nextTask
		}

		bp.TaskInd = curTask.TaskID + 1
		bp.FromDBType = curTask.FromDBType
		bp.FromDBName = curTask.FromDBName
		bp.FromTable = curTask.FromTable
		bp.OrgDBName = curTask.OrgDBName
		bp.OrgTbName = curTask.OrgTbName
		bp.PrimaryKeyType = curTask.PrimaryKeyType
		bp.PrimaryKey = curTask.PrimaryKey
		bp.StrStartId = curTask.StrStartId
		bp.StrEndId = curTask.StrEndId
		bp.StartId = curTask.StartId
		bp.EndId = curTask.EndId
	}
}

func InitBreakPoint(jobID int64) {
	dirname := moverconfig.GetConfig().DataDir
	BpFile = GetBreakPointFile(dirname, jobID)

	mapBPoints = make(map[string]*BreakPoint)

	if err := loadBreakPointsIfNeeded(); nil != err {
		logging.LogError(logging.EC_Runtime_ERR, "fail to load breakpoint, err: %v", err)
	}
}

func FlushBreakPoint() error {
	mtxBPoints.RLock()
	defer mtxBPoints.RUnlock()

	if len(mapBPoints) == 0 {
		seelog.Info("No breakpoints.")
		return nil
	}

	data, err := json.Marshal(mapBPoints)
	if nil != err {
		logging.LogError(logging.EC_Encode_ERR, "Marshal fail, err: %v", err)
		return err
	}

	f, err := os.Create(BpFile)
	if err != nil {
		logging.LogError(logging.EC_Runtime_ERR, "Create file %s, err: %v", BpFile, err)
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if nil != err {
		logging.LogError(logging.EC_Runtime_ERR, "flush breakpoint fail, err: %v", err)
		return err
	}

	return nil
}

func ExistFile(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func loadBreakPointsIfNeeded() error {
	if !ExistFile(BpFile) {
		return nil
	}

	data, err := utils.ReadData(BpFile)
	if nil != err {
		logging.LogError(logging.EC_Runtime_ERR, "read breakpoint file fail, err: %v", err)
		return err
	}

	m := &map[string]*BreakPoint{}
	err = json.Unmarshal(data, m)
	if nil != err {
		logging.LogError(logging.EC_Decode_ERR, "Json unmarshal, err=%v", err)
		return err
	}

	mapBPoints = *m
	return nil
}

func FinishTask(DBName, TableName string, task *moverconfig.Task) error {
	tablekey := fmt.Sprintf("%s.%s", DBName, TableName)

	mtxBPoints.Lock()
	defer mtxBPoints.Unlock()

	if nil == task {
		if _, ok := mapBPoints[tablekey]; !ok {
			mapBPoints[tablekey] = &BreakPoint{
				FromDBName: DBName,
				FromTable:  TableName,
				TaskMap:    make(map[string]*moverconfig.Task),
			}
		}
	} else {
		bp, ok := mapBPoints[tablekey]
		if !ok {
			mapBPoints[tablekey] = &BreakPoint{
				FromDBName: DBName,
				FromTable:  TableName,
			}
			return errors.New(fmt.Sprintf("No breakpoint %s in cache.", tablekey))
		}
		bp.freshBreakPoint(task)
	}

	return nil
}

func FindBreakPoint(DBName, TableName string) *BreakPoint {
	tablekey := fmt.Sprintf("%s.%s", DBName, TableName)

	mtxBPoints.Lock()
	defer mtxBPoints.Unlock()

	if bp, ok := mapBPoints[tablekey]; ok {
		return bp
	}

	return nil
}
