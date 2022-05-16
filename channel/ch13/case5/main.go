package main

//第4条规则是，如果Channel的容量是m（m>0），那么，第n个receive一定happens before 第 n+m 个 send的完成。

/*
前一条规则是针对unbuffered channel的，这里给出了更广泛的针对buffered channel的保证。
利用这个规则，我们可以实现信号量（Semaphore）的并发原语。
Channel的容量相当于可用的资源，发送一条数据相当于请求信号量，接收一条数据相当于释放信号。
关于信号量这个并发原语，我会在下一讲专门给你介绍一下，这里你只需要知道它可以控制多个资源的并发访问，就可以了。
*/
