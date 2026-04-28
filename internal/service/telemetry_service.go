package service

import (
	"context"
	"fmt"
	"time"

	"github.com/VLGFoxRU/smartic-home/internal/repository"
)

type TelemetryService struct {
	repo repository.TelemetryRepository
}

func NewTelemetryService(repo repository.TelemetryRepository) *TelemetryService {
	return &TelemetryService{repo: repo}
}

// Ingest сохраняет одно показание
func (s *TelemetryService) Ingest(ctx context.Context, deviceID, telemetryType string, value float64, timestamp time.Time) error {
	if deviceID == "" || telemetryType == "" {
		return fmt.Errorf("deviceID и type обязательны")
	}
	rec := repository.TelemetryRecord{
		DeviceID:  deviceID,
		Type:      telemetryType,
		Value:     value,
		Timestamp: timestamp,
	}
	return s.repo.Save(ctx, rec)
}

// GetHistory возвращает последние записи за период (по умолчанию 1000)
func (s *TelemetryService) GetHistory(ctx context.Context, deviceID string, from, to time.Time) ([]repository.TelemetryRecord, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("deviceID обязателен")
	}
	return s.repo.GetHistory(ctx, deviceID, from, to, 1000)
}