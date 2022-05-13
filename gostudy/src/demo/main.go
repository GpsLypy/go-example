package main

func main() {
	server := &http.Server{
		// 请求监听地址
		Addr: ":8080",
		// 自定义的请求核心处理函数
		Handler: framework.NewCore(),
	}
	server.ListenAndServe()

}
