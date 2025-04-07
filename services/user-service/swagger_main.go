package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title User Service API
// @version 1.0
// @description User Service for SparkFund
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
// @BasePath /
// @schemes http

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UsersResponse struct {
	Users    []User `json:"users"`
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	r := mux.NewRouter()

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API Routes
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/v1/users", usersHandler).Methods("GET")
	r.HandleFunc("/api/v1/users/register", registerUserHandler).Methods("POST")

	// Serve static Swagger files
	fs := http.FileServer(http.Dir("./docs"))
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	log.Printf("User Service starting on port %s...", port)
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
		Service: "user-service",
		Version: "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get users
// @Description Get list of users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UsersResponse
// @Router /api/v1/users [get]
func usersHandler(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{
			ID:        "user-123456",
			Email:     "john.doe@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Status:    "ACTIVE",
			CreatedAt: "2025-04-01T00:00:00Z",
			UpdatedAt: "2025-04-01T00:00:00Z",
		},
		{
			ID:        "user-789012",
			Email:     "jane.smith@example.com",
			FirstName: "Jane",
			LastName:  "Smith",
			Status:    "ACTIVE",
			CreatedAt: "2025-04-02T00:00:00Z",
			UpdatedAt: "2025-04-02T00:00:00Z",
		},
	}

	response := UsersResponse{
		Users:    users,
		Total:    2,
		Page:     1,
		PageSize: 10,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Register user
// @Description Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Router /api/v1/users/register [post]
func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		ID:        "user-345678",
		Email:     "new.user@example.com",
		FirstName: "New",
		LastName:  "User",
		Status:    "PENDING",
		CreatedAt: "2025-04-06T00:00:00Z",
		UpdatedAt: "2025-04-06T00:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
