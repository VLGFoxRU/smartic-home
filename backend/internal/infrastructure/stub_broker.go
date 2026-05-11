package infrastructure

import (
	"context"
	"log"
)

type StubBroker struct{}

func NewStubBroker() *StubBroker {
	return &StubBroker{}
}

func (b *StubBroker) Publish(ctx context.Context, topic string, message []byte) error {
	log.Printf("[STUB BROKER] topic=%s message=%s", topic, string(message))
	return nil
}