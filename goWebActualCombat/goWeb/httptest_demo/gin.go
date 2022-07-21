package httptestdemo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//在web开发场景下的单元测试，若涉及到HTTP请求，推荐使用Go标准库net/http/httptest进行测试

//param请求参数
type Param struct {
	Name string `json:"name"`
}

//helloHandler /hello请求处理函数

func helloHandler(c *gin.Context) {
	var p Param
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "we need a name",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("hello %s", p.Name),
	})
}

//SetupRouter路由
func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/hello", helloHandler)
	return router
}
