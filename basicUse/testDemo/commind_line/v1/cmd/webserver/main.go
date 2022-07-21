// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"testDemo/commind_line/v1"
// )

// const dbFileName = "game.db.json"

// func main() {
// 	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

// 	if err != nil {
// 		log.Fatalf("problem opening %s %v", dbFileName, err)
// 	}

// 	//store := &FileSystemPlayerStore{db}
// 	store, _ := poker.NewFileSystemPlayerStore(db)
// 	server := poker.NewPlayerServer(store)

// 	if err := http.ListenAndServe(":5000", server); err != nil {
// 		log.Fatalf("could not listen on port 5000 %v", err)
// 	}
// }

package main

import (
	"log"
	"net/http"
	"testDemo/commind_line/v1"
)

const dbFileName = "game.db.json"

func main() {
	store, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
