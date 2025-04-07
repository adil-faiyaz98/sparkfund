package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

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

	// API Routes
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/v1/kyc/verify", kycVerifyHandler).Methods("POST")
	r.HandleFunc("/api/v1/kyc/status", kycStatusHandler).Methods("GET")

	// Serve static Swagger UI files
	fs := http.FileServer(http.Dir("./docs"))
	r.PathPrefix("/").Handler(fs)

	log.Printf("KYC Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "UP",
		Service: "kyc-service",
		Version: "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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
