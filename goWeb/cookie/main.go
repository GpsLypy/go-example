package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	CopeHandle("GET", "https://www.baidu.com", "")
}

//http请求处理
func CopeHandle(method, urlVal, data string) {
	fmt.Println(urlVal)
	client := &http.Client{}
	var req *http.Request

	if data == "" {
		urlArr := strings.Split(urlVal, "?")
		fmt.Println("urlArr", urlArr)
		if len(urlArr) == 2 {
			urlVal = urlArr[0] + "?" + getParseParam(urlArr[1])
			fmt.Println(urlArr[0])
			fmt.Println(urlVal)
		}
		req, _ = http.NewRequest(method, urlVal, nil)
	} else {
		req, _ = http.NewRequest(method, urlVal, strings.NewReader(data))
	}

	cookie := &http.Cookie{Name: "X-Xsrftoken", Value: "abccadf41ba5fasfasjijalkjaqezgbea3ga", HttpOnly: true}
	req.AddCookie(cookie)

	//添加header
	//这个是添加请求头的一个字段的值，用于模拟跨域请求的token[令牌]
	// token 放在 HTTP Header 里则可能导致产生一个预检请求
	req.Header.Add("X-Xsrftoken", "aaab6d695bbdcd111e8b681002324e63af81")
	fmt.Println(req)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
}

//将get请求的参数进行转义
func getParseParam(param string) string {
	return url.PathEscape(param)
}
