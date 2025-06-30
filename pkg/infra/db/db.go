package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func ExecContext(context context.Context, execer sqlx.ExecerContext, query string, args ...interface{}) (sql.Result, error) {
	return execer.ExecContext(context, query, args...)
}
