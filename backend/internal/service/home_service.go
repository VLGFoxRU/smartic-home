package service

import (
    "context"
    "fmt"
    "github.com/VLGFoxRU/smartic-home/internal/domain"
    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/google/uuid"
)

type HomeService struct {
    repo repository.HomeRepository
}

func NewHomeService(repo repository.HomeRepository) *HomeService {
    return &HomeService{repo: repo}
}

func (s *HomeService) Create(ctx context.Context, name, ownerID string) (*domain.Home, error) {
    id := uuid.New().String()
    home := domain.NewHome(id, name, ownerID)
    if err := s.repo.Create(ctx, home); err != nil {
        return nil, fmt.Errorf("ошибка создания дома: %w", err)
    }
    return home, nil
}

func (s *HomeService) GetByID(ctx context.Context, id string) (*domain.Home, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *HomeService) ListByOwner(ctx context.Context, ownerID string) ([]domain.Home, error) {
    return s.repo.FindByOwner(ctx, ownerID)
}

func (s *HomeService) Update(ctx context.Context, id, name string) error {
    home, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    home.SetName(name)
    return s.repo.Update(ctx, home)
}

func (s *HomeService) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}