package main

import "fmt"

//not  ok yet
const (
	mutexLocked      = 1 << iota // 表示互斥锁的锁定状态
	mutexWoken                   // 表示从正常模式被从唤醒
	mutexStarving                // 当前的互斥锁进入饥饿状态
	mutexWaiterShift = iota      // 当前互斥锁上等待者的数量
)

func main() {
	fmt.Println(mutexLocked)      //1
	fmt.Println(mutexWoken)       //2
	fmt.Println(mutexStarving)    //4
	fmt.Println(mutexWaiterShift) //3
}
