package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/liuzheran/lockInRaft/pkg/entity"
)

type LockTypeRepo struct {
}

func NewLockTypeRepo() *LockTypeRepo {
	return &LockTypeRepo{}
}

func (r *LockTypeRepo) List(ctx context.Context, queryer sqlx.QueryerContext) ([]*entity.LockType, error) {
	lockTypes := make([]*entity.LockType, 0)

	err := sqlx.SelectContext(ctx, queryer, &lockTypes, "SELECT * FROM lock_type;")
	if err != nil {
		fmt.Printf("query DB failed, err:%v\n", err)
		return nil, err
	}

	return lockTypes, nil
}
