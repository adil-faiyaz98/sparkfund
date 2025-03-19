package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yourusername/money-pulse/pkg/discovery"
	// ...existing code...
)

func main() {
	// ...existing code...

	// Initialize service discovery
	consulClient, err := discovery.NewConsulClient(getEnv("CONSUL_ADDR", "consul:8500"))
	if err != nil {
		log.Fatalf("Failed to create consul client: %v", err)
	}

	// Register service with Consul
	serviceHost := getEnv("SERVICE_HOST", "user-service")
	servicePort := 8080 // Update according to your configuration
	err = consulClient.Register("user-service", serviceHost, servicePort, []string{"user", "api"})
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	defer consulClient.Deregister()

	// ...existing code...

	// Create a health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// ...existing code...
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
