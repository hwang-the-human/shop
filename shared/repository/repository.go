package repository

import (
	"context"
	"database/sql"
)

type Repository interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
	Close() error
}
