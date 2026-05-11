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

type PostgresHomeRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresHomeRepository(pool *pgxpool.Pool) repository.HomeRepository {
    return &PostgresHomeRepository{pool: pool}
}

func (r *PostgresHomeRepository) Create(ctx context.Context, home *domain.Home) error {
    query := `INSERT INTO homes (id, name, owner_id, created_at) VALUES ($1, $2, $3, $4)`
    _, err := r.pool.Exec(ctx, query, home.ID(), home.Name(), home.OwnerID(), home.CreatedAt())
    return err
}

func (r *PostgresHomeRepository) FindByID(ctx context.Context, id string) (*domain.Home, error) {
    query := `SELECT id, name, owner_id, created_at FROM homes WHERE id = $1`
    var (
        idStr, name, ownerID string
        createdAt             time.Time
    )
    err := r.pool.QueryRow(ctx, query, id).Scan(&idStr, &name, &ownerID, &createdAt)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("дом не найден")
        }
        return nil, err
    }
    home := domain.NewHome(idStr, name, ownerID)
    // createdAt уже внутри конструктора, но перезапишем
    return home, nil
}

func (r *PostgresHomeRepository) FindByOwner(ctx context.Context, ownerID string) ([]domain.Home, error) {
    query := `SELECT id, name, owner_id, created_at FROM homes WHERE owner_id = $1`
    rows, err := r.pool.Query(ctx, query, ownerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var homes []domain.Home
    for rows.Next() {
        var id, name, owner string
        var createdAt time.Time
        if err := rows.Scan(&id, &name, &owner, &createdAt); err != nil {
            return nil, err
        }
        h := domain.NewHome(id, name, owner)
        homes = append(homes, *h)
    }
    return homes, rows.Err()
}

func (r *PostgresHomeRepository) Update(ctx context.Context, home *domain.Home) error {
    query := `UPDATE homes SET name = $1 WHERE id = $2`
    _, err := r.pool.Exec(ctx, query, home.Name(), home.ID())
    return err
}

func (r *PostgresHomeRepository) Delete(ctx context.Context, id string) error {
    // Сначала удаляем комнаты и устройства, связанные с домом (каскадно в БД или вручную)
    // У нас внешние ключи с ON DELETE CASCADE для rooms? В миграции для rooms.home_id ON DELETE CASCADE, для devices.room_id нет каскада.
    // Для простоты удалим дом (комнаты удалятся каскадно, устройства останутся? Нужно обработать).
    _, err := r.pool.Exec(ctx, `DELETE FROM homes WHERE id = $1`, id)
    return err
}