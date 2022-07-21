package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/gocron")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	//连接数据库查询
	for i := 0; i < 100; i++ {
		go func(i int) {
			mSql := "select * from user"
			rows, _ := db.Query(mSql)
			rows.Close() //这里如果不释放连接到池里，执行5次后其他并发就会阻塞
			fmt.Println("第 ", i)
		}(i)
	}

	for {
		time.Sleep(time.Second)
	}
}
