package ai

import (
    "context"
    "sync"
    "time"
)

type ModelManager struct {
    models     map[string]Model
    versions   map[string]string
    mu         sync.RWMutex
    updateChan chan ModelUpdate
}

func NewModelManager(ctx context.Context) *ModelManager {
    mm := &ModelManager{
        models:     make(map[string]Model),
        versions:   make(map[string]string),
        updateChan: make(chan ModelUpdate),
    }
    go mm.modelUpdateLoop(ctx)
    return mm
}

func (mm *ModelManager) GetModel(name string) (Model, error) {
    mm.mu.RLock()
    defer mm.mu.RUnlock()
    if model, exists := mm.models[name]; exists {
        return model, nil
    }
    return nil, ErrModelNotFound
}