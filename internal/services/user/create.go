package user

import (
	"context"
	"fmt"

	"wireguard-api/internal/repository/user"
	"wireguard-api/internal/repository/users2servers"
	"wireguard-api/internal/utils"
)

const defaultServerId = 1

func (u *ServiceUser) Create(
	ctx context.Context,
	serverId int,
	userId int64,
	username,
	firstName,
	lastname string,
) error {
	privateKey, publicKey, err := utils.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("[user_service.create] %w", err)
	}

	if serverId == 0 {
		serverId = defaultServerId
	}

	err = u.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		userModel := &user.Model{
			Id:         userId,
			Username:   username,
			FirstName:  firstName,
			LastName:   lastname,
			Role:       1,
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		}

		if err := u.userRepo.CreateUser(ctx, userModel); err != nil {
			return fmt.Errorf("[user_service.create] %w", err)
		}

		users2serversModel := &users2servers.Users2Servers{
			UserId:   userId,
			ServerId: serverId,
			Address:  "10.0.0.2/24",
		}

		if err = u.users2serversRepo.CreateUsers2Servers(ctx, users2serversModel); err != nil {
			return fmt.Errorf("[user_service.create] %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("[user_service.create] %w", err)
	}

	return nil
}
