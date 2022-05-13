package main

import (
	"fmt"
)

func main() {
	var arr [5]int = [...]int{10, 20, 30, 40, 50}
	slice := arr[1:4]
	//使用for--range 方式遍历切片
	for i := range slice {
		//fmt.Printf("i=%v v=%v \n", i, v)
		fmt.Printf("i=%v\n", i)
	}
}
