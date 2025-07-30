package entity

type LockType struct {
	ID                   int64  `db:"id"`
	Namespace            string `db:"namespace"`
	LockTypeName         string `db:lock_type_name`
	Mutex                int32  `db:mutex`
	MaxConcurrency       int64  `db:max_concurrency`
	ActualMaxConcurrency int64  `db:actual_max_concurrency`
	CreateTime           string `db:"create_time"`
	UpdateTime           string `db:"update_time"`
}
