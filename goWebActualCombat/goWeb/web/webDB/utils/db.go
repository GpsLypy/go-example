package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//全局变量
var (
	Db  *sql.DB
	err error
)

func init() {
	Db, err = sql.Open("mysql", "root:password#dbr@tcp(10.100.156.210:3306)/shop_lyp")
	if err != nil {
		panic(err.Error())
	}
}
