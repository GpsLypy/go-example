package main

import "fmt"

//解锁了rune类型的第一个功能，即统计字符串长度。
func main() {
	//out:14 輸出底层占用字节长度，而不是字符串长度
	fmt.Println(len("Go語言編程"))
	//轉換成rune數組后統計字符串長度
	fmt.Println(len([]rune("Go語言編程")))

}
