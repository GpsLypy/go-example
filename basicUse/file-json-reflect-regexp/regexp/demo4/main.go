package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "regexp"
    "strconv"
)

func savToFile(index int, filmName, filmScore [][]string) {
    f, err := os.Create("第" + strconv.Itoa(index) + "页.txt")
    if err != nil {
        fmt.Println("os create err", err)
        return
    }
    defer f.Close()
    // 查出有多少条
    n := len(filmName)
    // 先写抬头 名称     评分
    f.WriteString("电影名称" + "\t\t\t" + "评分" + "\n")
    for i := 0; i < n; i++ {
        f.WriteString(filmName[i][1] + "\t\t\t" + filmScore[i][1] + "\n")
    }
}

func main() {
    var start, end int
    fmt.Print("请输入要爬取的起始页")
    fmt.Scan(&start)
    fmt.Print("请输入要爬取的终止页")
    fmt.Scan(&end)
    working(start, end)
}

func working(start int, end int) {
    fmt.Printf("正在爬取%d到%d页", start, end)
    for i := start; i <= end; i++ {
        SpiderPage(i)
    }
}

// 爬取一个豆瓣页面数据信息保存到文档
func SpiderPage(index int) {
    // 获取url
    url := "https://movie.douban.com/top250?start=" + strconv.Itoa((index-1)*25) + "&filter="

    // 爬取url对应页面
    result, err := HttpGet(url)
    if err != nil {
        fmt.Println("httpget err", err)
        return
    }
    //fmt.Println("result=", result)
    // 解析，编译正则表达式  ---电影名称
    ret := regexp.MustCompile(`<img width="100" alt="(?s:(.*?))"`)
    filmName := ret.FindAllStringSubmatch(result, -1)
    for _, name := range filmName {
        fmt.Println("name", name[1])
    }

    ret2 := regexp.MustCompile(`<span class="rating_num" property="v:average">(?s:(.*?))<`)
    filmScore := ret2.FindAllStringSubmatch(result, -1)
    for _, score := range filmScore {
        fmt.Println("score", score[1])
    }

    savToFile(index, filmName, filmScore)

}

// 爬取指定url页面，返回result
func HttpGet(url string) (result string, err error) {
    req, _ := http.NewRequest("GET", url, nil)
    // 设置头部信息
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.115")
    resp, err1 := (&http.Client{}).Do(req)
    //resp, err1 := http.Get(url)  //此方法已经被豆瓣视为爬虫，返回状态吗为418，所以必须伪装头部用上述办法
    if err1 != nil {
        err = err1
        return
    }
    defer resp.Body.Close()

    buf := make([]byte, 4096)

    //循环爬取整页数据
    for {
        n, err2 := resp.Body.Read(buf)
        if n == 0 {
            break
        }
        if err2 != nil && err2 != io.EOF {
            err = err2
            return
        }
        result += string(buf[:n])
    }

    return

}