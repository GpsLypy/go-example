package broker

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

//Payload使用字节数组，是因为不管数据是什么，只当做字节数组来处理即可。
//Msg承载着生产者生产的消息，消费者消费的消息，ACK、和错误消息.
//前两者会有负载，而后两者负载和长度都为空。
type Msg struct {
	Id       int64
	TopicLen int64
	Topic    string
	// 1-consumer 2-producer 3-comsumer-ack 4-error
	MsgType int64  // 消息类型
	Len     int64  // 消息长度
	Payload []byte // 消息
}

//协议的编解码处理，就是对字节的处理，接下来有从字节转为Msg，和从Msg转为字节两个函数
func BytesToMsg(reader io.Reader) Msg {

	m := Msg{}
	var buf [128]byte
	n, err := reader.Read(buf[:])
	if err != nil {
		fmt.Println("read failed, err:", err)
	}
	fmt.Println("read bytes:", n)
	// id
	buff := bytes.NewBuffer(buf[0:8])
	binary.Read(buff, binary.LittleEndian, &m.Id)
	// topiclen
	buff = bytes.NewBuffer(buf[8:16])
	binary.Read(buff, binary.LittleEndian, &m.TopicLen)
	// topic
	msgLastIndex := 16 + m.TopicLen
	m.Topic = string(buf[16:msgLastIndex])
	// msgtype
	buff = bytes.NewBuffer(buf[msgLastIndex : msgLastIndex+8])
	binary.Read(buff, binary.LittleEndian, &m.MsgType)

	buff = bytes.NewBuffer(buf[msgLastIndex : msgLastIndex+16])
	binary.Read(buff, binary.LittleEndian, &m.Len)

	if m.Len <= 0 {
		return m
	}

	m.Payload = buf[msgLastIndex+16:]
	return m
}

func MsgToBytes(msg Msg) []byte {
	msg.TopicLen = int64(len([]byte(msg.Topic)))
	msg.Len = int64(len([]byte(msg.Payload)))

	var data []byte
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, msg.Id)
	data = append(data, buf.Bytes()...)

	buf = bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, msg.TopicLen)
	data = append(data, buf.Bytes()...)

	data = append(data, []byte(msg.Topic)...)

	buf = bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, msg.MsgType)
	data = append(data, buf.Bytes()...)

	buf = bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, msg.Len)
	data = append(data, buf.Bytes()...)
	data = append(data, []byte(msg.Payload)...)

	return data
}

// 队列
// 使用container/list，实现先入先出，生产者在队尾写，消费者在队头读取
//这里使用Queue结构体对List进行封装，其实是有必要的，List作为底层的数据结构，
//我们希望隐藏更多的底层操作，只给客户提供基本的操作。
type Queue struct {
	len  int
	data list.List
}

var lock sync.Mutex

//方法offer往队列里插入数据
func (queue *Queue) offer(msg Msg) {
	queue.data.PushBack(msg)
	queue.len = queue.data.Len()
}

//poll从队列头读取数据素
func (queue *Queue) poll() Msg {
	if queue.len == 0 {
		return Msg{}
	}
	msg := queue.data.Front()
	return msg.Value.(Msg)
}

//delete根据消息ID从队列删除数据
// delete操作是在消费者消费成功且发送ACK后，对消息从队列里移除的，因为消费者可以多个同时消费，
// 所以这里进入临界区时加锁（em，加锁是否就一定会影响性能呢？？）。
func (queue *Queue) delete(id int64) {
	lock.Lock()
	for msg := queue.data.Front(); msg != nil; msg = msg.Next() {
		if msg.Value.(Msg).Id == id {
			queue.data.Remove(msg)
			queue.len = queue.data.Len()
			break
		}
	}
	lock.Unlock()
}

//broker作为服务器角色，负责接收连接，接收和响应请求。
var topics = sync.Map{}

func handleErr(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			println(err.(string))
			conn.Write(MsgToBytes(Msg{MsgType: 4}))
		}
	}()
}

// MsgType等于1时，直接消费消息；
//MsgType等于2时是生产者生产消息，
// 如果队列为空，那么还需创建一个新的队列，放在对应的topic下；
// MsgType等于3时，代表消费者成功消费，可以删除
func Process(conn net.Conn) {
	handleErr(conn)
	reader := bufio.NewReader(conn)
	msg := BytesToMsg(reader)
	queue, ok := topics.Load(msg.Topic)
	var res Msg
	if msg.MsgType == 1 {
		// comsumer
		if queue == nil || queue.(*Queue).len == 0 {
			return
		}
		msg = queue.(*Queue).poll()
		msg.MsgType = 1
		res = msg
	} else if msg.MsgType == 2 {
		// producer
		if !ok {
			queue = &Queue{}
			queue.(*Queue).data.Init()
			topics.Store(msg.Topic, queue)
		}
		queue.(*Queue).offer(msg)
		res = Msg{Id: msg.Id, MsgType: 2}
	} else if msg.MsgType == 3 {
		// consumer ack
		if queue == nil {
			return
		}
		queue.(*Queue).delete(msg.Id)

	}
	conn.Write(MsgToBytes(res))

}

// 删除消息
// 我们说消息不丢失，这里实现不完全，我就实现了持久化（持久化也没全部实现）。
//思路就是该topic对应的队列里的消息，按协议格式进行序列化，当broker启动时，从文件恢复。
// 持久化需要考虑的是增量还是全量，需要保存多久，这些都会影响实现的难度和性能（想想Kafka和Redis的持久化）
//这里表示简单实现就好：定时器定时保存

func Save() {
	ticker := time.NewTicker(60)
	for {
		select {
		case <-ticker.C:
			topics.Range(func(key, value interface{}) bool {
				if value == nil {
					return false
				}
				file, _ := os.Open(key.(string))
				if file == nil {
					file, _ = os.Create(key.(string))
				}
				for msg := value.(*Queue).data.Front(); msg != nil; msg = msg.Next() {
					file.Write(MsgToBytes(msg.Value.(Msg)))
				}
				res := file.Close()
				fmt.Println(res)
				return false
			})
		default:
			time.Sleep(1)
		}
	}
}
