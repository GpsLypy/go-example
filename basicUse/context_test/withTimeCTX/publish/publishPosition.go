package publish

import (
	"context"
	"fmt"
	"time"

	"github.com/GpyLypy/go-example/basicUse/context_test/withTimeCTX/flight"
)

type publisher2 interface {
	Publish(ctx context.Context, position flight.Position) (string, error)
}

type flyPublisher struct{}

func (f flyPublisher) Publish(ctx context.Context, position flight.Position) (string, error) {
	ch := make(chan string)
	go func() {
		// 模拟超时工作
		time.Sleep(time.Second * 3)
		select {
		case ch <- "result":
		default:
			return
		}
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case result := <-ch:
		return result, nil
	}
}

type publisher struct {
}

func (p publisher) Publish(ctx context.Context, position flight.Position) (string, error) {
	ch := make(chan string)
	go func() {
		// 模拟超时工作
		time.Sleep(time.Second * 5)
		select {
		case ch <- "result":
		default:
			return
		}
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case result := <-ch:
		return result, nil
	}
}

type PublishHandler struct {
	pub publisher
	//接口作为成员类型
	pub2 publisher2
}

func NewPublishHandler() PublishHandler {
	return PublishHandler{
		pub:  publisher{},
		pub2: flyPublisher{},
	}
}

func (h PublishHandler) PublishPosition(position flight.Position) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	fmt.Printf("[%v] begin start handle \n", time.Second*4)
	res, err := h.pub2.Publish(ctx, position)
	if err != nil {
		fmt.Printf("[%v] return error:%v \n", time.Now(), err)
		return
	}
	fmt.Printf("[%v] %s\n", time.Now(), res)
}
