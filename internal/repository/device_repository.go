package repository

import (
	"context"
	"github.com/VLGFoxRU/smartic-home/internal/domain"
)

type DeviceRepository interface {
	FindAll(ctx context.Context) ([]domain.Device, error)
	FindByID(ctx context.Context, id string) (*domain.Device, error)
	Save(ctx context.Context, device *domain.Device) error
}