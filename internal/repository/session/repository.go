package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"wireguard-bot/internal/db"
)

const table = "sessions"

const (
	colPk        = "id"
	colUserID    = "user_id"
	colExpiredAt = "expired_at"
	colCreatedAt = "created_at"
	colUpdatedAt = "updated_at"
)

type Repository struct {
	db db.Client
}

func NewRepository(db db.Client) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, session *Session) error {
	q, args, err := squirrel.
		Insert(table).
		PlaceholderFormat(squirrel.Dollar).
		Columns(
			colPk,
			colUserID,
			colExpiredAt,
			colCreatedAt,
		).
		Values(
			session.ID,
			session.UserID,
			session.ExpiredAt,
			time.Now(),
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("session_repository.create_session: %w", err)
	}

	query := db.Query{
		Name:     "session_repository.create_session",
		QueryRaw: q,
	}

	row, err := r.db.DB().QueryContext(ctx, query, args...)
	defer func() {
		if clsErr := row.Close(); clsErr != nil {
			err = errors.Join(err, clsErr)
		}
	}()

	if err != nil {
		return fmt.Errorf("session_repository.create_session: %w", err)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, session *Session) error {
	q, args, err := squirrel.
		Update(table).
		PlaceholderFormat(squirrel.Dollar).
		Set(colExpiredAt, session.ExpiredAt).
		Set(colUpdatedAt, time.Now()).
		Where(squirrel.Eq{colPk: session.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("session_repository.update_session: %w", err)
	}

	query := db.Query{
		Name:     "session_repository.update_session",
		QueryRaw: q,
	}

	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("session_repository.update_session: %w", err)
	}

	return nil
}

func (r *Repository) FindByID(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	return r.FindBy(ctx, squirrel.Eq{colPk: sessionID})
}

func (r *Repository) FindByUserID(ctx context.Context, userID int64) (*Session, error) {
	return r.FindBy(ctx, squirrel.Eq{colUserID: userID})
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*Session, error) {
	q, args, err := squirrel.
		Select(table+"."+colPk, colUserID, table+"."+colExpiredAt, table+"."+colCreatedAt, table+"."+colUpdatedAt).
		PlaceholderFormat(squirrel.Dollar).
		From(table).
		Join("users ON sessions.user_id = users.id").
		Where(squirrel.Eq{"username": username}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("session_repository.find_one_by_username: %w", err)
	}

	query := db.Query{
		Name:     "session_repository.find_one_by_username",
		QueryRaw: q,
	}

	var model Session
	if err := r.db.DB().GetContext(ctx, &model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//nolint:nilnil
			return nil, nil
		}
		return nil, fmt.Errorf("session_repository.find_one_by_username: %w", err)
	}

	return &model, nil
}

func (r *Repository) FindBy(ctx context.Context, eq squirrel.Eq) (*Session, error) {
	q, args, err := squirrel.
		Select(colPk, colUserID, colExpiredAt, colCreatedAt, colUpdatedAt).
		PlaceholderFormat(squirrel.Dollar).
		From(table).
		Where(eq).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("session_repository.find_one_by: %w", err)
	}

	query := db.Query{
		Name:     "session_repository.find_one_by",
		QueryRaw: q,
	}

	var model Session
	if err := r.db.DB().GetContext(ctx, &model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//nolint:nilnil
			return nil, nil
		}
		return nil, fmt.Errorf("session_repository.find_one_by: %w", err)
	}

	return &model, nil
}
