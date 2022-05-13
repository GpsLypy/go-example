package info

import (
	"app/tools/mover/help"
	"fmt"
)

var EtcdEndpointMap = map[string][]string{
	"dev":     []string{"http://10.100.156.210:2379"},
	"test":    []string{"http://10.100.38.15:2379"},
	"uat":     []string{"http://etcdgateway.sqlproxy.uat.17usoft.com"},
	"qa":      []string{"http://etcdgateway.sqlproxy.qa.17usoft.com"},
	"product": []string{"http://etcdgateway.sqlproxy.17usoft.com"},
	"stage":   []string{"http://etcdgateway.sqlproxy.17usoft.com"},
}

var CtrlServerAddrMap = map[string][]string{
	"dev":     []string{"http://10.100.156.207:7000"},
	"test":    []string{"http://10.100.156.207:7001"},
	"uat":     []string{"http://10.100.156.207:7002"},
	"qa":      []string{"http://10.100.156.207:7002"},
	"product": []string{"http://ctrl.sqlproxy.17usoft.com"},
	"stage":   []string{"http://ctrl.sqlproxy.17usoft.com"},
}

//EtcdEndpoint存储多个节点的服务地址
type EnvInfo struct {
	EnvType       string
	Host          string
	JobID         int64
	EtcdEndpoints []string
}

var envInfo EnvInfo

func InitEnvInfo(flagEnv, flagJobID string) error {
	envInfo.EnvType = help.GetEnvType(flagEnv)
	envInfo.Host = help.GetHost()
	envInfo.JobID = help.GetJobID(flagJobID)
	endpoints, hasKey := EtcdEndpointMap[envInfo.EnvType]
	if !hasKey {
		return fmt.Errorf("Can't get matched etcd endpoints")
	}
	envInfo.EtcdEndpoints = endpoints
	return nil
}

func GetEnvInfo() *EnvInfo {
	return &envInfo
}

func GetJobID() int64 {
	return envInfo.JobID
}
