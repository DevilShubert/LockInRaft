package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/liuzheran/lockInRaft/pkg/entity"
	"github.com/liuzheran/lockInRaft/pkg/repository"
)

type LockCache struct {
	// key: lock_uuid value: lock
	LocksByUUID sync.Map
	// key: lock_type_id value: lock
	LocksByTypeID sync.Map
	// key: lock_type_id value: lock_type
	LockTypesByID sync.Map
	// key: lock_resource value: lock_type
	LockTypeByResource sync.Map
}

type CacheManager struct {
	// DB查询对象
	DB *sqlx.DB
	// 众多Repo
	lockRecordRepo *repository.LockRecordRepo
	lockTypeRepo   *repository.LockTypeRepo
	// 缓存实体
	LockCache *LockCache
	// Available 主要是在重建缓存时不可用
	Available bool
}

func NewCacheManager(
	db *sqlx.DB,
	lockRecordRepo *repository.LockRecordRepo,
	lockTypeRepo *repository.LockTypeRepo) *CacheManager {
	return &CacheManager{
		DB:             db,
		lockRecordRepo: lockRecordRepo,
		lockTypeRepo:   lockTypeRepo,
		LockCache:      &LockCache{},
	}
}

func (m *CacheManager) RebuildCache(ctx context.Context) error {
	// 1.清空缓存
	// 2.重建manager中的锁信息（所有locks以及lockType）
	m.cleanCache()

	err := m.initLocks(ctx)
	if err != nil {

	}

	err = m.initLockTypes(ctx)
	if err != nil {

	}
	return nil
}

func (m *CacheManager) cleanCache() {
	m.LockCache.LockTypeByResource = sync.Map{}
	m.LockCache.LocksByTypeID = sync.Map{}
	m.LockCache.LockTypesByID = sync.Map{}
	m.LockCache.LockTypeByResource = sync.Map{}
}

func (m *CacheManager) initLocks(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	// 从数据库中查询所有lockRecord并包装到cacheManager中
	lockRecords, err := m.lockRecordRepo.List(ctxWithTimeout, m.DB)
	if err != nil {
		// ...
	}
	for _, lock := range lockRecords {
		m.lockStore(lock)
	}
	return nil
}

func (m *CacheManager) initLockTypes(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	lockTypes, err := m.lockTypeRepo.List(ctxWithTimeout, m.DB)
	if err != nil {
		// ...
	}
	for _, lock_type := range lockTypes {
		fmt.Println(lock_type)
	}
	return nil
}

func (m *CacheManager) lockStore(lock *entity.LockRecord) {
	// 按照两类进行划分：lock_uuid和lock_type_id
	m.LockCache.LocksByUUID.Store(lock.LockUUID, lock)
	// 	如果 key（这里是 lock.LockTypeID）已经存在，返回已存在的 value，并且 ok 为 true。
	// 	如果 key 不存在，则存入第二个参数（这里是 []*entity.LockEntity{lock}），并且 ok 为 false。
	if value, ok := m.LockCache.LocksByTypeID.LoadOrStore(lock.LockTypeID, []*entity.LockRecord{lock}); ok {
		data := value.([]*entity.LockRecord)
		data = append(data, lock)
		m.LockCache.LockTypesByID.Store(lock.LockTypeID, data)
	}
}
