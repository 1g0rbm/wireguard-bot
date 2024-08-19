// Package config provides functionality for generating and managing vpn config for users.
package config

import (
	"context"
	"fmt"

	"wireguard-api/internal/services"
	"wireguard-api/internal/utils"
)

const defaultServerName = "astana_1"

const (
	keepAlive = 25
)

type vpnConfig struct {
	UserPrivateKey      string
	UserAddress         string
	DNS                 string
	ServerPublicKey     string
	ServerEndpoint      string
	AllowedIPs          string
	PersistentKeepalive int
}

// GenerateConf method to generate vpn config for userId.
func (s *ServiceConfig) GenerateConf(ctx context.Context, userID int64) ([]byte, error) {
	userModel, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate] %w", err)
	}
	if userModel == nil {
		return nil, fmt.Errorf("[config_service.generate] %w", services.ErrUserNotFound)
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
		PersistentKeepalive: keepAlive,
	}

	configBytes, err := utils.Render("static/user_vpn_config.tmp", config)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate] %w", err)
	}

	return configBytes, nil
}
