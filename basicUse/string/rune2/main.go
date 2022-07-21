package main

import "fmt"

func main() {
	//截取带中文字符串
	s := "Go语言编程"
	// 转成 rune 数组，需要几个字符，取几个字符
	fmt.Println(string([]rune(s)[:4])) // 输出：Go语言
}

/*
func MysqlQuoteString(i string) string {
	buf := bytes.NewBuffer(nil)
	rs := []rune(i)

	for i, v := range rs {
		var escape rune

		switch v {
		case 0:
			{
				escape = '0'
			}
		case '\n':
			{
				escape = 'n'
			}
		case '\r':
			{
				escape = 'r'
			}
		case '\\':
			{
				// Check if has the next rune
				if i == len(rs)-1 {
					escape = '\\'
				} else {
					if rs[i+1] != '%' && rs[i+1] != '_' {
						escape = '\\'
					}
				}
			}
		case '\'':
			{
				escape = '\''
			}
		case '"':
			{
				escape = '"'
			}
		case '\032':
			{
				escape = 'Z'
			}
		}

		if escape != 0 {
			buf.WriteRune('\\')
			buf.WriteRune(escape)
		} else {
			buf.WriteRune(v)
		}
	}

	return buf.String()
}

*/
