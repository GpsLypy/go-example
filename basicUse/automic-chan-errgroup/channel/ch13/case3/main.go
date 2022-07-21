package main

//第2条规则是，close一个Channel的调用，肯定happens before 从关闭的Channel中读取出一个零值。
var ch = make(chan struct{}, 10) // buffered或者unbuffered
var s string

func f() {
	s = "hello, world"
	//ch <- struct{}{}
	close(ch)
}

func main() {
	go f()
	<-ch
	print(s)
}
