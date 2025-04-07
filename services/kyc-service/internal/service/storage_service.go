package service

import (
	"context"
	"io"
)

// StorageService defines the interface for document storage operations
type StorageService interface {
	Store(ctx context.Context, path string, content io.Reader) error
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}
