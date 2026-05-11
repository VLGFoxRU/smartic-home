package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/google/uuid"
)

type ControlService struct {
	deviceRepo repository.DeviceRepository
	cache      repository.CacheRepository
	audit      repository.AuditRepository
	broker     repository.MessageBroker
	eventPub   repository.EventPublisher
}

func NewControlService(
	deviceRepo repository.DeviceRepository,
	cache repository.CacheRepository,
	audit repository.AuditRepository,
	broker repository.MessageBroker,
	eventPub repository.EventPublisher,
) *ControlService {
	return &ControlService{
		deviceRepo: deviceRepo,
		cache:      cache,
		audit:      audit,
		broker:     broker,
		eventPub:   eventPub,
	}
}

// SendCommand выполняет проверки, пишет аудит, публикует команду в брокер
func (s *ControlService) SendCommand(ctx context.Context, deviceID, command string, params map[string]interface{}, userID, role string) error {
	// 1. Получить устройство
	dev, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("устройство не найдено: %w", err)
	}

	// 2. Проверка прав (controller/admin/owner)
	if role != "owner" && role != "admin" && role != "controller" {
		return fmt.Errorf("недостаточно прав для управления устройством")
	}

	// 3. Проверка статуса устройства
	if !dev.CanAcceptCommand() {
		return fmt.Errorf("устройство %s не может принять команду (статус: %s)", dev.ID(), dev.Status())
	}

	// 4. Валидация команды
	if err := dev.ValidateCommand(command, params); err != nil {
		return fmt.Errorf("невалидная команда: %w", err)
	}

	// 5. Генерация correlationId
	corrID := uuid.New().String()

	// 6. Запись в аудит
	paramsJSON, _ := json.Marshal(params)
	auditEntry := repository.AuditLogEntry{
		UserID:        userID,
		Action:        "command",
		EntityType:    "device",
		EntityID:      deviceID,
		Params:        string(paramsJSON),
		Timestamp:     time.Now(),
		CorrelationID: corrID,
	}
	if err := s.audit.Log(ctx, auditEntry); err != nil {
		return fmt.Errorf("ошибка записи аудита: %w", err)
	}

	// 7. Публикация в брокер (заглушка)
	if s.eventPub != nil {
		eventPayload := map[string]interface{}{
			"deviceId": deviceID,
			"command":  command,
			"params":   params,
		}
		_ = s.eventPub.Publish(ctx, "device.command_sent", eventPayload)
	}

	// Формируем новое состояние на основе команды
	newState := map[string]interface{}{}
	switch command {
	case "set_power":
		on, _ := params["on"].(bool)
		newState["power"] = map[bool]string{true: "on", false: "off"}[on]
	case "set_brightness":
		// сохраняем яркость, если применимо
		if brightness, ok := params["brightness"]; ok {
			newState["brightness"] = brightness
		}
	default:
		newState["lastCommand"] = command
		newState["params"] = params
	}
	// Добавляем всегда статус pending для отслеживания
	newState["status"] = "pending"
	stateJSON, _ := json.Marshal(newState)
	_ = s.cache.Set(ctx, "device:"+deviceID+":state", string(stateJSON), 5*time.Minute)

	// Публикация события через WebSocket
	

	return nil
}