package main

import (
	"fmt"
	"log"

	"github.com/colinmarc/hdfs"
)

func main() {
	client, err := hdfs.New("127.0.0.1:9870")
	//client, err := hdfs.New("")

	fmt.Println(111)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(222)
	file, _ := client.Open("/test.txt")

	buf := make([]byte, 59)
	file.ReadAt(buf, 48847)

	fmt.Println(string(buf))
}
