package infrastructure

import (
	"context"
	"fmt"
	"time"
	"database/sql"

	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTelemetryRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTelemetryRepository(pool *pgxpool.Pool) repository.TelemetryRepository {
	return &PostgresTelemetryRepository{pool: pool}
}

func (r *PostgresTelemetryRepository) Save(ctx context.Context, record repository.TelemetryRecord) error {
	query := `INSERT INTO telemetry (device_id, time, value, type) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, record.DeviceID, record.Timestamp, record.Value, record.Type)
	if err != nil {
		return fmt.Errorf("PostgresTelemetryRepository.Save: %w", err)
	}
	return nil
}

func (r *PostgresTelemetryRepository) GetHistory(ctx context.Context, deviceID string, from, to time.Time, limit int) ([]repository.TelemetryRecord, error) {
	query := `SELECT device_id, time, value, type FROM telemetry
	          WHERE device_id = $1 AND time BETWEEN $2 AND $3
	          ORDER BY time DESC
	          LIMIT $4`
	rows, err := r.pool.Query(ctx, query, deviceID, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("PostgresTelemetryRepository.GetHistory: %w", err)
	}
	defer rows.Close()

	var records []repository.TelemetryRecord
	for rows.Next() {
		var rec repository.TelemetryRecord
		if err := rows.Scan(&rec.DeviceID, &rec.Timestamp, &rec.Value, &rec.Type); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

func (r *PostgresTelemetryRepository) GetStats(ctx context.Context, deviceID, telemetryType string, days int) (*repository.StatsResult, error) {
    query := `SELECT COUNT(*), AVG(value), STDDEV_SAMP(value)
              FROM telemetry
              WHERE device_id = $1 AND type = $2 AND time >= NOW() - ($3 || ' days')::INTERVAL`
    var count int
    var mean, stddev sql.NullFloat64
    err := r.pool.QueryRow(ctx, query, deviceID, telemetryType, fmt.Sprintf("%d", days)).Scan(&count, &mean, &stddev)
    if err != nil {
        return nil, fmt.Errorf("PostgresTelemetryRepository.GetStats: %w", err)
    }
    result := &repository.StatsResult{Count: count}
    if mean.Valid {
        result.Mean = mean.Float64
    }
    if stddev.Valid {
        result.StdDev = stddev.Float64
    }
    return result, nil
}