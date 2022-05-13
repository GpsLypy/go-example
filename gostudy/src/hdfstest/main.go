package main

import (
	"fmt"
	"log"

	"github.com/colinmarc/hdfs"
)

func main() {
	//client, err := hdfs.New("192.168.65.130:9870")
	client, err := hdfs.New("diy2.bigdata.ly:50070")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Remove("/test/1.txt")
	fmt.Println(err)
}
