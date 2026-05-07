package infrastructure

import (
    "context"
    "fmt"

    "github.com/VLGFoxRU/smartic-home/internal/domain"
    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRoomRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresRoomRepository(pool *pgxpool.Pool) repository.RoomRepository {
    return &PostgresRoomRepository{pool: pool}
}

func (r *PostgresRoomRepository) Create(ctx context.Context, room *domain.Room) error {
    query := `INSERT INTO rooms (id, home_id, name) VALUES ($1, $2, $3)`
    _, err := r.pool.Exec(ctx, query, room.ID(), room.HomeID(), room.Name())
    return err
}

func (r *PostgresRoomRepository) FindByID(ctx context.Context, id string) (*domain.Room, error) {
    query := `SELECT id, home_id, name FROM rooms WHERE id = $1`
    var idStr, homeID, name string
    err := r.pool.QueryRow(ctx, query, id).Scan(&idStr, &homeID, &name)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("комната не найдена")
        }
        return nil, err
    }
    return domain.NewRoom(idStr, homeID, name), nil
}

func (r *PostgresRoomRepository) FindByHome(ctx context.Context, homeID string) ([]domain.Room, error) {
    query := `SELECT id, home_id, name FROM rooms WHERE home_id = $1`
    rows, err := r.pool.Query(ctx, query, homeID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var rooms []domain.Room
    for rows.Next() {
        var id, hID, name string
        if err := rows.Scan(&id, &hID, &name); err != nil {
            return nil, err
        }
        rooms = append(rooms, *domain.NewRoom(id, hID, name))
    }
    return rooms, rows.Err()
}

func (r *PostgresRoomRepository) Update(ctx context.Context, room *domain.Room) error {
    query := `UPDATE rooms SET name = $1 WHERE id = $2`
    _, err := r.pool.Exec(ctx, query, room.Name(), room.ID())
    return err
}

func (r *PostgresRoomRepository) Delete(ctx context.Context, id string) error {
    _, err := r.pool.Exec(ctx, `DELETE FROM rooms WHERE id = $1`, id)
    return err
}