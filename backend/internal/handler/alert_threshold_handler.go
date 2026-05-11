package handler

import (
    "encoding/json"
    "net/http"

    "github.com/VLGFoxRU/smartic-home/internal/middleware"
    "github.com/VLGFoxRU/smartic-home/internal/service"
    "github.com/gorilla/mux"
)

type AlertThresholdHandler struct {
    svc *service.AlertThresholdService
}

func NewAlertThresholdHandler(svc *service.AlertThresholdService) *AlertThresholdHandler {
    return &AlertThresholdHandler{svc: svc}
}

func (h *AlertThresholdHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        DeviceID      string   `json:"device_id"`
        TelemetryType string   `json:"telemetry_type"`
        MinValue      *float64 `json:"min_value,omitempty"`
        MaxValue      *float64 `json:"max_value,omitempty"`
        Severity      string   `json:"severity"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
        return
    }
    userID, _ := r.Context().Value(middleware.UserIDKey).(string)
    threshold, err := h.svc.Create(r.Context(), req.DeviceID, req.TelemetryType, req.Severity, userID, req.MinValue, req.MaxValue)
    if err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "id":             threshold.ID(),
            "device_id":      threshold.DeviceID(),
            "telemetry_type": threshold.TelemetryType(),
            "min_value":      threshold.MinValue(),
            "max_value":      threshold.MaxValue(),
            "severity":       threshold.Severity(),
        },
    })
}

func (h *AlertThresholdHandler) ListByDevice(w http.ResponseWriter, r *http.Request) {
    deviceID := r.URL.Query().Get("device_id")
    if deviceID == "" {
        http.Error(w, `{"error":"device_id required"}`, http.StatusBadRequest)
        return
    }
    thresholds, err := h.svc.ListByDevice(r.Context(), deviceID)
    if err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    dtos := make([]map[string]interface{}, len(thresholds))
    for i, t := range thresholds {
        dtos[i] = map[string]interface{}{
            "id":             t.ID(),
            "device_id":      t.DeviceID(),
            "telemetry_type": t.TelemetryType(),
            "min_value":      t.MinValue(),
            "max_value":      t.MaxValue(),
            "severity":       t.Severity(),
        }
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data":    dtos,
    })
}

func (h *AlertThresholdHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if err := h.svc.Delete(r.Context(), id); err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}