package main

import (
	"fmt"
	"strings"
)

func main() {
	var s string = "言念君子，温其如玉"
	fmt.Println(len(s)) // 27,汉字占3个字节长度

	//字符串拼接
	var s1 = "python"
	var s2 = "go"
	fmt.Println(s1 + s2) //pythongo
	s3 := fmt.Sprintf("%s---%s", s1, s2)
	fmt.Println(s3) //python---go

	//分割字符串
	res := strings.Split(s, ",")
	for _, v := range res {
		fmt.Println(v)
	}

	res = append(res, "hahh")
	fmt.Println(res) //[言念君子，温其如玉]

	// 判断是否包含
	res2 := strings.Contains(s1, "on")
	fmt.Println(res2) // true

	// 判断前缀后缀
	res3 := strings.HasPrefix(s1, "py")
	res4 := strings.HasSuffix(s1, "on")
	fmt.Println(res3, res4) // true  true

	// 子串出现的位置
	var s4 string = "applepen"
	fmt.Println(strings.Index(s4, "p"))     // 1
	fmt.Println(strings.LastIndex(s4, "p")) // 5

	// join()
	a1 := []string{"python", "php", "go"}
	fmt.Println(strings.Join(a1, "-")) // python-php-go
}
