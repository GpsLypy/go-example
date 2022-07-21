package main

// 扇出模式
// 有扇入模式，就有扇出模式，扇出模式是和扇入模式相反的。

// 扇出模式只有一个输入源Channel，有多个目标Channel，扇出比就是1比目标Channel数的值，
// 经常用在设计模式中的观察者模式中（观察者设计模式定义了对象间的一种一对多的组合关系。
// 这样一来，一个对象的状态发生变化时，所有依赖于它的对象都会得到通知并自动刷新）。
// 在观察者模式中，数据变动后，多个观察者都会收到这个变更信号。

//下面是一个扇出模式的实现。从源Channel取出一个数据后，依次发送给目标Channel。
//在发送给目标Channel的时候，可以同步发送，也可以异步发送：

func fanOut(ch <-chan interface{}, out []chan interface{}, async bool) {
	go func() {
		defer func() { //退出时关闭所有的输出chan
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()
		for v := range ch { //从输入chan 中读取数据
			v := v
			for i := 0; i < len(out); i++ {
				i := i
				if async { //异步
					go func() {
						out[i] <- v //放入到输出chan 中，异步方式
					}()
				} else {
					out[i] <- v //放入到输出chan中，同步方式
				}
			}
		}
	}()
}

func main() {

}
