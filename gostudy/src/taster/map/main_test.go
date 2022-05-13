package main

import (
	"fmt"
	"testing"
)

func TestFilterIgnoreTables(t *testing.T) {
	s1 := []string{"t1", "t2"}
	s2 := "t1,t2,t3"
	ret := filterIgnoreTables(s2, s1)
	fmt.Println(ret)
}
