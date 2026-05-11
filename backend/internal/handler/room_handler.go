package handler

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/VLGFoxRU/smartic-home/internal/service"
)

type RoomHandler struct {
    svc *service.RoomService
}

func NewRoomHandler(svc *service.RoomService) *RoomHandler {
    return &RoomHandler{svc: svc}
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
    homeID := mux.Vars(r)["homeId"]
    var req struct{ Name string `json:"name"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
        return
    }
    room, err := h.svc.Create(r.Context(), homeID, req.Name)
    if err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": map[string]string{
            "id":      room.ID(),
            "name":    room.Name(),
            "home_id": room.HomeID(),
        },
    })
}

func (h *RoomHandler) ListByHome(w http.ResponseWriter, r *http.Request) {
    homeID := mux.Vars(r)["homeId"]
    rooms, err := h.svc.ListByHome(r.Context(), homeID)
    if err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    dtos := make([]map[string]string, len(rooms))
    for i, room := range rooms {
        dtos[i] = map[string]string{
            "id":      room.ID(),
            "name":    room.Name(),
            "home_id": room.HomeID(),
        }
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data":    dtos,
    })
}

func (h *RoomHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    room, err := h.svc.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, `{"error":"room not found"}`, http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": map[string]string{
            "id":      room.ID(),
            "name":    room.Name(),
            "home_id": room.HomeID(),
        },
    })
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    var req struct{ Name string `json:"name"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
        return
    }
    if err := h.svc.Update(r.Context(), id, req.Name); err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if err := h.svc.Delete(r.Context(), id); err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}