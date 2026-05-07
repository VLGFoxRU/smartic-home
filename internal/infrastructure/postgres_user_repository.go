package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/VLGFoxRU/smartic-home/internal/domain"
	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, username, email, password_hash, role, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query,
		user.ID(), user.Username(), user.Email(),
		user.PasswordHash(), string(user.Role()), user.CreatedAt(),
	)
	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at FROM users WHERE id = $1`
	var (
		idStr, username, email, passwordHash, roleStr string
		createdAt                                     time.Time
	)
	err := r.pool.QueryRow(ctx, query, id).Scan(&idStr, &username, &email, &passwordHash, &roleStr, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, err
	}
	user := domain.NewUser(idStr, username, email, passwordHash, domain.UserRole(roleStr))
	return user, nil
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at FROM users WHERE username = $1`
	var (
		idStr, uname, email, passwordHash, roleStr string
		createdAt                                 time.Time
	)
	err := r.pool.QueryRow(ctx, query, username).Scan(&idStr, &uname, &email, &passwordHash, &roleStr, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, err
	}
	user := domain.NewUser(idStr, uname, email, passwordHash, domain.UserRole(roleStr))
	return user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at FROM users WHERE email = $1`
	var (
		idStr, uname, userEmail, passwordHash, roleStr string
		createdAt                                     time.Time
	)
	err := r.pool.QueryRow(ctx, query, email).Scan(&idStr, &uname, &userEmail, &passwordHash, &roleStr, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, err
	}
	user := domain.NewUser(idStr, uname, userEmail, passwordHash, domain.UserRole(roleStr))
	return user, nil
}