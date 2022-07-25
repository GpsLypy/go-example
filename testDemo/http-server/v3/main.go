package main

import (
	"log"
	"net/http"
)

//之前我们探讨过 Handler 接口是为创建服务器而需要实现的。
//一般来说，我们通过创建 struct 来实现接口。然而，struct 的用途是用于存储数据，
//但是目前没有状态可存储，因此创建一个 struct 感觉不太对。
//HandlerFunc 可以让我们避免这样。
//HandlerFunc 类型是一个允许将普通函数用作 HTTP handler 的适配器。如果 f 是具有适当签名的函数，
//则 HandlerFunc(f) 是一个调用 f 的 Handler。
//所以我们用它来封装 PlayerServer 函数，使它现在符合 Handler。
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func (i *InMemoryPlayerStore) RecordWin(name string) {}

// func (p* PlayerServer) processWin(w http.ResponseWriter){
// 	p.store.
// }

func main() {
	//handler := http.HandlerFunc(PlayerServer)
	server := &PlayerServer{&InMemoryPlayerStore{}}
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
