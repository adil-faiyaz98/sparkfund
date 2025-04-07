package service

import (
	"context"
)

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	Publish(ctx context.Context, eventType string, data interface{}) error
}
