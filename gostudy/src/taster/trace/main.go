package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {
	// NewProcess 会返回一个持有PID的Process对象，方法会检查PID是否存在，如果不存在会返回错误
	// 通过Process对象上定义的其他方法我们可以获取关于进程的各种信息。
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic(err)
	}

	// 返回指定时间内进程占用CPU时间的比例
	cpuPercent, err := p.Percent(time.Second)
	if err != nil {
		panic(err)
	}
	// 上面返回的是占所有CPU核心时间的比例，如果想更直观的看占比，可以算一下占单个核心的比例
	cp := cpuPercent / float64(runtime.NumCPU())

	// 获取进程占用内存的比例
	mp, _ := p.MemoryPercent()

	// 创建的线程数
	threadCount := pprof.Lookup("threadcreate").Count()

	// Goroutine数

	gNum := runtime.NumGoroutine()

	time.Sleep(time.Second * 2)

	fmt.Println(cpuPercent, cp, mp, threadCount, gNum)

}
