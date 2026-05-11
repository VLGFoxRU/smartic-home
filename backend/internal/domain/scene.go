package domain

import (
	"encoding/json"
	"errors"
	"time"
)

type SceneStatus string

const (
	SceneStatusDraft   SceneStatus = "draft"
	SceneStatusActive  SceneStatus = "active"
	SceneStatusPaused  SceneStatus = "paused"
)

// Scene – доменный агрегат "Сценарий автоматизации"
type Scene struct {
	id          string
	homeID      string
	name        string
	condition   map[string]interface{} // JSON-условие
	action      map[string]interface{} // JSON-действие
	isActive    bool
	createdBy   string
	createdAt   time.Time
	version     int
}

func NewScene(id, homeID, name, createdBy string, condition, action map[string]interface{}) *Scene {
	return &Scene{
		id:        id,
		homeID:    homeID,
		name:      name,
		condition: condition,
		action:    action,
		createdBy: createdBy,
		createdAt: time.Now(),
		version:   1,
		// изначально не активен (draft)
	}
}

// Геттеры
func (s *Scene) ID() string                          { return s.id }
func (s *Scene) HomeID() string                      { return s.homeID }
func (s *Scene) Name() string                        { return s.name }
func (s *Scene) Condition() map[string]interface{}   { return s.condition }
func (s *Scene) Action() map[string]interface{}      { return s.action }
func (s *Scene) IsActive() bool                      { return s.isActive }
func (s *Scene) CreatedBy() string                   { return s.createdBy }
func (s *Scene) CreatedAt() time.Time                { return s.createdAt }
func (s *Scene) Version() int                        { return s.version }

// SetStatus для восстановления из БД
func (s *Scene) SetStatus(status SceneStatus) {
	switch status {
	case SceneStatusDraft, SceneStatusActive, SceneStatusPaused:
		s.isActive = (status == SceneStatusActive)
	}
}

// Активировать сценарий
func (s *Scene) Activate() error {
	if s.isActive {
		return errors.New("сценарий уже активен")
	}
	// Здесь можно добавить проверку, что condition и action валидны
	s.isActive = true
	return nil
}

// Приостановить
func (s *Scene) Pause() error {
	if !s.isActive {
		return errors.New("сценарий не активен")
	}
	s.isActive = false
	return nil
}

// Обновить имя/условия (для простоты)
func (s *Scene) Update(name string, condition, action map[string]interface{}) {
	s.name = name
	s.condition = condition
	s.action = action
}

// Для индикации статуса как перечисления
func (s *Scene) Status() SceneStatus {
	if s.isActive {
		return SceneStatusActive
	}
	return SceneStatusDraft // или paused – упростим
}

// Валидация JSON (может быть расширена)
func ValidateConditionJSON(cond map[string]interface{}) error {
	// Проверяем наличие обязательных полей
	typ, _ := cond["type"].(string)
	if typ == "" {
		return errors.New("условие должно содержать поле 'type'")
	}
	return nil
}

// Сериализация/десериализация для БД
func (s *Scene) ConditionJSON() string {
	b, _ := json.Marshal(s.condition)
	return string(b)
}

func (s *Scene) ActionJSON() string {
	b, _ := json.Marshal(s.action)
	return string(b)
}

// Восстановление из JSON (используется репозиторием)
func NewSceneFromDB(id, homeID, name, createdBy, conditionJSON, actionJSON string, isActive bool, version int, createdAt time.Time) (*Scene, error) {
	var cond, act map[string]interface{}
	if err := json.Unmarshal([]byte(conditionJSON), &cond); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(actionJSON), &act); err != nil {
		return nil, err
	}
	s := &Scene{
		id:        id,
		homeID:    homeID,
		name:      name,
		condition: cond,
		action:    act,
		createdBy: createdBy,
		createdAt: createdAt,
		version:   version,
	}
	s.SetStatus(func() SceneStatus {
		if isActive { return SceneStatusActive }
		return SceneStatusDraft
	}())
	return s, nil
}