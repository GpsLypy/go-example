package main

import (
	"fmt"
	"time"
)

func printones() {
	logTicker := time.Tick(1 * time.Second)
	channel := make(chan struct{}, 1)
	go func() {
		for {
			//1s中之内只能选择其中一个两个分支之一
			select {
			case <-logTicker:
				channel <- struct{}{}

			case <-channel:
				fmt.Println("got it")
			}
		}
	}()
	<-make(chan struct{})
}

func main() {
	fmt.Println("begin")
	printones()
	fmt.Println("commit")
}
