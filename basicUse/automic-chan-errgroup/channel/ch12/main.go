package main

import "fmt"

// map-reduce
// map-reduce是一种处理数据的方式，最早是由Google公司研究提出的一种面向大规模数据处理的并行计算模型和方法，开源的版本是hadoop，前几年比较火。

// 不过，我要讲的并不是分布式的map-reduce，而是单机单进程的map-reduce方法。

// map-reduce分为两个步骤，第一步是映射（map），处理队列中的数据
//第二步是规约（reduce），把列表中的每一个元素按照一定的处理方式处理成结果，放入到结果队列中。

// 就像做汉堡一样，map就是单独处理每一种食材，reduce就是从每一份食材中取一部分，做成一个汉堡。

// 我们先来看下map函数的处理逻辑:

func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{}) //创建一个输出chan
	if nil == in {
		close(out)
		return out
	}
	//启动一个goroutine，实现map的主要逻辑
	go func() {
		defer close(out)
		//从输入chan读取数据，执行业务操作，也就是map操作
		for v := range in {
			out <- fn(v)
		}
	}()
	return out
}

//reduce函数的处理逻辑如下：

func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if nil == in {
		return nil
	}
	out := <-in //先读取第一个元素
	//实现reduce的主要逻辑
	for v := range in {
		out = fn(out, v)
	}
	return out
}

//我们可以写一个程序，这个程序使用map-reduce模式处理一组整数，map函数就是为每个整数乘以10，reduce函数就是把map处理的结果累加起来：

func asStream(done <-chan struct{}) <-chan interface{} {
	s := make(chan interface{})
	values := []int{1, 2, 3, 4, 5}
	go func() {
		defer close(s)
		for _, v := range values {
			select {
			case <-done:
				return
			case s <- v:
			}
		}
	}()
	return s
}

func main() {
	in := asStream(nil)
	//map操作，乘以10
	mapFn := func(v interface{}) interface{} {
		return v.(int) * 10
	}
	//reduce 操作，对map的结果进行累加
	reduceFn := func(r, v interface{}) interface{} {
		return r.(int) + v.(int)
	}

	sum := reduce(mapChan(in, mapFn), reduceFn) //返回累加结果
	fmt.Println(sum)
}
