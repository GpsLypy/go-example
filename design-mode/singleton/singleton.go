package singleton

import (
	"sync"
	"sync/atomic"
)

//饿汉式单例
//注意定义非导出类型

type databaseConn struct {
	//
}

var dbconn *databaseConn

func init() {
	dbconn = &databaseConn{}
}

func GetInstance1() *databaseConn {
	return dbconn
}

//懒汉模式
//java用双重锁实现，volatile
//go中实现可以考虑定义一个实例的状态变量，然后用原子操作atmic.Load() 、 atomic.Store() 去读写这个状态变量

var initialized uint32

type singleton struct {
	//
}

func GetInstance() *singleton {
	if atomic.LoadUint32(&initialized) == 1 { //原子操作
		return instance
	}
	mu.Lock()
	defer mu.Unlock()
	if initialized == 0 {
		instance = &singleton{}
		atomic.StoreUint32(&initialized, 1)
	}
	return instance
}

//懒汉模式写法二

type singleton struct{}

var instance *singleton

var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}
