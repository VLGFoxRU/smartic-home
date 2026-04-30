package repository

import (
	"context"
	"github.com/VLGFoxRU/smartic-home/internal/domain"
)

type AnomalyRepository interface {
	Create(ctx context.Context, anomaly *domain.Anomaly) error
	FindByID(ctx context.Context, id string) (*domain.Anomaly, error)
	FindByDevice(ctx context.Context, deviceID string, statusFilter string) ([]domain.Anomaly, error)
	Update(ctx context.Context, anomaly *domain.Anomaly) error
}