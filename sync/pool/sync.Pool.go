package main

import (
	"fmt"
	"sync"
)

//sync.Pool是一个并发池，负责安全的保存一组对象，他有两个导出方法
//1、当我们必须重用共享的和长期存在的对象时(比如数据库连接) 、
//2、用于内存分配优化时
func main() {
	pool := &sync.Pool{}
	pool.Put(NewConnection(1))
	pool.Put(NewConnection(2))
	pool.Put(NewConnection(3))

	//Get方法会随机存取对象，不会以固定顺序
	connection := pool.Get().(*Connection)
	fmt.Printf("%d\n", connection.id)
	connection = pool.Get().(*Connection)
	fmt.Printf("%d\n", connection.id)
	connection = pool.Get().(*Connection)
	fmt.Printf("%d\n", connection.id)
}

/*
还可以为sync.Pool指定一个创建者方法
pool:=&sync.Pool{
	New :func() interface{}{
		return NewConnection()
	},
}

connection :=sync.Get().(*Connection)
*/

/*
让我们考虑一个写入缓冲区并将结果持久保存到文件中的函数示例。使用sync.Pool，我们可以通过在不同的函数调用之间重用同一对象来重用为缓冲区分配的空间。
第一步是检索先前分配的缓冲区（如果是第一个调用，则创建一个缓冲区，但这是抽象的）。然后，defer操作是将缓冲区放回sync.Pool中。

func writeFile(pool *sync.Pool, filename string) error {
    buf := pool.Get().(*bytes.Buffer)

  defer pool.Put(buf)

    // Reset 缓存区，不然会连接上次调用时保存在缓存区里的字符串foo
    // 编程foofoo 以此类推
    buf.Reset()

    buf.WriteString("foo")

    return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}
*/
