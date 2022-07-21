package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func helloHandleFunc(w http.ResponseWriter, r *http.Request) {
	//1、解析模板
	t, err := template.ParseFiles("src/taster/template1/template_example.tmpl")
	if err != nil {
		fmt.Println("template parseFiles failed,err:", err)
		return
	}
	//2、渲染模板
	name := "我爱GO语言"
	t.Execute(w, name)
}
func main() {
	http.HandleFunc("/", helloHandleFunc)
	http.ListenAndServe(":8080", nil)
}
