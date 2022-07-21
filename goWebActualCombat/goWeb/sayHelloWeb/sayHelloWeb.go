package main

import (
	_ "fmt"
	"net/http"
)

//此函数实现了Net/http包中的提供的公共方法，即实现了type HandlerFunc func(responsWritee,*request)类型，
//该类型绑定了ServerHTTP方法，所以说间接实现了Handler接口（处理器）
func SayHello(w http.ResponseWriter, req *http.Request) {
	//将string 转换为字节切片
	w.Write([]byte("hello"))
}

func main() {
	http.HandleFunc("/", SayHello)
	//第二个参数为事件处理器Handler
	http.ListenAndServe(":8080", nil)
}
