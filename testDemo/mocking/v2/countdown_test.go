package main

import (
	"bytes"
	"testing"
)

func TestCountdown(t *testing.T) {
	//我们的目的是让 Countdown 函数将数据写到某处，io.writer 就是作为 Go 的一个接口来抓取数据的一种方式。
	buffer := &bytes.Buffer{}
	Countdown(buffer)

	got := buffer.String()
	want := `3
2
1
Go!`
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}

}
