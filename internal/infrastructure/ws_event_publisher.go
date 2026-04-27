package infrastructure

import (
	"context"
	"time"

	"github.com/VLGFoxRU/smartic-home/internal/repository"
	"github.com/VLGFoxRU/smartic-home/internal/ws"
)

type WSEventPublisher struct {
	hub *ws.Hub
}

func NewWSEventPublisher(hub *ws.Hub) repository.EventPublisher {
	return &WSEventPublisher{hub: hub}
}

func (p *WSEventPublisher) Publish(ctx context.Context, eventType string, payload interface{}) error {
	message := map[string]interface{}{
		"type":      eventType,
		"payload":   payload,
		"timestamp": time.Now().UTC(),
		"status":    "success",
	}
	return p.hub.Broadcast(message)
}