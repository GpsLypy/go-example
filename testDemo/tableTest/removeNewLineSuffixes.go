package main

import (
	"fmt"
	"strings"
)

//s="123\r\n"
func removeNewLineSuffixes(s string) string {
	if s == "" {
		return s
	}
	if strings.HasSuffix(s, "\r\n") {
		return removeNewLineSuffixes(s[:len(s)-2])
	}
	if strings.HasSuffix(s, "\n") {
		return removeNewLineSuffixes(s[:len(s)-1])
	}
	return s
}

func main() {
	fmt.Println(removeNewLineSuffixes("123\r\n"))
}

/*
上面函数采用了递归实现。现在，假设我们要全面地测试这个函数，至少要覆盖以下几种情况：

输入的是空串
输入的字符串以\n结尾
输入的字符串以\r\n结尾
输入的字符串以多个\n结尾
输入的字符串不含换行符
*/
