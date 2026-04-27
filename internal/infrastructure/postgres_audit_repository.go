package infrastructure

import (
    "context"
    "fmt"

    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAuditRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresAuditRepository(pool *pgxpool.Pool) repository.AuditRepository {
    return &PostgresAuditRepository{pool: pool}
}

func (r *PostgresAuditRepository) Log(ctx context.Context, entry repository.AuditLogEntry) error {
    query := `INSERT INTO audit_log (user_id, action, entity_type, entity_id, params, timestamp, correlation_id)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
    _, err := r.pool.Exec(ctx, query,
        entry.UserID, entry.Action, entry.EntityType, entry.EntityID,
        entry.Params, entry.Timestamp, entry.CorrelationID,
    )
    if err != nil {
        return fmt.Errorf("PostgresAuditRepository.Log: %w", err)
    }
    return nil
}