package service

import (
	"context"

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
	ListLockRecords(ctx context.Context) ([]*entity.LockRecord, error)
	LockAcquire(ctx context.Context) error
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

func (s *lockService) ListLockRecords(ctx context.Context) ([]*entity.LockRecord, error) {
	return s.lockRecordRepo.List(ctx)
}

func (s *lockService) LockAcquire(ctx context.Context) error {
	// 加锁逻辑如下
	// 1. 验证当前节点是否是raft的leader 并且 Leader可用
	// 2. 查看要加的锁是否合规
	// 	2.1 检查对应的锁类型是否存在
	//  2.2 检查对应加锁内容是否合规
	//  2.3 检查lock_uuid是否已经存在
	// 3. 检查锁的互斥性
	//  3.1 是否已经存在
	//  3.2 对应锁类型是否互斥
	//  3.3 锁是否达到最大并发度
	// 4. 对CacheManager添加全局Mutex锁
	// 5. CacheManager中添加锁记录
	// 6. DB中添加锁记录
	// 7. 对CacheManager解锁全局Mutex锁
	// 8. 返回加锁成功
	return nil
}
