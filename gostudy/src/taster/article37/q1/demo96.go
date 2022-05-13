package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/pprof"
	"taster/article37/common"
	"taster/article37/common/op"
)

var (
	profileName = "cpuprofile.out"
)

func main() {
	f, err := common.CreateFile("", profileName)
	if err != nil {
		fmt.Printf("CPU profile creation error: %v\n", err)
		return
	}
	defer f.Close()
	//在我们想让程序开始对CPU概要信息进行采样的时候，需要调用StartCPUProfile函数
	if err := startCPUProfile(f); err != nil {
		fmt.Printf("CPU profile start error: %v\n", err)
		return
	}
	if err = common.Execute(op.CPUProfile, 10); err != nil {
		fmt.Printf("execute error: %v\n", err)
		return
	}
	stopCPUProfile()
}

func startCPUProfile(f *os.File) error {
	if f == nil {
		return errors.New("nil file")
	}
	return pprof.StartCPUProfile(f)
}

func stopCPUProfile() {
	pprof.StopCPUProfile()
}
