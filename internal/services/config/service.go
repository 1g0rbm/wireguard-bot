package config

import "wireguard-api/internal/repository"

// ServiceConfig is business logic to handle config command.
type ServiceConfig struct {
	users2serversRepo repository.Users2Servers
}

// NewConfigService creates new ServiceConfig instance.
func NewConfigService(users2serversRepo repository.Users2Servers) *ServiceConfig {
	return &ServiceConfig{
		users2serversRepo: users2serversRepo,
	}
}
