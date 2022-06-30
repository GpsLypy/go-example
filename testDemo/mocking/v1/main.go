package main

import (
	"fmt"
	"io"
	"os"
)

//现在需要你写一个程序，从 3 开始依次向下，当到 0 时打印 「GO!」 并退出，要求每次打印从新的一行开始且打印间隔一秒的停顿

//所谓迭代是指：确保我们采取最小的步骤让软件可用。
//尽你所能拆分需求是一项很重要的技能

//下面是我们如何划分工作和迭代的方法：
//打印 3
//打印 3 到 Go!
//在每行中间等待一秒
func Countdown(out io.Writer) {
	fmt.Fprintf(out, "3")
}

func main() {
	Countdown(os.Stdout)
}
