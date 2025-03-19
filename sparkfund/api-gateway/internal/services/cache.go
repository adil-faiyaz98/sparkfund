package services

import (
	"github.com/adil-faiyaz98/sparkfund/api-gateway/internal/models"
)

type Cache interface {
	Get(req *models.Request) ([]byte, error)
	Set(req *models.Request, data []byte)
	Delete(req *models.Request)
	Clear()
}

type cache struct {
	// TODO: Add dependencies (e.g., Redis for distributed caching)
}

func NewCache() Cache {
	return &cache{}
}

func (c *cache) Get(req *models.Request) ([]byte, error) {
	// TODO: Implement cache retrieval
	return nil, nil
}

func (c *cache) Set(req *models.Request, data []byte) {
	// TODO: Implement cache setting
}

func (c *cache) Delete(req *models.Request) {
	// TODO: Implement cache deletion
}

func (c *cache) Clear() {
	// TODO: Implement cache clearing
}
