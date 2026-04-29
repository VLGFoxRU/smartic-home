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

type PostgresSceneRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSceneRepository(pool *pgxpool.Pool) repository.SceneRepository {
	return &PostgresSceneRepository{pool: pool}
}

func (r *PostgresSceneRepository) Create(ctx context.Context, scene *domain.Scene) error {
	query := `INSERT INTO scenes (id, home_id, name, condition_json, action_json, is_active, created_by, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query,
		scene.ID(), scene.HomeID(), scene.Name(),
		scene.ConditionJSON(), scene.ActionJSON(),
		scene.IsActive(), scene.CreatedBy(), scene.CreatedAt(),
	)
	return err
}

func (r *PostgresSceneRepository) FindByID(ctx context.Context, id string) (*domain.Scene, error) {
	query := `SELECT id, home_id, name, condition_json, action_json, is_active, created_by, created_at, version
	          FROM scenes WHERE id = $1`
	var (
		idStr, homeID, name, condJSON, actJSON, createdBy string
		isActive bool
		createdAt time.Time
		version int
	)
	err := r.pool.QueryRow(ctx, query, id).Scan(&idStr, &homeID, &name, &condJSON, &actJSON, &isActive, &createdBy, &createdAt, &version)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("сценарий не найден")
		}
		return nil, err
	}
	scene, err := domain.NewSceneFromDB(idStr, homeID, name, createdBy, condJSON, actJSON, isActive, version, createdAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка восстановления сценария: %w", err)
	}
	return scene, nil
}

func (r *PostgresSceneRepository) FindByHome(ctx context.Context, homeID string) ([]domain.Scene, error) {
	query := `SELECT id, home_id, name, condition_json, action_json, is_active, created_by, created_at, version
	          FROM scenes WHERE home_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, homeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scenes []domain.Scene
	for rows.Next() {
		var (
			idStr, homeID, name, condJSON, actJSON, createdBy string
			isActive bool
			createdAt time.Time
			version int
		)
		if err := rows.Scan(&idStr, &homeID, &name, &condJSON, &actJSON, &isActive, &createdBy, &createdAt, &version); err != nil {
			return nil, err
		}
		scene, err := domain.NewSceneFromDB(idStr, homeID, name, createdBy, condJSON, actJSON, isActive, version, createdAt)
		if err != nil {
			return nil, err
		}
		scenes = append(scenes, *scene)
	}
	return scenes, rows.Err()
}

func (r *PostgresSceneRepository) Update(ctx context.Context, scene *domain.Scene) error {
	query := `UPDATE scenes SET name=$1, condition_json=$2, action_json=$3, is_active=$4, version=version+1
	          WHERE id=$5 AND version=$6`
	tag, err := r.pool.Exec(ctx, query,
		scene.Name(), scene.ConditionJSON(), scene.ActionJSON(),
		scene.IsActive(), scene.ID(), scene.Version(),
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("сценарий не найден или версия устарела")
	}
	return nil
}

func (r *PostgresSceneRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM scenes WHERE id = $1`, id)
	return err
}