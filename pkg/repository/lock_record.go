package repository

import "github.com/liuzheran/lockInRaft/pkg/entity"

/*repository层是对表建立对象，也就是Entity的确定操作*/

type LockRecordRepository interface {
	List() ([]entity.LockRecord, error)
}

type lockRecordRepo struct {
}

// 用于创建lockRecordRepo的函数
func NewLockRecordRepository() LockRecordRepository {
	return &lockRecordRepo{}
}

func (r *lockRecordRepo) List() ([]entity.LockRecord, error) {
	// TODO 从数据库中获取lock_record表中的所有数据
	data := []entity.LockRecord{
		{LockUUID: "123", Namespace: "namespace1", Username: "username1", LockResource: "resource1", LockTypeID: 1, LockType: "type1", ExpireTime: "expire1", Comment: "comment1", CreateTime: "create1", UpdateTime: "update1"},
		{LockUUID: "456", Namespace: "namespace2", Username: "username2", LockResource: "resource2", LockTypeID: 2, LockType: "type2", ExpireTime: "expire2", Comment: "comment2", CreateTime: "create2", UpdateTime: "update2"},
		{LockUUID: "789", Namespace: "namespace3", Username: "username3", LockResource: "resource3", LockTypeID: 3, LockType: "type3", ExpireTime: "expire3", Comment: "comment3", CreateTime: "create3", UpdateTime: "update3"},
	}
	return data, nil
}
