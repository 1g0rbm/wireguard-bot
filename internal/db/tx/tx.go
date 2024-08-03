package tx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"wireguard-api/internal/db"
)

type ctxKey string

const TxKey ctxKey = "tx"

type Manager struct {
	transactor db.Transactor
}

func NewTxManager(transactor db.Transactor) *Manager {
	return &Manager{
		transactor: transactor,
	}
}

func (m *Manager) ReadCommited(ctx context.Context, fn db.TxHandler) error {
	return m.transaction(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted}, fn)
}

func (m *Manager) transaction(ctx context.Context, opts *sql.TxOptions, fn db.TxHandler) error {
	tx, ok := ctx.Value(TxKey).(db.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err := m.transactor.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("transaction panic recover: %v", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				err = errors.Join(err, fmt.Errorf("transaction rollback err: %w", errRollback))
			}

			return
		}

		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("transaction commit err: %w", err)
		}
	}()

	ctx = WithTransactionCtx(ctx, tx)
	if err = fn(ctx); err != nil {
		err = fmt.Errorf("failed execute transaction: %w", err)
	}

	return err
}

func WithTransactionCtx(ctx context.Context, tx db.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
