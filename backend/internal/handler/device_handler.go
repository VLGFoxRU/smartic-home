package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"context"

	"github.com/VLGFoxRU/smartic-home/internal/middleware"
	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/VLGFoxRU/smartic-home/internal/service"
	"github.com/gorilla/mux"
)

type DeviceHandler struct {
    deviceSvc  *service.DeviceService
    controlSvc *service.ControlService
    cache      repository.CacheRepository
}

func NewDeviceHandler(
    deviceSvc *service.DeviceService,
    controlSvc *service.ControlService,
    cache repository.CacheRepository,
) *DeviceHandler {
    return &DeviceHandler{
        deviceSvc:  deviceSvc,
        controlSvc: controlSvc,
        cache:      cache,
    }
}

// getDeviceState читает состояние устройства из кэша или возвращает пустой объект
func (h *DeviceHandler) getDeviceState(ctx context.Context, deviceID string) json.RawMessage {
    stateJSON, err := h.cache.Get(ctx, "device:"+deviceID+":state")
    if err != nil || stateJSON == "" {
        return json.RawMessage(`{}`)
    }
    return json.RawMessage(stateJSON)
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
	dtos := make([]DeviceDTO, len(devices))
    for i, d := range devices {
        dtos[i] = DeviceDTO{
            ID:     d.ID(),
            Name:   d.Name(),
            Type:   d.Type(),
            Status: string(d.Status()),
            RoomID: d.RoomID(),
            State: h.getDeviceState(r.Context(), d.ID()),
        }
    }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    dtos,
	})
}

func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	dev, err := h.deviceSvc.GetDevice(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"device not found"}`, http.StatusNotFound)
		return
	}
	dto := DeviceDTO{
        ID:     dev.ID(),
        Name:   dev.Name(),
        Type:   dev.Type(),
        Status: string(dev.Status()),
        RoomID: dev.RoomID(),
        State: h.getDeviceState(r.Context(), dev.ID()),
    }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    dto,
	})
}

func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		RoomID string `json:"room_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Type == "" {
		http.Error(w, `{"error":"name and type required"}`, http.StatusBadRequest)
		return
	}
	dev, err := h.deviceSvc.CreateDevice(r.Context(), req.Name, req.Type, req.RoomID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	dto := DeviceDTO{
		ID:     dev.ID(),
		Name:   dev.Name(),
		Type:   dev.Type(),
		Status: string(dev.Status()),
		RoomID: dev.RoomID(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    dto,
	})
}

func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if err := h.deviceSvc.UpdateDevice(r.Context(), id, req.Name, req.Type); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.deviceSvc.DeleteDevice(r.Context(), id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func (h *DeviceHandler) SendCommand(w http.ResponseWriter, r *http.Request) {
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

func (h *DeviceHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req struct {
		Status string `json:"status"` // "enable" или "disable"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	var err error
	switch req.Status {
	case "enable":
		err = h.deviceSvc.EnableDevice(r.Context(), id)
	case "disable":
		err = h.deviceSvc.DisableDevice(r.Context(), id)
	default:
		err = fmt.Errorf("недопустимый статус: %s", req.Status)
	}
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}