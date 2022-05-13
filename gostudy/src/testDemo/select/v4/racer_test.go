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
	got, err := Racer(slowURL, fastURL)

	if err != nil {
		t.Fatalf("did not expect an erro but got one %v", err)
	}
	if got != want {
		t.Errorf("got '%s',want '%s'", got, want)
	}

	t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
		serverA := makeDelayedServer(25 * time.Second)
		//serverB := makeDelayedServer(12 * time.Second)
		defer serverA.Close()
		//defer serverB.Close()

		_, err := ConfigurableRacer(serverA.URL, serverA.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})

}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
