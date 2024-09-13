package repository

import (
	"context"

	"github.com/google/uuid"

	"wireguard-bot/internal/repository/server"
	"wireguard-bot/internal/repository/session"
	"wireguard-bot/internal/repository/user"
	"wireguard-bot/internal/repository/users2servers"
)

type SessionRepository interface {
	Create(ctx context.Context, session *session.Session) error
	Update(ctx context.Context, session *session.Session) error
	FindByID(ctx context.Context, sessionID uuid.UUID) (*session.Session, error)
	FindByUserID(ctx context.Context, userID int64) (*session.Session, error)
	FindByUsername(ctx context.Context, username string) (*session.Session, error)
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*user.Model, error)
	GetUserByUsername(ctx context.Context, username string) (*user.Model, error)
	CreateUser(ctx context.Context, user *user.Model) error
	List(ctx context.Context) ([]user.Model, error)
}

type ServerRepository interface {
	GetByName(ctx context.Context, name string) (*server.Server, error)
}

type Users2Servers interface {
	CreateUsers2Servers(ctx context.Context, u2s *users2servers.Users2Servers) error
	GetAllAllocatedIPsByServerAlias(ctx context.Context, alias string) ([]string, error)
	GetFullInfo(ctx context.Context, userID int64) (*users2servers.UsersServers, error)
}
