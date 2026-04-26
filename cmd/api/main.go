package main

import (
    "encoding/json"
    "log"
    "net/http"
)

// Структура для ответа (пока заглушка)
type Device struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Status string `json:"status"`
}

func main() {
    // Регистрируем обработчик по пути /api/v1/devices
    http.HandleFunc("/api/v1/devices", func(w http.ResponseWriter, r *http.Request) {
        // Заглушка списка устройств
        devices := []Device{
            {ID: "1", Name: "Лампа в гостиной", Status: "online"},
            {ID: "2", Name: "Датчик температуры", Status: "online"},
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": true,
            "data":    devices,
        })
    })

    log.Println("API Gateway запущен на :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }
}