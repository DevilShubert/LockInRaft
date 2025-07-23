package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/liuzheran/lockInRaft/pkg/setting"
)

var lockDbInitOnce sync.Once
var lockDb *sqlx.DB

func ExecContext(context context.Context, execer sqlx.ExecerContext, query string, args ...interface{}) (sql.Result, error) {
	return execer.ExecContext(context, query, args...)
}

func ProvideLockDB(dbConfig *setting.DBConfig) (*sqlx.DB, error) {
	var err error
	lockDbInitOnce.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName)
		fmt.Println("Connecting to DB with DSN:", dsn)
		lockDb, err = sqlx.Connect("mysql", dsn)
	})
	return lockDb, err
}
