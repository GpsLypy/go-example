package main

import "reflect"

//扇入模式
//扇入借鉴了数字电路的概念，它定义了单个逻辑门能够接受的数字信号输入最大量的术语。一个逻辑门可以有多个输入，一个输出。
//在软件工程中，模块的扇入是指有多少个上级模块调用它。
//而对于我们这里的Channel扇入模式来说，就是指有多个源Channel输入、一个目的Channel输出的情况。扇入比就是源Channel数量比1。
//每个源Channel的元素都会发送给目标Channel，相当于目标Channel的receiver只需要监听目标Channel，
//就可以接收所有发送给源Channel的数据。

func fanInReflect(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		//构造SelectCase slice
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		//循环，从cases中选择一个可用的
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok { //此channel已经close
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}
	}()

	return out
}

func main() {

}
