package repository

import (
    "context"
    "github.com/VLGFoxRU/smartic-home/internal/domain"
)

type AlertThresholdRepository interface {
    Create(ctx context.Context, threshold *domain.AlertThreshold) error
    FindByID(ctx context.Context, id string) (*domain.AlertThreshold, error)
    FindByDevice(ctx context.Context, deviceID string) ([]domain.AlertThreshold, error)
    GetByDeviceAndType(ctx context.Context, deviceID, telemetryType string) (*domain.AlertThreshold, error)
    Update(ctx context.Context, threshold *domain.AlertThreshold) error
    Delete(ctx context.Context, id string) error
}