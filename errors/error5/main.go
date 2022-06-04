package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func ReadFile(path string) {
	f, err := os.Open(path)
	errors.Wrap(err, "opem faild") //stack
	errors.Wrapf(err, "failed to open %q", path)
	defer f.Close()
}

func ReadConfig() (byte, error) {
	//home :=os.Getenv("HOME")
	config, err := ReadFile("/temp")
	return config, errors.WithMessage(err, "could not read config")
}

//第一次出现错误的时候wrap
func main() {
	_, err := ReadConfig()
	if err != nil {
		fmt.Printf("originnal error:%T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace:\n%+v\n", err)
	}
}
