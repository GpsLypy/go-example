package main

import (
	"fmt"
	"time"
)

//演示数据传递=>击鼓传花
//任务编排的题 其实它就可以用数据传递的方式实现。

//有4个goroutine，编号为1、2、3、4。每秒钟会有一个goroutine打印出它自己的编号，要求你编写程序，
//让输出的编号总是按照1、2、3、4、1、2、3、4……这个顺序打印出来。

type Token struct{}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch         //取得令牌
		fmt.Println((id + 1)) //id 从1开始
		time.Sleep(time.Second)
		nextCh <- token
	}
}
func main() {
	//为了实现顺序的数据传递，我们可以定义一个令牌的变量，
	//谁得到令牌，谁就可以打印一次自己的编号，同时将令牌传递给下一个goroutine，我们尝试使用chan来实现，可以看下下面的代码。
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}

	//创建4个worker
	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}

	//首先把令牌交给第一个worker
	chs[0] <- struct{}{}
	select {}
}

//这类场景有一个特点，就是当前持有数据的goroutine都有一个信箱，信箱使用chan实现，
//goroutine只需要关注自己的信箱中的数据，处理完毕后，就把结果发送到下一家的信箱中。
