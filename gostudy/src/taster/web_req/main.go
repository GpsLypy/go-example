package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webapp/webDB/model"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "你发送的请求的地址是:", r.URL.Path)
	fmt.Fprintln(w, "你发送的请求的地址后的查询字符串是:", r.URL.RawQuery)
	fmt.Fprintln(w, "请求中的所有信息:", r.Header)
	fmt.Fprintln(w, "请求头中Accept-Encoding 的信息是:", r.Header.Get("Accept-Encoding"))
	fmt.Fprintln(w, "请求头中Accept-Encoding 的信息是:", r.Header["Accept-Encoding"])

	//获取请求体内容长度
	//len := r.ContentLength
	//创建一个body切片准备将内容读取到body中
	//body := make([]byte, len)

	//r.Body.Read(body)
	//fmt.Fprintln(w, "请求体中的内容：", string(body))

	///解析表单
	// r.ParseForm()
	// fmt.Fprintln(w, "请求参数有：", r.Form)
	// fmt.Fprintln(w, "POST请求的form表单中的参数有：", r.PostForm)

	//通过直接调用FormValue 和 PostForm Value方法直接获取请求参数的值
	fmt.Fprintln(w, "URL中的USER请求参数的值是：", r.FormValue("user"))
	fmt.Fprintln(w, "FORM表单中的USERname请求参数的值是：", r.PostFormValue("username"))

}

func testJsonRes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := model.User{
		ID:   1,
		Name: "tom",
	}
	//将结构体转为序列化的json格式
	json, _ := json.Marshal(user)
	w.Write(json)
}

func testRedir(w http.ResponseWriter, r *http.Request) {
	//设置响应头中的Location属性
	w.Header().Set("Location", "https://www.baidu.com")

	//设置响应状态码
	w.WriteHeader(302)
}

func main() {
	http.HandleFunc("/hello", handler)
	http.HandleFunc("/testjson", testJsonRes)
	http.HandleFunc("/testRedir", testRedir)
	http.ListenAndServe(":8080", nil)
}
