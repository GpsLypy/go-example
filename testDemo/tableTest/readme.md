使用t.Parallel运行的测试项将与其他并行测试项并行执行。在执行的时候，Go语言测试工具(go test)会一个接一个地运行完所有的顺序测试项，然后，一旦顺序测试完成，将执行并行测试。例如，下面的测试代码包含3个测试项，其中只有两个被标记为并行运行。

func TestA(t *testing.T) {
        t.Parallel()
        // ...
}

func TestB(t *testing.T) {
        t.Parallel()
        // ...
}

func TestC(t *testing.T) {
        // ...
}
运行上述测试文件(example1_test.go)得到的输出日志信息如下：


 go test -v example1_test.go          


=== RUN   TestA
=== PAUSE TestA
=== RUN   TestB
=== PAUSE TestB
=== RUN   TestC
--- PASS: TestC (0.00s)
=== CONT  TestA
--- PASS: TestA (0.00s)
=== CONT  TestB
--- PASS: TestB (0.00s)
PASS
ok      command-line-arguments  0.007s

通过上面输出日志可以看到，go test 会按文件中编写测试项的先后顺序运行，但是第一个真正被执行的是TestC。TestA和TestB在按顺序运行到的时候，它们被暂停实际并没有执行。等待TestC执行完成之后，TestA和TestB都被恢复并且并行执行。

默认情况下，可以同时并行运行的最大测试数等于GOMAXPROCS值。例如，如果我们想要进行序列化测试或者要在长时间运行的测试中执行大量I/O操作，可以增大并行运行的最大测试数。具体我们可以使用-parallel参数设置并行度值，像下面这样将最大并行测试数设置为16.
go test -parallel 16 .

打乱顺序(Shuffle)

从Go1.17版本开始，可以随机化测试和性能测试的执行顺序。为什么要进行随机化测试呢？编写测试的最佳实践是各个测试项之间隔离。例如，它们不应该依赖于执行顺序或共享变量。这些隐藏的依赖关系可能意味着一个可能的测试错误，更糟糕的是一个在测试中不会被发现的错误。为了防止这种情况，我们应该使用-shuffle参数设置要进行随机化测试，该参数设置为on表示启用随机化测试，设置off表示关闭随机化测试，默认是禁用的。

$ go test -shuffle=on -v .
但是，在某些情况下，我们希望以相同的顺序再次运行测试。例如，在CI期间测试失败，我们可能希望在本地重现错误。这时候，我们可以传递用于随机化测试的种子值给-shuffle参数。在执行测试的命令中携带-v参数，可以获得运行shuffled时种子值。下一次运行的时候将此种子值给-shuffle参数,可以使得运行的顺序跟获得种子值那次一样。

go test -shuffle=on -v example1_test.go                                                                           
-test.shuffle 1658273859224698000
=== RUN   TestC
--- PASS: TestC (0.00s)
=== RUN   TestA
=== PAUSE TestA
=== RUN   TestB
=== PAUSE TestB
=== CONT  TestA
--- PASS: TestA (0.00s)
=== CONT  TestB
--- PASS: TestB (0.00s)
PASS
ok      command-line-arguments  0.006s
上面的随机化测试通过加-v参数打印了种子值1658273859224698000。下面测试时通过将-shuffle设置为1658273859224698000以保持运行的顺序与上面的一样。通过输出信息可以看到，运行顺序与上面是一样的。

go test -shuffle=1658273859224698000 -v example1_test.go                                                             
-test.shuffle 1658273859224698000
=== RUN   TestC
--- PASS: TestC (0.00s)
=== RUN   TestA
=== PAUSE TestA
=== RUN   TestB
=== PAUSE TestB
=== CONT  TestA
--- PASS: TestA (0.00s)
=== CONT  TestB
--- PASS: TestB (0.00s)
PASS
ok      command-line-arguments  0.009s
总结，我们应该对现在的测试参数/标志熟悉使用，并随时了解Go最新版本提供功能。并行运行测试是减少运行所有测试的整体执行时间的好方法。此外，shuffle模式可以帮助我们发现隐藏的依赖关系。这些依赖关系可能意味着以相同顺序运行测试暴露不出来问题，但是通过随机打乱执行顺序可以提高暴露问题的机会。