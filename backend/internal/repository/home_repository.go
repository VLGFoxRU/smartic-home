package repository

import (
    "context"
    "github.com/VLGFoxRU/smartic-home/internal/domain"
)

type HomeRepository interface {
    Create(ctx context.Context, home *domain.Home) error
    FindByID(ctx context.Context, id string) (*domain.Home, error)
    FindByOwner(ctx context.Context, ownerID string) ([]domain.Home, error)
    Update(ctx context.Context, home *domain.Home) error
    Delete(ctx context.Context, id string) error
}