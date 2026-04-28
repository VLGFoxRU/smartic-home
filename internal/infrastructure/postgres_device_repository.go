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

func NewPostgresDeviceRepository(pool *pgxpool.Pool) repository.DeviceRepository {
	return &PostgresDeviceRepository{pool: pool}
}

func (r *PostgresDeviceRepository) FindAll(ctx context.Context) ([]domain.Device, error) {
	query := `SELECT id, room_id, name, type, status, last_seen, version FROM devices`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PostgresDeviceRepository.FindAll: %w", err)
	}
	defer rows.Close()

	var devices []domain.Device
	for rows.Next() {
		var (
			id       string
			roomID   *string
			name     string
			devType  string
			status   domain.DeviceStatus
			lastSeen *time.Time
			version  int
		)
		if err := rows.Scan(&id, &roomID, &name, &devType, &status, &lastSeen, &version); err != nil {
			return nil, fmt.Errorf("PostgresDeviceRepository.FindAll scan: %w", err)
		}
		device := domain.NewDevice(id, name, devType)
		if roomID != nil {
			device.SetRoomID(*roomID)
		}
		_ = device.SetStatus(status)
		if lastSeen != nil {
			device.SetLastSeen(*lastSeen)
		}
		device.SetVersion(version)
		devices = append(devices, *device)
	}
	return devices, rows.Err()
}

func (r *PostgresDeviceRepository) FindByID(ctx context.Context, id string) (*domain.Device, error) {
	query := `SELECT id, room_id, name, type, status, last_seen, version FROM devices WHERE id = $1`
	var (
		devID    string
		roomID   *string
		name     string
		devType  string
		status   domain.DeviceStatus
		lastSeen *time.Time
		version  int
	)
	err := r.pool.QueryRow(ctx, query, id).Scan(&devID, &roomID, &name, &devType, &status, &lastSeen, &version)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("устройство не найдено")
		}
		return nil, fmt.Errorf("PostgresDeviceRepository.FindByID: %w", err)
	}
	device := domain.NewDevice(devID, name, devType)
	if roomID != nil {
		device.SetRoomID(*roomID)
	}
	_ = device.SetStatus(status)
	if lastSeen != nil {
		device.SetLastSeen(*lastSeen)
	}
	device.SetVersion(version)
	return device, nil
}

func (r *PostgresDeviceRepository) Create(ctx context.Context, device *domain.Device) error {
	query := `INSERT INTO devices (id, room_id, type, name, status) VALUES ($1, $2, $3, $4, $5)`
	var roomID interface{}
	if device.RoomID() != "" {
		roomID = device.RoomID()
	}
	_, err := r.pool.Exec(ctx, query, device.ID(), roomID, device.Type(), device.Name(), device.Status())
	if err != nil {
		return fmt.Errorf("PostgresDeviceRepository.Create: %w", err)
	}
	return nil
}

func (r *PostgresDeviceRepository) Save(ctx context.Context, device *domain.Device) error {
	query := `UPDATE devices SET name=$1, type=$2, status=$3, last_seen=$4, version=version+1
	          WHERE id=$5 AND version=$6`
	tag, err := r.pool.Exec(ctx, query, device.Name(), device.Type(), device.Status(), device.LastSeen(), device.ID(), device.Version())
	if err != nil {
		return fmt.Errorf("PostgresDeviceRepository.Save: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("устройство не найдено или версия устарела")
	}
	return nil
}

func (r *PostgresDeviceRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM devices WHERE id = $1`, id)
	return err
}