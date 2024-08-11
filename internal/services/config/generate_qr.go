package config

import (
	"context"
	"fmt"
	"github.com/skip2/go-qrcode"
)

func (s *ServiceConfig) GenerateQR(ctx context.Context, userId int64) ([]byte, error) {
	conf, err := s.GenerateConf(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate_qr] %w", err)
	}

	qrCode, err := qrcode.Encode(string(conf), qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate_qr] %w", err)
	}

	return qrCode, nil
}
