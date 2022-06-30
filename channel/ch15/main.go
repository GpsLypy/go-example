package main

import (
	"fmt"
	_"time"
)

var tableChan chan int

func main(){
	tableChan=make(chan int,100)
	tableChan<-1
	for val :=range tableChan{
		fmt.Println(val)
	}
	
	//time.Sleep(10*time.Second)
}