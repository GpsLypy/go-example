package service

import (
	"app/tools/mover/help"
	"app/tools/mover/info"
	"app/tools/mover/logging"
	"app/tools/mover/moverconfig"
	"app/tools/mover/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var (
	firstUploadStatus bool
)

type UploadStatus struct {
	JobID         string   `json:"jobID"`
	Time          int64    `json:"time"`
	Healthy       int      `json:"health"`
	ProgressRatio int      `json:"ratio"`
	Prepared      int      `json:"Prepared"`
	Finish        int      `json:"Finish"`
	Tool          int      `json:"tool"`
	SumTasks      int      `json:"SumTasks"`
	DoneTasks     int32    `json:"DoneTasks"`
	DoingTasks    int32    `json:"DoingTasks"`
	ErrMsg        string   `json:"message"`
	RdsMessage    string   `json:"rdsMessage"`
	FromHost      []string `json:"fromHost"`

	Stats []moverconfig.TableStat `json:"Stats"`
	Tasks []string                `json:"tasks"`
}

type JobMonitorStatus struct {
	JobID string `json:"jobid"`
	Time  int64  `json:"time"`
}

func init() {
	firstUploadStatus = true
}

func GetUploaderTime() int {
	if firstUploadStatus {
		logging.LogDebug("First status not sent")
		return 1
	}

	return info.DefaultUploadTime
}

func UploadStatusInfo(errMsg string) {
	cfg := moverconfig.GetConfig()
	if nil == cfg {
		logging.LogWarn("Skip upload, service config nil")
		return
	}

	if "" != errMsg {
		//currState.ErrMsg = errMsg
		(*moverconfig.GetCurrState()).ErrMsg = errMsg
	}

	// get redis status
	var status UploadStatus
	status.JobID = strconv.FormatInt(info.GetEnvInfo().JobID, 10)
	status.Time = time.Now().Unix()
	if moverconfig.GetState().IsHealthy {
		status.Healthy = 1
	} else {
		status.Healthy = 0
	}
	status.Prepared = utils.Bool2int(moverconfig.GetState().Prepared)
	status.ProgressRatio = moverconfig.GetState().ProgressRatio
	status.Tool = info.ToolType
	status.ErrMsg = moverconfig.GetState().ErrMsg
	status.SumTasks = int(moverconfig.GetState().SumTasks)
	status.DoneTasks = moverconfig.GetState().DoneTasks
	status.DoingTasks = moverconfig.GetState().DoingTasks

	if moverconfig.GetState().IsHealthy && status.DoneTasks >= int32(status.SumTasks) && moverconfig.GetState().Prepared {
		status.Finish = 1
	} else {
		status.Finish = 0
	}

	dataLog, err := json.Marshal(&status)
	if nil != err {
		logging.LogError(logging.EC_Encode_ERR, "Marshal err, status: %v err: %v", status, err)
		(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Marshal err, status: %v err: %v", status, err)
		(*moverconfig.GetCurrState()).IsHealthy = false
		return
	}
	status.RdsMessage = moverconfig.GetState().RdsMessage
	hosts := moverconfig.GetState().FromHost
	for host, _ := range hosts {
		status.FromHost = append(status.FromHost, host)
	}
	status.Tasks = moverconfig.GetState().Tasks
	status.Stats = GetTableInfos()
	data, err := json.Marshal(&status)
	if nil != err {
		logging.LogError(logging.EC_Encode_ERR, "Marshal err, status: %v err: %v", status, err)
		(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Marshal err, status: %v err: %v", status, err)
		(*moverconfig.GetCurrState()).IsHealthy = false
		return
	}

	// etcd status
	etcdKey := help.GetEtcdJobStatusKey(info.GetEnvInfo().JobID)
	err = help.SetEtcdKey(etcdKey, string(data), info.EtcdStatusTtl, false)
	if nil != err {
		logging.LogError(logging.EC_Ignore_ERR, "Failed to upload status to etcd, err: %s", err.Error())
		(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Failed to upload status to etcd, err: %s", err.Error())
		(*moverconfig.GetCurrState()).IsHealthy = false

	}

	logging.LogInfo("Upload status: %s", string(dataLog))
	firstUploadStatus = false
}

func UploadSummary(errMsg string) {
	cfg := moverconfig.GetConfig()
	if nil == cfg {
		logging.LogWarn("Skip upload, service config nil")
		return
	}

	if "" != errMsg {
		//currState.ErrMsg = errMsg
		(*moverconfig.GetCurrState()).ErrMsg = errMsg
	}

	// get redis status
	var status UploadStatus
	status.JobID = strconv.FormatInt(info.GetEnvInfo().JobID, 10)
	status.Time = time.Now().Unix()
	status.Healthy = utils.Bool2int(moverconfig.GetState().IsHealthy)
	status.Prepared = utils.Bool2int(moverconfig.GetState().Prepared)
	status.ProgressRatio = moverconfig.GetState().ProgressRatio
	status.Tool = info.ToolType
	status.ErrMsg = moverconfig.GetState().ErrMsg
	status.SumTasks = int(moverconfig.GetState().SumTasks)
	status.DoneTasks = moverconfig.GetState().DoneTasks
	status.DoingTasks = moverconfig.GetState().DoingTasks
	status.Stats = GetTableInfos()

	if moverconfig.GetState().IsHealthy && status.DoneTasks >= int32(status.SumTasks) && moverconfig.GetState().Prepared {
		status.Finish = 1
	} else {
		status.Finish = 0
	}

	status.RdsMessage = moverconfig.GetState().RdsMessage
	hosts := moverconfig.GetState().FromHost
	for host, _ := range hosts {
		status.FromHost = append(status.FromHost, host)
	}

	data, err := json.Marshal(&status)
	if nil != err {
		logging.LogError(logging.EC_Encode_ERR, "Marshal err, status: %v err: %v", status, err)
		(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Marshal err, status: %v err: %v", status, err)
		(*moverconfig.GetCurrState()).IsHealthy = false
		return
	}

	// etcd status
	etcdKey := help.GetEtcdJobStatusKey(info.GetEnvInfo().JobID)
	err = help.SetEtcdKey(etcdKey, string(data), info.EtcdStatusTtl, false)
	if nil != err {
		logging.LogError(logging.EC_Ignore_ERR, "Failed to upload status to etcd, err: %s", err.Error())
		(*moverconfig.GetCurrState()).ErrMsg = fmt.Sprintf("Failed to upload status to etcd, err: %s", err.Error())
		(*moverconfig.GetCurrState()).IsHealthy = false
	}

	firstUploadStatus = false
}
