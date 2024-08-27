// Package config provides functionality for generating and managing vpn config for users.
package config

import (
	"context"
	"fmt"

	"wireguard-bot/internal/services"
	"wireguard-bot/internal/utils"
)

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
	userServer, err := s.users2serversRepo.GetFullInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate] %w", err)
	}
	if userServer == nil {
		return nil, fmt.Errorf("[config_service.generate] %w", services.ErrUserNotFound)
	}

	config := vpnConfig{
		UserPrivateKey:      userServer.UserPrivateKey,
		UserAddress:         userServer.UserAddress,
		DNS:                 "8.8.8.8",
		ServerPublicKey:     userServer.ServerPublicKey,
		ServerEndpoint:      userServer.ServerAddress,
		AllowedIPs:          "0.0.0.0/0",
		PersistentKeepalive: keepAlive,
	}

	configBytes, err := utils.Render("static/configs/user_vpn.tmp", config)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate] %w", err)
	}

	return configBytes, nil
}
