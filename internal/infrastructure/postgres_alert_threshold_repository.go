package infrastructure

import (
    "context"
    "time"

    "github.com/VLGFoxRU/smartic-home/internal/domain"
    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAlertThresholdRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresAlertThresholdRepository(pool *pgxpool.Pool) repository.AlertThresholdRepository {
    return &PostgresAlertThresholdRepository{pool: pool}
}

func (r *PostgresAlertThresholdRepository) Create(ctx context.Context, t *domain.AlertThreshold) error {
    query := `INSERT INTO alert_thresholds (id, device_id, telemetry_type, min_value, max_value, severity, created_by, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
    _, err := r.pool.Exec(ctx, query,
        t.ID(), t.DeviceID(), t.TelemetryType(), t.MinValue(), t.MaxValue(), t.Severity(), t.CreatedBy(), t.CreatedAt(),
    )
    return err
}

func (r *PostgresAlertThresholdRepository) FindByID(ctx context.Context, id string) (*domain.AlertThreshold, error) {
    query := `SELECT id, device_id, telemetry_type, min_value, max_value, severity, created_by, created_at
              FROM alert_thresholds WHERE id = $1`
    var (
        idStr, devID, tType, severity, createdBy string
        minValue, maxValue *float64
        createdAt time.Time
    )
    err := r.pool.QueryRow(ctx, query, id).Scan(&idStr, &devID, &tType, &minValue, &maxValue, &severity, &createdBy, &createdAt)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return domain.NewAlertThreshold(idStr, devID, tType, severity, createdBy, minValue, maxValue), nil
}

func (r *PostgresAlertThresholdRepository) FindByDevice(ctx context.Context, deviceID string) ([]domain.AlertThreshold, error) {
    query := `SELECT id, device_id, telemetry_type, min_value, max_value, severity, created_by, created_at
              FROM alert_thresholds WHERE device_id = $1`
    rows, err := r.pool.Query(ctx, query, deviceID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var thresholds []domain.AlertThreshold
    for rows.Next() {
        var id, devID, telemetryType, severity, createdBy string
        var minValue, maxValue *float64
        var createdAt time.Time
        if err := rows.Scan(&id, &devID, &telemetryType, &minValue, &maxValue, &severity, &createdBy, &createdAt); err != nil {
            return nil, err
        }
        threshold := domain.NewAlertThreshold(id, devID, telemetryType, severity, createdBy, minValue, maxValue)
        thresholds = append(thresholds, *threshold)
    }
    return thresholds, rows.Err()
}

func (r *PostgresAlertThresholdRepository) GetByDeviceAndType(ctx context.Context, deviceID, telemetryType string) (*domain.AlertThreshold, error) {
    query := `SELECT id, device_id, telemetry_type, min_value, max_value, severity, created_by, created_at
              FROM alert_thresholds WHERE device_id = $1 AND telemetry_type = $2`
    var id, devID, tType, severity, createdBy string
    var minValue, maxValue *float64
    var createdAt time.Time
    err := r.pool.QueryRow(ctx, query, deviceID, telemetryType).Scan(&id, &devID, &tType, &minValue, &maxValue, &severity, &createdBy, &createdAt)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return domain.NewAlertThreshold(id, devID, tType, severity, createdBy, minValue, maxValue), nil
}

func (r *PostgresAlertThresholdRepository) Update(ctx context.Context, t *domain.AlertThreshold) error {
    query := `UPDATE alert_thresholds SET min_value=$1, max_value=$2, severity=$3 WHERE id=$4`
    _, err := r.pool.Exec(ctx, query, t.MinValue(), t.MaxValue(), t.Severity(), t.ID())
    return err
}

func (r *PostgresAlertThresholdRepository) Delete(ctx context.Context, id string) error {
    _, err := r.pool.Exec(ctx, `DELETE FROM alert_thresholds WHERE id = $1`, id)
    return err
}