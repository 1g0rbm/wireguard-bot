package session

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID    `db:"id"`
	UserID    int64        `db:"user_id"`
	ExpiredAt time.Time    `db:"expired_at"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
