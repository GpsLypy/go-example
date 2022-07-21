package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func operation1(ctx context.Context) error {
	//让我们假设这个操作会因为某种原因失败
	//我们使用time.Sleep来模拟一个资源密集型操作
	time.Sleep(1000 * time.Millisecond)
	return errors.New("failed")
}

func operation2(ctx context.Context) {
	//我们使用在前面HTTP服务器例子里使用过的类似模式
	select {
	case <-time.After(5000 * time.Millisecond):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("halted operation2")
	}
}

func main() {
	//新建一个上下文
	ctx := context.Background()
	//在初始上下文的基础上创建一个有取消功能的上下文
	ctx, cancel := context.WithCancel(ctx)
	//在不同即goroutine中运行operation2
	go func() {
		operation2(ctx)
	}()

	err := operation1(ctx)
	if err != nil {
		cancel()
	}

	time.Sleep(1 * time.Second)
}
