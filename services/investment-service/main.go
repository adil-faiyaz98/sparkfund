package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

type Investment struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	InvestmentType string  `json:"investment_type"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type InvestmentsResponse struct {
	Investments []Investment `json:"investments"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PageSize    int          `json:"page_size"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/v1/investments", investmentsHandler)
	http.HandleFunc("/api/v1/investments/create", createInvestmentHandler)
	
	log.Printf("Investment Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "UP",
		Service: "investment-service",
		Version: "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func investmentsHandler(w http.ResponseWriter, r *http.Request) {
	investments := []Investment{
		{
			ID:            "inv-123456",
			UserID:        "user-123456",
			Amount:        1000.00,
			Currency:      "USD",
			Status:        "ACTIVE",
			InvestmentType: "STOCK",
			CreatedAt:     "2025-04-06T00:00:00Z",
			UpdatedAt:     "2025-04-06T00:00:00Z",
		},
		{
			ID:            "inv-789012",
			UserID:        "user-123456",
			Amount:        5000.00,
			Currency:      "USD",
			Status:        "PENDING",
			InvestmentType: "BOND",
			CreatedAt:     "2025-04-05T00:00:00Z",
			UpdatedAt:     "2025-04-05T00:00:00Z",
		},
	}

	response := InvestmentsResponse{
		Investments: investments,
		Total:       2,
		Page:        1,
		PageSize:    10,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	investment := Investment{
		ID:            "inv-345678",
		UserID:        "user-123456",
		Amount:        2000.00,
		Currency:      "USD",
		Status:        "PENDING",
		InvestmentType: "ETF",
		CreatedAt:     "2025-04-06T00:00:00Z",
		UpdatedAt:     "2025-04-06T00:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(investment)
}
