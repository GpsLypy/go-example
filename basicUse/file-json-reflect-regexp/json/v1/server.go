package main

import (
	"fmt"
	"net/http"
)

type PlayerScore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

//为了让 PlayerServer 能够使用 PlayerStore，它需要一个引用。现在是改变架构的时候了，将 PlayerServer 改成一个 struct。
type PlayerServer struct {
	store PlayerScore
}

// func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		w.WriteHeader(http.StatusAccepted)
// 		return
// 	}
// 	player := r.URL.Path[len("/players/"):]
// 	score := p.store.GetPlayerScore(player)
// 	if score == 0 {
// 		w.WriteHeader(http.StatusNotFound)
// 	}

// 	fmt.Fprint(w, score)
// }

//把一个路由作为一个请求来处理并调用它挺奇怪的（并且效率低下）。
//我们想要的理想情况是有一个 NewPlayerServer 这样的函数，它可以取得依赖并进行一次创建路由的设置。每个请求都可以使用该路由的一个实例。
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	router.ServeHTTP(w, r)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
