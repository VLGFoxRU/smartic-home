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

type SceneHandler struct {
	svc *service.SceneService
}

func NewSceneHandler(svc *service.SceneService) *SceneHandler {
	return &SceneHandler{svc: svc}
}

func (h *SceneHandler) CreateScene(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string                 `json:"name"`
		HomeID    string                 `json:"home_id"`
		Condition map[string]interface{} `json:"condition"`
		Action    map[string]interface{} `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	scene, err := h.svc.CreateScene(r.Context(), req.Name, req.HomeID, userID, req.Condition, req.Action)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":        scene.ID(),
			"name":      scene.Name(),
			"is_active": scene.IsActive(),
		},
	})
}

func (h *SceneHandler) GetScene(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	scene, err := h.svc.GetScene(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"scene not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    sceneToDTO(scene),
	})
}

func (h *SceneHandler) ListScenes(w http.ResponseWriter, r *http.Request) {
	homeID := r.URL.Query().Get("home_id")
	if homeID == "" {
		homeID = "b0000000-0000-0000-0000-000000000001" // дефолтный дом
	}
	scenes, err := h.svc.ListScenes(r.Context(), homeID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	dtos := make([]map[string]interface{}, len(scenes))
	for i, s := range scenes {
		dtos[i] = sceneToDTO(&s)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    dtos,
	})
}

func (h *SceneHandler) UpdateScene(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		Name      string                 `json:"name"`
		Condition map[string]interface{} `json:"condition"`
		Action    map[string]interface{} `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if err := h.svc.UpdateScene(r.Context(), id, req.Name, req.Condition, req.Action); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *SceneHandler) ActivateScene(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.Activate(r.Context(), id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "activated"})
}

func (h *SceneHandler) PauseScene(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.Pause(r.Context(), id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "paused"})
}

func (h *SceneHandler) DeleteScene(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.DeleteScene(r.Context(), id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func sceneToDTO(s *domain.Scene) map[string]interface{} {
	return map[string]interface{}{
		"id":         s.ID(),
		"home_id":    s.HomeID(),
		"name":       s.Name(),
		"condition":  s.Condition(),
		"action":     s.Action(),
		"is_active":  s.IsActive(),
		"created_at": s.CreatedAt().Format(time.RFC3339),
	}
}