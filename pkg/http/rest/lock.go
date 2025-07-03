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

func (l *lockApi) List(c *gin.Context) {
	lockRecords, err := l.service.ListLockRecords(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lockRecords)
}
