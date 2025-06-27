package service

import (
	"github.com/liuzheran/lockInRaft/pkg/entity"
	"github.com/liuzheran/lockInRaft/pkg/repository"
)

/*
service层的lock.go文件 主要负责“最上层的”逻辑交互，在这个文件中这些“函数”会相互调用

	加锁、解锁、查询、续期
	查询leader、触发重新选举、重建缓存、查看当前节点是否是Leader
*/

// 接口
type LockService interface {
	ListLockRecords() ([]entity.LockRecord, error)
}

// 类（结构体）
type lockService struct {
	lockRecordRepo repository.LockRecordRepository
}

func NewLockService(repo repository.LockRecordRepository) LockService {
	return &lockService{
		lockRecordRepo: repo,
	}
}

func (s *lockService) ListLockRecords() ([]entity.LockRecord, error) {
	return s.lockRecordRepo.List()
}
