package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Investment Service API
// @version 1.0
// @description Investment Service for SparkFund
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /
// @schemes http

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

type Investment struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	InvestmentType string  `json:"investment_type"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
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

	r := mux.NewRouter()

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API Routes
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/v1/investments", investmentsHandler).Methods("GET")
	r.HandleFunc("/api/v1/investments/create", createInvestmentHandler).Methods("POST")

	// Serve static Swagger files
	fs := http.FileServer(http.Dir("./docs"))
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	log.Printf("Investment Service starting on port %s...", port)
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
		Service: "investment-service",
		Version: "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get investments
// @Description Get list of investments
// @Tags investments
// @Accept json
// @Produce json
// @Success 200 {object} InvestmentsResponse
// @Router /api/v1/investments [get]
func investmentsHandler(w http.ResponseWriter, r *http.Request) {
	investments := []Investment{
		{
			ID:             "inv-123456",
			UserID:         "user-123456",
			Amount:         1000.00,
			Currency:       "USD",
			Status:         "ACTIVE",
			InvestmentType: "STOCK",
			CreatedAt:      "2025-04-06T00:00:00Z",
			UpdatedAt:      "2025-04-06T00:00:00Z",
		},
		{
			ID:             "inv-789012",
			UserID:         "user-123456",
			Amount:         5000.00,
			Currency:       "USD",
			Status:         "PENDING",
			InvestmentType: "BOND",
			CreatedAt:      "2025-04-05T00:00:00Z",
			UpdatedAt:      "2025-04-05T00:00:00Z",
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

// @Summary Create investment
// @Description Create a new investment
// @Tags investments
// @Accept json
// @Produce json
// @Success 200 {object} Investment
// @Router /api/v1/investments/create [post]
func createInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	investment := Investment{
		ID:             "inv-345678",
		UserID:         "user-123456",
		Amount:         2000.00,
		Currency:       "USD",
		Status:         "PENDING",
		InvestmentType: "ETF",
		CreatedAt:      "2025-04-06T00:00:00Z",
		UpdatedAt:      "2025-04-06T00:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(investment)
}
