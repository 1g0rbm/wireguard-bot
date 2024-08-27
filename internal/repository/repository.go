package repository

import (
	"context"
	"wireguard-bot/internal/repository/users2servers"

	"wireguard-bot/internal/repository/server"
	"wireguard-bot/internal/repository/user"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*user.Model, error)
	CreateUser(ctx context.Context, user *user.Model) error
}

type ServerRepository interface {
	GetByName(ctx context.Context, name string) (*server.Server, error)
}

type Users2Servers interface {
	CreateUsers2Servers(ctx context.Context, u2s *users2servers.Users2Servers) error
	GetAllAllocatedIPsByServerAlias(ctx context.Context, alias string) ([]string, error)
	GetFullInfo(ctx context.Context, userID int64) (*users2servers.UsersServers, error)
}
