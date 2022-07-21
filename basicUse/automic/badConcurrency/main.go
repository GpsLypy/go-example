package main

import (
	"time"
)

func badConcurrency() {
	batchSize := 500
	for {
		data, _ := queryDataWithSizeN(batchSize)
		if len(data) == 0 {
			break
		}
		for _, item := range data {
			go func(i int) {
				doSomething(i)
			}(item)
		}
		time.Sleep(1 * time.Second)
	}
}
