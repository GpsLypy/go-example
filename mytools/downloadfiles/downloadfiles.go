package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/qianlnk/pgbar"
)

/**
* 需求:
1. 多协程下载文件
2.断点续连
**/

func main() {
	//获取要下载文件
	DownloadFileName := "/home/go-example/demo/d1/main.go"
	//copy的文件
	copyFileName := "./copy.go"
	storgeFileName := "./current.go" //记录读取源文件的偏移量
	//打开文件
	sfile, err := os.Open(DownloadFileName)
	if err != nil {
		panic(err)
	}
	defer sfile.Close()
	//获取文件大小
	info, _ := sfile.Stat()
	downloadSize := info.Size()
	var scount int64 = 1
	if downloadSize%5 == 0 {
		scount *= 5
	} else {
		scount *= 10
	}
	//分给每个协程的大小
	si := downloadSize / scount
	//分片数10 每个分片交给一个协程去做，需要10个协程
	fmt.Printf("文件总大小：%v, 分片数：%v,每个分片大小：%v\n", downloadSize, scount, si)
	//open copy file
	copyFile, err := os.OpenFile(copyFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	storgeFile, err := os.OpenFile(storgeFileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer copyFile.Close()

	var currentIndex int64 = 0
	wg := sync.WaitGroup{}
	fmt.Println("协程进度条")
	pgb := pgbar.New("")
	for ; currentIndex < scount; currentIndex++ {
		wg.Add(1)
		go func(current int64) {
			p := pgb.NewBar(fmt.Sprint((current+1))+"st", int(si))
			//p.SetSpeedSection(900, 100)
			b := make([]byte, 1024)
			bs := make([]byte, 16)
			//返回读取的字节数放到currentIndex里面
			currentIndex, _ := storgeFile.ReadAt(bs, current*16)
			//取出所有整数
			//MustCompile类似Compile但会在解析失败时panic，主要用于全局正则表达式变量的安全初始化。
			//解析并返回一个正则表达式。如果成功返回，该Regexp就可用于匹配文本。
			reg := regexp.MustCompile(`\d+`) //查找连续的数字
			countStr := reg.FindString(string(bs[:currentIndex]))
			total, _ := strconv.ParseInt(countStr, 10, 0)
			progressBar := 1
			for {
				if total >= si {
					wg.Done()
					break
				}
				//从指定位置开始读
				n, err := sfile.ReadAt(b, current*si+total)
				if err == io.EOF {
					wg.Done()
					break
				}
				//从指定位置开始写
				copyFile.WriteAt(b, current*si+total)
				//FoematInt 返回total的10进制的字符串表示
				storgeFile.WriteAt([]byte(strconv.FormatInt(total, 10)+" "), current*16)
				total += int64(n)
				if total >= si/10*int64(progressBar) {
					progressBar += 1
					p.Add(int(si / 10))
				}

			}

		}(currentIndex)
	}
	wg.Wait()
	storgeFile.Close()
	os.Remove(storgeFileName)
	fmt.Println("下载完成")
}
