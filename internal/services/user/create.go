package user

import (
	"context"
	"fmt"

	"wireguard-bot/internal/repository/user"
	"wireguard-bot/internal/repository/users2servers"
	"wireguard-bot/internal/utils"
)

const defaultServerID = 1

func (u *ServiceUser) Create(
	ctx context.Context,
	serverID int,
	userID int64,
	username,
	firstName,
	lastname string,
) (*user.Model, error) {
	privateKey, publicKey, err := utils.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("[user_service.create] %w", err)
	}

	if serverID == 0 {
		serverID = defaultServerID
	}

	ip, err := u.dhcp.Reserve()
	if err != nil {
		return nil, fmt.Errorf("[user_service.create] %w", err)
	}

	userModel := &user.Model{
		ID:         userID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastname,
		Role:       1,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	err = u.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		if err := u.userRepo.CreateUser(ctx, userModel); err != nil {
			return fmt.Errorf("[user_service.create] %w", err)
		}

		users2serversModel := &users2servers.Users2Servers{
			UserID:   userID,
			ServerID: serverID,
			Address:  ip.String(),
		}

		if err = u.users2serversRepo.CreateUsers2Servers(ctx, users2serversModel); err != nil {
			return fmt.Errorf("[user_service.create] %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("[user_service.create] %w", err)
	}

	return userModel, nil
}
