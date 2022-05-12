package main

import (
	"fmt"
	"sync"
)

//muetx正常姿势应该嵌入到struct中使用
type Counter struct {
	mu    sync.Mutex
	Count uint64
}

func main() {
	//var count = 0
	//使用WaitGroup等待10个goroutine完成
	var wg sync.WaitGroup
	//var mu sync.Mutex
	var counter Counter
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			//对变量count执行10次加1
			for j := 0; j < 100000; j++ {
				counter.mu.Lock()
				counter.Count++
				counter.mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println(counter.Count)
}
