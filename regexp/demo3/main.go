package main

import (
 "fmt"
 "io"
 "net/http"
 "regexp"
 "strconv"
)

var readCount int = 0
var commentCount int = 0
var diggCount int = 0

//http读取网页数据写入result返回
func HttpGet(url string) (result string, err error) {
 resp, err1 := http.Get(url)
 if err1 != nil {
  err = err1
  return
 }
 defer resp.Body.Close()

 buf := make([]byte, 4096)

 for {
  n, err2 := resp.Body.Read(buf)
  //fmt.Println(url)
  if n == 0 {
   break
  }
  if err2 != nil && err2 != io.EOF {
   err = err2
   return
  }
  result += string(buf[:n])
 }
 return result, err
}

//横向纵向爬取文章标题数据，并累计数值
func SpiderPageDB(index int, page chan int) {
 url := "https://www.cnblogs.com/littleperilla/default.html?page=" + strconv.Itoa(index)

 result, err := HttpGet(url)

 if err != nil {
  fmt.Println("HttpGet err:", err)
  return
 }

 str := regexp.MustCompile("post-view-count\">阅读[(](?s:(.*?))[)]</span>")
 alls := str.FindAllStringSubmatch(result, -1)
 for _, j := range alls {
  temp, err := strconv.Atoi(j[1])
  if err != nil {
   fmt.Println("string2int err:", err)
  }
  readCount += temp
 }

 str = regexp.MustCompile("post-comment-count\">评论[(](?s:(.*?))[)]</span>")
 alls = str.FindAllStringSubmatch(result, -1)
 for _, j := range alls {
  temp, err := strconv.Atoi(j[1])
  if err != nil {
   fmt.Println("string2int err:", err)
  }
  commentCount += temp
 }

 str = regexp.MustCompile("post-digg-count\">推荐[(](?s:(.*?))[)]</span>")
 alls = str.FindAllStringSubmatch(result, -1)
 for _, j := range alls {
  temp, err := strconv.Atoi(j[1])
  if err != nil {
   fmt.Println("string2int err:", err)
  }
  diggCount += temp
 }

 page <- index
}

//主要工作方法
func working(start, end int) {
 fmt.Printf("正在从%d到%d爬取中...\n", start, end)

 //channel通知主线程是否所有go都结束
 page := make(chan int)

 //多线程go程同时爬取
 for i := start; i <= end; i++ {
  //go SpiderPageDB(i, page)
  SpiderPageDB(i, page)
 }

 for i := start; i <= end; i++ {
  fmt.Printf("拉取到%d页\n", <-page)
 }
}

//入口函数
func main() {
 //输入爬取的起始页
// var start, end int
 fmt.Print("startPos:")
 //fmt.Scan(&start)
 fmt.Print("endPos:")
 //fmt.Scan(&end)

 //working(start, end)
 working(1, 3)

 fmt.Println("阅读:", readCount)
 fmt.Println("评论:", commentCount)
 fmt.Println("推荐:", diggCount)
}