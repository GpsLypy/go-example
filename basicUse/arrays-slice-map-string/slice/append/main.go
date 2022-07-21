package main

import "fmt"

// //包全局变量
// var a []int

// func fn(b []int) []int {
// 	//ab共享同一块底层数组
// 	//b结束生命后，a只引用一小段，造成内存泄漏
// 	a = b[:2]
// 	return a
// }

// func main() {
// 	//....
// }

var a []int
var c []int

func f(b []int) []int {
	a = b[:2]
	//c 实现了申请新的底层数组
	c = append(c, b[:2]...)
	fmt.Printf("a: %p\nc: %p\nb: %p\n", &a[0], &c[0], &b[0])
	return a
}

func main() {
	b := []int{1, 2, 3}
	f(b)

}
