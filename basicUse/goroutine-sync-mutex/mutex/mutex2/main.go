package main

//此种方式对外不会暴露锁逻辑
import (
	"fmt"
	"sync"
)

//muetx正常姿势应该嵌入到struct中使用
//线程安全的计数器类型
type Counter struct {
	CounterType int
	Name        string
	mu          sync.Mutex
	count       uint64
}

//加一的方法，内部使用互斥锁保护
func (c *Counter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

//读取计数器值的操作也需要互斥锁保护
func (c *Counter) Count() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
func main() {

	//使用WaitGroup等待10个goroutine完成
	//var wg sync.WaitGroup
	wg := &sync.WaitGroup{}
	var counter Counter
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			//对变量count执行10次加1
			for j := 0; j < 100000; j++ {
				counter.Incr()
			}
		}()
	}
	wg.Wait()
	fmt.Println(counter.Count())
}
