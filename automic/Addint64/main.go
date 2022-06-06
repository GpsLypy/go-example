package main

import (
	"fmt"
	"sync/atomic"
)

var value int32

//于 atomic.AddUint32() 和 atomic.AddUint64() 的第二个参数为 uint32 与 uint64，
//因此无法直接传递一个负的数值进行减法操作，Go语言提供了另一种方法来迂回实现：使用二进制补码的特性

//注意：unsafe.Pointer类型的值无法被加减。
func main() {
	var counter int64 = 23
	atomic.AddInt64(&counter, -3)
	fmt.Println(counter)

	AddValue(20)

	fmt.Println(value)

}

//比较并交换（Compare And Swap）
//简称CAS，在标准库代码包sync/atomic中以”Compare And Swap“为前缀的若干函数就是CAS操作函数，
//比如下面这个
//func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
//第一个参数的值是这个变量的指针，第二个参数是这个变量的旧值，第三个参数指的是这个变量的新值。

//运行过程：调用CompareAndSwapInt32 后，会先判断这个指针上的值是否跟旧值相等，
//若相等，就用新值覆盖掉这个值，若不相等，那么后面的操作就会被忽略掉。
//返回一个 swapped 布尔值，表示是否已经进行了值替换操作。

//与锁有不同之处：锁总是假设会有并发操作修改被操作的值，而CAS总是假设值没有被修改，
//因此CAS比起锁要更低的性能损耗，锁被称为悲观锁，而CAS被称为乐观锁。

//由示例可以看出，我们需要多次使用for循环来判断该值是否已被更改，
//为了保证CAS操作成功，仅在 CompareAndSwapInt32 返回为 true时才退出循环，这跟自旋锁的自旋行为相似。
func AddValue(delta int32) {
	for {
		v := value
		if atomic.CompareAndSwapInt32(&value, v, (v + delta)) {
			break
		}
	}
}

//载入与存储
//对一个值进行读或写时，并不代表这个值是最新的值，也有可能是在在读或写的过程中进行了并发的写操作导致原值改变。
//为了解决这问题，
//Go语言的标准库代码包sync/atomic提供了原子的读取（Load为前缀的函数）或写入（Store为前缀的函数）某个值
//将上面的示例改为原子读取
func AddValue2(delta int32) {
	for {
		v := atomic.LoadInt32(&value)
		if atomic.CompareAndSwapInt32(&value, v, (v + delta)) {
			break
		}
	}
}

// 原子写入总会成功，因为它不需要关心原值是什么，而CAS中必须关注旧值，因此原子写入并不能代替CAS，原子写入包含两个参数，以下面的StroeInt32为例：

// //第一个参数是被操作值的指针，第二个是被操作值的新值
// func StoreInt32(addr *int32, val int32)
// 交换
// 这类操作都以”Swap“开头的函数，称为”原子交换操作“，功能与之前说的CAS操作与原子写入操作有相似之处。

// 以 SwapInt32 为例，第一个参数是int32类型的指针，第二个是新值。原子交换操作不需要关心原值，而是直接设置新值，但是会返回被操作值的旧值。
// func SwapInt32(addr *int32, new int32) (old int32)
