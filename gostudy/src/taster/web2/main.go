package main

import (
	"fmt"
	"net/http"
	"time"
)

type Myhandler struct {
}

func (h *Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "~hello world!", r.URL.Path)
}

func main() {
	Myhandler := Myhandler{}

	server := http.Server{
		Addr:        ":8080",
		Handler:     &Myhandler,
		ReadTimeout: 2 * time.Second,
	}

	server.ListenAndServe()

}
