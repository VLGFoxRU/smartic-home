package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VLGFoxRU/smartic-home/internal/service"
	"github.com/gorilla/mux"
)

type TelemetryHandler struct {
	svc *service.TelemetryService
}

func NewTelemetryHandler(svc *service.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{svc: svc}
}

// IngestTelemetry принимает одно или массив показаний
func (h *TelemetryHandler) IngestTelemetry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["id"]

	var req struct {
		Type      string  `json:"type"`
		Value     float64 `json:"value"`
		Timestamp string  `json:"timestamp"` // RFC3339, опционально
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Type == "" || req.Value == 0 { // значение 0 допустимо? допустим, пусть проверяем на наличие
		http.Error(w, `{"error":"type and value required"}`, http.StatusBadRequest)
		return
	}

	ts := time.Now()
	if req.Timestamp != "" {
		var err error
		ts, err = time.Parse(time.RFC3339, req.Timestamp)
		if err != nil {
			http.Error(w, `{"error":"invalid timestamp format, use RFC3339"}`, http.StatusBadRequest)
			return
		}
	}

	if err := h.svc.Ingest(r.Context(), deviceID, req.Type, req.Value, ts); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

// GetTelemetryHistory возвращает историю показаний
func (h *TelemetryHandler) GetTelemetryHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["id"]

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	if fromStr == "" || toStr == "" {
		http.Error(w, `{"error":"from and to query params required"}`, http.StatusBadRequest)
		return
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		http.Error(w, `{"error":"invalid from format"}`, http.StatusBadRequest)
		return
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		http.Error(w, `{"error":"invalid to format"}`, http.StatusBadRequest)
		return
	}

	records, err := h.svc.GetHistory(r.Context(), deviceID, from, to)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Преобразуем в DTO
	type TelemetryDTO struct {
		DeviceID  string  `json:"device_id"`
		Timestamp string  `json:"timestamp"`
		Value     float64 `json:"value"`
		Type      string  `json:"type"`
	}
	dtos := make([]TelemetryDTO, len(records))
	for i, rec := range records {
		dtos[i] = TelemetryDTO{
			DeviceID:  rec.DeviceID,
			Timestamp: rec.Timestamp.Format(time.RFC3339),
			Value:     rec.Value,
			Type:      rec.Type,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    dtos,
	})
}

func (h *TelemetryHandler) GetLatestTelemetry(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    deviceID := vars["id"]
    if deviceID == "" {
        http.Error(w, `{"error":"device id required"}`, http.StatusBadRequest)
        return
    }

    // Запрашиваем последнюю запись (самую свежую) из сервиса
    records, err := h.svc.GetHistory(r.Context(), deviceID,
        time.Now().Add(-7*24*time.Hour), time.Now())
    if err != nil || len(records) == 0 {
        http.Error(w, `{"error":"no telemetry found"}`, http.StatusNotFound)
        return
    }

    // Берём самую свежую запись (последний элемент после сортировки DESC)
    latest := records[0]
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "device_id": latest.DeviceID,
            "timestamp": latest.Timestamp.Format(time.RFC3339),
            "value":     latest.Value,
            "type":      latest.Type,
        },
    })
}