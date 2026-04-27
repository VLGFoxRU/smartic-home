package repository

import "context"

type MessageBroker interface {
	Publish(ctx context.Context, topic string, message []byte) error
}