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
    eventPub      repository.EventPublisher
}

func NewAnomalyDetector(
    telemetryRepo repository.TelemetryRepository,
    anomalySvc *service.AnomalyService,
    eventPub repository.EventPublisher,
) *AnomalyDetector {
    return &AnomalyDetector{
        telemetryRepo: telemetryRepo,
        anomalySvc:    anomalySvc,
        eventPub:      eventPub,
    }
}

func (d *AnomalyDetector) ProcessTelemetry(ctx context.Context, deviceID, telemetryType string, value float64, timestamp time.Time) {
    // Получаем статистику за 7 дней
    stats, err := d.telemetryRepo.GetStats(ctx, deviceID, telemetryType, 7)
    if err != nil {
        log.Printf("AnomalyDetector: ошибка получения статистики: %v", err)
        return
    }
    if stats.Count < 10 { // нужно хотя бы 10 записей для надёжной статистики
        return
    }
    if stats.StdDev == 0 {
        return // недостаточная вариативность
    }
    threshold := 3.0
    upper := stats.Mean + threshold*stats.StdDev
    lower := stats.Mean - threshold*stats.StdDev

    if value > upper || value < lower {
        // Определяем severity
        severity := "warning"
        if value > stats.Mean+5*stats.StdDev || value < stats.Mean-5*stats.StdDev {
            severity = "critical"
        }
        anomaly, err := d.anomalySvc.CreateAnomaly(ctx, deviceID, value, &stats.Mean, severity)
        if err != nil {
            log.Printf("AnomalyDetector: ошибка создания аномалии: %v", err)
            return
        }
        log.Printf("AnomalyDetector: обнаружена аномалия %s (значение %.2f, ожидалось %.2f±%.2f)", anomaly.ID(), value, stats.Mean, stats.StdDev)

        // Публикуем событие
        _ = d.eventPub.Publish(ctx, "anomaly.detected", map[string]interface{}{
            "id":          anomaly.ID(),
            "device_id":   deviceID,
            "value":       value,
            "expected":    stats.Mean,
            "severity":    severity,
            "timestamp":   timestamp,
        })
    }
}