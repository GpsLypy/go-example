package main

import "fmt"

func main() {
	//带缓冲的channel 缺点：容量不好控制
	var ch = make(chan int, 10)
	for i := 0; i < 10; i++ {
		select {
		case ch <- i:
			//len返回chan中缓存的还未被取走的元素数量
			fmt.Printf("len1=%d\n", len(ch))
		case v := <-ch:
			fmt.Printf("len2=%d\n", len(ch))
			fmt.Println(v)
		}
	}
}
