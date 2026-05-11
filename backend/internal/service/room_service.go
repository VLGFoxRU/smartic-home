package service

import (
    "context"
    "fmt"
    "github.com/VLGFoxRU/smartic-home/internal/domain"
    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/google/uuid"
)

type RoomService struct {
    repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) *RoomService {
    return &RoomService{repo: repo}
}

func (s *RoomService) Create(ctx context.Context, homeID, name string) (*domain.Room, error) {
    id := uuid.New().String()
    room := domain.NewRoom(id, homeID, name)
    if err := s.repo.Create(ctx, room); err != nil {
        return nil, fmt.Errorf("ошибка создания комнаты: %w", err)
    }
    return room, nil
}

func (s *RoomService) GetByID(ctx context.Context, id string) (*domain.Room, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *RoomService) ListByHome(ctx context.Context, homeID string) ([]domain.Room, error) {
    return s.repo.FindByHome(ctx, homeID)
}

func (s *RoomService) Update(ctx context.Context, id, name string) error {
    room, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    room.SetName(name)
    return s.repo.Update(ctx, room)
}

func (s *RoomService) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}