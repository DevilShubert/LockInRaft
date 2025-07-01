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
	List(ctx context.Context) ([]*entity.LockRecord, error)
}

type lockRecordRepo struct {
}

// 用于创建lockRecordRepo的函数
func NewLockRecordRepository() LockRecordRepository {
	return &lockRecordRepo{}
}

func (r *lockRecordRepo) List(ctx context.Context) ([]*entity.LockRecord, error) {
	// TODO 从数据库中获取lock_record表中的所有数据
	dsn := "root:123456xx@tcp(127.0.0.1:3306)/lock_in_raft?charset=utf8mb4&parseTime=True"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return nil, err
	}

	locks := make([]*entity.LockRecord, 0)

	// 测试代码：单独的 goroutine 监听 context 取消
	go func() {
		<-ctx.Done()
		fmt.Println("[goroutine] ctx is done, reason:", ctx.Err())
	}()

	// 使用sqlx.SelectContext 查询数据库，这里的ctx是gin.Context的上下文，如果用户手动取消请求，则ctx会取消，从而取消数据库查询
	err = sqlx.SelectContext(ctx, db, &locks, "SELECT * FROM lock_record;")

	if err != nil {
		fmt.Printf("query DB failed, err:%v\n", err)
		return nil, err
	}

	return locks, nil
}
