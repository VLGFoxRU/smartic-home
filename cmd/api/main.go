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
	"github.com/VLGFoxRU/smartic-home/internal/ws"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Подключение к БД
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
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	log.Println("Подключено к Redis")

	hub := ws.NewHub()
	go hub.Run()

	// Собираем зависимости
	deviceRepo := infrastructure.NewPostgresDeviceRepository(pool)
	cacheRepo := infrastructure.NewRedisCacheRepository(redisClient)
	auditRepo := infrastructure.NewPostgresAuditRepository(pool)
	wsEventPub := infrastructure.NewWSEventPublisher(hub)
	stubBroker := infrastructure.NewStubBroker()
	deviceSvc := service.NewDeviceService(deviceRepo, cacheRepo)
	controlSvc := service.NewControlService(deviceRepo, cacheRepo, auditRepo, stubBroker, wsEventPub)
	deviceHandler := handler.NewDeviceHandler(deviceSvc, controlSvc)
	telemetryRepo := infrastructure.NewPostgresTelemetryRepository(pool)
	telemetrySvc := service.NewTelemetryService(telemetryRepo)
	telemetryHandler := handler.NewTelemetryHandler(telemetrySvc)
	sceneRepo := infrastructure.NewPostgresSceneRepository(pool)
	sceneSvc := service.NewSceneService(sceneRepo)
	sceneHandler := handler.NewSceneHandler(sceneSvc)
	anomalyRepo := infrastructure.NewPostgresAnomalyRepository(pool)
	anomalySvc := service.NewAnomalyService(anomalyRepo)
	anomalyHandler := handler.NewAnomalyHandler(anomalySvc)

	// Секрет для JWT (в реальности из переменной окружения)
	jwtSecret := []byte(getEnv("JWT_SECRET", "super-secret-key"))
	authHandler := handler.NewAuthHandler(jwtSecret)

	// Создаём роутер
	r := mux.NewRouter()

	// Публичные маршруты
	r.HandleFunc("/api/v1/login", authHandler.Login).Methods("POST")

	// WebSocket handler
	wsHandler := handler.NewWSHandler(hub, jwtSecret)
	r.HandleFunc("/ws", wsHandler.ServeWS) // не защищаем middleware, т.к. проверка внутри

	// Защищённые маршруты (требуют JWT)
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.AuthMiddleware(jwtSecret))

	api.HandleFunc("/devices", deviceHandler.ListDevices).Methods("GET")
	api.HandleFunc("/devices", deviceHandler.CreateDevice).Methods("POST")
	api.HandleFunc("/devices/{id}", deviceHandler.GetDevice).Methods("GET")
	api.HandleFunc("/devices/{id}", deviceHandler.UpdateDevice).Methods("PUT")
	api.HandleFunc("/devices/{id}", deviceHandler.DeleteDevice).Methods("DELETE")
	api.HandleFunc("/devices/{id}/command", deviceHandler.SendCommand).Methods("POST")
	api.HandleFunc("/devices/{id}/status", deviceHandler.UpdateStatus).Methods("PATCH")
	api.HandleFunc("/telemetry/{id}", telemetryHandler.IngestTelemetry).Methods("POST")
	api.HandleFunc("/telemetry/{id}", telemetryHandler.GetTelemetryHistory).Methods("GET")
	api.HandleFunc("/scenes", sceneHandler.ListScenes).Methods("GET")
	api.HandleFunc("/scenes", sceneHandler.CreateScene).Methods("POST")
	api.HandleFunc("/scenes/{id}", sceneHandler.GetScene).Methods("GET")
	api.HandleFunc("/scenes/{id}", sceneHandler.UpdateScene).Methods("PUT")
	api.HandleFunc("/scenes/{id}", sceneHandler.DeleteScene).Methods("DELETE")
	api.HandleFunc("/scenes/{id}/activate", sceneHandler.ActivateScene).Methods("PATCH")
	api.HandleFunc("/scenes/{id}/pause", sceneHandler.PauseScene).Methods("PATCH")
	api.HandleFunc("/anomalies", anomalyHandler.ListAnomalies).Methods("GET")
	api.HandleFunc("/anomalies", anomalyHandler.CreateAnomaly).Methods("POST")
	api.HandleFunc("/anomalies/{id}", anomalyHandler.UpdateAnomalyStatus).Methods("PATCH")

	// Применяем логгирование ко всем маршрутам
	wrappedR := middleware.CORSMiddleware(middleware.LoggingMiddleware(r))

	log.Println("API Gateway запущен на :8080")
	if err := http.ListenAndServe(":8080", wrappedR); err != nil {
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
