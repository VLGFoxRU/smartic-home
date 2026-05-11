package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/VLGFoxRU/smartic-home/internal/domain"
	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAnomalyRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAnomalyRepository(pool *pgxpool.Pool) repository.AnomalyRepository {
	return &PostgresAnomalyRepository{pool: pool}
}

func (r *PostgresAnomalyRepository) Create(ctx context.Context, a *domain.Anomaly) error {
	query := `INSERT INTO anomalies (id, device_id, detected_at, value, expected_value, severity, status)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query,
		a.ID(), a.DeviceID(), a.DetectedAt(), a.Value(),
		a.ExpectedValue(), a.Severity(), string(a.Status()),
	)
	return err
}

func (r *PostgresAnomalyRepository) FindByID(ctx context.Context, id string) (*domain.Anomaly, error) {
	query := `SELECT id, device_id, detected_at, value, expected_value, severity, status,
	          acknowledged_by, acknowledged_at, resolved_by, resolved_at, version
	          FROM anomalies WHERE id = $1`
	var (
		idStr, deviceID, severity, statusStr string
		detectedAt                           time.Time
		value                                float64
		expectedValue                        *float64
		ackBy, resBy                         *string
		ackAt, resAt                         *time.Time
		version                              int
	)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&idStr, &deviceID, &detectedAt, &value, &expectedValue, &severity, &statusStr,
		&ackBy, &ackAt, &resBy, &resAt, &version,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("аномалия не найдена")
		}
		return nil, err
	}
	anomaly := domain.NewAnomaly(idStr, deviceID, value, expectedValue, severity)
	anomaly.SetStatus(domain.AnomalyStatus(statusStr))
	anomaly.AcknowledgedBy = ackBy
	anomaly.AcknowledgedAt = ackAt
	anomaly.ResolvedBy = resBy
	anomaly.ResolvedAt = resAt
	anomaly.Version = version
	return anomaly, nil
}

func (r *PostgresAnomalyRepository) FindByDevice(ctx context.Context, deviceID string, statusFilter string) ([]domain.Anomaly, error) {
	query := `SELECT id, device_id, detected_at, value, expected_value, severity, status,
	          acknowledged_by, acknowledged_at, resolved_by, resolved_at, version
	          FROM anomalies WHERE device_id = $1`
	args := []interface{}{deviceID}
	if statusFilter != "" {
		query += " AND status = $2"
		args = append(args, statusFilter)
	}
	query += " ORDER BY detected_at DESC LIMIT 100"
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var anomalies []domain.Anomaly
	for rows.Next() {
		var (
			idStr, deviceID, severity, statusStr string
			detectedAt                           time.Time
			value                                float64
			expectedValue                        *float64
			ackBy, resBy                         *string
			ackAt, resAt                         *time.Time
			version                              int
		)
		if err := rows.Scan(&idStr, &deviceID, &detectedAt, &value, &expectedValue, &severity, &statusStr,
			&ackBy, &ackAt, &resBy, &resAt, &version); err != nil {
			return nil, err
		}
		a := domain.NewAnomaly(idStr, deviceID, value, expectedValue, severity)
		a.SetStatus(domain.AnomalyStatus(statusStr))
		a.AcknowledgedBy = ackBy
		a.AcknowledgedAt = ackAt
		a.ResolvedBy = resBy
		a.ResolvedAt = resAt
		a.Version = version
		anomalies = append(anomalies, *a)
	}
	return anomalies, rows.Err()
}

func (r *PostgresAnomalyRepository) Update(ctx context.Context, a *domain.Anomaly) error {
	query := `UPDATE anomalies SET status=$1, acknowledged_by=$2, acknowledged_at=$3,
	          resolved_by=$4, resolved_at=$5, version=version+1
	          WHERE id=$6 AND version=$7`
	tag, err := r.pool.Exec(ctx, query,
		string(a.Status()), a.AcknowledgedBy, a.AcknowledgedAt,
		a.ResolvedBy, a.ResolvedAt, a.ID(), a.Version,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("аномалия не найдена или версия устарела")
	}
	return nil
}