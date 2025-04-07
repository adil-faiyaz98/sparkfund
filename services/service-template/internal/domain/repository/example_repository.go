package repository

import (
	"context"

	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/domain/model"
)

// ExampleRepository defines the interface for example repository
type ExampleRepository interface {
	// GetAll returns all examples
	GetAll(ctx context.Context) ([]*model.Example, error)

	// GetByID returns an example by ID
	GetByID(ctx context.Context, id string) (*model.Example, error)

	// Create creates a new example
	Create(ctx context.Context, example *model.Example) error

	// Update updates an example
	Update(ctx context.Context, example *model.Example) error

	// Delete deletes an example
	Delete(ctx context.Context, id string) error
}
