package main

import (
	"fmt"
)

var curId int64

func main() {
	curId = 1
	for i := 0; i < 3; i++ {
		if curId >= 1 {
			break
		}
		fmt.Println("1\n")
	}

	fmt.Println("2\n")
}
