package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adil-faiyaz98/sparkfund/api-gateway/internal/models"
	"github.com/adil-faiyaz98/sparkfund/api-gateway/internal/services"
	"github.com/gorilla/mux"
)

type Handler struct {
	authService  services.AuthService
	rateLimiter  services.RateLimiter
	cache        services.Cache
	loadBalancer services.LoadBalancer
}

func NewHandler(
	authService services.AuthService,
	rateLimiter services.RateLimiter,
	cache services.Cache,
	loadBalancer services.LoadBalancer,
) *Handler {
	return &Handler{
		authService:  authService,
		rateLimiter:  rateLimiter,
		cache:        cache,
		loadBalancer: loadBalancer,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Health check endpoint
	mux.HandleFunc("/health", h.handleHealthCheck)

	// Protected routes
	mux.HandleFunc("/api/v1/kyc", h.handleKYC)
	mux.HandleFunc("/api/v1/blockchain", h.handleBlockchain)
	mux.HandleFunc("/api/v1/email", h.handleEmail)
	mux.HandleFunc("/api/v1/security", h.handleSecurity)
	mux.HandleFunc("/api/v1/notification", h.handleNotification)
	mux.HandleFunc("/api/v1/analytics", h.handleAnalytics)
	mux.HandleFunc("/api/v1/audit", h.handleAudit)
	mux.HandleFunc("/api/v1/identity", h.handleIdentity)
	mux.HandleFunc("/api/v1/transaction", h.handleTransaction)
	mux.HandleFunc("/api/v1/risk", h.handleRisk)
	mux.HandleFunc("/api/v1/compliance", h.handleCompliance)
}

func (h *Handler) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) handleKYC(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement KYC request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleBlockchain(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement blockchain request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleEmail(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement email request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleSecurity(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement security request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleNotification(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement notification request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleAnalytics(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement analytics request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleAudit(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement audit request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleIdentity(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement identity request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleTransaction(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement transaction request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleRisk(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement risk request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleCompliance(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement compliance request handling
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, models.Response{
		Status:  status,
		Message: message,
	})
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Extract service and path from URL
	vars := mux.Vars(r)
	service := vars["service"]
	path := vars["path"]

	// Create request object
	req := &models.Request{
		ID:            r.Header.Get("X-Request-ID"),
		Method:        r.Method,
		Path:          path,
		Headers:       make(map[string]string),
		Body:          nil, // TODO: Read body
		QueryParams:   r.URL.Query(),
		PathParams:    vars,
		ClientIP:      r.RemoteAddr,
		UserAgent:     r.UserAgent(),
		ServiceTarget: service,
	}

	// Check cache first
	if r.Method == "GET" {
		if cached, err := h.cache.Get(req); err == nil {
			w.Write(cached)
			return
		}
	}

	// Get service instance from load balancer
	serviceInstance, err := h.loadBalancer.GetInstance(service)
	if err != nil {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	// Forward request to service
	resp, err := serviceInstance.ForwardRequest(req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Cache response if it's a GET request
	if r.Method == "GET" {
		h.cache.Set(req, resp.Body)
	}

	// Write response
	for k, v := range resp.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(resp.Status)
	w.Write(resp.Body)
}

func (h *Handler) HandleAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify token
		user, err := h.authService.VerifyToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) HandleRateLimiting(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.User)
		if !h.rateLimiter.Allow(user.ID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) HandleLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement logging
		next.ServeHTTP(w, r)
	})
}
