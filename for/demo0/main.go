package main

import (
	"fmt"
	"sync"
)

//detail := make([]int,10)
var wg sync.WaitGroup
var lock sync.Mutex
// func append(detail *[]int){
// 	for i:=0;i<10;i++{
// 		*detail=append(*detail,i)
// 	}
// }



func main(){
	Detail :=[]int{1,2,3,4,5}
	//append(detail)
	for _,job :=range Detail{
		wg.Add(1)
		go func(jobId int){
			defer wg.Done()
			fmt.Println(jobId)
			lock.Lock()
			defer lock.Unlock()
		}(job)
	}
	wg.Wait()
	fmt.Println("all done")
}