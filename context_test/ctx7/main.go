package main

import (
	"context"
	"fmt"
)

//WithValue基于parent Context生成一个新的Context，保存了一个key-value键值对。它常常用来传递上下文。
//WithValue方法其实是创建了一个类型为valueCtx的Context，它的类型定义如下：
/*
type valueCtx struct {
    Context
    key, val interface{}
}
*/
func main() {
	ctx1 := context.TODO()
	ctx2 := context.WithValue(ctx1, "key1", "0001")
	ctx3 := context.WithValue(ctx2, "key2", "0002")
	ctx4 := context.WithValue(ctx3, "key3", "0003")
	ctx5 := context.WithValue(ctx4, "key4", "0004")

	fmt.Println(ctx5.Value("key2"))
}
