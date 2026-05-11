package service

import (
    "context"
    "errors"
    "fmt"

    "github.com/VLGFoxRU/smartic-home/internal/domain"
    "github.com/VLGFoxRU/smartic-home/internal/repository"
)

type HomeMemberService struct {
    memberRepo repository.HomeMemberRepository
    homeRepo   repository.HomeRepository
    userRepo   repository.UserRepository
}

func NewHomeMemberService(
    memberRepo repository.HomeMemberRepository,
    homeRepo repository.HomeRepository,
    userRepo repository.UserRepository,
) *HomeMemberService {
    return &HomeMemberService{
        memberRepo: memberRepo,
        homeRepo:   homeRepo,
        userRepo:   userRepo,
    }
}

// AddMember добавляет участника в дом. Только владелец (owner) может приглашать.
func (s *HomeMemberService) AddMember(ctx context.Context, homeID, userID string, role domain.HomeRole, requesterID string) error {
    // Проверяем, что requester является владельцем дома
    home, err := s.homeRepo.FindByID(ctx, homeID)
    if err != nil {
        return fmt.Errorf("дом не найден: %w", err)
    }
    if home.OwnerID() != requesterID {
        // Также можно разрешить админам, но пока строго владелец
        return errors.New("только владелец дома может приглашать участников")
    }

    // Проверяем, существует ли пользователь
    _, err = s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return fmt.Errorf("пользователь не найден: %w", err)
    }

    // Проверяем, не состоит ли уже в доме
    if _, err := s.memberRepo.FindByHomeAndUser(ctx, homeID, userID); err == nil {
        return errors.New("пользователь уже является участником дома")
    }

    member := domain.NewHomeMember(homeID, userID, role)
    return s.memberRepo.Add(ctx, member)
}

// ChangeRole изменяет роль участника (требует права владельца)
func (s *HomeMemberService) ChangeRole(ctx context.Context, homeID, userID string, newRole domain.HomeRole, requesterID string) error {
    home, err := s.homeRepo.FindByID(ctx, homeID)
    if err != nil {
        return fmt.Errorf("дом не найден: %w", err)
    }
    if home.OwnerID() != requesterID {
        return errors.New("только владелец дома может изменять роли")
    }
    member, err := s.memberRepo.FindByHomeAndUser(ctx, homeID, userID)
    if err != nil {
        return fmt.Errorf("участник не найден: %w", err)
    }
    member.SetRole(newRole)
    return s.memberRepo.Update(ctx, member)
}

// RemoveMember удаляет участника из дома (требует права владельца)
func (s *HomeMemberService) RemoveMember(ctx context.Context, homeID, userID string, requesterID string) error {
    home, err := s.homeRepo.FindByID(ctx, homeID)
    if err != nil {
        return fmt.Errorf("дом не найден: %w", err)
    }
    if home.OwnerID() != requesterID {
        return errors.New("только владелец дома может удалять участников")
    }
    // Нельзя удалить самого владельца? Владелец не должен быть в таблице home_members, он и так владелец по полю owner_id в homes.
    // Если он добавлен как участник, можно удалить, но он останется владельцем. Пока разрешим.
    return s.memberRepo.Remove(ctx, homeID, userID)
}

// ListMembers возвращает список участников дома
func (s *HomeMemberService) ListMembers(ctx context.Context, homeID string, requesterID string) ([]domain.HomeMember, error) {
    // Проверим, что запрашивающий является участником или владельцем
    // Если владелец, то ок. Если нет, проверим членство.
    home, _ := s.homeRepo.FindByID(ctx, homeID)
    if home != nil && home.OwnerID() != requesterID {
        _, err := s.memberRepo.FindByHomeAndUser(ctx, homeID, requesterID)
        if err != nil {
            return nil, errors.New("доступ запрещён: вы не участник дома")
        }
    }
    return s.memberRepo.FindByHome(ctx, homeID)
}