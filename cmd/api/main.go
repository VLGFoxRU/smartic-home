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
	"github.com/VLGFoxRU/smartic-home/internal/engine"
	"github.com/VLGFoxRU/smartic-home/internal/detector"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	// БД
	pool, err := pgxpool.New(context.Background(), dbURL())
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Нет ответа от БД: %v", err)
	}
	log.Println("Подключено к PostgreSQL")

	// Redis
	redisAddr := getEnv("REDIS_ADDR", "redis:6379")
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	log.Println("Подключено к Redis")

	// Репозитории
	deviceRepo := infrastructure.NewPostgresDeviceRepository(pool)
	cacheRepo := infrastructure.NewRedisCacheRepository(redisClient)
	auditRepo := infrastructure.NewPostgresAuditRepository(pool)
	telemetryRepo := infrastructure.NewPostgresTelemetryRepository(pool)
	sceneRepo := infrastructure.NewPostgresSceneRepository(pool)
	anomalyRepo := infrastructure.NewPostgresAnomalyRepository(pool)
	homeRepo := infrastructure.NewPostgresHomeRepository(pool)
	roomRepo := infrastructure.NewPostgresRoomRepository(pool)
	userRepo := infrastructure.NewPostgresUserRepository(pool)
	homeMemberRepo := infrastructure.NewPostgresHomeMemberRepository(pool)

	// Инфраструктурные адаптеры
	stubBroker := infrastructure.NewStubBroker() // замена на RabbitMQ позже

	// EventPublisher (WebSocket)
	hub := ws.NewHub()
	go hub.Run()
	wsEventPub := infrastructure.NewWSEventPublisher(hub)

	// Сервисы
	deviceSvc := service.NewDeviceService(deviceRepo, cacheRepo)
    controlSvc := service.NewControlService(deviceRepo, cacheRepo, auditRepo, stubBroker, wsEventPub)
    sceneSvc := service.NewSceneService(sceneRepo)
    anomalySvc := service.NewAnomalyService(anomalyRepo, wsEventPub) // получает publisher
	homeSvc := service.NewHomeService(homeRepo)
	roomSvc := service.NewRoomService(roomRepo)
	userSvc := service.NewUserService(userRepo)
	homeMemberSvc := service.NewHomeMemberService(homeMemberRepo, homeRepo, userRepo)

	// Движки и детекторы
	sceneEngine := engine.NewSceneEngine(sceneSvc, controlSvc, cacheRepo, wsEventPub)
	anomalyDetector := detector.NewAnomalyDetector(telemetryRepo, anomalySvc)

	// TelemetryService (с процессорами)
	telemetrySvc := service.NewTelemetryService(telemetryRepo, wsEventPub, anomalyDetector, sceneEngine)

	// Обработчики
	jwtSecret  := []byte(getEnv("JWT_SECRET", "super-secret-key"))
	authHandler := handler.NewAuthHandler(userSvc, jwtSecret)
	deviceHandler := handler.NewDeviceHandler(deviceSvc, controlSvc)
	telemetryHandler := handler.NewTelemetryHandler(telemetrySvc)
	sceneHandler := handler.NewSceneHandler(sceneSvc)
	anomalyHandler := handler.NewAnomalyHandler(anomalySvc)
	wsHandler := handler.NewWSHandler(hub, jwtSecret )
	homeHandler := handler.NewHomeHandler(homeSvc)
	roomHandler := handler.NewRoomHandler(roomSvc)
	homeMemberHandler := handler.NewHomeMemberHandler(homeMemberSvc)

	// Роутер
	r := mux.NewRouter()
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.LoggingMiddleware)

	// Публичные
	r.HandleFunc("/api/v1/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/ws", wsHandler.ServeWS)

	// Защищённые (JWT)
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.AuthMiddleware(jwtSecret ))

	// Devices
	api.HandleFunc("/devices", deviceHandler.ListDevices).Methods("GET")
	api.HandleFunc("/devices", deviceHandler.CreateDevice).Methods("POST")
	api.HandleFunc("/devices/{id}", deviceHandler.GetDevice).Methods("GET")
	api.HandleFunc("/devices/{id}", deviceHandler.UpdateDevice).Methods("PUT")
	api.HandleFunc("/devices/{id}", deviceHandler.DeleteDevice).Methods("DELETE")
	api.HandleFunc("/devices/{id}/command", deviceHandler.SendCommand).Methods("POST")
	api.HandleFunc("/devices/{id}/status", deviceHandler.UpdateStatus).Methods("PATCH")

	// Telemetry
	api.HandleFunc("/telemetry/{id}", telemetryHandler.IngestTelemetry).Methods("POST")
	api.HandleFunc("/telemetry/{id}", telemetryHandler.GetTelemetryHistory).Methods("GET")

	// Scenes
	api.HandleFunc("/scenes", sceneHandler.ListScenes).Methods("GET")
	api.HandleFunc("/scenes", sceneHandler.CreateScene).Methods("POST")
	api.HandleFunc("/scenes/{id}", sceneHandler.GetScene).Methods("GET")
	api.HandleFunc("/scenes/{id}", sceneHandler.UpdateScene).Methods("PUT")
	api.HandleFunc("/scenes/{id}", sceneHandler.DeleteScene).Methods("DELETE")
	api.HandleFunc("/scenes/{id}/activate", sceneHandler.ActivateScene).Methods("PATCH")
	api.HandleFunc("/scenes/{id}/pause", sceneHandler.PauseScene).Methods("PATCH")

	// Anomalies
	api.HandleFunc("/anomalies", anomalyHandler.ListAnomalies).Methods("GET")
	api.HandleFunc("/anomalies", anomalyHandler.CreateAnomaly).Methods("POST")
	api.HandleFunc("/anomalies/{id}", anomalyHandler.UpdateAnomalyStatus).Methods("PATCH")

	// Homes
	api.HandleFunc("/homes", homeHandler.List).Methods("GET")
	api.HandleFunc("/homes", homeHandler.Create).Methods("POST")
	api.HandleFunc("/homes/{id}", homeHandler.GetByID).Methods("GET")
	api.HandleFunc("/homes/{id}", homeHandler.Update).Methods("PUT")
	api.HandleFunc("/homes/{id}", homeHandler.Delete).Methods("DELETE")

	// Rooms
	api.HandleFunc("/homes/{homeId}/rooms", roomHandler.ListByHome).Methods("GET")
	api.HandleFunc("/homes/{homeId}/rooms", roomHandler.Create).Methods("POST")
	api.HandleFunc("/rooms/{id}", roomHandler.GetByID).Methods("GET")
	api.HandleFunc("/rooms/{id}", roomHandler.Update).Methods("PUT")
	api.HandleFunc("/rooms/{id}", roomHandler.Delete).Methods("DELETE")

	// Home members
	api.HandleFunc("/homes/{homeId}/members", homeMemberHandler.ListMembers).Methods("GET")
	api.HandleFunc("/homes/{homeId}/members", homeMemberHandler.AddMember).Methods("POST")
	api.HandleFunc("/homes/{homeId}/members/{userId}", homeMemberHandler.ChangeRole).Methods("PUT")
	api.HandleFunc("/homes/{homeId}/members/{userId}", homeMemberHandler.RemoveMember).Methods("DELETE")

	// Старт сервера
	log.Println("API Gateway запущен на :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
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
