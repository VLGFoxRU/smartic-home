package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VLGFoxRU/smartic-home/internal/domain"
	"github.com/VLGFoxRU/smartic-home/internal/middleware"
	"github.com/VLGFoxRU/smartic-home/internal/service"
	"github.com/gorilla/mux"
)

type AnomalyHandler struct {
	svc *service.AnomalyService
}

func NewAnomalyHandler(svc *service.AnomalyService) *AnomalyHandler {
	return &AnomalyHandler{svc: svc}
}

func (h *AnomalyHandler) CreateAnomaly(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID      string  `json:"device_id"`
		Value         float64 `json:"value"`
		ExpectedValue *float64 `json:"expected_value,omitempty"`
		Severity      string  `json:"severity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.DeviceID == "" || req.Severity == "" {
		http.Error(w, `{"error":"device_id and severity required"}`, http.StatusBadRequest)
		return
	}
	anomaly, err := h.svc.CreateAnomaly(r.Context(), req.DeviceID, req.Value, req.ExpectedValue, req.Severity)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    anomalyToDTO(anomaly),
	})
}

func (h *AnomalyHandler) ListAnomalies(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("device_id")
	status := r.URL.Query().Get("status")
	anomalies, err := h.svc.ListAnomalies(r.Context(), deviceID, status)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	dtos := make([]map[string]interface{}, len(anomalies))
	for i, a := range anomalies {
		dtos[i] = anomalyToDTO(&a)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    dtos,
	})
}

func (h *AnomalyHandler) UpdateAnomalyStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	if err := h.svc.UpdateStatus(r.Context(), id, req.Status, userID); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func anomalyToDTO(a *domain.Anomaly) map[string]interface{} {
	return map[string]interface{}{
		"id":              a.ID(),
		"device_id":       a.DeviceID(),
		"detected_at":     a.DetectedAt().Format(time.RFC3339),
		"value":           a.Value(),
		"expected_value":  a.ExpectedValue(),
		"severity":        a.Severity(),
		"status":          string(a.Status()),
		"acknowledged_by": a.AcknowledgedBy,
		"acknowledged_at": func() string {
			if a.AcknowledgedAt != nil {
				return a.AcknowledgedAt.Format(time.RFC3339)
			}
			return ""
		}(),
		"resolved_by": a.ResolvedBy,
		"resolved_at": func() string {
			if a.ResolvedAt != nil {
				return a.ResolvedAt.Format(time.RFC3339)
			}
			return ""
		}(),
	}
}