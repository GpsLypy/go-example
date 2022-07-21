package main

import (
	"fmt"
	"net/http"
)

//创建处理器

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world!", r.URL.Path)
}
func main() {
	http.HandleFunc("/http", handler)

	//1创建路由
	http.ListenAndServe(":8080", nil)

}
