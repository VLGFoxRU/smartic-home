package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/VLGFoxRU/smartic-home/internal/domain"
	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/google/uuid"
)

type DeviceService struct {
	repo  repository.DeviceRepository
	cache repository.CacheRepository
	audit repository.AuditRepository
}

func NewDeviceService(repo repository.DeviceRepository, cache repository.CacheRepository, audit repository.AuditRepository) *DeviceService {
	return &DeviceService{
		repo:  repo,
		cache: cache,
		audit: audit,
	}
}

// ListDevices возвращает список всех устройств
func (s *DeviceService) ListDevices(ctx context.Context) ([]domain.Device, error) {
	devices, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("DeviceService.ListDevices: %w", err)
	}
	return devices, nil
}

// SendCommand выполняет команду управления устройством
func (s *DeviceService) SendCommand(ctx context.Context, deviceID string, command string, params map[string]interface{}) error {
	// 1. Получаем устройство
	dev, err := s.repo.FindByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("устройство не найдено: %w", err)
	}

	// 2. Проверяем, может ли принять команду
	if !dev.CanAcceptCommand() {
		return fmt.Errorf("устройство %s не может принять команду (статус: %s)", dev.ID(), dev.Status())
	}

	// 3. Эмулируем выполнение (позже через брокер)
	// Пока просто обновляем состояние устройства в кэше и БД
	// и пишем аудит.

	// Генерируем корреляционный ID
	corrID := uuid.New().String()

	// Записываем в аудит
	paramsJSON, _ := json.Marshal(params) // упрощённо
	entry := repository.AuditLogEntry{
		UserID:        "a0000000-0000-0000-0000-000000000001", // пока заглушка
		Action:        "command",
		EntityType:    "device",
		EntityID:      deviceID,
		Params:        string(paramsJSON),
		Timestamp:     time.Now(),
		CorrelationID: corrID,
	}
	if err := s.audit.Log(ctx, entry); err != nil {
		return fmt.Errorf("ошибка записи аудита: %w", err)
	}

	// Обновляем состояние в кэше
	newStatus := map[string]interface{}{"status": "executed", "command": command}
	statusJSON, _ := json.Marshal(newStatus)
	_ = s.cache.Set(ctx, "device:"+deviceID+":state", string(statusJSON), 5*time.Minute)

	// Обновление в БД (заглушка, позже реализуем)
	return nil
}