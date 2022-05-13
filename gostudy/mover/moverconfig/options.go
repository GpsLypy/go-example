package moverconfig

import (
	"app/platform/etcd-gateway/rdsdef"
	"app/tools/mover/logging"
	"app/tools/mover/utils"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func GetConfig() *MoverConfig {
	return currConfig
}

func GetConfigData(configType int) ([]byte, error) {
	var cfgByte []byte
	var err error

	switch configType {
	case CT_JSON:
		{
			cfg := GetConfig()
			if nil == cfg {
				return nil, errors.New("Nil config")
			}

			cfgByte, err = json.Marshal(*cfg)
			if nil != err {
				//LogError(EC_Config_ERR, "Json marshal in getConfigData, err=%v", err)
				logging.LogError(logging.EC_Config_ERR, "Json marshal in getConfigData, err=%v", err)
				currState.IsHealthy = false
				currState.ErrMsg = fmt.Sprintf("Json marshal in getConfigData, err=%v", err)
				return nil, err
			}
		}
	case CT_BINARY:
		{
			return nil, errors.New("Don't support this format(binary)")
		}
	default:
		{
			err := errors.New("Wrong config type")
			return nil, err
		}
	}

	return cfgByte, nil
}

func CheckConfig(configType int, configData []byte) error {
	var err error
	switch configType {
	case CT_JSON:
		{
			config := &MoverConfig{}
			err = json.Unmarshal(configData, config)
			if nil != err {
				//LogError(EC_Decode_ERR, "Json unmarshal, err=%v", err)
				logging.LogError(logging.EC_Decode_ERR, "Json unmarshal, err=%v", err)
				currState.IsHealthy = false
				currState.ErrMsg = fmt.Sprintf("Json unmarshal, err=%v", err)
				return err
			}
			if len(config.TaskList) == 0 || config.WorkerSize == 0 {
				//LogError(EC_Config_ERR, "Invalid config, config: %s", string(configData))
				logging.LogError(logging.EC_Config_ERR, "Invalid config, config: %s", string(configData))
				currState.IsHealthy = false
				currState.ErrMsg = fmt.Sprintf("Invalid config, config: %s", string(configData))
				return errors.New("Invalid config")
			}
		}
	case CT_BINARY:
		{
			return errors.New("Don't support this format(binary)")
		}
	default:
		{
			return errors.New("Wrong config type")
		}
	}

	return nil
}

func SetConfig(configType int, configData []byte) error {
	configMutex.Lock()
	defer func() {
		configMutex.Unlock()
	}()

	var err error
	switch configType {
	case CT_JSON:
		{
			config := &MoverConfig{}
			err = json.Unmarshal(configData, config)
			if nil != err {
				//LogError(EC_Decode_ERR, "Json unmarshal, err=%v", err)
				logging.LogError(logging.EC_Decode_ERR, "Json unmarshal, err=%v", err)
				currState.IsHealthy = false
				currState.ErrMsg = fmt.Sprintf("Json unmarshal, err=%v", err)
				return err
			}
			if len(config.TaskList) <= 0 || config.WorkerSize <= 0 {
				//LogError(EC_Config_ERR, "Invalid config, config: %s,len(config.TaskList)=%v,WorkerSize=%v", string(configData), len(config.TaskList), config.WorkerSize)
				logging.LogError(logging.EC_Config_ERR, "Invalid config, config: %s,len(config.TaskList)=%v,WorkerSize=%v", string(configData), len(config.TaskList), config.WorkerSize)
				currState.IsHealthy = false
				currState.ErrMsg = fmt.Sprintf("Invalid config, config: %s,len(config.TaskList)=%v,WorkerSize=%v", string(configData), len(config.TaskList), config.WorkerSize)
				return errors.New("Invalid config")
			}

			for i, _ := range config.TaskList {
				dsp := &config.TaskList[i]
				dsp.DestTo = config.Target
			}

			currConfig = config
		}
	case CT_BINARY:
		{
			return errors.New("Don't support this format(binary)")
		}
	default:
		{
			return errors.New("Wrong config type")
		}
	}

	return nil
}

func CheckRdsRsp(rdsConfig *rdsdef.DsGetRsp) error {
	if nil == rdsConfig {
		return errors.New("Rds config nil")
	}

	if len(rdsConfig.InstanceList) == 0 {
		seelog.Errorf("Invalid rds config, config: %v", *rdsConfig)
		return errors.New("Invalid rds config")
	}

	return nil
}

func SetRdsConfig(rdsData []byte) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	rdsRsp := &rdsdef.DsGetRsp{}
	err := json.Unmarshal(rdsData, rdsRsp)
	if nil != err {
		//LogError(EC_Decode_ERR, "Json unmarshal rds Response, err=%v", err)
		logging.LogError(logging.EC_Decode_ERR, "Json unmarshal rds Response, err=%v", err)
		currState.IsHealthy = false
		currState.ErrMsg = fmt.Sprintf("Json unmarshal rds Response, err=%v", err)
		return err
	}

	err = CheckRdsRsp(rdsRsp)
	if nil != err {
		return err
	}

	currRdsDS = rdsRsp
	return nil
}

func getDetailedEndpoint(host string, port int, isDs bool) Endpoint {
	var user, pwd string

	if isDs {
		user = currConfig.RdsConf.User.UserName
		pwd = currConfig.RdsConf.User.Password
	} else {
		user = currConfig.RdsConf.BackendUser.UserName
		pwd = currConfig.RdsConf.BackendUser.Password
	}

	if len(currConfig.RdsConf.DetailedEndpoints) == 0 {
		return Endpoint{
			Host:     host,
			Port:     port,
			User:     user,
			Password: pwd,
		}
	}

	for _, ep := range currConfig.RdsConf.DetailedEndpoints {
		if ep.Host == host && ep.Port == port {
			return ep
		}
	}

	return Endpoint{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pwd,
	}
}

func getEndpointsFromInstanceGroup(instGroup rdsdef.RdsInstances, isSlave, isDs bool) ([]Endpoint, error) {
	var endpoints []Endpoint
	var endpoint, slaveToolEP, slaveBakEP Endpoint
	var user, pwd string

	if isDs {
		user = currConfig.RdsConf.User.UserName
		pwd = currConfig.RdsConf.User.Password
	} else {
		user = currConfig.RdsConf.BackendUser.UserName
		pwd = currConfig.RdsConf.BackendUser.Password
	}

	// get data source from instance group
	if len(instGroup.Instances) == 0 {
		return endpoints, errors.New("no instances in RdsInstances group.")
	}

	for _, inst := range instGroup.Instances {
		port, _ := strconv.Atoi(inst.Port)

		if isSlave && strings.ToUpper(inst.Role) == utils.MasterRole {
			continue
		}

		if len(currConfig.RdsConf.DetailedEndpoints) == 0 {
			endpoint = Endpoint{
				Host:     inst.IP,
				Port:     port,
				User:     user,
				Password: pwd,
			}
		} else {
			endpoint = getDetailedEndpoint(inst.IP, port, isDs)
		}

		if utils.IsToolSlave(inst.Role) {
			slaveToolEP = endpoint
		} else if utils.IsBakSlave(inst.Role) {
			slaveBakEP = endpoint
		} else if utils.IsMaster(inst.Role) {
			endpoints = append(endpoints, endpoint)
		} else {
			s := []Endpoint{endpoint}
			endpoints = append(s, endpoints...)
		}
	}

	// slavetool, slavebak, slave..., master
	h := []Endpoint{}
	if slaveToolEP.Validate() {
		h = append(h, slaveToolEP)
	}
	if slaveBakEP.Validate() {
		h = append(h, slaveBakEP)
	}

	if len(h) > 0 {
		endpoints = append(h, endpoints...)
	}

	return endpoints, nil
}

func GetRdsDataSourceEndpoints() ([]Endpoint, error) {
	var endpoints []Endpoint
	//var endpoint Endpoint
	var err error = nil

	// pull data from backend mysql instance, so neither master or slave for dbrouter
	if len(currRdsDS.DBRNodes) > 0 {

		/*for _, n := range currRdsDS.DBRNodes {
			port, _ := strconv.Atoi(n.Port)

			if len(currConfig.RdsConf.DetailedEndpoints) == 0 {
				endpoint = Endpoint{
					Host:     n.IP,
					Port:     port,
					User:     currConfig.RdsConf.User.UserName,
					Password: currConfig.RdsConf.User.Password,
				}
			} else {
				endpoint = getDetailedEndpoint(n.IP, port, true)
			}

			endpoints = append(endpoints, endpoint)
		}*/
	} else {
		if len(currRdsDS.InstanceList) == 0 {
			return endpoints, errors.New("InstanceList is empty.")
		}

		// get data source from first group, if more than one group
		if len(currRdsDS.InstanceList[0].Instances) == 0 {
			return endpoints, errors.New("no instances in InstanceList[0].")
		}

		endpoints, err = getEndpointsFromInstanceGroup(currRdsDS.InstanceList[0], GetConfig().IsMaster == 0, false)
		if nil != err {
			return endpoints, err
		}
	}

	return endpoints, err
}

func GetRdsShardingBackendEndpoints() ([][]Endpoint, error) {
	var EndpointsSlice [][]Endpoint
	var err error = nil

	if len(currRdsDS.InstanceList) == 0 {
		return EndpointsSlice, errors.New("no instanceList in the RdsRsp.")
	}

	for _, instGroup := range currRdsDS.InstanceList {
		var Eps []Endpoint
		Eps, err = getEndpointsFromInstanceGroup(instGroup, GetConfig().IsMaster == 0, false)
		if nil != err {
			return EndpointsSlice, err
		}

		EndpointsSlice = append(EndpointsSlice, Eps)
	}

	return EndpointsSlice, nil
}
