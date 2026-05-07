package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/VLGFoxRU/smartic-home/internal/domain"
	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, username, email, password string, role domain.UserRole) (*domain.User, error) {
	// Проверка уникальности
	if _, err := s.repo.FindByUsername(ctx, username); err == nil {
		return nil, errors.New("пользователь с таким именем уже существует")
	}
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	id := uuid.New().String()
	user := domain.NewUser(id, username, email, string(hash), role)
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}
	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, username, password string) (*domain.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("неверное имя пользователя или пароль")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(password)); err != nil {
		return nil, errors.New("неверное имя пользователя или пароль")
	}
	return user, nil
}