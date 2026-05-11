package repository

import (
	"context"
	"github.com/VLGFoxRU/smartic-home/internal/domain"
)

type SceneRepository interface {
	Create(ctx context.Context, scene *domain.Scene) error
	FindByID(ctx context.Context, id string) (*domain.Scene, error)
	FindByHome(ctx context.Context, homeID string) ([]domain.Scene, error)
	Update(ctx context.Context, scene *domain.Scene) error
	Delete(ctx context.Context, id string) error
}