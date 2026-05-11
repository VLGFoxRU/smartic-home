package repository

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload interface{}) error
}