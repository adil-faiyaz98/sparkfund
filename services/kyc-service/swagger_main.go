package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title KYC Service API
// @version 1.0
// @description KYC Service for SparkFund
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @schemes http

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

type KYCResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := mux.NewRouter()

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API Routes
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/v1/kyc/verify", kycVerifyHandler).Methods("POST")
	r.HandleFunc("/api/v1/kyc/status", kycStatusHandler).Methods("GET")

	// Serve static Swagger files
	fs := http.FileServer(http.Dir("./docs"))
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	log.Printf("KYC Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// @Summary Health check
// @Description Get service health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "UP",
		Service: "kyc-service",
		Version: "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Verify KYC
// @Description Submit KYC verification request
// @Tags kyc
// @Accept json
// @Produce json
// @Success 200 {object} KYCResponse
// @Router /api/v1/kyc/verify [post]
func kycVerifyHandler(w http.ResponseWriter, r *http.Request) {
	response := KYCResponse{
		ID:        "kyc-123456",
		UserID:    "user-123456",
		Status:    "PENDING",
		Message:   "KYC verification request received and is being processed",
		CreatedAt: "2025-04-06T00:00:00Z",
		UpdatedAt: "2025-04-06T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get KYC Status
// @Description Get KYC verification status
// @Tags kyc
// @Accept json
// @Produce json
// @Success 200 {object} KYCResponse
// @Router /api/v1/kyc/status [get]
func kycStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := KYCResponse{
		ID:        "kyc-123456",
		UserID:    "user-123456",
		Status:    "VERIFIED",
		Message:   "KYC verification completed successfully",
		CreatedAt: "2025-04-06T00:00:00Z",
		UpdatedAt: "2025-04-06T00:00:00Z",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
