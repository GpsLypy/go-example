/*
gollback也是用来处理一组子任务的执行的，不过它解决了ErrGroup收集子任务返回结果的痛点。
使用ErrGroup时，如果你要收到子任务的结果和错误，你需要定义额外的变量收集执行结果和错误，但是这个库可以提供更便利的方式。

我刚刚在说官方扩展库ErrGroup的时候，举了一些例子（返回第一个错误的例子和返回所有子任务错误的例子），
在例子中，如果想得到每一个子任务的结果或者error，我们需要额外提供一个result slice进行收集。使用gollback的话，就不需要这些额外的处理了，因为它的方法会把结果和error信息都返回。

接下来，我们看一下它提供的三个方法，分别是All、Race和Retry。
*/

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/vardius/gollback"
	"time"
)

/*
All方法

All方法的签名如下：

func All(ctx context.Context, fns ...AsyncFunc) ([]interface{}, []error)
它会等待所有的异步函数（AsyncFunc）都执行完才返回，而且返回结果的顺序和传入的函数的顺序保持一致。第一个返回参数是子任务的执行结果，第二个参数是子任务执行时的错误信息。

其中，异步函数的定义如下：

type AsyncFunc func(ctx context.Context) (interface{}, error)
可以看到，ctx会被传递给子任务。如果你cancel这个ctx，可以取消子任务。

我们来看一个使用All方法的例子：
*/

func main() {
	rs, errs := gollback.Race(
		context.Background(),
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(3 * time.Second)
			return 1, nil
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("faild")
		},
		func(ctx context.Context) (interface{}, error) {
			return 3, nil
		},
	)
	fmt.Println(rs)   //输出子任务的结果
	fmt.Println(errs) //输出子任务的错误信息
}

/*
Race方法

Race方法跟All方法类似，只不过，在使用Race方法的时候，只要一个异步函数执行没有错误，就立马返回，而不会返回所有的子任务信息。如果所有的子任务都没有成功，就会返回最后一个error信息。

Race方法签名如下：

func Race(ctx context.Context, fns ...AsyncFunc) (interface{}, error)
如果有一个正常的子任务的结果返回，Race会把传入到其它子任务的Context cancel掉，这样子任务就可以中断自己的执行。

Race的使用方法也跟All方法类似，我就不再举例子了，你可以把All方法的例子中的All替换成Race方式测试下。
*/

/*
Retry方法

Retry不是执行一组子任务，而是执行一个子任务。如果子任务执行失败，它会尝试一定的次数，如果一直不成功 ，就会返回失败错误 ，如果执行成功，它会立即返回。如果retires等于0，它会永远尝试，直到成功。

func Retry(ctx context.Context, retires int, fn AsyncFunc) (interface{}, error)
再来看一个使用Retry的例子：

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/vardius/gollback"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试5次，或者超时返回
	res, err := gollback.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
		return nil, errors.New("failed")
	})

	fmt.Println(res) // 输出结果
	fmt.Println(err) // 输出错误信息
}

*/
