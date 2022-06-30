package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//channel导致死锁之一：unbuffered channel

//1、只要生产者，没有消费者，或者反过来，略

//2、生产者消费者出现在同一个goroutine中
func f3() {
	ch := make(chan int)
	ch <- 1 //由于消费者还没执行到，这里会一直阻塞住
	<-ch
}

//对于带缓冲buffer,则是buffer channel 已经满了，然后出现上述情况

//goroutine泄漏的核心是：生产者/消费者所在的goroutine已经退出，而其对应得消费者/生产者所在的goroutine会阻塞住，直到进程退出

//2、1 生产者阻塞导致内存泄漏

//我们一般会用channel来做一些超时控制
func leak1() {
	ch := make(chan int)
	//g1
	go func() {
		time.Sleep(2 * time.Second) //模拟IO操作
		ch <- 100
	}()
	//g2，阻塞住，直到超时或者返回
	select {
	case <-time.After(500 * time.Microsecond):
		fmt.Println("timeout exit...")
	case result := <-ch:
		fmt.Printf("result:%d\n", result)
	}
}

//消费者阻塞导致内存泄漏
//这种情况下，只需要增加close(ch) 的操作即可，for range 操作在收到close的信号后悔退出
func leak2() {
	ch := make(chan int)
	//消费者g1
	go func() {
		for result := range ch {
			fmt.Printf("result:%d\n", result)
		}
	}()
	//生产者g2
	ch <- 1
	ch <- 2
	time.Sleep(time.Second)
	fmt.Println("main goroutine g2 done ...")
}

// func main() {
// 	leak1()
// }

//预防goroutine泄露的核心是：创建goroutine时就要考虑清楚他什么时候被回收
//1、当goroutine退出时，需要考虑它使用的channel有没有可能阻塞对应的
//2、尽量使用buffered channel  能减少阻塞发生，即使疏忽了一些极端情况，也能减低goroutine泄漏的概率

//panic

// channel 导致的 panic 一般是以下几个原因：
//4.1向已经close掉的channel继续发送数据
//4.2 重复close

//如何优雅的关闭channel?
//1、除非必须关闭chan 否则不要关闭 ，当一个chan 没有sender和receiver ，即不在使用时，GC会在一段时间后标记
//清理掉这个chan 。也就是说将close作为一种通知机制，尤其是生产者和消费者之间是1：M的关系时
//通过close告诉下游，我收工了，你们别读啦
//2、channel关闭的原则;
//不要再消费者端关闭chan
//有多个并发写的生产者时也别关
//一写一读：生产者关闭即可
//一写多读：生产者关闭即可
//多写一读：多个生产者之间需要引入一个协调channel来处理信号***
//多写多读：引入一个中间层以及使用try_send 的套路来处理非阻塞的写入

//列如：
func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)
	const Max = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)
	stopCh := make(chan struct{})
	//stopCh 是额外引入的一个信号channel ,
	//他的生产者是下面的toStop channel
	//消费者是上面dataCh的生产者和消费者
	toStop := make(chan string, 1)
	//toStop是拿来关闭stopCh用的，由dataChan 的生产者和消费者写入
	//由下面这个匿名中介函数moderator消费
	//要注意，这个一定要是buffered channel（否则没办法用try_send来处理了）
	var stoppedBy string

	//moderator
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	//senders
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(Max)
				//准备结束时走的分支
				if value == 0 {
					//try-send操作
					//如果toStop满了，(不知道被同伙谁填充了)，就会走defalut分支，也不会阻塞
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}
				//try-receive 操作，尽快退出
				//如果没有这一步，下面的select操作可能会造成panic
				//发送数据时先检查退出标志
				select {
				case <-stopCh:
					return
				default:
				}

				//如果尝试从stopCh取数据的同时，也尝试向dataCh 写数据，则会命中select
				//的伪随机逻辑，可能会写入数据
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	//receivers
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()
			for {
				select {
				case <-stopCh:
					return
				default:
				}

				//尝试读取数据
				select {

				case <-stopCh:
					return
				case value := <-dataCh:
					if value == Max-1 {
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}
					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}
	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
