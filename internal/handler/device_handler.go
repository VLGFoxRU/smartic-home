package handler

import (
	"encoding/json"
	"net/http"

	"github.com/VLGFoxRU/smartic-home/internal/middleware"
	"github.com/VLGFoxRU/smartic-home/internal/service"
	"github.com/gorilla/mux"
)

type DeviceHandler struct {
	deviceSvc  *service.DeviceService
	controlSvc *service.ControlService
}

func NewDeviceHandler(deviceSvc *service.DeviceService, controlSvc *service.ControlService) *DeviceHandler {
	return &DeviceHandler{
		deviceSvc:  deviceSvc,
		controlSvc: controlSvc,
	}
}

func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	devices, err := h.deviceSvc.ListDevices(r.Context())
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

	// Извлекаем ID устройства из переменных маршрута
	vars := mux.Vars(r)
	deviceID := vars["id"]
	if deviceID == "" {
		http.Error(w, `{"error":"device id required"}`, http.StatusBadRequest)
		return
	}

	var req struct {
		Command string                 `json:"command"`
		Params  map[string]interface{} `json:"params"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Command == "" {
		http.Error(w, `{"error":"command field required"}`, http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста (установлен middleware)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role, _ := r.Context().Value(middleware.UserRoleKey).(string)

	err := h.controlSvc.SendCommand(r.Context(), deviceID, req.Command, req.Params, userID, role)
	if err != nil {
		http.Error(w, `{"error":"command failed: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}