package main

import (
	"fmt"
	"time"
)

// 使用chan也可以实现互斥锁。

// 在chan的内部实现中，就有一把互斥锁保护着它的所有字段。从外在表现上，chan的发送和接收之间也存在着happens-before的关系，保证元素放进去之后，receiver才能读取到（关于happends-before的关系，是指事件发生的先后顺序关系，我会在下一讲详细介绍，这里你只需要知道它是一种描述事件先后顺序的方法）。

// 要想使用chan实现互斥锁，至少有两种方式。一种方式是先初始化一个capacity等于1的Channel，然后再放入一个元素。这个元素就代表锁，谁取得了这个元素，就相当于获取了这把锁。另一种方式是，先初始化一个capacity等于1的Channel，它的“空槽”代表锁，谁能成功地把元素发送到这个Channel，谁就获取了这把锁。

// 这是使用Channel实现锁的两种不同实现方式，我重点介绍下第一种。理解了这种实现方式，第二种方式也就很容易掌握了，我就不多说了。

// 使用chan实现互斥锁
//创建带缓冲的通道。往通道发送数据时，只有当缓冲满时才会阻塞

type Mutex struct {
	ch chan struct{}
}

//使用锁需要初始化
func NewMutex() *Mutex {
	mu := &Mutex{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

//请求锁，直到获取到
func (m *Mutex) Lock() {
	<-m.ch
}

//解锁
func (m *Mutex) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

//尝试获取锁
func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

// 加入一个超时的设置
func (m *Mutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

// 锁是否已被持有
func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}

func main() {
	m := NewMutex()
	ok := m.TryLock()
	fmt.Printf("locked v %v\n", ok)
	ok = m.TryLock()
	fmt.Printf("locked %v\n", ok)
	// m.Unlock()

	// for i := 0; i < 5; i++ {
	// 	fmt.Println(i)
	// 	time.Sleep(1 * time.Second)
	// }
}

// 在这段代码中，还有一点需要我们注意下：利用select+chan的方式，
// 很容易实现TryLock、Timeout的功能。具体来说就是，在select语句中，我们可以使用default实现TryLock，使用一个Timer来实现Timeout的功能。
