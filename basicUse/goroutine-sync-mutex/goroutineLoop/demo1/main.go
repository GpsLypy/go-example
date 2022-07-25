package main

import "context"

/*
首先让我们对一个goroutine泄露产生的影响有一个量的概念。
在内存占用方面，goroutine以最小的2KB大小开始分配，可以根据需要增长或缩小，64位系统的最大堆栈为1GB,32位系统的最大堆栈为250MB.
goroutine上还可以保存分配给堆的变量引用，保存资源，例如HTTP或DB连接，打开的文件，最终应该正常关闭的网络套接字，如果一个goroutine泄露，
这些资源可能也会被泄露。
下面来看一个不清楚什么该停止goroutine运行的例子。程序中，父goroutine调用一个返回通道的函数foo,然后创建一个新的goroutine将从该通道中接收消息。
*/

/*

ch := foo()
go func() {
        for v := range ch {
                // ...
        }
}()

创建的子goroutine将在ch被关闭时退出，但是，我们是否确切知道该通道何时关闭？可能不明显，
因为ch是由foo函数创建的，如果通道从未被关闭，那么就会导致泄露。因此，我们应该始终对goroutine的退出点保持谨慎，并确保最终能够退出不会泄露。

*/

//现在通过一个具体的例子进行分析说明。我们将设计一个需要监视外部配置的应用程序，例如使用数据库连接，下面是实例代码：

func main() {
	newWatcher()
	//run the application
}

type watcher struct { /*some resources */
}

func newWatcher() {
	w := watcher{}
	go w.watch()
}

/*
程序调用newWatcher，它创建了一个watcher结构对象，并启动一个goroutine来负责监视配置变动。
这段代码的问题点是当main goroutine退出时（可能是因为操作系统信号或者是有限的工作被处理完），
应用程序将停止。这会导致观察者创建的资源不会被优雅地关闭。那我们应该才能防止这种情况产生呢？
*/

//1、一种处理方法是向newWatcher传递一个上下文，该上下文将在main函数返回时被取消，代码如下。

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newWatcher(ctx)
	//run the application
}

func newWatcher() {
	w := watcher{}
	go w.watch()
}

/*
我们将创建的上下文传递给watch方法，当上下文被取消时，观察者应该关闭它的资源，
但是，我们能保证观察者有时间完成关闭资源操作吗？
我们不能保证，不过这是一个设计的问题。问题的原因是使用信号来传达一个goroutine必须停止，
在资源关闭之前，我们没有阻塞父goroutine，下面是一个改进的版本
*/

func main() {
	w := newWatcher()
	defer w.close()
}

func newWatcher() watcher {
	w := watcher{}
	go w.watch()
	return w
}

func (w watcher) close() {
	//close the resources
}

/*
watcher对象有一个close方法，现在不是通过向watcher方法发出信号来关闭它的资源，而是使用defer调用close方法来保证应用程序退出之前资源已经关闭。
*/
