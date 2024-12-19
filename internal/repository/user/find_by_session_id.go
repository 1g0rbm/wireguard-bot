package user

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

func (r *Repository) FIndUserByActiveSessionID(ctx context.Context, id uuid.UUID) (*Model, error) {
	q, args, err := squirrel.Select(
		"users."+colPk,
		"users."+colUsername,
		"users."+colFirstname,
		"users."+colLastname,
		"users."+colRole,
		"users."+colState,
		"users."+colPublicKey,
		"users."+colPrivateKey,
		"users."+colCreatedAt,
		"users."+colUpdatedAt,
	).
		PlaceholderFormat(squirrel.Dollar).
		From(table).
		Join("sessions ON sessions.user_id = users.id").
		Where(squirrel.Eq{"sessions.id": id}).
		Where(squirrel.Gt{"sessions.expired_at": time.Now()}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("user_repository.find_user_by_session_id: %w", err)
	}

	query := db.Query{
		Name:     "user_repository.FIndUserBySessionID",
		QueryRaw: q,
	}

	var model Model
	if err = r.db.DB().GetContext(ctx, &model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//nolint:nilnil
			return nil, nil
		}
		return nil, fmt.Errorf("user_repository.find_user_by_session_id: %w", err)
	}

	return &model, nil
}
