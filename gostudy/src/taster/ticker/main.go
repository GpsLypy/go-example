package main

import (
	"fmt"
	"time"
)

func foo() {
	fmt.Printf("foo() start.")
	time.Sleep(time.Second * 3)
	fmt.Println("foo() end")
}

//当定时任务执行时间过长且超过定时的间隔时间时，
//定时的间隔时间到了之后会等待定时任务执行完才会进行下一轮的定时任务.
