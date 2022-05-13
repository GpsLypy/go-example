package concurrency

import (
	_ "net/http"
)

type WebsiteChecker func(string) bool

// func CheckWebsite(url string)bool{
// 	response,err :=http.Head(url)
// 	if err!=nil{
// 		return false
// 	}
// 	if response.StatusCode!=http.StatusOK{
// 		return false
// 	}
// 	return true
// }

// func CheckWebsite(wc WebsiteChecker, urls []string) map[string]bool {
// 	results := make(map[string]bool)
// 	for _, url := range urls {
// 		//这里的问题是变量 url 被重复用于 for 循环的每次迭代 —— 每次都会从 urls 获取新值。
// 		//但是我们的每个 goroutine 都是 url 变量的引用 —— 它们没有自己的独立副本。
// 		//所以他们都会写入在迭代结束时的 url —— 最后一个 url。这就是为什么我们得到的结果是最后一个 url。
// 		go func() {
// 			results[url] = wc(url)
// 		}()
// 	}
// 	return results
// }

type result struct {
	string
	bool
}

//我们可以通过使用 channels 协调我们的 goroutines 来解决这个数据竞争
func CheckWebsite(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)
	for _, url := range urls {
		//这里的问题是变量 url 被重复用于 for 循环的每次迭代 —— 每次都会从 urls 获取新值。
		//但是我们的每个 goroutine 都是 url 变量的引用 —— 它们没有自己的独立副本。
		//所以他们都会写入在迭代结束时的 url —— 最后一个 url。这就是为什么我们得到的结果是最后一个 url。
		go func(u string) {
			//results[u] = wc(u)
			resultChannel <- result{u, wc(u)}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		result := <-resultChannel
		results[result.string] = result.bool
	}
	return results
}
