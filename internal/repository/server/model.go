package server

import (
	"database/sql"
	"time"
)

type Server struct {
	ID         int          `db:"id"`
	Name       string       `db:"name"`
	Address    string       `db:"address"`
	PublicKey  string       `db:"public_key"`
	PrivateKey string       `db:"private_key"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
}
