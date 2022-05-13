package main

import (
	_ "io"
	"log"
	"net/http"
	_ "net/http/pprof"
	_ "os"
)

func main() {
	log.Println(http.ListenAndServe("localhost:8082", nil))
}

// func main() {
// 	myString := ""
// 	arguments := os.Args
// 	if len(arguments) == 1 {
// 		myString = "Please give me one argument!"
// 	} else {
// 		myString = arguments[1]
// 	}

// 	io.WriteString(os.Stdout, "This is Standard output\n")
// 	io.WriteString(os.Stderr, myString)
// 	io.WriteString(os.Stderr, "\n")
// }
