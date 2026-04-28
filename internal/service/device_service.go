package service

import (
	"context"
	"fmt"

	"github.com/VLGFoxRU/smartic-home/internal/domain"
	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/google/uuid"
)

type DeviceService struct {
	repo  repository.DeviceRepository
	cache repository.CacheRepository
}

func NewDeviceService(repo repository.DeviceRepository, cache repository.CacheRepository) *DeviceService {
	return &DeviceService{repo: repo, cache: cache}
}

func (s *DeviceService) ListDevices(ctx context.Context) ([]domain.Device, error) {
	devices, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("DeviceService.ListDevices: %w", err)
	}
	return devices, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, id string) (*domain.Device, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *DeviceService) CreateDevice(ctx context.Context, name, devType, roomID string) (*domain.Device, error) {
	id := uuid.New().String()
	device := domain.NewDeviceWithRoom(id, name, devType, roomID)
	if err := s.repo.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("ошибка создания устройства: %w", err)
	}
	return device, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, id, name, devType string) error {
	dev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	dev.SetName(name)
	dev.SetType(devType)
	return s.repo.Save(ctx, dev)
}

func (s *DeviceService) DeleteDevice(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *DeviceService) EnableDevice(ctx context.Context, id string) error {
	dev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := dev.Enable(); err != nil {
		return err
	}
	return s.repo.Save(ctx, dev)
}

func (s *DeviceService) DisableDevice(ctx context.Context, id string) error {
	dev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := dev.Disable(); err != nil {
		return err
	}
	return s.repo.Save(ctx, dev)
}