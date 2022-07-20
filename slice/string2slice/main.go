package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

//字符串转成切片，会产生拷贝。有没有什么办法在字符串转为切片的时候不用发生拷贝呢？

func main() {
	a := "aaa"
	//(unsafe.Pointer(&a))可以得到a的地址
	//(*reflect.StringHeader)(unsafe.Pointer(&a)) 可以把字符串转成底层结构的形式
	ssh := *(*reflect.StringHeader)(unsafe.Pointer(&a))
	//把底层结构题转化为byte切片的指针
	b := *(*[]byte)(unsafe.Pointer(&ssh))
	//转为指针指向的实际内容
	fmt.Printf("%v", b)
}
