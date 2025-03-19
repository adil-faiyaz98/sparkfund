package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/adil-faiyaz98/sparkfund/aml-service/config"
    "github.com/adil-faiyaz98/sparkfund/aml-service/internal/handlers"
    "github.com/adil-faiyaz98/sparkfund/aml-service/internal/services"
    "github.com/adil-faiyaz98/sparkfund/aml-service/internal/repositories"
    "github.com/adil-faiyaz98/sparkfund/aml-service/pkg/database"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    // Initialize database connection
    dbConfig := database.NewConfig()
    db, err := database.NewConnection(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize repositories
    repo := repositories.NewRepository(db)

    // Initialize services
    svc := services.NewService(repo)

    // Initialize handlers
    handler := handlers.NewHandler(svc)

    // Set up router
    r := mux.NewRouter()
    handler.RegisterRoutes(r)

    // Start server
    log.Printf("Starting aml-service on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
