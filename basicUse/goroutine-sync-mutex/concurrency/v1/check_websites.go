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

func CheckWebsite(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	for _, url := range urls {
		results[url] = wc(url)
	}
	return results
}
