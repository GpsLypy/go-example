package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("hello")

	//to do....
	go func() {
		panic("exit")
	}()

	//to do....

	time.Sleep(5 * time.Second)
}

// func main() {
// 	Go(func() {
// 		fmt.Println("hello")
// 		panic("exit")
// 	})
// 	time.Sleep(5 * time.Second)
// }

// func Go(x func()) {
// 	go func() {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				fmt.Println(err)
// 			}
// 		}()
// 		x()
// 	}()
// }

// type Message struct{

// }

// ch <- &Message{}

// //强依赖 配置 ->panic

// blocking 阻塞
// nonblocking 流量进来-》资源未准备好，进行连接，报一小波错，然后正常
// nonblocking + 超时 ：流量进来，资源未准备好，尝试连接几次，超时了在报错
