package main

import "fmt"

func Positive(n int) bool {
	return n > -1
}

func Check(n int) {
	if Positive(n) {
		fmt.Println(n, "is positive")
	} else {
		fmt.Println(n, "is negative")
	}
}

func main() {
	Check(1)
	Check(0)
	Check(-1)
}


sentinel error 预定义的，比如错误码