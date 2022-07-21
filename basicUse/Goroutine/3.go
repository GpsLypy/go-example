package main
import(
	"fmt"
	"time"
)
func gohead() {
//	go func() {
	 fmt.Println(1111)
	 panic("煎鱼下班了")
//	}()
}
   
   func main() {
	go func() {
	 defer func() {
	  if r := recover(); r != nil {
	   fmt.Println(r)
	   fmt.Println(11)
	  }
	 }()
   
	 gohead()
	}()
   
	time.Sleep(5*time.Second)
   }