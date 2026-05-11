package infrastructure

import (
    "context"
    "fmt"
    "time"

    "github.com/VLGFoxRU/smartic-home/internal/repository"
    "github.com/redis/go-redis/v9"
)

type RedisCacheRepository struct {
    client *redis.Client
}

func NewRedisCacheRepository(client *redis.Client) repository.CacheRepository {
    return &RedisCacheRepository{client: client}
}

func (r *RedisCacheRepository) Get(ctx context.Context, key string) (string, error) {
    val, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return "", fmt.Errorf("key not found")
    }
    return val, err
}

func (r *RedisCacheRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
    return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCacheRepository) Del(ctx context.Context, key string) error {
    return r.client.Del(ctx, key).Err()
}

// SetNX пытается установить ключ, если его ещё нет. Возвращает true, если установка успешна.
func (r *RedisCacheRepository) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
    result, err := r.client.SetNX(ctx, key, value, ttl).Result()
    if err != nil {
        return false, fmt.Errorf("RedisCacheRepository.SetNX: %w", err)
    }
    return result, nil
}