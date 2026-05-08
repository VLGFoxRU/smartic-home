package domain

import "time"

type AlertThreshold struct {
    id             string
    deviceID       string
    telemetryType  string
    minValue       *float64
    maxValue       *float64
    severity       string
    createdBy      string
    createdAt      time.Time
}

func NewAlertThreshold(id, deviceID, telemetryType, severity, createdBy string, minValue, maxValue *float64) *AlertThreshold {
    return &AlertThreshold{
        id:            id,
        deviceID:      deviceID,
        telemetryType: telemetryType,
        minValue:      minValue,
        maxValue:      maxValue,
        severity:      severity,
        createdBy:     createdBy,
        createdAt:     time.Now(),
    }
}

func (a *AlertThreshold) ID() string            { return a.id }
func (a *AlertThreshold) DeviceID() string      { return a.deviceID }
func (a *AlertThreshold) TelemetryType() string { return a.telemetryType }
func (a *AlertThreshold) MinValue() *float64    { return a.minValue }
func (a *AlertThreshold) MaxValue() *float64    { return a.maxValue }
func (a *AlertThreshold) Severity() string      { return a.severity }
func (a *AlertThreshold) CreatedBy() string     { return a.createdBy }
func (a *AlertThreshold) CreatedAt() time.Time  { return a.createdAt }

func (a *AlertThreshold) SetMin(val *float64) { a.minValue = val }
func (a *AlertThreshold) SetMax(val *float64) { a.maxValue = val }
func (a *AlertThreshold) SetSeverity(s string) { a.severity = s }

func (a *AlertThreshold) IsExceeded(value float64) bool {
    if a.minValue != nil && value < *a.minValue {
        return true
    }
    if a.maxValue != nil && value > *a.maxValue {
        return true
    }
    return false
}