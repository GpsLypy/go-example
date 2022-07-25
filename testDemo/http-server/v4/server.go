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

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
