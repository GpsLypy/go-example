package main

import (
	"fmt"
	"reflect"
)

func main() {
	//1. 获取变量类型
	fmt.Println("获取变量类型")
	fmt.Println(reflect.TypeOf(10))                          //int
	fmt.Println(reflect.TypeOf(10.0))                        //float64
	fmt.Println(reflect.TypeOf(struct{ age int }{10}))       //struct { age int }
	fmt.Println(reflect.TypeOf(map[string]string{"a": "a"})) //map[string]string
	fmt.Println("")
	//2. 获取变量值
	fmt.Println("获取变量值")
	fmt.Println(reflect.ValueOf("hello word"))                //hello word
	fmt.Println(reflect.ValueOf(struct{ age int }{10}))       //{10}
	fmt.Println(reflect.TypeOf(struct{ age int }{10}).Kind()) //struct
	//类型判断
	if t := reflect.TypeOf(struct{ age int }{10}).Kind(); t == reflect.Struct {
		fmt.Println("是结构体")
	} else {
		fmt.Println("不是结构体")
	}
	//修改目标对象
	str := "hello word"
	//普通变量修改
	reflect.ValueOf(&str).Elem().SetString("张三")
	fmt.Println(str)
	//结构体变量修改
	user := User{Name: "张三", Age: 10}
	//Elem() 获取user原始的值
	elem := reflect.ValueOf(&user).Elem()
	//FieldByName() 通过Name返回具有给定名称的结构字段 通过SetString 修改原始的值
	elem.FieldByName("Name").SetString("李四")
	elem.FieldByName("Age").SetInt(18)
	fmt.Println(user)
	//获取结构体的标签的值
	fmt.Println(reflect.TypeOf(&user).Elem().Field(0).Tag.Get("name"))
	//调用无参方法
	reflect.ValueOf(&user).MethodByName("Say").Call([]reflect.Value{})
	reflect.ValueOf(user).MethodByName("Say").Call(make([]reflect.Value, 0))
	//调用有参方法
	reflect.ValueOf(user).MethodByName("SayContent").Call([]reflect.Value{reflect.ValueOf("该说话了"), reflect.ValueOf(1)})
	//调用本地的方法
	reflect.ValueOf(Hello).Call([]reflect.Value{})
	reflect.ValueOf(Hello).Call(nil)
	fmt.Printf("%#v\n", reflect.TypeOf(user).Field(0))
}
func Hello() {
	fmt.Println("hello")
}

type Person struct {
	Name string
}
type User struct {
	Person        // //反射会将匿名字段作为一个独立字段来处理
	Name   string `json:"name" name:"张三"`
	Age    int
}

func (_ User) Say() {
	fmt.Println("user 说话")
}
func (_ User) SayContent(content string, a int) {
	fmt.Println("user", content, a)
}
