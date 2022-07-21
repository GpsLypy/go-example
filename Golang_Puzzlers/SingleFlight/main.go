//SingleFlight的作用是将并发请求合并成一个请求，以减少对下层服务的压力
//SingleFlight是Go开发组提供的一个扩展并发原语。它的作用是，在处理多个goroutine同时调用同一个函数的时候，
//只让一个goroutine去调用这个函数，等到这个goroutine返回结果的时候，再把结果返回给这几个同时调用的goroutine，这样可以减少并发调用的数量。

//sync.Once不是只在并发的时候保证只有一个goroutine执行函数f，而是会保证永远只执行一次，
//而SingleFlight是每次调用都重新执行，并且在多个请求同时调用的时候只有一个执行。
//它们两个面对的场景是不同的，sync.Once主要是用在单次初始化场景中，而SingleFlight主要用在合并并发请求的场景中，尤其是缓存场景

//如果你学会了SingleFlight，在面对秒杀等大并发请求的场景，而且这些请求都是读请求时，你就可以把这些请求合并为一个请求，
//这样，你就可以将后端服务的压力从n降到1。尤其是在面对后端是数据库这样的服务的时候，采用 SingleFlight可以极大地提高性能。
package main

import (
	"singleflight"

	"golang.org/x/sync/singleflight"
)

func main() {
	singleflight.Group{}
}
