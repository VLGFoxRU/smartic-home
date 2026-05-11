package service

import (
    "context"
    "fmt"

    "github.com/VLGFoxRU/smartic-home/internal/domain"
    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/google/uuid"
)

type AlertThresholdService struct {
    repo repository.AlertThresholdRepository
}

func NewAlertThresholdService(repo repository.AlertThresholdRepository) *AlertThresholdService {
    return &AlertThresholdService{repo: repo}
}

func (s *AlertThresholdService) Create(ctx context.Context, deviceID, telemetryType, severity, createdBy string, minValue, maxValue *float64) (*domain.AlertThreshold, error) {
    id := uuid.New().String()
    threshold := domain.NewAlertThreshold(id, deviceID, telemetryType, severity, createdBy, minValue, maxValue)
    if err := s.repo.Create(ctx, threshold); err != nil {
        return nil, fmt.Errorf("не удалось создать порог: %w", err)
    }
    return threshold, nil
}

func (s *AlertThresholdService) GetByDeviceAndType(ctx context.Context, deviceID, telemetryType string) (*domain.AlertThreshold, error) {
    return s.repo.GetByDeviceAndType(ctx, deviceID, telemetryType)
}

func (s *AlertThresholdService) ListByDevice(ctx context.Context, deviceID string) ([]domain.AlertThreshold, error) {
    return s.repo.FindByDevice(ctx, deviceID)
}

func (s *AlertThresholdService) Update(ctx context.Context, id string, minValue, maxValue *float64, severity string) error {
    threshold, err := s.repo.FindByID(ctx, id)
    if err != nil || threshold == nil {
        return fmt.Errorf("порог не найден")
    }
    if minValue != nil {
        threshold.SetMin(minValue)
    }
    if maxValue != nil {
        threshold.SetMax(maxValue)
    }
    if severity != "" {
        threshold.SetSeverity(severity)
    }
    return s.repo.Update(ctx, threshold)
}

func (s *AlertThresholdService) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}