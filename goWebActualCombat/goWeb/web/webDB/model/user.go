package model

import (
	"fmt"
	"taster/webDB/utils"
)

type User struct {
	Id  string
	Uid int
}

//ADD User 可以防止SQL注入
func (User *User) AddUser() error {
	sqlStr := "insert into t1(Id,Uid) values(?,?)"
	inStmt, err := utils.Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("预编译出现异常...", err)
	}

	//Sql = fmt.Sprintf("SELECT `%s` FROM `%s` WHERE `%s` >= '%s' order by `%s`  LIMIT 1;", priKey, TableName, priKey, curId, priKey)

	//interpolateParams("\"'\"", args []driver.Value) (string, error)
	i := "1"
	j := 900000
	for {

		_, err2 := inStmt.Exec(i, j)
		if err2 != nil {
			fmt.Println("执行出现异常1..", err2)
			return err2
		}
		j++
		if j == 20000001 {
			break
		}

	}
	return nil
}

//ADD User
// func (User *User) AddUser2() error {
// 	sqlStr := "insert into girls(id,name,weight,btime) values(?,?,?,?)"

// 	_, err2 := utils.Db.Exec(sqlStr, 11, "momo", 2, "2022-03-01 00:25:31")
// 	if err2 != nil {
// 		fmt.Println("执行出现异常2..", err2)
// 		return err2
// 	}
// 	return nil
// }

// func (user *User) GetUserById() (*User, error) {
// 	sqlStr := "select id,name,weight from girls where id =?"
// 	row := utils.Db.QueryRow(sqlStr, user.ID)
// 	var username string
// 	var weight float32
// 	var id int
// 	err := row.Scan(&id, &username, &weight)
// 	if err != nil {
// 		return nil, err
// 	}
// 	u := &User{
// 		ID:     id,
// 		Name:   username,
// 		Weight: weight,
// 	}
// 	return u, nil
// }

// func (user *User) GetUsers() ([]*User, error) {
// 	sqlStr := "select id,name,weight from girls"
// 	rows, err := utils.Db.Query(sqlStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var users []*User
// 	for rows.Next() {
// 		var username string
// 		var weight float32
// 		var id int
// 		err := rows.Scan(&id, &username, &weight)
// 		if err != nil {
// 			return nil, err
// 		}
// 		u := &User{
// 			ID:     id,
// 			Name:   username,
// 			Weight: weight,
// 		}
// 		users = append(users, u)
// 	}
// 	return users, nil
// }
