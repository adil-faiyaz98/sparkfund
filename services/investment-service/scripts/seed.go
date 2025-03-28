package main

import (
	"log"
	"time"

	"investment-service/internal/database"
	"investment-service/internal/models"
)

func main() {
	// Initialize database
	database.InitDB()

	// Create test user ID (using uint)
	userID := uint(1)

	// Create test portfolio
	portfolio := models.Portfolio{
		UserID:      userID,
		Name:        "Test Portfolio",
		Description: "A test portfolio for development",
		TotalValue:  10000.00,
		LastUpdated: time.Now(),
	}
	if err := database.DB.Create(&portfolio).Error; err != nil {
		log.Fatalf("Failed to create portfolio: %v", err)
	}

	// Create test investments
	investments := []models.Investment{
		{
			UserID:        userID,
			PortfolioID:   portfolio.ID,
			Amount:        1000.00,
			Type:          "STOCK",
			Status:        "ACTIVE",
			PurchaseDate:  time.Now(),
			PurchasePrice: 150.50,
			Symbol:        "AAPL",
			Quantity:      10,
			Notes:         "Test investment 1",
		},
		{
			UserID:        userID,
			PortfolioID:   portfolio.ID,
			Amount:        2000.00,
			Type:          "STOCK",
			Status:        "ACTIVE",
			PurchaseDate:  time.Now(),
			PurchasePrice: 280.75,
			Symbol:        "GOOGL",
			Quantity:      5,
			Notes:         "Test investment 2",
		},
	}

	for _, investment := range investments {
		if err := database.DB.Create(&investment).Error; err != nil {
			log.Fatalf("Failed to create investment: %v", err)
		}
	}

	log.Println("Seed data inserted successfully!")
}
