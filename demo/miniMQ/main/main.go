package main

import (
	"fmt"
	"net"
	"taster/miniMQ/broker"
)

//消费者这里ack时重新创建了连接，
// 如果不创建连接的话，那服务端那里就需要一直从conn读取数据，直到结束。
// 思考一下，像RabbitMQ的ack就有自动和手工的ack，
// 如果是手工的ack，必然需要一个新的连接，因为不知道客户端什么时候发送ack，
// 自动的话，当然可以使用同一个连接，but这里就简单创建一条新连接吧

// 启动
// 先启动broker，再启动producer，然后启动comsumer能实现发送消息到队列

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Print("listen failed, err:", err)
		return
	}
	go broker.Save()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Print("accept failed, err:", err)
			continue
		}

		go broker.Process(conn)

	}
}
