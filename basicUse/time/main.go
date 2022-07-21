package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool)

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 6)
			c <- false
		}

		time.Sleep(time.Second * 6)
		c <- true
	}()

	go func() {
		// try to read from channel, block at most 5s.
		// if timeout, print time event and go on loop.
		// if read a message which is not the type we want(we want true, not false),
		// retry to read.
		timer := time.NewTimer(time.Second * 5)
		for {
			// timer is active , not fired, stop always returns true, no problems occurs.
			//如果计时器已经过期或已经停止，则返回false。 停止不关闭通道，计时器过期
			//维护者所有活跃计时器的最小堆中已经不包含该计时器了，而此时time.C没有数据就被阻塞了
			// if !timer.Stop() {
			// 	select {
			// 	//计时器过期或者停止时执行以下代码
			// 	case <-timer.C: //try to drain from the cahnnel 里面不在有过期的值，避免意外情况
			// 	default:
			// 	}
			// }
			timer.Reset(time.Second * 5)
			select {
			case b := <-c:
				if b == false {
					fmt.Println(time.Now(), ":recv false. continue")
					continue
				}
				//we want true, not false
				fmt.Println(time.Now(), ":recv true. return")
				return
			case <-timer.C:
				fmt.Println(time.Now(), ":timer expired")
				continue
			}
		}
	}()

	//to avoid that all goroutine blocks.
	var s string
	fmt.Scanln(&s)
}
