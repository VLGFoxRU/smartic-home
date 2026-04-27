package handler

import (
	"encoding/json"
	"net/http"

	"github.com/VLGFoxRU/smartic-home/internal/service"
)

type DeviceHandler struct {
	svc *service.DeviceService
}

func NewDeviceHandler(svc *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{svc: svc}
}

func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	devices, err := h.svc.ListDevices(r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Преобразуем доменные объекты в DTO
	dtos := make([]DeviceDTO, 0, len(devices))
	for _, d := range devices {
		dtos = append(dtos, DeviceDTO{
			ID:     d.ID(),
			Name:   d.Name(),
			Type:   d.Type(),
			Status: string(d.Status()),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    dtos,
	})
}

func (h *DeviceHandler) SendCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		DeviceID string                 `json:"device_id"`
		Command  string                 `json:"command"`
		Params   map[string]interface{} `json:"params"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.DeviceID == "" || req.Command == "" {
		http.Error(w, `{"error":"device_id and command are required"}`, http.StatusBadRequest)
		return
	}

	err := h.svc.SendCommand(r.Context(), req.DeviceID, req.Command, req.Params)
	if err != nil {
		http.Error(w, `{"error":"command failed: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}