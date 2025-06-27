package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/liuzheran/lockInRaft/pkg/service"
)

// List 因为在 gin 源码中，HandlerFunc 的定义明确要求处理函数的签名必须是 func(c *gin.Context)
func List(c *gin.Context, service service.LockService) {
	// TODO 访问数据库查询数据
	fmt.Println("执行List API")
	data := []string{"item1", "item2", "item3"} // 示例数据
	// TODO 访问service层的 LockService.List() 方法
	c.JSON(200, gin.H{
		"status": "success",
		"data":   data,
	})
}
