package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liuzheran/lockInRaft/pkg/service"
)

type lockApi struct {
	service service.LockService
}

func NewLockApi(service service.LockService) *lockApi {
	return &lockApi{service: service}
}

// 查询
func (l *lockApi) List(c *gin.Context) {
	lockRecords, err := l.service.ListLockRecords(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lockRecords)
}

// 加锁
func (l *lockApi) Acquire(c *gin.Context) {
	err := l.service.LockAcquire(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lock acquired"})
}

// 释放锁
func (l *lockApi) Release(c *gin.Context) {}

// 续期
func (l *lockApi) KeepAlive(c *gin.Context) {}
