package repository

import (
    "context"
    "time"
)

type AuditLogEntry struct {
    UserID      string
    Action      string
    EntityType  string
    EntityID    string
    Params      string
    Timestamp   time.Time
    CorrelationID string
}

type AuditRepository interface {
    Log(ctx context.Context, entry AuditLogEntry) error
}