//goroutine泄漏,没有人会对此通道发信号

/*
func leak() {
	ch := make(chan int)
	go func() {
		val := <-ch
		fmt.Println("We received a value:", val)
	}()
}


//search 函数是一个模拟实现，用于模拟长时间运行的操作，如数据库查询或rpc调用，本例子中硬编码为200ms
func search(term string) (string,error){
	time.Sleep(200*time.Millisecond)
	return "some value",nil
}

//定义一个process函数，接受字符串参数，传递给search。对于某些应用程序，顺序调用产生的延迟可能是不可接受的
func process(term string) error{
	record,err:=search(term)
	if err!=nil{
		return err
	}
	fmt.Println("Received:",record)
	return nil
}

*/

/*
type result struct{
	record string
	err error
}

func process(term string) error{
	ctx,cancel :=context.WithTimeout(context.Background(),100*time.Millisecond)
	defer cancel()

	ch :=make(chan result)
	go func(){
		record,err:=search(term)
		ch <-result{record,err}
	}()

	select{
		case <- ctx.Done():
		return errors.New("search canceled")
		case <- ch :
		if result.err!=nil{
			return result.err
		}
		fmt.Println("Receive:",result.record)
		return nil
	}
}

*/


