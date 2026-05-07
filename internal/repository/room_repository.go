package repository

import (
    "context"
    "github.com/VLGFoxRU/smartic-home/internal/domain"
)

type RoomRepository interface {
    Create(ctx context.Context, room *domain.Room) error
    FindByID(ctx context.Context, id string) (*domain.Room, error)
    FindByHome(ctx context.Context, homeID string) ([]domain.Room, error)
    Update(ctx context.Context, room *domain.Room) error
    Delete(ctx context.Context, id string) error
}