package entity

type LockRecord struct {
	ID           int64  `db:"id"`
	Namespace    string `db:"namespace"`
	Username     string `db:"username"`
	LockUUID     string `db:"lock_uuid"`
	LockResource string `db:"lock_resource"`
	LockTypeID   int64  `db:"lock_type_id"`
	LockType     string `db:"lock_type"`
	ExpireTime   string `db:"expire_time"`
	Comment      string `db:"comment"`
	CreateTime   string `db:"create_time"`
	UpdateTime   string `db:"update_time"`
}
