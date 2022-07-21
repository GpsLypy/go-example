package main

import "fmt"

//版本一
func BinaryFind(arr *[6]int, leftIndex, rightIndex, findVal int) {
	if leftIndex > rightIndex {
		fmt.Println("not find")
		return
	}

	mid := leftIndex + ((rightIndex - leftIndex) >> 2)

	if findVal < (*arr)[mid] {
		BinaryFind(arr, leftIndex, mid-1, findVal)
	} else if findVal > (*arr)[mid] {
		BinaryFind(arr, mid+1, rightIndex, findVal)
	} else {
		fmt.Printf("find %v\n", mid)
	}
}

//版本二
func BinaryFind2(nums []int, target int) int {
	low := 0
	high := len(nums)

	for low < high {
		mid := low + (high-low)/2
		if nums[mid] > target {
			high = mid
		} else if nums[mid] < target {
			low = mid + 1
		} else {
			return mid
		}
	}

	return -1
}

func main() {
	arr := []int{15, 25, 35, 45, 55, 65}
	fmt.Println(BinaryFind2(arr, 12))
	//BinaryFind(&arr, 0, len(arr)-1, 15)
}
