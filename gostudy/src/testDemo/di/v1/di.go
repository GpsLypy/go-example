package main

import (
	"fmt"
	"io"
	"net/http"
	_ "os"
)

//因为有了注入依赖，我们可以控制数据向哪儿写入

// func Greet(writer *bytes.Buffer, name string) {
//     fmt.Fprintf(writer, "Hello, %s", name)
// }

// Greet sends a personalised greeting to writer.
// func Greet(writer io.Writer, name string) {
// 	fmt.Fprintf(writer, "Hello, %s", name)
// }

// func main() {
// 	Greet(os.Stdout, "Elodie")
// }

//测试代码。如果你不能很轻松地测试函数，这通常是因为有依赖硬链接到了函数或全局状态。例如，如果某个服务层使用了全局的数据库连接池，
//这通常难以测试，并且运行速度会很慢。DI 提倡你注入一个数据库依赖（通过接口），然后就可以在测试中控制你的模拟数据了。

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	Greet(w, "world")
}

func main() {
	http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler))
}
