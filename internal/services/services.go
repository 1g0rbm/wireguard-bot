package services

import "context"

type ConfigService interface {
	GenerateConf(ctx context.Context, userId int64) ([]byte, error)
	GenerateQR(ctx context.Context, userId int64) ([]byte, error)
}

type UserService interface {
	Create(ctx context.Context, serverId int, userId int64, username, firstName, lastname string) error
	IsUserExist(ctx context.Context, userId int64) (bool, error)
}
