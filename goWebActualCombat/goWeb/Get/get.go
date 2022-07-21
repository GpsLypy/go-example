package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resq, err := http.Get("https://www.baidu.com")
	if err != nil {
		fmt.Print("err", err)
	}
	defer resq.Body.Close()
	closer := resq.Body
	bytes, err := ioutil.ReadAll(closer)
	fmt.Println(string(bytes))
}
