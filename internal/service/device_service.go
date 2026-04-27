package service

import (
	"context"
	"fmt"

	"github.com/VLGFoxRU/smartic-home/internal/domain"
	"github.com/VLGFoxRU/smartic-home/internal/repository"
)

type DeviceService struct {
	repo  repository.DeviceRepository
	cache repository.CacheRepository
}

func NewDeviceService(repo repository.DeviceRepository, cache repository.CacheRepository) *DeviceService {
	return &DeviceService{
		repo:  repo,
		cache: cache,
	}
}

func (s *DeviceService) ListDevices(ctx context.Context) ([]domain.Device, error) {
	devices, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("DeviceService.ListDevices: %w", err)
	}
	return devices, nil
}