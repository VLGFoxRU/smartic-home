package domain

import (
	"errors"
	"time"
)

// DeviceStatus – перечисление статусов (как в спецификации)
type DeviceStatus string

const (
	DeviceStatusOnline   DeviceStatus = "online"
	DeviceStatusOffline  DeviceStatus = "offline"
	DeviceStatusDisabled DeviceStatus = "disabled"
)

// Device – доменная сущность "Устройство"
type Device struct {
	id      string       // внутренний идентификатор (UUID)
	name    string
	devType string       // тип устройства (light, temperature_sensor ...)
	status  DeviceStatus
	lastSeen *time.Time  // nullable время последней телеметрии
	version int          // оптимистичная блокировка
}

// NewDevice – конструктор, создаёт устройство с начальным состоянием
func NewDevice(id, name, devType string) *Device {
	return &Device{
		id:      id,
		name:    name,
		devType: devType,
		status:  DeviceStatusOffline, // по умолчанию офлайн
		version: 1,
	}
}

// ID, Name, Type, Status – геттеры (доступ на чтение)
func (d *Device) ID() string          { return d.id }
func (d *Device) Name() string        { return d.name }
func (d *Device) Type() string        { return d.devType }
func (d *Device) Status() DeviceStatus { return d.status }
func (d *Device) Version() int        { return d.version }

// Enable – включает устройство, если оно было disabled
func (d *Device) Enable() error {
	if d.status != DeviceStatusDisabled {
		return errors.New("устройство не отключено, включение невозможно")
	}
	d.status = DeviceStatusOnline
	return nil
}

// Disable – выключает устройство (перевод в disabled)
func (d *Device) Disable() error {
	if d.status == DeviceStatusDisabled {
		return errors.New("устройство уже отключено")
	}
	d.status = DeviceStatusDisabled
	return nil
}

// MarkOnline – вызывается при получении heartbeat, переводит в online
func (d *Device) MarkOnline(timestamp time.Time) error {
	if d.status == DeviceStatusDisabled {
		return errors.New("нельзя перевести disabled устройство в online")
	}
	d.status = DeviceStatusOnline
	d.lastSeen = &timestamp
	return nil
}

// MarkOffline – переводит в offline (при таймауте)
func (d *Device) MarkOffline() error {
	if d.status == DeviceStatusDisabled {
		return errors.New("устройство disabled, игнорируем offine")
	}
	d.status = DeviceStatusOffline
	return nil
}

// CanAcceptCommand – проверяет, можно ли отправить команду (инвариант из документации)
func (d *Device) CanAcceptCommand() bool {
	return d.status == DeviceStatusOnline
}

// SetStatus – для восстановления из БД (ТОЛЬКО для репозитория)
// В реальном проекте можно сделать экспортируемым методом Reconstruct.
func (d *Device) SetStatus(status DeviceStatus) error {
	switch status {
	case DeviceStatusOnline, DeviceStatusOffline, DeviceStatusDisabled:
		d.status = status
		return nil
	default:
		return errors.New("недопустимый статус")
	}
}