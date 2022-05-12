package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	once := &sync.Once{}
	for i := 0; i < 4; i++ {
		i := i
		go func() {
			once.Do(func() {
				fmt.Printf("first %d\n", i)
			})
		}()
	}
	time.Sleep(1 * time.Second)
}
