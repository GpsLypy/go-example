package main

/*

Facebook提供的这个ErrGroup，其实并不是对Go扩展库ErrGroup的扩展，而是对标准库WaitGroup的扩展。不过，因为它们的名字一样，处理的场景也类似，所以我把它也列在了这里。

标准库的WaitGroup只提供了Add、Done、Wait方法，而且Wait方法也没有返回子goroutine的error。而Facebook提供的ErrGroup提供的Wait方法可以返回error，而且可以包含多个error。子任务在调用Done之前，可以把自己的error信息设置给ErrGroup。接着，Wait在返回的时候，就会把这些error信息返回给调用者。

我们来看下Group的方法：

type Group
  func (g *Group) Add(delta int)
  func (g *Group) Done()
  func (g *Group) Error(e error)
  func (g *Group) Wait() error
关于Wait方法，我刚刚已经介绍了它和标准库WaitGroup的不同，我就不多说了。这里还有一个不同的方法，就是Error方法，
*/

import (
	"errors"
	"fmt"
	"time"
	//"github.com/facebookgo/errgroup"
)

func main() {
	var g errgroup.Group
	g.Add(3)

	// 启动第一个子任务,它执行成功
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("exec #1")
		g.Done()
	}()

	// 启动第二个子任务，它执行失败
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("exec #2")
		g.Error(errors.New("failed to exec #2"))
		g.Done()
	}()

	// 启动第三个子任务，它执行成功
	go func() {
		time.Sleep(15 * time.Second)
		fmt.Println("exec #3")
		g.Done()
	}()

	// 等待所有的goroutine完成，并检查error
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully exec all")
	} else {
		fmt.Println("failed:", err)
	}
}
