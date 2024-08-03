package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
)

type Client interface {
	DB() DB
	Close() error
}

type DB interface {
	SQLExecutor
	Pinger
	Transactor
	Close() error
}

type TxManager interface {
	ReadCommited(ctx context.Context, f TxHandler) error
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions *sql.TxOptions) (Tx, error)
}

type Tx interface {
	Commit() error
	Rollback() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type SQLExecutor interface {
	NamedExecutor
	QueryExecutor
}

type NamedExecutor interface {
	SelectContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type QueryExecutor interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Query struct {
	Name     string
	QueryRaw string
}

type TxHandler func(ctx context.Context) error
