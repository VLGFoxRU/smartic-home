package infrastructure

import (
    "context"
    "fmt"

    "github.com/VLGFoxRU/smartic-home/internal/domain"
    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHomeMemberRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresHomeMemberRepository(pool *pgxpool.Pool) repository.HomeMemberRepository {
    return &PostgresHomeMemberRepository{pool: pool}
}

func (r *PostgresHomeMemberRepository) Add(ctx context.Context, member *domain.HomeMember) error {
    query := `INSERT INTO home_members (home_id, user_id, role) VALUES ($1, $2, $3)`
    _, err := r.pool.Exec(ctx, query, member.HomeID(), member.UserID(), string(member.Role()))
    return err
}

func (r *PostgresHomeMemberRepository) FindByHomeAndUser(ctx context.Context, homeID, userID string) (*domain.HomeMember, error) {
    query := `SELECT home_id, user_id, role FROM home_members WHERE home_id = $1 AND user_id = $2`
    var hID, uID, roleStr string
    err := r.pool.QueryRow(ctx, query, homeID, userID).Scan(&hID, &uID, &roleStr)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("участник не найден")
        }
        return nil, err
    }
    return domain.NewHomeMember(hID, uID, domain.HomeRole(roleStr)), nil
}

func (r *PostgresHomeMemberRepository) FindByHome(ctx context.Context, homeID string) ([]domain.HomeMember, error) {
    query := `SELECT home_id, user_id, role FROM home_members WHERE home_id = $1`
    rows, err := r.pool.Query(ctx, query, homeID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var members []domain.HomeMember
    for rows.Next() {
        var hID, uID, roleStr string
        if err := rows.Scan(&hID, &uID, &roleStr); err != nil {
            return nil, err
        }
        members = append(members, *domain.NewHomeMember(hID, uID, domain.HomeRole(roleStr)))
    }
    return members, rows.Err()
}

func (r *PostgresHomeMemberRepository) Update(ctx context.Context, member *domain.HomeMember) error {
    query := `UPDATE home_members SET role = $1 WHERE home_id = $2 AND user_id = $3`
    _, err := r.pool.Exec(ctx, query, string(member.Role()), member.HomeID(), member.UserID())
    return err
}

func (r *PostgresHomeMemberRepository) Remove(ctx context.Context, homeID, userID string) error {
    _, err := r.pool.Exec(ctx, `DELETE FROM home_members WHERE home_id = $1 AND user_id = $2`, homeID, userID)
    return err
}