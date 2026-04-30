package service

import (
    "context"
    "fmt"
    "time"

    "github.com/VLGFoxRU/smartic-home/internal/repository"
)

type TelemetryProcessor interface {
    ProcessTelemetry(ctx context.Context, deviceID, telemetryType string, value float64, timestamp time.Time)
}

type TelemetryService struct {
    repo      repository.TelemetryRepository
    processor TelemetryProcessor
}

func NewTelemetryService(repo repository.TelemetryRepository, processor TelemetryProcessor) *TelemetryService {
    return &TelemetryService{repo: repo, processor: processor}
}

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
    if err := s.repo.Save(ctx, rec); err != nil {
        return err
    }
    if s.processor != nil {
        s.processor.ProcessTelemetry(ctx, deviceID, telemetryType, value, timestamp)
    }
    return nil
}

func (s *TelemetryService) GetHistory(ctx context.Context, deviceID string, from, to time.Time) ([]repository.TelemetryRecord, error) {
    if deviceID == "" {
        return nil, fmt.Errorf("deviceID обязателен")
    }
    return s.repo.GetHistory(ctx, deviceID, from, to, 1000)
}