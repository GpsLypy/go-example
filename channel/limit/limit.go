package main

import (
	"fmt"
	"time"
)

func Run(task_id, sleeptime, timeout int, ch chan string) {
	ch_run := make(chan string)
	go run(task_id, sleeptime, ch_run)
	select {
	case res := <-ch_run:
		ch <- res
	case <-time.After(time.Duration(timeout) * time.Second):
		res := fmt.Sprintf("task id %d , timeout", task_id)
		ch <- res
	}
}
func run(task_id, sleeptime int, ch chan string) {
	time.Sleep(time.Duration(sleeptime) * time.Second)
	ch <- fmt.Sprintf("task id %d , sleep %d second", task_id, sleeptime)
}

func main() {
	input := []int{300, 200, 100} //准备起3个goroutine
	timeout := 2
	//可以用一个 bool 类型的带缓冲 channel 作为并发限制的计数器
	chLimit := make(chan bool, 1)
	chstring := make([]chan string, len(input))
	limitFunc := func(chLimit chan bool, ch chan string, task_id, sleeptime, timeout int) {
		//time.Sleep(30 * time.Second)
		Run(task_id, sleeptime, timeout, ch)
		//让并发的 goroutine在执行完成后把这个 channel 里的东西给读走。
		//这样整个并发的数量就控制在这个 channel的缓冲区大小上。
		//他在执行完后，会把 chLimit的缓冲区里给消费掉一个。
		<-chLimit
	}

	startTime := time.Now()
	fmt.Println("Multirun start")

	for i, sleeptime := range input {
		chstring[i] = make(chan string, 1)
		//然后在并发执行的地方，每创建一个新的 goroutine，都往 chLimit 里塞个东西。
		//这样一来，当创建的 goroutine 数量到达 chLimit 的缓冲区上限后。主 goroutine 就挂起阻塞了，
		//直到这些 goroutine 执行完毕，
		//消费掉了 chLimit 缓冲区中的数据，程序才会继续创建新的 goroutine 。我们并发数量限制的目的也就达到了。
		chLimit <- true
		fmt.Println(111)
		//chstring监控执行结果，可能是产品可能是超时信息
		go limitFunc(chLimit, chstring[i], i, sleeptime, timeout)
	}

	for _, chresString := range chstring {
		fmt.Println(<-chresString)
	}

	endTime := time.Now()
	fmt.Printf("Multissh finished.Process time %s.Number of task is %d", endTime.Sub(startTime), len(input))
}
