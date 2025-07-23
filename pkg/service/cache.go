package service

import (
	"context"
	"sync"
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
	LockCache *LockCache
	// Available 主要是在重建缓存时不可用
	Available bool
	// TODO
	// 1、后续将众多repository添加进入CacheManager，而不是放在lockService中
	// 2、将sqlx准备好的DB对象添加进每个repo中（repo实例化时）
}

func (m *CacheManager) RebuildCache(ctx context.Context) error {
	// 1.清空缓存
	// 2.重建manager中的锁信息（所有locks以及lockType）
	m.cleanCache()
	return nil
}

func (m *CacheManager) cleanCache() {
	m.LockCache.LockTypeByResource = sync.Map{}
	m.LockCache.LocksByTypeID = sync.Map{}
	m.LockCache.LockTypesByID = sync.Map{}
	m.LockCache.LockTypeByResource = sync.Map{}
}
