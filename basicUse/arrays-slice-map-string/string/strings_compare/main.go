package main

import (
	"fmt"
	"strings"
)

//golang字符串操作
func main() {
	s := "Hello world hello world"
	str := "Hello"

	//比较字符串，区分大小写，比”==”速度快。相等为0，a<b 为-1.a>b为1。
	ret := strings.Compare(s, str)
	fmt.Println(ret) //  1
}
