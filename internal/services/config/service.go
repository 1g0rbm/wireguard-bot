package config

import "wireguard-api/internal/repository"

type ServiceConfig struct {
	userRepo   repository.UserRepository
	serverRepo repository.ServerRepository
}

func NewConfigService(userRepo repository.UserRepository, serverRepo repository.ServerRepository) *ServiceConfig {
	return &ServiceConfig{
		userRepo:   userRepo,
		serverRepo: serverRepo,
	}
}
