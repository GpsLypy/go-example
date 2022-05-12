package main

import (
	"errors"
	"fmt"
)

type errorString string

func (e errorString) Error() string {
	return string(e)
}

func New(text string) error {
	return errorString(text)
}

var ErrNamedType = New("EOF")
var ErrStructType = errors.New("EOF")

func main() {
	if ErrNamedType == New("EOF") {
		fmt.Println("Named Type Error")
	}
	//struct 包裹，并取地址，避免了两个人定义了相同错误造成错误判断。比较的是内存地址，而不是像浅拷贝那样比较结构体的每个字段值是否相等
	if ErrStructType == errors.New("EOF") {
		fmt.Println("Struct Type Error")
	}
}
