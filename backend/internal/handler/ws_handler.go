package handler

import (
	"log"
	"net/http"

	"github.com/VLGFoxRU/smartic-home/internal/middleware"
	"github.com/VLGFoxRU/smartic-home/internal/ws"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // разрешаем все origins (для разработки)
	},
}

type WSHandler struct {
	hub    *ws.Hub
	secret []byte
}

func NewWSHandler(hub *ws.Hub, secret []byte) *WSHandler {
	return &WSHandler{hub: hub, secret: secret}
}

// ServeWS обрабатывает запрос на апгрейд до WebSocket
func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Проверяем токен из query-параметра (?token=...)
	tokenString := r.URL.Query().Get("token")
	log.Printf("WebSocket connection attempt, token=%s", tokenString)

	if tokenString == "" {
		http.Error(w, "Token required", http.StatusUnauthorized)
		return
	}

	// Здесь нужно валидировать токен (аналогично AuthMiddleware). Для простоты вызовем общую функцию.
	_, _, err := middleware.ValidateToken(tokenString, h.secret)
	if err != nil {
        log.Printf("WebSocket token validation failed: %v", err)
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := ws.NewClient(h.hub, conn)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}