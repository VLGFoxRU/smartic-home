package repository

import (
    "context"
    "github.com/VLGFoxRU/smartic-home/internal/domain"
)

type HomeMemberRepository interface {
    Add(ctx context.Context, member *domain.HomeMember) error
    FindByHomeAndUser(ctx context.Context, homeID, userID string) (*domain.HomeMember, error)
    FindByHome(ctx context.Context, homeID string) ([]domain.HomeMember, error)
    Update(ctx context.Context, member *domain.HomeMember) error
    Remove(ctx context.Context, homeID, userID string) error
}