package main

import (
	"fmt"
	"webapp/webDB/model"
	_ "webapp/webDB/utils"
)

func main() {
	//init()
	// user := &model.User{}
	// user.AddUser()
	//str := "\"'\""
	//带有反斜杠
	//str := "'a\\L'"
	str := "'strkey1 + \\中文"
	str = b(str)
	fmt.Println(str)
}

// escapeStringBackslash is similar to escapeBytesBackslash but for string.
//也会自动加上单引号
func EscapeStringBackslash(buf []byte, v string) []byte {
	pos := len(buf)
	buf = ReserveBuffer(buf, len(v)*2)

	for i := 0; i < len(v); i++ {
		c := v[i]
		switch c {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = '0'
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		default:
			buf[pos] = c
			pos++
		}
	}

	return buf[:pos]
}

func a(v string) string {
	var buf []byte
	buf = append(buf, '\'')
	buf = EscapeStringQuotes(buf, v)

	buf = append(buf, '\'')

	return string(buf)
}

func b(v string) string {
	var buf []byte
	buf = append(buf, '\'')
	buf = EscapeStringBackslash(buf, v)

	buf = append(buf, '\'')

	return string(buf)
}

//不带反斜杠,自动加单引号
func EscapeStringQuotes(buf []byte, v string) []byte {
	pos := len(buf)
	buf = ReserveBuffer(buf, len(v)*2)

	for i := 0; i < len(v); i++ {
		c := v[i]
		if c == '\'' {
			buf[pos] = '\''
			buf[pos+1] = '\''
			pos += 2
		} else {
			buf[pos] = c
			pos++
		}
	}

	return buf[:pos]
}

func ReserveBuffer(buf []byte, appendSize int) []byte {
	newSize := len(buf) + appendSize
	if cap(buf) < newSize {
		// Grow buffer exponentially
		newBuf := make([]byte, len(buf)*2+appendSize)
		copy(newBuf, buf)
		buf = newBuf
	}
	return buf[:newSize]
}
