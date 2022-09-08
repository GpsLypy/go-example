package main

import (
	"time"

	"github.com/juju/errors"
)

func (s *source) start() error {
	timer := time.NewTimer(time.Second * 3)
	tickerC := 0
	for {
		select {
		case receiveTask, isOpen := <-TablesChan:
			if !isOpen {
				return nil
			} else {
				receiveTask := receiveTask
				go s.HandleTask(receiveTask)
			}
		case <-timer.C:
			tickerC++
			timer.Reset(time.Second)
			if tickerC%2 == 0 {

			}
		default:
			continue
		}
	}
	return nil
}

func (s *source) HandleTask(task syncTable) error {
	//atomic.AddInt32(&currState.DoingTables, 1)
	//pull数据时在实例化
	mysqlPolicy := &MysqlPolicy{}
	var FromField string
	LogInfo("table_Task: >>>> %+v <<<<", task)
	_, _, _, _, FieldItems, _ := GetFromFieldListfromMysql(task.queryDB, task.FromDBName, task.FromTable, FromField)
	ib := inbound.NewTableBuilder(task.FromDBName, task.FromTable)
	for _, v := range FieldItems {
		ib.AddColumn(v.Field, v.Type, v.PrimaryKey)
	}
	inboundEvent, err := ib.Convert()
	if err != nil {
		errors.Details(err)
	}
	s.parent.pushEvent(inboundEvent)

	err = GetAmoutOfTableData(&task)
	if err != nil {
		return err
	}
	// pull data to eventCh
	return s.pullDatas(task, mysqlPolicy)
}
