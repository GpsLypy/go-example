package main

import "sync"

//Once
//它提供的保证是：对于once.Do(f)调用，f函数的那个单次调用一定happens before 任何once.Do(f)调用的返回。换句话说，就是函数f一定会在Do方法返回之前执行。

var s string

var once sync.Once

func foo() {
	s = "hello"
}

func main() {
	once.Do(foo)
	print(s)
}
