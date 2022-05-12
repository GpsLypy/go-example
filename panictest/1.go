package main

// // import (
// // 	"fmt"
// // 	"time"
// // )

// // func main() {
// // 	fmt.Println("vim-go")
// // 	Go(func() {
// // 		fmt.Println("hello")
// // 		panic("一路向西")
// // 	})

// // 	time.Sleep(5 * time.Second)
// // }

// func Go(x func()) {
// 	// go func() {
// 	// 	defer func() {
// 	// 		if err := recover(); err != nil {
// 	// 			fmt.Println(err)
// 	// 		}
// 	// 	}()
// 	// 	x()
// 	// }()
// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()
// 	x()
// }
