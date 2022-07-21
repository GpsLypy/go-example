package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//为了访问数据库而定义的全局变量
var DB *sql.DB

type User struct {
	Uid   int
	Name  string
	Phone string
}

//初始化数据库连接
func init() {
	DB, _ = sql.Open("mysql",
		"root:123456@tcp(127.0.0.1:3306)/chapter3")
}

//获取用户信息
func GetUser(uid int) (u User) {
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := DB.QueryRow("select uid,name,phone from `user` where uid=?", uid).Scan(&u.Uid, &u.Name, &u.Phone)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	return u
}

// //为了访问数据库而定义的全局变量
// var db *sql.DB

// type User struct{
// 	Uid int
// 	Name string
// 	Phone string
// }

// //初始化数据库连接
// func init(){
// 	db,_ =sql.Open("mysql","root:Mysql123..@tcp(47.99.176.238:3306)")
// }

// func GetUser(uid int)(u User){
// 	//确保在QueryRow()方法之后调用Scan方法，否则持有的数据库连接不会释放

// }
