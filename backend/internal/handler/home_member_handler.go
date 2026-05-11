package handler

import (
    "encoding/json"
    "net/http"

    "github.com/VLGFoxRU/smartic-home/internal/domain"
    "github.com/VLGFoxRU/smartic-home/internal/middleware"
    "github.com/VLGFoxRU/smartic-home/internal/service"
    "github.com/gorilla/mux"
)

type HomeMemberHandler struct {
    svc *service.HomeMemberService
}

func NewHomeMemberHandler(svc *service.HomeMemberService) *HomeMemberHandler {
    return &HomeMemberHandler{svc: svc}
}

func (h *HomeMemberHandler) AddMember(w http.ResponseWriter, r *http.Request) {
    homeID := mux.Vars(r)["homeId"]
    var req struct {
        UserID string `json:"user_id"`
        Role   string `json:"role"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
        return
    }
    requesterID, _ := r.Context().Value(middleware.UserIDKey).(string)

    err := h.svc.AddMember(r.Context(), homeID, req.UserID, domain.HomeRole(req.Role), requesterID)
    if err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "added"})
}

func (h *HomeMemberHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
    homeID := mux.Vars(r)["homeId"]
    requesterID, _ := r.Context().Value(middleware.UserIDKey).(string)
    members, err := h.svc.ListMembers(r.Context(), homeID, requesterID)
    if err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusForbidden)
        return
    }
    dtos := make([]map[string]string, len(members))
    for i, m := range members {
        dtos[i] = map[string]string{
            "home_id": m.HomeID(),
            "user_id": m.UserID(),
            "role":    string(m.Role()),
        }
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data":    dtos,
    })
}

func (h *HomeMemberHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    homeID := vars["homeId"]
    userID := vars["userId"]
    var req struct {
        Role string `json:"role"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
        return
    }
    requesterID, _ := r.Context().Value(middleware.UserIDKey).(string)
    if err := h.svc.ChangeRole(r.Context(), homeID, userID, domain.HomeRole(req.Role), requesterID); err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *HomeMemberHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    homeID := vars["homeId"]
    userID := vars["userId"]
    requesterID, _ := r.Context().Value(middleware.UserIDKey).(string)
    if err := h.svc.RemoveMember(r.Context(), homeID, userID, requesterID); err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}