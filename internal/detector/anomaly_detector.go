package detector

import (
    "context"
    "log"
    "time"

    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/VLGFoxRU/smartic-home/internal/service"
)

type AnomalyDetector struct {
    telemetryRepo repository.TelemetryRepository
    anomalySvc    *service.AnomalyService
    thresholdSvc  *service.AlertThresholdService
}

func NewAnomalyDetector(
    telemetryRepo repository.TelemetryRepository,
    anomalySvc *service.AnomalyService,
    thresholdSvc *service.AlertThresholdService,
) *AnomalyDetector {
    return &AnomalyDetector{
        telemetryRepo: telemetryRepo,
        anomalySvc:    anomalySvc,
        thresholdSvc:  thresholdSvc,
    }
}

func (d *AnomalyDetector) ProcessTelemetry(ctx context.Context, deviceID, telemetryType string, value float64, timestamp time.Time) {
    // 1. Проверяем пороги
    threshold, err := d.thresholdSvc.GetByDeviceAndType(ctx, deviceID, telemetryType)
    if err == nil && threshold != nil && threshold.IsExceeded(value) {
        _, err := d.anomalySvc.CreateAnomaly(ctx, deviceID, value, nil, threshold.Severity())
        if err != nil {
            log.Printf("AnomalyDetector: ошибка создания аномалии по порогу: %v", err)
        }
        return
    }

    // 2. Статистический детектор (среднее ± 3σ)
    stats, err := d.telemetryRepo.GetStats(ctx, deviceID, telemetryType, 7)
    if err != nil {
        log.Printf("AnomalyDetector: ошибка получения статистики: %v", err)
        return
    }
    if stats.Count < 10 || stats.StdDev == 0 {
        return
    }
    upper := stats.Mean + 3*stats.StdDev
    lower := stats.Mean - 3*stats.StdDev
    if value > upper || value < lower {
        severity := "warning"
        if value > stats.Mean+5*stats.StdDev || value < stats.Mean-5*stats.StdDev {
            severity = "critical"
        }
        _, err := d.anomalySvc.CreateAnomaly(ctx, deviceID, value, &stats.Mean, severity)
        if err != nil {
            log.Printf("AnomalyDetector: ошибка создания аномалии: %v", err)
        }
    }
}