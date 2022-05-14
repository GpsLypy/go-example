package main

import (
	"fmt"
	"time"
)

// 首先来看Or-Done模式。Or-Done模式是信号通知模式中更宽泛的一种模式。这里提到了“信号通知模式”，我先来解释一下。

// 我们会使用“信号通知”实现某个任务执行完成后的通知机制，在实现时，我们为这个任务定义一个类型为chan struct{}类型的done变量，等任务结束后，我们就可以close这个变量，然后，其它receiver就会收到这个通知。

// 这是有一个任务的情况，如果有多个任务，只要有任意一个任务执行完，我们就想获得这个信号，这就是Or-Done模式。

// 比如，你发送同一个请求到多个微服务节点，只要任意一个微服务节点返回结果，就算成功，这个时候，就可以参考下面的实现：

func or(channels ...<-chan interface{}) <-chan interface{} {
	//特殊情况只有一个或零个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	//无缓冲channel
	orDone := make(chan interface{})

	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2: //2个也是一种特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: //超过两个，二分法递归处理
			m := len(channels) / 2
			select {
			case <-or(channels[:m]...):
			case <-or(channels[m:]...):
			}
		}
	}()

	return orDone
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()
	<-or(
		sig(10*time.Second),
		sig(20*time.Second),
		sig(30*time.Second),
		sig(40*time.Second),
		sig(50*time.Second),
		sig(01*time.Minute),
	)
	fmt.Printf("done  after %v", time.Since(start))
}

//这里的实现使用了一个巧妙的方式，当chan的数量大于2时，使用递归的方式等待信号。

//在chan数量比较多的情况下，递归并不是一个很好的解决方式，根据这一讲最开始介绍的反射的方法，我们也可以实现Or-Done模式：
