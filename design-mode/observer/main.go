package main

import (
	"fmt"
	"time"
)

//观察者模式也叫发布-订阅模式，待发布的状态变更信息会被包装到一个对象里，这个对象被称为事件
//事件发布给订阅者的过程，其实就是遍历一下已经注册的事件订阅者，逐个去调用订阅者实现的接口，比如叫handleEvent之类的

//Subject 接口 它相当于是发布者的定义
type Subject interface {
	Subscribe(observer Observer)
	Notify(msg string)
}

type Observer interface {
	Update(msg string)
}

//Subject 实现
type SubjectImpl struct {
	observers []Observer
}

//发布者会调用这个接口把想要订阅到自己主题下的订阅者加入到自己专门存放的队列里（添加观察者-订阅者）
func (sub *SubjectImpl) Subscribe(observer Observer) {
	sub.observers = append(sub.observers, observer)
}

//发布通知
func (sub *SubjectImpl) Notify(msg string) {
	for _, o := range sub.observers {
		o.Update(msg)
	}
}

type Observer1 struct{}

//Updata 实现观察者接口

func (Observer1) Update(msg string) {
	fmt.Printf("Observer1:%s\n", msg)
}

type Observer2 struct{}

func (Observer2) Update(msg string) {
	fmt.Printf("Observer2:%s\n", msg)
}

//在实际应用中，一般会定一个事件总线,EventBus 或Event_Dispatcher来管理事件和订阅者间的关系，以及分发事件
//ex1:
//1、异步不阻塞
//2、支持任意参数值

func main() {
	// s := &SubjectImpl{}
	// o1 := Observer1{}
	// o2 := Observer2{}
	// s.Subscribe(o1)
	// s.Subscribe(o2)
	// s.Notify("hello")

	bus := NewAsyncEventBus()
	bus.Subscribe("topic:1", sub1)
	bus.Subscribe("topic:1", sub2)
	bus.Publish("topic:1", "test1", "test2")
	bus.Publish("topic:1", "testA", "testB")
	time.Sleep(1 * time.Second)
}
