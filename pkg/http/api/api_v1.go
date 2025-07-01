package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/liuzheran/lockInRaft/pkg/entity"
	"github.com/liuzheran/lockInRaft/pkg/service"
)

// List 因为在 gin 源码中，HandlerFunc 的定义明确要求处理函数的签名必须是 func(c *gin.Context)
func List(c *gin.Context, service service.LockService) {
	// TODO 访问数据库查询数据
	fmt.Println("执行List API")

	lockRecords, err := service.ListLockRecords(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
	}
	data := []*entity.LockRecord{}
	for _, lockRecord := range lockRecords {
		//  TODO 打印对象
		data = append(data, lockRecord)
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   data,
	})
}
