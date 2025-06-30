package repository

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/liuzheran/lockInRaft/pkg/entity"
)

/*repository层是对表建立对象，也就是Entity的确定操作*/

type LockRecordRepository interface {
	List() ([]*entity.LockRecord, error)
}

type lockRecordRepo struct {
}

// 用于创建lockRecordRepo的函数
func NewLockRecordRepository() LockRecordRepository {
	return &lockRecordRepo{}
}

func (r *lockRecordRepo) List() ([]*entity.LockRecord, error) {
	// TODO 从数据库中获取lock_record表中的所有数据
	dsn := "root:123456xx@tcp(127.0.0.1:3306)/lock_in_raft?charset=utf8mb4&parseTime=True"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return nil, err
	}

	locks := make([]*entity.LockRecord, 0)
	err = sqlx.SelectContext(context.Background(), db, &locks, "SELECT * FROM lock_record")

	if err != nil {
		fmt.Printf("query DB failed, err:%v\n", err)
		return nil, err
	}

	return locks, nil
}
