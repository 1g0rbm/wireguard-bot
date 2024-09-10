package user

import (
	"github.com/go-telegram/bot"

	"wireguard-bot/internal/db"
	"wireguard-bot/internal/repository"
	"wireguard-bot/internal/utils/dhcp"
)

type ServiceUser struct {
	userRepo          repository.UserRepository
	users2serversRepo repository.Users2Servers
	txManager         db.TxManager
	dhcp              *dhcp.DHCP
	outTxtMsgChan     chan<- *bot.SendMessageParams
}

func NewServiceUser(
	userRepo repository.UserRepository,
	users2serversRepo repository.Users2Servers,
	txManager db.TxManager,
	dhcp *dhcp.DHCP,
	outTxtMsgChan chan<- *bot.SendMessageParams,
) *ServiceUser {
	return &ServiceUser{
		userRepo:          userRepo,
		users2serversRepo: users2serversRepo,
		txManager:         txManager,
		dhcp:              dhcp,
		outTxtMsgChan:     outTxtMsgChan,
	}
}
