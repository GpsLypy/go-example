package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

//使用计数器实现请求限流
//限流的要求是在指定的时间间隔内，server 最多只能服务指定数量的请求。
//实现的原理是我们启动一个计数器，每次服务请求会把计数器加一，同时到达指定的时间间隔后会把计数器清零
type RequestLimitService struct {
	Interval time.Duration //时间间隔
	MaxCount int           //最大连接数
	ReqCount int           //目前请求数
	Lock     sync.Mutex    //同步锁
}

func NewRequestLimitService(interval time.Duration, maxCnt int) *RequestLimitService {
	reqLimit := &RequestLimitService{
		Interval: interval,
		MaxCount: maxCnt,
	}

	go func() {
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			reqLimit.Lock.Lock()
			fmt.Println("Reset Count...")
			reqLimit.ReqCount = 0
			reqLimit.Lock.Unlock()
		}
	}()

	return reqLimit
}

func (reqLimit *RequestLimitService) Increase() {
	reqLimit.Lock.Lock()
	defer reqLimit.Lock.Unlock()

	reqLimit.ReqCount += 1
}

func (reqLimit *RequestLimitService) IsAvailable() bool {
	reqLimit.Lock.Lock()
	defer reqLimit.Lock.Unlock()
	return reqLimit.ReqCount < reqLimit.MaxCount
}

var RequestLimit = NewRequestLimitService(10*time.Second, 5)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if RequestLimit.IsAvailable() {
		RequestLimit.Increase()
		fmt.Println(RequestLimit.ReqCount)
		io.WriteString(w, "Hello world!\n")
	} else {
		fmt.Println("Reach request limiting!")
		io.WriteString(w, "Reach request limit!\n")
	}
}

func main() {
	fmt.Println("Server Started!")
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":8000", nil)
}

/*

使用golang来编写httpserver时，可以使用官方已经有实现好的包：

import(
  "fmt"
  "net"
  "golang.org/x/net/netutil"
)

func main() {
  l, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    fmt.Fatalf("Listen: %v", err)
  }
  defer l.Close()
  l = LimitListener(l, max)

  http.Serve(l, http.HandlerFunc())

  //bla bla bla.................
}
*/
