package help

import (
	"fmt"
	"time"

	"app/tools/mover/logging"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var (
	etcdApi client.KeysAPI
)

func GetEtcdJobConfigKey(jobID int64) string {
	return fmt.Sprintf("/data_transfer/config/DT_CFG_%d", jobID)
}

func GetEtcdJobStatusKey(jobID int64) string {
	return fmt.Sprintf("/data_transfer/status/DT_STATUS_%d/status", jobID)
}

func CreateEtcd(etcdEndpointList []string) error {
	//要创建client，它需要传入一个Config配置
	//Endpoints:etcd的多个节点服务地址
	etcdCfg := client.Config{
		Endpoints:               etcdEndpointList,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second * 5,
	}
	//创建访问etcd服务的客户端
	etcdClient, err := client.New(etcdCfg)
	if err != nil {
		return err
	}
	etcdApi = client.NewKeysAPI(etcdClient)
	return nil
}

func SetEtcdKey(k string, v string, ttl int, first bool) error {
	logging.LogDebug("Set k=%s, v=%s", k, v)
	var opt client.SetOptions
	if !first {
		opt.PrevExist = client.PrevExist
	}
	if ttl > 0 {
		opt.TTL = time.Second * time.Duration(ttl)
	}
	_, err := etcdApi.Set(context.Background(), k, v, &opt)
	if nil != err {
		if !first {
			// If key not found, register to etcd again
			if etcdErr, ok := err.(client.Error); ok {
				if etcdErr.Code == client.ErrorCodeKeyNotFound {
					// Register again
					opt.PrevExist = client.PrevIgnore
					_, err := etcdApi.Set(context.Background(), k, v, &opt)
					return err
				}
			}
		}
		return err
	}

	return nil
}

func GetEtcdKey(k string) (string, error) {
	resp, err := etcdApi.Get(context.Background(), k, nil)
	if nil != err {
		return "", err
	}

	return resp.Node.Value, nil
}

func DelEtcdKey(k string) error {
	_, err := etcdApi.Delete(context.Background(), k, nil)
	if nil != err {
		return err
	}

	return nil
}
