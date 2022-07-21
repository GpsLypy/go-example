package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

func main() {
	times := time.Now().Unix()
	fmt.Println(times)
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%v", times))
	fmt.Printf("%x\n", h.Sum(nil))

	a := 5
	fmt.Println("····" + strconv.Itoa(a) + "主服务器-----")
}
