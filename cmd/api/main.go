package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/VLGFoxRU/smartic-home/internal/handler"
	"github.com/VLGFoxRU/smartic-home/internal/infrastructure"
	"github.com/VLGFoxRU/smartic-home/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
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

	// Собираем зависимости
	deviceRepo := infrastructure.NewPostgresDeviceRepository(pool)
	deviceSvc := service.NewDeviceService(deviceRepo)
	deviceHandler := handler.NewDeviceHandler(deviceSvc)

	// Регистрируем эндпоинт
	http.HandleFunc("/api/v1/devices", deviceHandler.ListDevices)

	log.Println("API Gateway запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
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