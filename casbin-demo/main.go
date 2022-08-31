package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//e, err := casbin.NewEnforcer("path/to/model.conf", "path/to/policy.csv")

	// 使用MySQL数据库初始化一个Xorm适配器
	a, err := xormadapter.NewAdapter("mysql", "root:password#dbr@tcp(10.100.156.210:3306)/")
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}

	m, err := model.NewModelFromString(`
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
	`)
	if err != nil {
		log.Fatalf("error: model: %s", err)
	}

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
	}

	sub := "alice" // the user that wants to access a resource.
	obj := "data1" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.

	ok, err := e.Enforce(sub, obj, act)

	if err != nil {
		// handle err
		fmt.Println(err)
	}

	if ok {
		// permit alice to read data1
		fmt.Println("yes")
	} else {
		// deny the request, show an error
		fmt.Println("no")
	}
	//roles, err := e.GetRolesForUser("alice")
	//fmt.Println(roles)
	// You could use BatchEnforce() to enforce some requests in batches.
	// This method returns a bool slice, and this slice's index corresponds to the row index of the two-dimensional array.
	// e.g. results[0] is the result of {"alice", "data1", "read"}
	//results, err := e.BatchEnforce([][]interface{}{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"jack", "data3", "read"}})

}
