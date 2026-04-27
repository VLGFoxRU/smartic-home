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