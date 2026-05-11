package handler

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/VLGFoxRU/smartic-home/internal/middleware"
    "github.com/VLGFoxRU/smartic-home/internal/service"
)

type HomeHandler struct {
    svc *service.HomeService
}

func NewHomeHandler(svc *service.HomeService) *HomeHandler {
    return &HomeHandler{svc: svc}
}

func (h *HomeHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name string `json:"name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
        return
    }
    userID, _ := r.Context().Value(middleware.UserIDKey).(string)
    home, err := h.svc.Create(r.Context(), req.Name, userID)
    if err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": map[string]string{
            "id":   home.ID(),
            "name": home.Name(),
        },
    })
}

func (h *HomeHandler) List(w http.ResponseWriter, r *http.Request) {
    userID, _ := r.Context().Value(middleware.UserIDKey).(string)
    homes, err := h.svc.ListByOwner(r.Context(), userID)
    if err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    dtos := make([]map[string]string, len(homes))
    for i, home := range homes {
        dtos[i] = map[string]string{"id": home.ID(), "name": home.Name()}
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data":    dtos,
    })
}

func (h *HomeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    home, err := h.svc.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, `{"error":"home not found"}`, http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": map[string]string{
            "id":      home.ID(),
            "name":    home.Name(),
            "owner_id": home.OwnerID(),
        },
    })
}

func (h *HomeHandler) Update(w http.ResponseWriter, r *http.Request) {
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

func (h *HomeHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if err := h.svc.Delete(r.Context(), id); err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}