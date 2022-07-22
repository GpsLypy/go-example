package main

import (
	"fmt"
)

// sql注入示例
func sqlInject(name string) {
	sqlStr := fmt.Sprintf("select uid, name, phone from user where name='%s'", name)
	fmt.Printf("SQL:%s\n", sqlStr)
	// ret, err := db.Exec(sqlStr)
	// if err != nil {
	// 	fmt.Printf("update failed, err:%v\n", err)
	// 	return
	// }
	// n, err := ret.RowsAffected() // 操作影响的行数
	// if err != nil {
	// 	fmt.Printf("get RowsAffected failed, err:%v\n", err)
	// 	return
	// }
	// fmt.Printf("get success, affected rows:%d\n", n)
}

func main() {

	sqlInject("xxx' or 1=1#")
}
