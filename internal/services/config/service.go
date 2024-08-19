package config

import "wireguard-api/internal/repository"

// ServiceConfig is business logic to handle config command.
type ServiceConfig struct {
	userRepo   repository.UserRepository
	serverRepo repository.ServerRepository
}

// NewConfigService creates new ServiceConfig instance.
func NewConfigService(userRepo repository.UserRepository, serverRepo repository.ServerRepository) *ServiceConfig {
	return &ServiceConfig{
		userRepo:   userRepo,
		serverRepo: serverRepo,
	}
}
