package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"log"

	"github.com/VLGFoxRU/smartic-home/internal/service"
	"github.com/VLGFoxRU/smartic-home/internal/domain"
	"github.com/VLGFoxRU/smartic-home/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	svc    *service.UserService
	secret []byte
}

func NewAuthHandler(svc *service.UserService, secret []byte) *AuthHandler {
	return &AuthHandler{svc: svc, secret: secret}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"username, email, password required"}`, http.StatusBadRequest)
		return
	}

	// По умолчанию регистрируем владельца (owner)
	user, err := h.svc.Register(r.Context(), req.Username, req.Email, req.Password, domain.RoleOwner)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":       user.ID(),
		"username": user.Username(),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	user, err := h.svc.Authenticate(r.Context(), creds.Username, creds.Password)
	if err != nil {
		log.Printf("Login failed for user %s: %v", creds.Username, err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID(),
		"role":    string(user.Role()),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(h.secret)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
    var req struct {
        OldPassword string `json:"old_password"`
        NewPassword string `json:"new_password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
        return
    }
    userID, _ := r.Context().Value(middleware.UserIDKey).(string)
    err := h.svc.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
    if err != nil {
        http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Me возвращает информацию о текущем авторизованном пользователе
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
    userID, _ := r.Context().Value(middleware.UserIDKey).(string)
    if userID == "" {
        http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
        return
    }

    user, err := h.svc.GetByID(r.Context(), userID)
    if err != nil {
        http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "id":       user.ID(),
            "username": user.Username(),
            "email":    user.Email(),
            "role":     string(user.Role()),
        },
    })
}