package repository

import (
	"context"
	"time"
)

type TelemetryRecord struct {
	DeviceID  string
	Timestamp time.Time
	Value     float64
	Type      string // temperature, humidity, power, energy, water_flow
}

type TelemetryRepository interface {
	Save(ctx context.Context, record TelemetryRecord) error
	GetHistory(ctx context.Context, deviceID string, from, to time.Time, limit int) ([]TelemetryRecord, error)
}