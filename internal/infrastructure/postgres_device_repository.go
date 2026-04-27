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

type PostgresDeviceRepository struct {
	pool *pgxpool.Pool
}

// Конструктор репозитория
func NewPostgresDeviceRepository(pool *pgxpool.Pool) repository.DeviceRepository {
	return &PostgresDeviceRepository{pool: pool}
}

func (r *PostgresDeviceRepository) FindAll(ctx context.Context) ([]domain.Device, error) {
	query := `SELECT id, name, type, status, last_seen, version FROM devices`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PostgresDeviceRepository.FindAll: %w", err)
	}
	defer rows.Close()

	var devices []domain.Device
	for rows.Next() {
		var (
			id       string
			name     string
			devType  string
			status   domain.DeviceStatus
			lastSeen *time.Time
			version  int
		)
		if err := rows.Scan(&id, &name, &devType, &status, &lastSeen, &version); err != nil {
			return nil, fmt.Errorf("PostgresDeviceRepository.FindAll scan: %w", err)
		}
		// Создаём доменный объект и восстанавливаем состояние
		device := domain.NewDevice(id, name, devType)
		// Восстанавливаем статус напрямую, но лучше добавить метод Reconstruct в domain.
		// Для простоты пока сделаем через SetStatus (добавим его в domain).
		// Добавим ниже.
		_ = device.SetStatus(status) // добавим метод в domain чуть позже
		devices = append(devices, *device)
	}
	return devices, rows.Err()
}

func (r *PostgresDeviceRepository) FindByID(ctx context.Context, id string) (*domain.Device, error) {
	query := `SELECT id, name, type, status, last_seen, version FROM devices WHERE id = $1`
	var (
		devID    string
		name     string
		devType  string
		status   domain.DeviceStatus
		lastSeen *time.Time
		version  int
	)
	err := r.pool.QueryRow(ctx, query, id).Scan(&devID, &name, &devType, &status, &lastSeen, &version)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("устройство не найдено")
		}
		return nil, fmt.Errorf("PostgresDeviceRepository.FindByID: %w", err)
	}
	device := domain.NewDevice(devID, name, devType)
	_ = device.SetStatus(status)
	return device, nil
}

func (r *PostgresDeviceRepository) Save(ctx context.Context, device *domain.Device) error {
	// TODO: обновление состояния в БД
	return fmt.Errorf("not implemented")
}

