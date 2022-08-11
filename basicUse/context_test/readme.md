# 截止日期

截止日期

截止日期是指通过下面的方式确定的特定时间点：

```cpp
time.Duration:从现在开始持续的时间值，例如250毫秒
time.Time:一个具体的日期时间，例如 2023-02-07 00:00:00 UTC
```

截止日期（deadline）想表达的语义是如果到了该截止日期，则应该**停止正在进行的活动**。例如，一个I/O请求，或是一个等待从channel中接收消息的goroutine.

现在有这样一个应用程序，它每隔4秒从雷达接收一次飞行位置，一旦收到位置信息，会将位置信息共享给对飞机最新位置感兴趣的应用程序。这个应用程序有一个接口，接口中包含一个Publish方法，代码如下：

```cpp
type publisher interface{
    Publish(ctx context.Context,position flight.Position)
}
```

Publish函数有上下文context和位置position两个参数，假设具体实现将调用一个函数来向代理发布消息，例如使用Sarama发布到Kafka。这个函数是上下文感知的，也就是说一旦上下文被取消，它就会取消请求，完整代码见https://github.com/ThomasMing0915/100-go-mistakes-code/tree/main/60。

假设现在我们没有已有的context上下文对象，那怎么构造一个context传递给Publish方法呢？前面已经提到，所有应用程序只对最新位置感兴趣，因此，构造的上下文context应该能够表达，四秒后如果我们不能发布新位置，应该停止它。


type publishHandler struct{
    pub publisher
}

func (h publishHandler) publishPosition(position flight.Position) error{
    ctx,cancel:=context.WithTimeout(context.Backgroud(),4*time.Second)
    defer cancel()
    return h.pub.Publish(ctx,position)
}


上面的程序使用context.WithTimeout函数创建一个context对象，该函数接收一个超时时间和一个context对象。由于publishPosition没有收到现有的上下文context，即它的入参没有context，只有一个position变量。我们使用context.Background从一个空的上下文创建一个，同时，context.WithTimeout返回两个变量，创建的上下文和一个取消func()函数，调用取消函数后将取消上下文，创建的上下文ctx传递给h.pub.Publish之后，使得Publish方法最长在4秒内返回。

为什么通过defer调用cancel函数，context.WithTimeout内部创建了一个goroutine, 这个goroutine将存活4秒中或者被调用取消。因此通过defer调用cancel意味着当父函数退出时，上下文被取消，创建的goroutine将被销毁，这是一种将无效垃圾对象不留在内存中的保护措施。

