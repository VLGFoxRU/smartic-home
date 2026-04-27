package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Device struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

var db *pgxpool.Pool

func main() {
	var err error
	db, err = pgxpool.New(context.Background(), dbURL())
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer db.Close()

	// Проверим соединение
	if err = db.Ping(context.Background()); err != nil {
		log.Fatalf("Нет ответа от БД: %v", err)
	}
	log.Println("Подключено к PostgreSQL")

	http.HandleFunc("/api/v1/devices", devicesHandler)

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

func devicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(context.Background(),
		"SELECT id, name, type, status FROM devices")
	if err != nil {
		log.Printf("Ошибка запроса устройств: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Type, &d.Status); err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			continue
		}
		devices = append(devices, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    devices,
	})
}