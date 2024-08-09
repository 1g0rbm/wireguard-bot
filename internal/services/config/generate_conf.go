package config

import (
	"context"
	"fmt"
	"wireguard-api/internal/services"
	"wireguard-api/internal/utils"
)

const defaultServerName = "test"

type vpnConfig struct {
	UserPrivateKey      string
	UserAddress         string
	DNS                 string
	ServerPublicKey     string
	ServerEndpoint      string
	AllowedIPs          string
	PersistentKeepalive int
}

func (s *ServiceConfig) GenerateConf(ctx context.Context, userId int64) ([]byte, error) {
	userModel, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate] %w", err)
	}
	if userModel == nil {
		return nil, fmt.Errorf("[config_service.generate] %w", services.UserNotFound)
	}

	serverModel, err := s.serverRepo.GetByName(ctx, defaultServerName)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate] %w", err)
	}

	config := vpnConfig{
		UserPrivateKey:      userModel.PrivateKey,
		UserAddress:         "10.0.0.2/24",
		DNS:                 "8.8.8.8",
		ServerPublicKey:     serverModel.PublicKey,
		ServerEndpoint:      serverModel.Address,
		AllowedIPs:          "0.0.0.0/0",
		PersistentKeepalive: 25,
	}

	configBytes, err := utils.Render("static/user_vpn_config.tmp", config)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate] %w", err)
	}

	return configBytes, nil
}
