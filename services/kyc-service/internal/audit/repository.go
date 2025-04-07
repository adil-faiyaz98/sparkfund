package audit

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for audit event storage
type Repository interface {
	// Create creates a new audit event
	Create(ctx context.Context, event *Event) error
	
	// GetByID gets an audit event by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	
	// GetByResourceID gets audit events by resource ID
	GetByResourceID(ctx context.Context, resourceType ResourceType, resourceID string, limit, offset int) ([]*Event, int, error)
	
	// GetByUserID gets audit events by user ID
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Event, int, error)
	
	// Search searches for audit events
	Search(ctx context.Context, query map[string]interface{}, limit, offset int) ([]*Event, int, error)
}
