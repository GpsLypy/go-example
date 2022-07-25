谨慎使用goroutine和循环变量

错误处理goroutine和循环变量可能是Go开发人员在编写并发应用程序时最常见的问题之一。下面通过一个具体的例子来说明，然后分析此类问题产生的原因以及如何防止它。

下面的程序中初始化一个切片s,然后循环遍历s，在循环中启动goroutine,通过闭包访问s中的元素. 代码如下。

s := []int{1, 2, 3}

for _, i := range s {
        go func() {
                fmt.Println(i)
        }()
}


也许你预期的结果是没有特定顺序的输出123，因为各个goroutine的执行先后顺序并不能保证，所以说是顺序不固定的。但是，实际的输出结果并不是确定性的包含1、2和3个数字，比如有时候打印233，有时候打印333. 这是为什么呢？

上面的程序中，新启动的goroutine引用了外部的变量i，这是函数闭包，其定义是函数内部引用了函数外部的变量。

====
有一点需要知道，当一个闭包goroutine被执行时，它不会立即处理闭包变量的值，所以上面的所有的goroutine都引用的是同一个变量i.当goroutine真正被执行时，它会在fmt.Println执行时打印i的值，这个时候i的值可能已经被修改。
====


上述程序输出结果为:

go run main.go                                                                                                                                                                                                                   
3
3
3
现在来看看打印233时可能的执行情况：


随着时间的流逝，i的值从1到2到3，在每次迭代中，都会启动一个新的goroutine，由于无法保证每个goroutine何时启动和完成，因此打印的结果也会有所不同。在上图中，当i等于2时，第一个goroutine打印i. 然后，当i的值已经为3时，其他goroutine打印，因此输出结果为233.

如果我们希望每个goroutine在创建时访问此时i的值，有什么解决方法？

一种解决方法是，如果想继续使用闭包，需要创建一个新变量，代码如下

for _, i := range s {
        val := i
        go func() {
                fmt.Println(val)
        }()
}

为什么上面这段代码工作是正常的。因为在每次迭代中，我们都会创建一个新的局部变量val, 此变量会在创建goroutine之前被赋值为i的当前值，当每个闭包goroutine在执行println语句时，会使用预期值执行，所以会输出123（当然123顺序可能是不固定的).

另一种解决方法是不使用闭包，而是通过函数参数传递的方式，代码如下：

for _, i := range s {
        go func(val int) {
                fmt.Print(val)
        }(i)
}

上述程序仍然采用匿名执行goroutine方式（而不是go f(i)), 但是没有采用闭包方式，内部打印的val不是函数外部的变量，而是函数的入参，通过这样处理，在每次迭代中保证val的值固定为当时i的值，使得程序如预期工作。

总结，在使用goroutine和循环变量时必须谨慎。如果一个goroutine访问的是函数外部的变量，这种闭包处理会引发问题。我们可以通过创建一个局部变量来修改它，或者不使用闭包操作，而是通过参数传递的方式。这两种方法都是有效的，不应墨守成规的只使用一种，你可能会发现使用闭包的方式处理更方便，采用函数参数传递的方式更容易理解。