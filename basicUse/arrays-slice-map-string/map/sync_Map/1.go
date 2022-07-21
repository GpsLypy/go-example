package main

import (
	"fmt"
	"sync"
)

//什么时候使用sync.Map而不是在普通的map上使用sync.Mutex？
//1、当我们对map读多写少时
//2、当多个goroutine读取、写入、覆盖不相交的键时

func main() {
	m := &sync.Map{}
	m.Store(1, "one")
	m.Store(2, "two")

	val, ok := m.Load(1)
	if ok {
		fmt.Println(val.(string))
	}

	//原来存在返回ture，不存在返回false
	val, loaded := m.LoadOrStore(3, "three")
	fmt.Println(loaded) //false
	if !loaded {
		fmt.Printf("%s\n", val.(string))
	}

	m.Delete(3)

	//迭代所有元素,若函数返回了false,则停止迭代
	m.Range(func(key, value interface{}) bool {
		fmt.Printf("%d: %s\n", key.(int), value.(string))
		return true
	})
}
