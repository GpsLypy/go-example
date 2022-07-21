package main

import (
	"fmt"
	_ "github.com/pkg/errors"
	"time"
)

func main() {
	fmt.Println("vim-go")
	Go(func() {
		fmt.Println("hello")
		panic("panic:一路向西")
	})

	time.Sleep(5 * time.Second)
}

func Go(x func()) {
	//以下为并行调用
	// go func() {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	}()
	// 	x()
	// }()
	//以下为串行调用方式，会在panic(致命错误)之后正常运行睡眠函数
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	x()
}
