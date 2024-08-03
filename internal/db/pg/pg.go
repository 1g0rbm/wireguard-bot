package pg

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"

	"wireguard-api/internal/db"
	"wireguard-api/internal/db/tx"
)

type pg struct {
	dbc *sqlx.DB
}

func NewDB(dbc *sqlx.DB) db.DB {
	return &pg{
		dbc: dbc,
	}
}

func (p *pg) SelectContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	if err := p.dbc.SelectContext(ctx, dest, q.QueryRaw, args...); err != nil {
		return err
	}

	return nil
}

func (p *pg) GetContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	if err := p.dbc.GetContext(ctx, dest, q.QueryRaw, args...); err != nil {
		return err
	}

	return nil
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (sql.Result, error) {
	t, ok := ctx.Value(tx.TxKey).(sql.Tx)
	if ok {
		return t.ExecContext(ctx, q.QueryRaw, args...)
	}

	return p.dbc.ExecContext(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (*sql.Rows, error) {
	t, ok := ctx.Value(tx.TxKey).(sql.Tx)
	if ok {
		return t.QueryContext(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryContext(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	t, ok := ctx.Value(tx.TxKey).(sql.Tx)
	if ok {
		return t.QueryRowContext(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRowContext(ctx, q.QueryRaw, args...)
}

func (p *pg) BeginTx(ctx context.Context, txOptions *sql.TxOptions) (db.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.PingContext(ctx)
}

func (p *pg) Close() error {
	return p.dbc.Close()
}
