package main

import (
	"fmt"
	"log"
	"os"
	"testDemo/commind_line/v1"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	// db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	// if err != nil {
	// 	log.Fatalf("problem opening %s %v", dbFileName, err)
	// }
	// store, err := poker.NewFileSystemPlayerStore(db)
	store, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}

	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()
}
