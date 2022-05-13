//++++++++++++++++++++++++++++++++++++++++
// 《Go Web编程实战派从入门到精通》源码
//++++++++++++++++++++++++++++++++++++++++
// Author:廖显东（ShirDon）
// Blog:https://www.shirdon.com/
// 仓库地址：https://gitee.com/shirdonl/goWebActualCombat
// 仓库地址：https://github.com/shirdonl/goWebActualCombat
//++++++++++++++++++++++++++++++++++++++++

package main

import (
	"fmt"
	"regexp"
)

func main()  {
	re := regexp.MustCompile(`Go(\w+)`)
	fmt.Println(re.ReplaceAllString("Hello Gopher，Hello GoLang", "Java$1"))

	text := "Hello Gopher，Hello Go Web"
	reg := regexp.MustCompile(`\w+`)
	fmt.Println(reg.MatchString(text))


}
