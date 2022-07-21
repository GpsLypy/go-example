//通过反射的方式执行select语句，在处理很多的case clause，尤其是不定长的case clause的时候，非常有用

package main

import (
	"fmt"
	"reflect"
)

// 通过reflect.Select函数，你可以将一组运行时的case clause传入，当作参数执行。
// Go的select是伪随机的，它可以在执行的case中随机选择一个case，并把选择的这个case的索引（chosen）返回，
// 如果没有可用的case返回，会返回一个bool类型的返回值，这个返回值用来表示是否有case成功被选择。
// 如果是recv case，还会返回接收的元素。Select的方法签名如下：

//func Select(cases []SelectCase) (chosen int, recv Value, recvOK bool)

//演示动态处理两个chan的场景

// 首先，createCases函数分别为每个chan生成了recv case和send case，并返回一个reflect.SelectCase数组。

// 然后，通过一个循环10次的for循环执行reflect.Select，这个方法会从cases中选择一个case执行。
// 第一次肯定是send case，因为此时chan还没有元素，recv还不可用。
// 等chan中有了数据以后，recv case就可以被选择了。这样，你就可以处理不定数量的chan了。

func main() {
	var ch1 = make(chan int, 10)
	var ch2 = make(chan int, 10)

	//创建SelectCase
	var cases = createCases(ch1, ch2)

	//执行10次select
	for i := 0; i < 10; i++ {
		chosen, recv, ok := reflect.Select(cases)
		if recv.IsValid() {
			fmt.Println("recv:", cases[chosen].Dir, recv, ok)
		} else {
			fmt.Println("send:", cases[chosen].Dir, ok)
		}
	}

}

func createCases(chs ...chan int) []reflect.SelectCase {
	var cases []reflect.SelectCase

	//创建recv case
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	//创建send case

	for i, ch := range chs {
		v := reflect.ValueOf(i)
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: v,
		})
	}

	return cases
}
