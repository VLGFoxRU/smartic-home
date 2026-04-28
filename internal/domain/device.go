package domain

import (
	"errors"
	"fmt"
	"time"
)

type DeviceStatus string

const (
	DeviceStatusOnline   DeviceStatus = "online"
	DeviceStatusOffline  DeviceStatus = "offline"
	DeviceStatusDisabled DeviceStatus = "disabled"
)

type Device struct {
	id       string
	name     string
	roomID   string
	devType  string
	status   DeviceStatus
	lastSeen *time.Time
	version  int
}

func NewDevice(id, name, devType string) *Device {
	return &Device{
		id:      id,
		name:    name,
		devType: devType,
		status:  DeviceStatusOffline,
		version: 1,
	}
}

func NewDeviceWithRoom(id, name, devType, roomID string) *Device {
	return &Device{
		id:      id,
		name:    name,
		roomID:  roomID,
		devType: devType,
		status:  DeviceStatusOffline,
		version: 1,
	}
}

// Геттеры
func (d *Device) ID() string            { return d.id }
func (d *Device) Name() string          { return d.name }
func (d *Device) Type() string          { return d.devType }
func (d *Device) Status() DeviceStatus  { return d.status }
func (d *Device) RoomID() string        { return d.roomID }
func (d *Device) LastSeen() *time.Time  { return d.lastSeen }
func (d *Device) Version() int          { return d.version }

// Сеттеры для полей, безопасных к изменению
func (d *Device) SetName(name string)    { d.name = name }
func (d *Device) SetType(t string)       { d.devType = t }
func (d *Device) SetRoomID(rid string)   { d.roomID = rid }
func (d *Device) SetLastSeen(t time.Time) { d.lastSeen = &t }
func (d *Device) SetVersion(v int)       { d.version = v }

// Восстановление статуса (используется репозиторием)
func (d *Device) SetStatus(status DeviceStatus) error {
	switch status {
	case DeviceStatusOnline, DeviceStatusOffline, DeviceStatusDisabled:
		d.status = status
		return nil
	default:
		return errors.New("недопустимый статус")
	}
}

// Включение устройства
func (d *Device) Enable() error {
	if d.status != DeviceStatusDisabled {
		return errors.New("устройство не отключено, включение невозможно")
	}
	d.status = DeviceStatusOnline
	return nil
}

// Отключение устройства
func (d *Device) Disable() error {
	if d.status == DeviceStatusDisabled {
		return errors.New("устройство уже отключено")
	}
	d.status = DeviceStatusDisabled
	return nil
}

// Пометить как онлайн (heartbeat)
func (d *Device) MarkOnline(timestamp time.Time) error {
	if d.status == DeviceStatusDisabled {
		return errors.New("нельзя перевести disabled устройство в online")
	}
	d.status = DeviceStatusOnline
	d.lastSeen = &timestamp
	return nil
}

// Пометить как офлайн
func (d *Device) MarkOffline() error {
	if d.status == DeviceStatusDisabled {
		return errors.New("устройство disabled, игнорируем offine")
	}
	d.status = DeviceStatusOffline
	return nil
}

// Может ли принять команду
func (d *Device) CanAcceptCommand() bool {
	return d.status == DeviceStatusOnline
}

// Валидация команды и параметров
type CommandDescriptor struct {
	Name       string
	ParamsDesc map[string]ParamConstraint
}

type ParamConstraint struct {
	Type     string // "bool", "int", "float", "string"
	Required bool
	Min      *float64
	Max      *float64
}

func (d *Device) SupportedCommands() map[string]CommandDescriptor {
	switch d.devType {
	case "light":
		return map[string]CommandDescriptor{
			"set_power": {Name: "set_power", ParamsDesc: map[string]ParamConstraint{
				"on": {Type: "bool", Required: true},
			}},
			"set_brightness": {Name: "set_brightness", ParamsDesc: map[string]ParamConstraint{
				"brightness": {Type: "int", Required: true, Min: ptr(0.0), Max: ptr(100.0)},
			}},
		}
	case "power_switch":
		return map[string]CommandDescriptor{
			"set_power": {Name: "set_power", ParamsDesc: map[string]ParamConstraint{
				"on": {Type: "bool", Required: true},
			}},
		}
	default:
		return map[string]CommandDescriptor{}
	}
}

func (d *Device) ValidateCommand(command string, params map[string]interface{}) error {
	cmds := d.SupportedCommands()
	desc, ok := cmds[command]
	if !ok {
		return fmt.Errorf("команда '%s' не поддерживается устройством типа %s", command, d.devType)
	}

	for paramName, constraint := range desc.ParamsDesc {
		val, exists := params[paramName]
		if constraint.Required && !exists {
			return fmt.Errorf("обязательный параметр '%s' отсутствует", paramName)
		}
		if !exists {
			continue
		}
		switch constraint.Type {
		case "bool":
			if _, ok := val.(bool); !ok {
				return fmt.Errorf("параметр '%s' должен быть boolean", paramName)
			}
		case "int":
			v, ok := val.(float64)
			if !ok {
				return fmt.Errorf("параметр '%s' должен быть числом", paramName)
			}
			if v != float64(int(v)) {
				return fmt.Errorf("параметр '%s' должен быть целым числом", paramName)
			}
			if constraint.Min != nil && v < *constraint.Min {
				return fmt.Errorf("параметр '%s' меньше минимума (%.0f)", paramName, *constraint.Min)
			}
			if constraint.Max != nil && v > *constraint.Max {
				return fmt.Errorf("параметр '%s' превышает максимум (%.0f)", paramName, *constraint.Max)
			}
		case "float":
			v, ok := val.(float64)
			if !ok {
				return fmt.Errorf("параметр '%s' должен быть числом", paramName)
			}
			if constraint.Min != nil && v < *constraint.Min {
				return fmt.Errorf("параметр '%s' меньше минимума (%f)", paramName, *constraint.Min)
			}
			if constraint.Max != nil && v > *constraint.Max {
				return fmt.Errorf("параметр '%s' превышает максимум (%f)", paramName, *constraint.Max)
			}
		}
	}
	return nil
}

func ptr(f float64) *float64 { return &f }