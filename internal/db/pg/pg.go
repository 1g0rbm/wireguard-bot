package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"

	"wireguard-bot/internal/db"
	"wireguard-bot/internal/db/tx"
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
		return fmt.Errorf("pg.select_context %w", err)
	}

	return nil
}

func (p *pg) GetContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	if err := p.dbc.GetContext(ctx, dest, q.QueryRaw, args...); err != nil {
		return fmt.Errorf("pg.get_context %w", err)
	}

	return nil
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (sql.Result, error) {
	t, ok := ctx.Value(tx.TxKey).(*sql.Tx)
	if ok {
		result, err := t.ExecContext(ctx, q.QueryRaw, args...)
		if err != nil {
			return nil, fmt.Errorf("pg.exec_context %w", err)
		}

		return result, nil
	}

	result, err := p.dbc.ExecContext(ctx, q.QueryRaw, args...)
	if err != nil {
		return nil, fmt.Errorf("pg.exec_context %w", err)
	}

	return result, nil
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (*sql.Rows, error) {
	t, ok := ctx.Value(tx.TxKey).(*sql.Tx)
	if ok {
		rows, err := t.QueryContext(ctx, q.QueryRaw, args...)
		if err != nil {
			return nil, fmt.Errorf("pg.query_context %w", err)
		}

		return rows, nil
	}

	rows, err := p.dbc.QueryContext(ctx, q.QueryRaw, args...)
	if err != nil {
		return nil, fmt.Errorf("pg.query_context %w", err)
	}

	return rows, nil
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	t, ok := ctx.Value(tx.TxKey).(*sql.Tx)
	if ok {
		return t.QueryRowContext(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRowContext(ctx, q.QueryRaw, args...)
}

func (p *pg) BeginTx(ctx context.Context, txOptions *sql.TxOptions) (db.Tx, error) {
	txInstance, err := p.dbc.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, fmt.Errorf("pg.begin_tx %w", err)
	}

	return txInstance, nil
}

func (p *pg) Ping(ctx context.Context) error {
	err := p.dbc.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("pg.ping %w", err)
	}

	return nil
}

func (p *pg) Close() error {
	err := p.dbc.Close()
	if err != nil {
		return fmt.Errorf("pg.close %w", err)
	}

	return nil
}
