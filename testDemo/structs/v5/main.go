package main

import "fmt"

type T struct {
	A string   //值类型
	B []string //引用类型
}

func main() {
	x := T{"剑鱼", []string{"上班"}}
	//拷贝的是指向对象的指针（浅拷贝）
	y := x
	y.A = "咸鱼"
	y.B[0] = "下班"
	fmt.Println(x)
	fmt.Println(y)

}
