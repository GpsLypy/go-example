package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/casbin/casbin/v2"
// 	"github.com/casbin/casbin/v2/model"
// 	xormadapter "github.com/casbin/xorm-adapter/v2"
// 	_ "github.com/go-sql-driver/mysql"
// )

// func main() {
// 	//e, err := casbin.NewEnforcer("path/to/model.conf", "path/to/policy.csv")

// 	// 使用MySQL数据库初始化一个Xorm适配器
// 	a, err := xormadapter.NewAdapter("mysql", "root:password#dbr@tcp(10.100.156.210:3306)/")
// 	if err != nil {
// 		log.Fatalf("error: adapter: %s", err)
// 	}

// 	//ACL模式
// 	m, err := model.NewModelFromString(`
// 	[request_definition]
// 	r = sub, obj, act

// 	[policy_definition]
// 	p = sub, obj, act

// 	[policy_effect]
// 	e = some(where (p.eft == allow))

// 	[matchers]
// 	m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
// 	`)
// 	if err != nil {
// 		log.Fatalf("error: model: %s", err)
// 	}

// 	e, err := casbin.NewEnforcer(m, a)
// 	if err != nil {
// 		log.Fatalf("error: enforcer: %s", err)
// 	}

// 	sub := "alice" // the user that wants to access a resource.
// 	obj := "data1" // the resource that is going to be accessed.
// 	act := "read"  // the operation that the user performs on the resource.

// 	ok, err := e.Enforce(sub, obj, act)

// 	if err != nil {
// 		// handle err
// 		fmt.Println(err)
// 	}

// 	if ok {
// 		// permit alice to read data1
// 		fmt.Println("yes")
// 	} else {
// 		// deny the request, show an error
// 		fmt.Println("no")
// 	}
// 	//roles, err := e.GetRolesForUser("alice")
// 	//fmt.Println(roles)
// 	// You could use BatchEnforce() to enforce some requests in batches.
// 	// This method returns a bool slice, and this slice's index corresponds to the row index of the two-dimensional array.
// 	// e.g. results[0] is the result of {"alice", "data1", "read"}
// 	//results, err := e.BatchEnforce([][]interface{}{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"jack", "data3", "read"}})

// }

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

//----简单实现权限验证系统。
//新增一个规则，就重新加载一次规则列表
//通过request 获取用户，查询其所属的角色是什么
//调用casbin提供的接口进行权限验证。
//增加用户所属角色接口
//1、调研
//2、测试
//3、封装
//4、验证
//5、上线

func main() {

	a, err := xormadapter.NewAdapter("mysql", "ro")
	if err != nil {
		fmt.Println("NewAdapter", err)
	}

	m, err := model.NewModelFromString(`
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`)
	if err != nil {
		fmt.Println("NewEnforcer", err)
	}
	//创建一个Casbin决策器需要有一个模型文件和策略文件为参数：
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		fmt.Println("NewEnforcer", err)
	}
	//从DB加载策略
	e.LoadPolicy()

	//获取router路由对象
	r := gin.New()
	//使用自定义拦截器中间件
	r.Use(LanjieqiHandler(e))
	//创建请求
	r.GET("/api/v1/aa", func(c *gin.Context) {
		var message string = "成功"
		var code int = 200
		var aa string = "data"
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": message,
			"data":    aa,
			"result":  "true",
		})
	})
	r.POST("/api/v1/aa", func(c *gin.Context) {
		var message string = "成功"
		var code int = 200
		var aa string = "data"
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": message,
			"data":    aa,
			"result":  "true",
		})
	})

	r.GET("/api/v1/bb", func(c *gin.Context) {
		var message string = "成功"
		var code int = 200
		var aa string = "data"
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": message,
			"data":    aa,
			"result":  "true",
		})
	})
	r.Run(":9090") //参数为空 默认监听8080端口
}

//拦截器
//拦截器
func LanjieqiHandler(e *casbin.Enforcer) gin.HandlerFunc {

	return func(c *gin.Context) {

		//获取请求的URI
		obj := c.Request.URL.RequestURI()
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色
		sub := "user"

		//判断策略中是否存在
		aa, err := e.Enforce(sub, obj, act)
		if err != nil {
			fmt.Println("e.Enforce error")
		}
		if aa {
			fmt.Println("通过权限")
			c.Next()
		} else {
			fmt.Println("权限没有通过")
			c.Abort()
		}
	}
}

// 访问地址：http://127.0.0.1:9090/api/v1/aa
// 返回结果：{"code":200,"data":"data","message":"成功","result":"true"}
