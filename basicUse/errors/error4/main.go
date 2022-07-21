package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	strCurId := "\\"
	isContain := strings.Contains(strCurId, "\\")
	if isContain {
		fmt.Println(11)
		strCurId = Stringescape(strCurId)
		fmt.Println(strCurId)
	}
}

//字符串转义不带反斜杠
func Stringescape(v string) string {
	var buf []byte
	buf = append(buf, '\'')
	isContain := strings.Contains(v, "\\")
	if isContain {
		EscapeStringBackslash(buf, v)
		buf = EscapeStringQuotes(buf, v)
	} else {
		//'\'
		buf = EscapeStringQuotes(buf, v)
	}
	buf = append(buf, '\'')

	return string(buf)
}

//带有反斜杠
// escapeStringBackslash is similar to escapeBytesBackslash but for string.
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

type Header struct {
	Key, Value string
}

type Status struct {
	Code   int
	Reason string
}

func WriteResponse(w io.Writer, st Status, headers []Header, body io.Reader) error {
	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", st.Code, st.Reason)
	//1
	if err != nil {
		return err
	}

	for _, h := range headers {
		_, err := fmt.Fprintf(w, "%s:%s \r\n", h.Key, h.Value)
		//2
		if err != nil {
			return err
		}
	}
	//3
	if _, err := fmt.Fprint(w, "\r\n"); err != nil {
		return err
	}
	//4
	_, err = io.Copy(w, body)
	return err
}

//带状态
type errWriter struct {
	io.Writer
	err error
}

func (e *errWriter) Write(buf []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}
	var n int
	n, e.err = e.Writer.Write(buf)
	return n, nil
}

func WriteResponse2(w io.Writer, st Status, headers []Header, body io.Reader) error {

	ew := &errWriter{Writer: w}
	fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", st.Code, st.Reason)

	for _, h := range headers {
		fmt.Fprintf(w, "%s:%s \r\n", h.Key, h.Value)
	}

	fmt.Fprint(w, "\r\n")

	io.Copy(w, body)
	return ew.err
}
