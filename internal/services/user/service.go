package user

import (
	"wireguard-api/internal/db"
	"wireguard-api/internal/repository"
)

type ServiceUser struct {
	userRepo          repository.UserRepository
	users2serversRepo repository.Users2Servers
	txManager         db.TxManager
}

func NewServiceUser(
	userRepo repository.UserRepository,
	users2serversRepo repository.Users2Servers,
	txManager db.TxManager,
) *ServiceUser {
	return &ServiceUser{
		userRepo:          userRepo,
		users2serversRepo: users2serversRepo,
		txManager:         txManager,
	}
}
