package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	//创建一个超时时间为100ms的上下文
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)

	//创建一个访问Google主页的请求
	req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)
	//将超时上下文关联到创建的请求上
	req = req.WithContext(ctx)

	//创建一个HTTP客户端并执行请求
	client := &http.Client{}
	res, err := client.Do(req)

	//如果请求失败了，记录到STDOUT
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	//请求成功后打印状态码
	fmt.Println("Response received,status code:", res.StatusCode)

}
