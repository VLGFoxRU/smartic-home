package domain

import (
	"errors"
	"time"
)

type AnomalyStatus string

const (
	AnomalyStatusOpen          AnomalyStatus = "open"
	AnomalyStatusAcknowledged  AnomalyStatus = "acknowledged"
	AnomalyStatusResolved      AnomalyStatus = "resolved"
	AnomalyStatusFalsePositive AnomalyStatus = "false_positive"
)

type Anomaly struct {
	id             string
	deviceID       string
	detectedAt     time.Time
	value          float64
	expectedValue  *float64
	severity       string
	status         AnomalyStatus
	AcknowledgedBy *string
	AcknowledgedAt *time.Time
	ResolvedBy     *string
	ResolvedAt     *time.Time
	Version        int
}

func NewAnomaly(id, deviceID string, value float64, expectedValue *float64, severity string) *Anomaly {
	return &Anomaly{
		id:            id,
		deviceID:      deviceID,
		detectedAt:    time.Now(),
		value:         value,
		expectedValue: expectedValue,
		severity:      severity,
		status:        AnomalyStatusOpen,
		Version:       1,
	}
}

// Геттеры
func (a *Anomaly) ID() string                   { return a.id }
func (a *Anomaly) DeviceID() string             { return a.deviceID }
func (a *Anomaly) DetectedAt() time.Time        { return a.detectedAt }
func (a *Anomaly) Value() float64              { return a.value }
func (a *Anomaly) ExpectedValue() *float64      { return a.expectedValue }
func (a *Anomaly) Severity() string             { return a.severity }
func (a *Anomaly) Status() AnomalyStatus        { return a.status }

func (a *Anomaly) SetStatus(status AnomalyStatus) {
	a.status = status
}

// Методы переходов
func (a *Anomaly) Acknowledge(userID string) error {
	if a.status != AnomalyStatusOpen {
		return errors.New("аномалия не в статусе open")
	}
	a.status = AnomalyStatusAcknowledged
	now := time.Now()
	a.AcknowledgedBy = &userID
	a.AcknowledgedAt = &now
	return nil
}

func (a *Anomaly) Resolve(userID string) error {
	if a.status == AnomalyStatusResolved || a.status == AnomalyStatusFalsePositive {
		return errors.New("аномалия уже завершена")
	}
	a.status = AnomalyStatusResolved
	now := time.Now()
	a.ResolvedBy = &userID
	a.ResolvedAt = &now
	return nil
}

func (a *Anomaly) MarkFalsePositive(userID string) error {
	if a.status == AnomalyStatusResolved || a.status == AnomalyStatusFalsePositive {
		return errors.New("аномалия уже завершена")
	}
	a.status = AnomalyStatusFalsePositive
	now := time.Now()
	a.ResolvedBy = &userID
	a.ResolvedAt = &now
	return nil
}

func (a *Anomaly) Reopen() error {
	if a.status != AnomalyStatusAcknowledged {
		return errors.New("переоткрыть можно только из acknowledged")
	}
	a.status = AnomalyStatusOpen
	a.AcknowledgedBy = nil
	a.AcknowledgedAt = nil
	return nil
}