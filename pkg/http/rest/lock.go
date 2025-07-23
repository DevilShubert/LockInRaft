package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
	"github.com/liuzheran/lockInRaft/pkg/schema"
	"github.com/liuzheran/lockInRaft/pkg/service"
)

type lockApi struct {
	lockService *service.LockService
}

func NewLockApi(lockService *service.LockService) *lockApi {
	return &lockApi{lockService: lockService}
}

// 查询
func (l *lockApi) List(c *gin.Context) {
	lockRecords, err := l.lockService.ListLockRecords(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lockRecords)
}

// 加锁
func (l *lockApi) Acquire(c *gin.Context) {
	// 包装入参
	_, err := l.lockService.LockAcquire(c.Request.Context())
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

// 添加节点
func (l *lockApi) AddNode(c *gin.Context) {
	var param schema.NodeParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("添加节点: %s, %s\n", param.Id, param.PeerAddr)
	err := l.lockService.RaftManager.AddVoter(raft.ServerID(param.Id), raft.ServerAddress(param.PeerAddr), 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Node added success"})
}

// 删除节点
func (l *lockApi) RemoveNode(c *gin.Context) {
	var param schema.NodeParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("删除节点: %s, %s\n", param.Id, param.PeerAddr)
	err := l.lockService.RaftManager.RemoveServer(raft.ServerID(param.Id), 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Node removed success"})
}

// 获取集群信息
func (l *lockApi) GetClusterInfo(c *gin.Context) {
	nodes, err := l.lockService.RaftManager.GetRaftClusterInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

// 获取Leader
func (l *lockApi) GetLeader(c *gin.Context) {
	leader, err := l.lockService.RaftManager.GetRaftLeader()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leader)
}

// 重建缓存
func (l *lockApi) RebuildCache(c *gin.Context) {

}
