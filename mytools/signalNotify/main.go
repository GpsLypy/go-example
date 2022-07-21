package main

import (
	"fmt"
	_ "os"
	_ "os/signal"
	"time"
)

// func main() {
// 	// 我们得使用带缓冲 channel
// 	// 否则，发送信号时我们还没有准备好接收，就有丢失信号的风险
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt)
// 	//time.Sleep(time.Second * 5)
// 	s := <-c
// 	fmt.Println("Got signal:", s)
// }

func main() {
	var waitFiveHundredMillisections time.Duration = 500 * time.Millisecond

	startingTime := time.Now().UTC()
	fmt.Println(startingTime)
	time.Sleep(600 * time.Millisecond)
	endingTime := time.Now().UTC()
	fmt.Println(endingTime)
	var duration time.Duration = endingTime.Sub(startingTime)

	if duration >= waitFiveHundredMillisections {
		//纳秒值除以 1^6 得到了 int64 类型表示的毫秒值
		fmt.Printf("Wait %v\nNative [%v]\nMilliseconds [%d]\nSeconds [%.3f]\n", waitFiveHundredMillisections, duration, duration.Nanoseconds()/1e6, duration.Seconds())
	}
}
