package service

import (
	"context"
	"fmt"

	"github.com/VLGFoxRU/smartic-home/internal/domain"
	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/google/uuid"
)

type SceneService struct {
	repo repository.SceneRepository
}

func NewSceneService(repo repository.SceneRepository) *SceneService {
	return &SceneService{repo: repo}
}

func (s *SceneService) CreateScene(ctx context.Context, name, homeID, createdBy string, condition, action map[string]interface{}) (*domain.Scene, error) {
	// Валидация условия
	if err := domain.ValidateConditionJSON(condition); err != nil {
		return nil, fmt.Errorf("некорректное условие: %w", err)
	}
	id := uuid.New().String()
	scene := domain.NewScene(id, homeID, name, createdBy, condition, action)
	if err := s.repo.Create(ctx, scene); err != nil {
		return nil, fmt.Errorf("ошибка создания сценария: %w", err)
	}
	return scene, nil
}

func (s *SceneService) GetScene(ctx context.Context, id string) (*domain.Scene, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *SceneService) ListScenes(ctx context.Context, homeID string) ([]domain.Scene, error) {
	return s.repo.FindByHome(ctx, homeID)
}

func (s *SceneService) UpdateScene(ctx context.Context, id, name string, condition, action map[string]interface{}) error {
	scene, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	scene.Update(name, condition, action)
	return s.repo.Update(ctx, scene)
}

func (s *SceneService) Activate(ctx context.Context, id string) error {
	scene, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := scene.Activate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, scene)
}

func (s *SceneService) Pause(ctx context.Context, id string) error {
	scene, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := scene.Pause(); err != nil {
		return err
	}
	return s.repo.Update(ctx, scene)
}

func (s *SceneService) DeleteScene(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}