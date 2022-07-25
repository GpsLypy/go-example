传递不合时宜的上下文

在Go语言编写的并发程序中，上下文context出现非常频繁。在很多种场合，传递context是一种推荐的做法。但是，有时候传递context会导致很细微的错误，以至于子功能不能正确执行。

现在有这样一个例子，我们提供了一个执行处理某些任务并将结果返回的HTTP程序，但是在返回结果之前，还希望将结果发送到Kafka消息队列。不希望发送操作影响HTTP处理，所以我们开启一个goroutine处理发送操作。假设这个发布函数接收一个context.Context类型参数，以便发布消息的操作可以在上下文取消时终止，下面是一个示例程序。

func handler(w http.ResponseWriter, r *http.Request) {
        response, err := doSomeTask(r.Context(), r)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        go func() {
                err := publish(r.Context(), response)
                // Do something with err
        }()

        writeResponse(response)
}
上述程序中doSomeTask的处理结果会传递给publish函数和writeResponse函数，用于发布消息和返回格式化的响应。在发布函数中我们传入了来自http的上下文r.Context(),你能看出这段代码有什么问题吗？

需要知道的一点就是，附加到HTTP请求的上下文可以在不同的条件下取消：1. 当客户端的连接关闭时 2. 在HTTP/2请求的情况下，当请求被取消时 3.当响应被写回客户端时。

在前两种情况下，我们可能会正确地处理。例如，如果我们刚刚收到来自doSomeTask的响应，但客户端已经关闭了连接，那么在上下文已经取消的情况下调用发布函数publish是可以的，这个时候消息是不会发布的。但是最后一种情况就无法确定了，当响应被写入客户端时，与请求关联的上下文将被取消，这时面临了竞争条件：

如果写响应操作是在Kafka发布之后完成的，都会返回响应成功并成功发布消息，这种情况，写响应和发布是一致的，没有问题
如果在Kafka发布之前或期间写入响应，则消息可能不会发布，因为写入响应之后,上下文已经取消了，这时候执行发布处理的时候，不会发布消息
在上面后一种情况下，调用publish会返回错误，因为在发布处理之前，写操作已经完成了，很快的响应了HTTP响应。如何解决这个问题呢？一个思路是不传递父上下文，而是使用空上下文调用发布函数：

err := publish(context.Background(), response)
像上面这样，传递一个空的context，不管HTTP响应需要多长时间，都可以调用publish，因为它不受HTTP中的context影响。然而，如果上下文包含了一些有用的值呢？例如，如果上下文包含有用于分布式追踪的ID,我们可以关联HTTP请求和Kafka发布。一个好的做法是有一个新的上下文，与潜在的父上下文取消分离，但包含父上下文的键值信息。标准库没有提供这样问题的解决方法，一种可能的解决方法是我们自己实现一个上下文，与父上下文类似，只是没有父上下文的取消信号。

一个context.Context接口包含有下面4个方法。上下文的截止日期由Deadline方法和通过Done和Err方法的取消信号管理。当上下文截止日期已过或上下文被取消时，Done应该返回一个关闭的通道，而Err应该返回一个错误，返回键的值是通过Value获取的。

type Context interface {
        Deadline() (deadline time.Time, ok bool)
        Done() <-chan struct{}
        Err() error
        Value(key any) any
}
现在，我们创建一个自定义的上下文，它将取消信号与父上下文分离, 除了调用父上下文来获取键值信息外，其他方法都返回一个默认值，以便上下文永远不会被视为过期或取消。

type detach struct {
        ctx context.Context
}

func (d detach) Deadline() (time.Time, bool) {
        return time.Time{}, false
}

func (d detach) Done() <-chan struct{} {
        return nil
}

func (d detach) Err() error {
        return nil
}

func (d detach) Value(key any) any {
        return d.ctx.Value(key)
}
通过上面的自定义上下文，可以在调用发布中使用以此取消父上下文的取消信号，像下面这样。传递给发布的上下文是一个永远不会过期也不会被取消的上下文，但它会携带父上下文的值。

err := publish(detach{ctx: r.Context()}, response)
总之，在传递上下文时应该谨慎处理。本节通过一个基于HTTP请求关联的上下文来处理异步操作的示例来说明。由于一旦返回响应上下文就会被取消，异步操作也可能会意外停止。留意传递上下文带来的问题影响，需要记住的是，如果有必要，始终为特定操作创建自定义上下文。