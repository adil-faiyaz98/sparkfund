package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status represents the health status of a component
type Status string

const (
	// StatusUp indicates the component is healthy
	StatusUp Status = "UP"
	// StatusDown indicates the component is unhealthy
	StatusDown Status = "DOWN"
	// StatusDegraded indicates the component is partially healthy
	StatusDegraded Status = "DEGRADED"
)

// Component represents a service component that can be health checked
type Component struct {
	Name    string      `json:"name"`
	Status  Status      `json:"status"`
	Details interface{} `json:"details,omitempty"`
}

// Checker defines the interface for component health checks
type Checker interface {
	Check(context.Context) Component
}

// Handler handles health check requests
type Handler struct {
	checkers map[string]Checker
	timeout  time.Duration
	mu       sync.RWMutex
}

// NewHandler creates a new health check handler
func NewHandler(timeout time.Duration) *Handler {
	return &Handler{
		checkers: make(map[string]Checker),
		timeout:  timeout,
	}
}

// AddChecker adds a health checker for a component
func (h *Handler) AddChecker(name string, checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers[name] = checker
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	status := h.checkHealth(ctx)
	w.Header().Set("Content-Type", "application/json")

	if status.Status != StatusUp {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(status)
}

// Health represents the overall health status
type Health struct {
	Status     Status      `json:"status"`
	Components []Component `json:"components"`
	Timestamp  time.Time   `json:"timestamp"`
}

// checkHealth performs health checks on all registered components
func (h *Handler) checkHealth(ctx context.Context) Health {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var (
		components = make([]Component, 0, len(h.checkers))
		wg         sync.WaitGroup
		mu         sync.Mutex
		status     = StatusUp
	)

	for name, checker := range h.checkers {
		wg.Add(1)
		go func(name string, checker Checker) {
			defer wg.Done()

			component := checker.Check(ctx)
			mu.Lock()
			components = append(components, component)
			if component.Status == StatusDown {
				status = StatusDown
			} else if component.Status == StatusDegraded && status != StatusDown {
				status = StatusDegraded
			}
			mu.Unlock()
		}(name, checker)
	}

	wg.Wait()

	return Health{
		Status:     status,
		Components: components,
		Timestamp:  time.Now(),
	}
}

// DatabaseChecker implements health checking for a database
type DatabaseChecker struct {
	db *sql.DB
}

// NewDatabaseChecker creates a new database health checker
func NewDatabaseChecker(db *sql.DB) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

// Check implements the Checker interface for database
func (c *DatabaseChecker) Check(ctx context.Context) Component {
	err := c.db.PingContext(ctx)
	details := map[string]interface{}{
		"maxOpenConnections": c.db.Stats().MaxOpenConnections,
		"openConnections":    c.db.Stats().OpenConnections,
		"inUse":              c.db.Stats().InUse,
		"idle":               c.db.Stats().Idle,
	}

	if err != nil {
		return Component{
			Name:    "database",
			Status:  StatusDown,
			Details: details,
		}
	}

	return Component{
		Name:    "database",
		Status:  StatusUp,
		Details: details,
	}
}

// Example usage:
// handler := health.NewHandler(5 * time.Second)
// handler.AddChecker("database", health.NewDatabaseChecker(db))
// http.Handle("/health", handler)
