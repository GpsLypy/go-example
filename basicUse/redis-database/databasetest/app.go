package main

import "database/sql"

//记录用户浏览产品信息
func recordStats(db *sql.DB, userID, productID int64) (err error) {
	//开启事务
	//操作views和product_viewers两张表
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()

		default:
			tx.Rollback()
		}
	}()

	//更新products表
	if _, err = tx.Exec("UPDATE products SET views = views + 1"); err != nil {
		return
	}

	//product_viewers表中插入一条数据
	if _, err = tx.Exec(
		"INSERT INTO product_viewers(user_id,product_id) VALUES(?,?)",
		userID, productID); err != nil {
		return
	}

	return
}

func main() {
	//测试过程不必建立真正的连接
	db, err := sql.Open("mysql", "root@/blog")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err = recordStats(db, 1, 5); err != nil {
		panic(err)
	}
}
