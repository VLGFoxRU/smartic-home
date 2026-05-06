package service

import (
    "context"
    "fmt"

    "github.com/VLGFoxRU/smartic-home/internal/domain"
    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/google/uuid"
)

type AnomalyService struct {
    repo     repository.AnomalyRepository
    eventPub repository.EventPublisher
}

func NewAnomalyService(repo repository.AnomalyRepository, eventPub repository.EventPublisher) *AnomalyService {
    return &AnomalyService{repo: repo, eventPub: eventPub}
}

func (s *AnomalyService) CreateAnomaly(ctx context.Context, deviceID string, value float64, expectedValue *float64, severity string) (*domain.Anomaly, error) {
    id := uuid.New().String()
    anomaly := domain.NewAnomaly(id, deviceID, value, expectedValue, severity)
    if err := s.repo.Create(ctx, anomaly); err != nil {
        return nil, fmt.Errorf("не удалось создать аномалию: %w", err)
    }
    // Публикуем событие
    if s.eventPub != nil {
        _ = s.eventPub.Publish(ctx, "anomaly.detected", map[string]interface{}{
            "id":          anomaly.ID(),
            "device_id":   deviceID,
            "value":       value,
            "expected":    expectedValue,
            "severity":    severity,
            "timestamp":   anomaly.DetectedAt(),
        })
    }
    return anomaly, nil
}

// остальные методы (GetAnomaly, ListAnomalies, UpdateStatus) без изменений
func (s *AnomalyService) GetAnomaly(ctx context.Context, id string) (*domain.Anomaly, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *AnomalyService) ListAnomalies(ctx context.Context, deviceID, statusFilter string) ([]domain.Anomaly, error) {
    return s.repo.FindByDevice(ctx, deviceID, statusFilter)
}

func (s *AnomalyService) UpdateStatus(ctx context.Context, id, newStatus, userID string) error {
	anomaly, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	switch domain.AnomalyStatus(newStatus) {
	case domain.AnomalyStatusAcknowledged:
		if err := anomaly.Acknowledge(userID); err != nil {
			return err
		}
	case domain.AnomalyStatusResolved:
		if err := anomaly.Resolve(userID); err != nil {
			return err
		}
	case domain.AnomalyStatusFalsePositive:
		if err := anomaly.MarkFalsePositive(userID); err != nil {
			return err
		}
	case domain.AnomalyStatusOpen:
		if err := anomaly.Reopen(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("недопустимый статус: %s", newStatus)
	}
	return s.repo.Update(ctx, anomaly)
}