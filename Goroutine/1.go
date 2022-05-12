//Never start a goroutine without knowning when it will stop
//when terminate
//what terminating

package main

import (
	"context"
	"fmt"
	_ "log"
	"net/http"
)

//版本一
// func main() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
// 		fmt.Fprintln(resp, "hello,QCon!")
// 	})
// 	go http.ListenAndServe("127.0.0.1:8001", http.DefaultServeMux) //Debug
// 	http.ListenAndServe("0.0.0.0:8000", mux)
// }

//版本二
// func serveApp() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
// 		fmt.Fprintln(resp, "hello,QCon")
// 	})
// 	http.ListenAndServe("0.0.0.0:8080", mux)
// }

// func serveDebug() {
// 	http.ListenAndServe("127.0.0.1:8001", http.DefaultServeMux)
// }

// func main() {
// 	go serveDebug()
// 	serveApp()
// }

//版本三
//log.Fatal 调用了os.Exit,会无条件终止程序，defers不会被调用到
//only use log.Fatal from main.main or init functions
// func serveApp() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
// 		fmt.Fprintln(resp, "hello,QCon")
// 	})
// 	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func serveDebug() {
// 	if err := http.ListenAndServe("127.0.0.1:8001", http.DefaultServeMux); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func main() {
// 	go serveDebug()
// 	go serveApp()
// 	select {}
// }

//版本四
func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-stop //wait for stop signal
		s.Shutdown(context.Background())
	}()
	return s.ListenAndServe()
}

func serveApp(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "hello,QCon")
	})
	return serve("0.0.0.0:8080", mux, stop)
}

func serveDebug(stop <-chan struct{}) error {
	return serve("127.0.0.1:8001", http.DefaultServeMux, stop)
}

func main() {
	done := make(chan error, 2)
	//用作信号
	stop := make(chan struct{})
	go func() {
		done <- serveDebug(stop)
	}()
	go func() {
		done <- serveApp(stop)
	}()

	var stopped bool
	for i := 0; i < cap(done); i++ {
		fmt.Println("ohuo")
		if err := <-done; err != nil {
			fmt.Println("error:%v", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}

	fmt.Println("lalla")

}

//func ListDirection(root string,walkFn walkfunc)
