package service

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	raft "github.com/hashicorp/raft"
	"github.com/liuzheran/lockInRaft/pkg/entity"
)

/*
service层的lock.go文件 主要负责“最上层的”逻辑交互，在这个文件中这些“函数”会相互调用

	加锁、解锁、查询、续期
	查询leader、触发重新选举、重建缓存、查看当前节点是否是Leader
*/

// 接口
// type LockService interface {
// 	ListLockRecords(ctx context.Context) ([]*entity.LockRecord, error)
// 	LockAcquire(ctx context.Context) error
// }

// 类（结构体）
type LockService struct {
	CacheManager *CacheManager
	RaftManager  *RaftManager
	mu           sync.Mutex
}

func NewLockService(
	cacheManager *CacheManager,
	raftManager *RaftManager) *LockService {
	return &LockService{
		CacheManager: cacheManager,
		RaftManager:  raftManager,
		mu:           sync.Mutex{},
	}
}

func (l *LockService) ListLockRecords(ctx context.Context) ([]*entity.LockRecord, error) {
	return l.CacheManager.lockRecordRepo.List(ctx, l.CacheManager.DB)
}

func (l *LockService) LockAcquire(ctx context.Context) (*entity.LockRecord, error) {
	// 加锁逻辑如下
	// 1. 验证当前节点是否是raft的leader 并且 CacheManager 可用
	// 2. 查看要加的锁是否合规
	// 	2.1 检查对应的锁类型是否存在
	//  2.2 检查对应加锁内容是否合规（对应namespace.lock_resource的格式是否正确，以及namespace是否存在）
	//  2.3 检查lock_uuid是否已经存在
	// 3. 检查锁的互斥性
	//  3.1 对应资源的锁是否已经存在（不存在则直接加锁，存在则进行下面的判断）
	//  3.2 对应锁类型是否互斥（互斥则加锁失败，不互斥则比较最大并发度）
	//  3.3 锁是否达到最大并发度（达到则加锁失败，没达到则加上）
	// 4. 对CacheManager添加全局Mutex锁（这里加锁是为了保证操作Cache和DB是原子）
	// 5. CacheManager中添加锁记录
	// 6. DB中添加锁记录
	// 7. 对CacheManager解锁全局Mutex锁
	// 8. 返回加锁成功

	if addr, _ := l.RaftManager.Raft.LeaderWithID(); addr == "" {
		return nil, errors.New("当前raft没有Leader节点")
	}

	if l.RaftManager.Raft.State() != raft.Leader {
		return nil, errors.New("当前raft不是Leader节点")
	}

	if !l.CacheManager.Available {
		return nil, errors.New("当前CacheManager正在重建缓存，不可用")
	}

	return nil, nil
}

func (l *LockService) RebuildCache(ctx context.Context) error {
	return nil
}

func (l *LockService) Elect(ctx context.Context) {
	// 这种range写法用于持续从信道 ch 读取数据，直到信道被关闭
	for leader := range l.RaftManager.Raft.LeaderCh() {
		if leader {
			atomic.StoreInt64(l.RaftManager.IsLeaderPtr(), 1)
			l.CacheManager.Available = false
			l.runLeader(ctx)
			l.CacheManager.Available = true
		} else {
			atomic.StoreInt64(l.RaftManager.IsLeaderPtr(), 0)
			l.CacheManager.Available = false
		}
	}
}

func (l *LockService) runLeader(ctx context.Context) {
	// 1.重建缓存
	err := l.CacheManager.RebuildCache(ctx)
	if err != nil {
		// TODO 清理缓存时报错
	}

	// 2.开启定时，定期清理超期的lock
	l.scheduleCleanExpiredLock(ctx)
}

func (l *LockService) scheduleCleanExpiredLock(ctx context.Context) {

}
