package user

import (
	"wireguard-bot/internal/db"
	"wireguard-bot/internal/repository"
	"wireguard-bot/internal/utils/dhcp"
)

type ServiceUser struct {
	userRepo          repository.UserRepository
	users2serversRepo repository.Users2Servers
	txManager         db.TxManager
	dhcp              *dhcp.DHCP
}

func NewServiceUser(
	userRepo repository.UserRepository,
	users2serversRepo repository.Users2Servers,
	txManager db.TxManager,
	dhcp *dhcp.DHCP,
) *ServiceUser {
	return &ServiceUser{
		userRepo:          userRepo,
		users2serversRepo: users2serversRepo,
		txManager:         txManager,
		dhcp:              dhcp,
	}
}
