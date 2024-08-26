package user

import (
	"database/sql"
	"time"
)

const (
	DisabledState = "disabled"
	EnabledState  = "enabled"
)

const (
	UserRole = iota + 1
	AdminRole
)

type Model struct {
	ID         int64        `db:"id"`
	Username   string       `db:"username"`
	FirstName  string       `db:"first_name"`
	LastName   string       `db:"last_name"`
	Role       int8         `db:"role"`
	PublicKey  string       `db:"public_key"`
	PrivateKey string       `db:"private_key"`
	State      string       `db:"state"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
}

func (m Model) RoleIs(role int8) bool {
	return m.Role == role
}

func (m Model) StateIs(state string) bool {
	return m.State == state
}
