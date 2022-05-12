package main

import (
	"sync/atomic"
	"time"
)

//典型的使用场景：应用程序定期的从外界获取最新的配置信息，然后更改自己内存中维护的配置变量。工作线程根据最新的配置来处理请求。

func loadConfig() map[string]string {
	//从数据库或文件系统中读取配置信息，然后以map的形式存放在内存里
	return make(map[string]string)
}

func requests() chan int {
	//将从外界中接收到的请求放到channel里
	return make(chan int)
}

func main() {
	//config 变量用来存放该服务的配置信息
	var config atomic.Value
	//初始化时从别的地方加载配置文件，并存到config变量里
	config.Store(loadConfig())
	go func() {
		for {
			//每10秒钟定时的拉取最新的配置信息，并且更新到config变量里
			time.Sleep(10 * time.Second)
			//对应于赋值操作config=loadConfig()
			config.Store(loadConfig())
		}
	}()

	//创建工作线程，每个工作线程都会根据他所读取到的最新的配置来处理请求
	for i := 0; i < 10; i++ {
		go func() {
			for r := range requests() {
				//对应于取值操作c :=config
				//由于Load()返回的是一个interface{}类型，所以我们要先强制转换一下
				c := config.Load().(map[string]string)
				//这里是根据配置信息处理请求的逻辑...
				_, _ = r, c
			}
		}()
	}
}

//type Value struct{
//automic.Value是一个结构体类型
