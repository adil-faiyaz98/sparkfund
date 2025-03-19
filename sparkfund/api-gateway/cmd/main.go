package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/adil-faiyaz98/sparkfund/api-gateway/internal/handlers"
	"github.com/adil-faiyaz98/sparkfund/api-gateway/internal/services"
)

func main() {
	// Initialize services
	authService := services.NewAuthService()
	rateLimiter := services.NewRateLimiter()
	cache := services.NewCache()
	loadBalancer := services.NewLoadBalancer()

	// Initialize handler
	handler := handlers.NewHandler(authService, rateLimiter, cache, loadBalancer)

	// Create router
	router := http.NewServeMux()
	handler.RegisterRoutes(router)

	// Create server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting API Gateway on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	log.Println("Shutting down server...")
	if err := server.Close(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
}
