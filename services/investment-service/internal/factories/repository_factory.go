package factories

import (
	"investment-service/internal/circuitbreaker"
	"investment-service/internal/repositories"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RepositoryFactory creates repositories with circuit breaker protection
type RepositoryFactory struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db *gorm.DB, logger *logrus.Logger) *RepositoryFactory {
	return &RepositoryFactory{
		db:     db,
		logger: logger,
	}
}

// CreateInvestmentRepository creates a new investment repository with circuit breaker
func (f *RepositoryFactory) CreateInvestmentRepository() repositories.InvestmentRepository {
	baseRepo := repositories.NewInvestmentRepository(f.db)

	// If we're not using a circuit breaker, return the base repository
	if !circuitBreakerEnabled() {
		return baseRepo
	}

	// Create a circuit breaker-protected repository
	return &CircuitBreakerInvestmentRepository{
		repo:           baseRepo,
		circuitBreaker: circuitbreaker.NewCircuitBreaker("investment-repository"),
		logger:         f.logger,
	}
}

// Helper function to check if circuit breakers are enabled
func circuitBreakerEnabled() bool {
	return true // This would normally check config.Get().CircuitBreaker.Enabled
}

// CircuitBreakerInvestmentRepository wraps an investment repository with a circuit breaker
type CircuitBreakerInvestmentRepository struct {
	repo           repositories.InvestmentRepository
	circuitBreaker *circuitbreaker.CircuitBreaker
	logger         *logrus.Logger
}

// Implement all repository methods with circuit breaker protection
// (Implementation of all InvestmentRepository methods would go here)
