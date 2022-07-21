package model

import (
	"fmt"
	"testing"
)

//在测试函数执行之前做准备工作
// func TestMain(m *testing.M) {
// 	fmt.Println("lalla")
// 	m.Run()
// }

// func TestUser(t *testing.T) {
// 	fmt.Println("开始测试User中的相关方法")
// 	//通过T.Run()来执行子测试函数
// 	//t.Run("准备开始执行", TestAddUser)
// 	//t.Run("查询测试", testGetUserById)
// 	t.Run("查询测试", testGetUsers)
// }

func TestAddUser(t *testing.T) {
	fmt.Println("测试添加用户：")
	user := &User{}

	fmt.Println("我开始执行了。。。")
	user.AddUser()
	//user.AddUser2()

}

// func testGetUserById(t *testing.T) {
// 	fmt.Println("测试查询")
// 	user := &User{
// 		ID: 1,
// 	}
// 	u, _ := user.GetUserById()
// 	fmt.Println("得到user:", u)
// }

// func testGetUsers(t *testing.T) {
// 	fmt.Println("测试查询数据库所有记录")
// 	user := &User{}

// 	//var users []*User
// 	users, _ := user.GetUsers()
// 	for _, v := range users {
// 		fmt.Println(v)
// 	}
// }
