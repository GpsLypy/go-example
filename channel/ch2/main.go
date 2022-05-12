package main

import (
	"fmt"
	_ "time"

	"github.com/eapache/channels"
)

//演示信息交流
//worker池
//在 go 中，是无法往一个 nil 的 channel 中发送元素的。例如
// func main() {
// 	var c chan interface{}
// 	select {
// 	case c <- 1:
// 	}
// }

// func main() {
// 	var c chan interface{}
// 	select {
// 	case c <- 1:
// 	default:
// 		fmt.Println("hello world")
// 	}
// }

// func main() {
// 	//all goroutines are asleep - deadlock!
// 	c := channels.NewInfiniteChannel()
// 	go func() {
// 		for i := 0; i < 20; i++ {
// 			c.In() <- i
// 			time.Sleep(time.Second)
// 		}
// 	}()

// 	for i := 0; i < 50; i++ {
// 		val := c.Out()
// 		fmt.Println(val)
// 		time.Sleep(time.Millisecond * 500)
// 	}

// 	select {}
// }

func main() {
	taskChan := channels.NewInfiniteChannel()
	select {
	case v := <-taskChan.Out():
		fmt.Println("拿到一个任务！")
		fmt.Println(v)
	default:
		fmt.Println("么有任务哦")
	}
	fmt.Println("要退出了.....")
	select {}
}
