package repository

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/liuzheran/lockInRaft/pkg/entity"
)

/*repository层是对表建立对象，也就是Entity的确定操作*/

// type LockRecordRepository interface {
// 	List(ctx context.Context) ([]*entity.LockRecord, error)
// }

type LockRecordRepo struct {
}

// 用于创建lockRecordRepo的函数
func NewLockRecordRepository() *LockRecordRepo {
	return &LockRecordRepo{}
}

func (r *LockRecordRepo) List(ctx context.Context, queryer sqlx.QueryerContext) ([]*entity.LockRecord, error) {
	// 从数据库中获取lock_record表中的所有数据

	locks := make([]*entity.LockRecord, 0)

	// 测试代码：单独的 goroutine 监听 context 取消
	// go func() {
	// 	<-ctx.Done()
	// 	fmt.Println("[goroutine] ctx is done, reason:", ctx.Err())
	// }()

	// 使用sqlx.SelectContext 查询数据库，这里的ctx是gin.Context的上下文，如果用户手动取消请求，则ctx会取消，从而取消数据库查询
	err := sqlx.SelectContext(ctx, queryer, &locks, "SELECT * FROM lock_record;")

	if err != nil {
		fmt.Printf("query DB failed, err:%v\n", err)
		return nil, err
	}

	return locks, nil
}

// 后续还有 Create Delete Update
