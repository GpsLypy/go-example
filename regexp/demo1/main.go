package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
)

//Compile——判断一个正则表达式是否合法，采用最左最短匹配

func TestCompile() {
	ret, err := regexp.Compile(`^[0-9]`)
	if err != nil {
		errors.WithMessage(err, "TestCompile faild")
	}
	//FindString——在传入字符串中查找正则表达式匹配的内容，并返回第一个匹配的对象
	fmt.Println(ret.FindString("2312eweqeqw"))
}

//CompilePOSIX——和complie一样，使用的是POSIX语法规则，采用最左最长匹配
func TestCompilePOSIX() {
	ret, err := regexp.CompilePOSIX(`[[:word:]]+`)
	if err != nil {
		errors.WithMessage(err, "TestCompilePOSIX faild")
	}
	fmt.Printf("%q\n", ret.FindString("hello world"))
}

//MustCompilePOSIX和CompilePOSIX作用一样，不一样的是当正则表达式不合法的时候会报出异常而不是错误
func TestMustCompilePOSIX() {
	ret := regexp.MustCompilePOSIX(`[[:word:]].+`)
	fmt.Printf("%q\n", ret.FindString("hello world"))
}

//统计正则表达式中分组的个数
//\x02和\x07表示16进制的ASCII代码。
//而\0表示结束.
//写出来是0x2
func TestNumSubexp() {
	ret := regexp.MustCompile(`(?U)(?:Hello)(\s+)(\w+)`)
	fmt.Printf("%q\n", ret.NumSubexp())
}

//SubexpNames——返回正则表达式中的分组名称,一个圆括号表示一个分组
//返回值[0]为整个表达式的名称
//返回值[1]为分组1的名称，其余依次向后推
func TestSubexpNames() {
	ret := regexp.MustCompile("(?:Hello)([321])(\t)()(World)")
	fmt.Printf("%q\n", ret.SubexpNames())
}

//LiteralPrefix——返回所有匹配项共同拥有的前缀（去除可变元素）
//返回值说明：第一个返回值表示共同拥有的前缀
//第二个返回值表示如果文字字符串包含整个正则表达式，则返回布尔值true。
func TestLiteralPrefix() {
	ret := regexp.MustCompile(`Hello[\w\s]`)
	fmt.Println(ret.LiteralPrefix())
	ret = regexp.MustCompile(`hello`)
	fmt.Println(ret.LiteralPrefix())
}

//MatchRead——判断在r中能否找到正则表达式所匹配的子串
//参数描述：第一个表示要查找的正则表达式，第二个表示要在其中查找的RuneReader接口
//返回值描述：第一个返回值表示是否匹配到，匹配到了返回true否则返回false，第二个返回值为err
func TestMatchReader() {
	ret := bytes.NewReader([]byte("jgdsafu gfjasgfgs"))
	fmt.Println(regexp.MatchReader("j.*", ret))
}

//MatchString——判断在给定字符串中能否找到匹配的串
//参数描述：第一个参数为给定的正则表达式，第二个为给定的要匹配的文本串
//返回值描述：第一个返回值表示是否匹配，匹配到了返回true，否则返回false，第二个返回值表示error
func TestMatchString() {
	fmt.Println(regexp.MatchString(`^[123]`, "13221")) //^表示开头必须匹配123其中一个啦
	fmt.Println(regexp.MatchString(`^[123]`, "456789"))
	fmt.Println(111)
}

//Match——判断在给定[]byte中能否找到正则表达式所匹配的子串
//参数说明：第一个参数为给定的正则表达式，第二个参数为给定的要查找的文本
//返回值说明：第一个表示正则表达试是否有匹配，第二个参数为error
func TestMatch() {
	fmt.Println(regexp.Match(`[e]`, []byte("321cwxfsfsdfs")))
}

//ReplaceAllString——将检索到的匹配项替换为给定的，并返回替换后的结果
func TestReplaceAllString() {
	str := "ewq 413, 4324, dsGuds Go"
	ret := regexp.MustCompile(`(Hell|G)o`)
	rep := "{1}ooo"
	fmt.Printf("%q\n", ret.ReplaceAllString(str, rep))
}

//ReplaceAllLiteralString——在给定文本中搜索，并将给定文本替换为第二个参数给定的内容
//参数描述：第一个参数表示给定的文本，第二个参数表示给的需要替换的内容(如果其中含有分组引用符，将分组引用符当做普通字符处理)
//返回值描述：返回替换后的文本，
func TestReplaceAllLiteralString() {
	src := "231 213 fdshof fgodsugi Go qhwie hifewir"
	reg := regexp.MustCompile(`(fg|hi|sh|G)o`)
	rep := "${1}AAA"
	fmt.Printf("%q\n", reg.ReplaceAllLiteralString(src, rep))
}

//QuoteMeta——将字符串中的特殊字符转换为其转义格式（特殊字符包括——\.+*?()|[]{}^$）
func TestQuoteMeta() {
	fmt.Println(regexp.QuoteMeta("^ewq[312]+(?:32){32|1}.*fsd$"))
}

//Find——在给定文本中查找正则表达式中匹配的内容，并返回第一个匹配的对象
func TestFind() {
	reg := regexp.MustCompile(`\w+`)
	fmt.Printf("%q\n", reg.Find([]byte("f11gdsgfs dfjdgsu")))
}

//FndIndex——在给定文本中查找匹配正则表达式的内容，并返回匹配内容的起始位置和结束位置[起始位置 结束位置]
func TestFindIndex() {
	reg := regexp.MustCompile(`\w+`)
	fmt.Println(reg.FindIndex([]byte("43242324 432432")))
}

//FindStringIndex——在给定文本中查找满足正则表达式的串并返回第一个匹配串的起始和结束位置

//FindSubmatch——在文本中查找正则表达式被匹配的第一个内容，同时返回子表达式匹配的内容{{完整匹配项}{起始子匹配项}{结束子匹配项}}
func TestFindSubmatch() {
	reg := regexp.MustCompile(`(\w)(\w)+`)
	fmt.Printf("%q\n", reg.FindSubmatch([]byte("ewewqx 423")))
}

//Expand——将template的内容处理之后追加到dst尾部，template中要有$1、$2、${name1}、${name2}这样的分组引用符
//match是由FindSubmatchIndex方法返回的结果，里面存放了各个位置的信息，如果template中有match信息，则以match为标准
//在src中取出相应的子串，替换掉template中的$1、$2等引用符号
func TestExpand() {
	reg := regexp.MustCompile(`(\w+),(\w+)`)
	src := []byte("Golang,World!")           // 源文本
	dst := []byte("Say: ")                   // 目标文本
	template := []byte("Hello $1, Hello $2") // 模板
	match := reg.FindSubmatchIndex(src)      // 解析源文本
	// 填写模板，并将模板追加到目标文本中
	fmt.Printf("%q\n", reg.Expand(dst, template, src, match))
	// "Say: Hello Golang, Hello World"
}

//ExpandString——功能和Expand一样，只不过参数类型换成了string

//Split——在给定文本项中搜索匹配项，并以匹配项为分割符，将s分割成多个子串，最多分割成多个子串，第n个子串不在进行分割
//如果n<0，则分割所有的子串，返回分割后的子串列表
func TestSplit() {
	src := "Hello World\t31321\nGolang"
	reg := regexp.MustCompile(`\s`)
	fmt.Printf("%q\n", reg.Split(src, -1))
}

func main() {
	TestCompile()
	TestCompilePOSIX()
	TestMustCompilePOSIX()
	TestNumSubexp()
	TestSubexpNames()
	TestLiteralPrefix()
	TestMatchReader()
	TestMatchString()
	TestMatch()
	TestReplaceAllString()
	TestReplaceAllLiteralString()
	TestQuoteMeta()
	TestFind()
	TestFindIndex()
	TestFindSubmatch()
	TestExpand()
	TestSplit()

	var shardingTableReg = regexp.MustCompile(`^(\\w+)_s_\\d+$`)
	var shardingTableReg2 = regexp.MustCompile(`table_1`)
	fmt.Println(shardingTableReg)
	fmt.Println(shardingTableReg2)
}
