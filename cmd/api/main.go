package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/VLGFoxRU/smartic-home/internal/handler"
    "github.com/VLGFoxRU/smartic-home/internal/infrastructure"
    "github.com/VLGFoxRU/smartic-home/internal/middleware"
    "github.com/VLGFoxRU/smartic-home/internal/service"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
)

func main() {
	pool, err := pgxpool.New(context.Background(), dbURL())
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Нет ответа от БД: %v", err)
	}
	log.Println("Подключено к PostgreSQL")

	// Подключение к Redis
    redisAddr := getEnv("REDIS_ADDR", "redis:6379")
    redisClient := redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })
    if err := redisClient.Ping(context.Background()).Err(); err != nil {
        log.Fatalf("Не удалось подключиться к Redis: %v", err)
    }
    log.Println("Подключено к Redis")

    // Собираем зависимости
    deviceRepo := infrastructure.NewPostgresDeviceRepository(pool)
    cacheRepo := infrastructure.NewRedisCacheRepository(redisClient)
    auditRepo := infrastructure.NewPostgresAuditRepository(pool)
    deviceSvc := service.NewDeviceService(deviceRepo, cacheRepo, auditRepo)
    deviceHandler := handler.NewDeviceHandler(deviceSvc)

    // Создаём маршрутизатор
    mux := http.NewServeMux()
    mux.HandleFunc("/api/v1/devices", deviceHandler.ListDevices)
    mux.HandleFunc("/api/v1/command", deviceHandler.SendCommand) // для /devices/{id}/command

    // Применяем middleware
    wrappedMux := middleware.LoggingMiddleware(mux)

    log.Println("API Gateway запущен на :8080")
    if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }
}

func dbURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "smarthome")
	pass := getEnv("DB_PASSWORD", "smarthome")
	dbname := getEnv("DB_NAME", "smarthome")
	return "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}