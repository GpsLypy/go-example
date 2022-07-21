package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

//我们经常会碰到需要将一个通用的父任务拆成几个小任务并发执行的场景，其实，将一个大的任务拆成几个小任务并发执行，可以有效地提高程序的并发度
//ErrGroup就是用来应对这种场景的。它和WaitGroup有些类似，但是它提供功能更加丰富：

// 和Context集成；
// error向上传播，可以把子任务的错误传递给Wait的调用者。

//1.WithContext

// 在创建一个Group对象时，需要使用WithContext方法：

// func WithContext(ctx context.Context) (*Group, context.Context)
// 这个方法返回一个Group实例，同时还会返回一个使用context.WithCancel(ctx)生成的新Context。一旦有一个子任务返回错误，或者是Wait调用返回，这个新Context就会被cancel。

// Group的零值也是合法的，只不过，你就没有一个可以监控是否cancel的Context了。

// 注意，如果传递给WithContext的ctx参数，是一个可以cancel的Context的话，那么，它被cancel的时候，并不会终止正在执行的子任务。

//2.Go

// 我们再来学习下执行子任务的Go方法：

// func (g *Group) Go(f func() error)
// 传入的子任务函数f是类型为func() error的函数，如果任务执行成功，就返回nil，否则就返回error，并且会cancel 那个新的Context。

// 一个任务可以分成好多个子任务，而且，可能有多个子任务执行失败返回error，不过，Wait方法只会返回第一个错误，所以，如果想返回所有的错误，需要特别的处理，我先留个小悬念，一会儿再讲。

//3.Wait

// 类似WaitGroup，Group也有Wait方法，等所有的子任务都完成后，它才会返回，否则只会阻塞等待。如果有多个子任务返回错误，它只会返回第一个出现的错误，如果所有的子任务都执行成功，就返回nil：

// func (g *Group) Wait() error

func main() {
	var g errgroup.Group

	//启动第一个子任务，它执行成功
	g.Go(func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("exec #1")
		return nil
	})

	//启动第二个子任务，它执行失败
	g.Go(func() error {
		time.Sleep(10 * time.Second)
		fmt.Println("exec #2")
		return errors.New("failed to exec #2")
	})

	// 启动第三个子任务，它执行成功
	g.Go(func() error {
		time.Sleep(15 * time.Second)
		fmt.Println("exec #3")
		return nil
	})

	// 等待三个任务都完成
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully exec all")
	} else {
		fmt.Println("failed:", err)
	}
}
