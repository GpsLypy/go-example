package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//httptest.NewServer 接受一个我们传入的 匿名函数 http.HandlerFunc。
// http.HandlerFunc 是一个看起来类似这样的类型：type HandlerFunc func(ResponseWriter, *Request)。
// 这些只是说它是一个需要接受一个 ResponseWriter 和 Request 参数的函数，这对于 HTTP 服务器来说并不奇怪。
func TestRacer(t *testing.T) {
	slowSever := makeDelayedServer(20 * time.Millisecond)
	fastServer := makeDelayedServer(0 * time.Millisecond)
	defer slowSever.Close()
	defer fastServer.Close()
	slowURL := slowSever.URL
	fastURL := fastServer.URL
	want := fastURL
	got := Racer(slowURL, fastURL)
	if got != want {
		t.Errorf("got '%s',want '%s'", got, want)
	}

}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
