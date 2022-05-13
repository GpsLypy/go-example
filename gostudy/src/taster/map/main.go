package main

import (
	"fmt"
	"strings"
)

// // ConvertStrSlice2Map 将字符串 slice 转为 map[string]struct{}
// func ConvertStrSlice2Map(sl []string) map[string]struct{} {
// 	set := make(map[string]struct{}, len(sl))
// 	for _, v := range sl {
// 		set[v] = struct{}{}
// 	}
// 	return set
// }

// // ContainsInMap 判断字符串是否在 map 中
// func ContainsInMap(m map[string]struct{}, s string) bool {
// 	_, ok := m[s]
// 	return ok
// }
// ConvertStrSlice2Map 将字符串 slice 转为 map[string]struct{}
func ConvertStrSlice2Map(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// 过滤函数
func filterIgnoreTables(tables string, ignoreTables []string) string {
	if len(ignoreTables) == 0 {
		return tables
	}

	//切割字符串request.Tables放到切片中tempTables
	tempTables := strings.Split(tables, ",")
	//遍历切片拿出表名比较是否是忽略表，是忽略，否加入新的切片中
	newTables := make([]string, 0, len(tempTables))
	m := ConvertStrSlice2Map(ignoreTables)
	for _, v := range tempTables {
		if _, ok := m[v]; ok {
			continue
		}
		newTables = append(newTables, v)
	}
	return strings.Join(newTables, ",")
}

func main() {
	s1 := []string{"t1", "t2"}
	s2 := "t1,t2,t3"
	ret := filterIgnoreTables(s2, s1)
	fmt.Println(ret)
}
