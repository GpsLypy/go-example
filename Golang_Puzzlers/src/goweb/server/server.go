package main

import (
	"flag"
	"net/http"
)

//路径处理函数
func Hello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	host := flag.String("host", "127.0.0.1", "listen host")
	port := flag.String("port", "80", "listen port")
	//用来注册路径处理函数，会根据给定的路径不同，调用不同的函数
	http.HandleFunc("/hello", Hello)

	//监听IP和端口，本机的话仅书写冒号加端口
	err := http.ListenAndServe(*host+":"+*port, nil)
	if err != nil {
		panic(err)
	}
}
