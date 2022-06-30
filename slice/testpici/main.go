package main

import (
	"fmt"
)

func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	size := 3
	var j int
	for i := 0; i < len(s); i += size {
		j += size
		if j > len(s) {
			j = len(s)
		}
		// do what do you want to with the sub-slice, here just printing the sub-slices
		fmt.Println(s[i:j])
	}
}
