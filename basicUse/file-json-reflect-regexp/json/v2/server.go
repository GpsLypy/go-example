package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PlayerScore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

//为了让 PlayerServer 能够使用 PlayerStore，它需要一个引用。现在是改变架构的时候了，将 PlayerServer 改成一个 struct。
//我们更改了 PlayerServer 的第二个属性，删除了命名属性 router http.ServeMux，并用 http.Handler 替换了它；这被称为 嵌入。
//这意味着我们的 PlayerServer 现在已经有了 http.Handler 所有的方法，也就是 ServeHTTP。
type PlayerServer struct {
	store PlayerScore
	//router *http.ServeMux
	http.Handler
}

func NewPlayerServer(store PlayerScore) *PlayerServer {
	// p := &PlayerServer{
	// 	store,
	// 	http.NewServeMux(),
	// }
	p := new(PlayerServer)
	p.store = store
	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	p.Handler = router
	return p
}

// func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	p.router.ServeHTTP(w, r)
// }

//把一个路由作为一个请求来处理并调用它挺奇怪的（并且效率低下）。
//我们想要的理想情况是有一个 NewPlayerServer 这样的函数，它可以取得依赖并进行一次创建路由的设置。每个请求都可以使用该路由的一个实例。
// func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

// 	router := http.NewServeMux()
// 	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
// 	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

// 	router.ServeHTTP(w, r)
// }

type Player struct {
	Name string
	Wins int
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	//json.NewEncoder(w).Encode(p.getLeagueTable())
	json.NewEncoder(w).Encode(p.store.GetLeague())
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) getLeagueTable() []Player {
	return []Player{
		{"Chris", 20},
	}
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
