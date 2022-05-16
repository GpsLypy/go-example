package main

//Channel是goroutine同步交流的主要方法。
//往一个Channel中发送一条数据，通常对应着另一个goroutine从这个Channel中接收一条数据。

//通用的Channel happens-before关系保证有4条规则

//第1条规则是，往Channel中的发送操作，happens before 从该Channel接收相应数据的动作完成之前，
//即第n个send一定happens before第n个receive的完成。

/*

在这个例子中，s的初始化（第5行）happens before 往ch中发送数据，
往ch发送数据 happens before从ch中读取出一条数据（第11行），
第12行打印s的值 happens after第11行，所以，打印的结果肯定是初始化后的s的值“hello world”。
*/
var ch = make(chan struct{}, 10)
var s string

func f() {
	s = "hello,world"
	ch <- struct{}{}
}

func main() {
	go f()
	<-ch
	print(s)
}
