package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

//利用channel实现协程池，缺点是造成携程的频繁开辟和和注销

//pool groutine

type Pool struct{
	queue chan int
	wg *sync.WaitGroup
}

func New(size int)* Pool{
	if size<=0{
		size=1
	}
	return &Pool{
		queue :make(chan int,size),
		wg: &sync.WaitGroup{},
	}
}

func (p *Pool) Add(delta int){
	for i:=0;i<delta;i++{
		p.queue<-1
	}
	for i:=0;i>delta;i--{
		<-p.queue
	}
	p.wg.Add(delta)
}

func (p *Pool) Done(){
	<-p.queue
	p.wg.Done()
}

func (p *Pool) Wait(){
	p.wg.Wait()
}


func main(){
	//这里限制100个并发
	pool :=New(100)
	for i:=0;i<100000;i++{
		pool.Add(1)
		go func(i int){
			resp,err:=http.Get("https://www.baidu.com")
			if err!=nil{
				fmt.Println(i,err)
			}else{
				defer resp.Body.Close()
				res,_:=ioutil.ReadAll(resp.Body)
				fmt.Println(i,string(res))
			}
			pool.Done()
		}(i)
	}
	pool.Wait()
}