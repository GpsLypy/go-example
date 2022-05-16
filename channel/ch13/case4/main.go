package main

import "fmt"

//第3条规则是，对于unbuffered的Channel，
//也就是容量是0的Channel，从此Channel中读取数据的调用一定happens before 往此Channel发送数据的调用完成。

var ch = make(chan struct{})

var s string

func f() {
	s = "hello,world"
	<-ch
	fmt.Println(111)
}

func main() {
	go f()
	ch <- struct{}{}
	print(s)
}
