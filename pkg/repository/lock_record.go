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
	return []entity.LockRecord{}, nil
}
