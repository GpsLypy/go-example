package main

import (
	"fmt"
)

type T struct{}

func (t *T) Hello() string {
	if t == nil {
		fmt.Println("hello")
		return ""
	}
	return "hi"
}

func a(s string) {
	switch s {
	case "1":
		fmt.Println(1)
	case "2", "3":
		fmt.Println(2)
	}
}

func main() {
	//类型决定其调用，而不是值
	var t *T
	t.Hello()
	a("3")

}
