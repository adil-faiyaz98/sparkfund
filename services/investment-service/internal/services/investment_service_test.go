package services_test

import (
	"context"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sparkfund/services/investment-service/internal/models"
	"github.com/sparkfund/services/investment-service/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockInvestmentRepository struct {
	mock.Mock
}

func (m *MockInvestmentRepository) Create(ctx context.Context, investment *models.Investment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *MockInvestmentRepository) GetByID(ctx context.Context, id uint) (*models.Investment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Investment, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) GetByPortfolioID(ctx context.Context, portfolioID uint) ([]models.Investment, error) {
	args := m.Called(ctx, portfolioID)
	return args.Get(0).([]models.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) Update(ctx context.Context, investment *models.Investment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *MockInvestmentRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInvestmentRepository) GetAll(ctx context.Context, page, pageSize int) ([]models.Investment, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]models.Investment), args.Get(1).(int64), args.Error(2)
}

// Tests
func TestCreateInvestment(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockInvestmentRepository)
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Silence logs

	// Create service
	service := services.NewInvestmentService(mockRepo, logger)

	// Test case: successful creation
	t.Run("Success", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		investment := &models.Investment{
			UserID:        1,
			PortfolioID:   1,
			Type:          "STOCK",
			Symbol:        "AAPL",
			Quantity:      10,
			PurchasePrice: 150.0,
		}

		// Mock expectations
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Investment")).Return(nil).Once()

		// Call service
		err := service.CreateInvestment(ctx, investment)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "ACTIVE", investment.Status)
		assert.Equal(t, 1500.0, investment.Amount)
		assert.False(t, investment.PurchaseDate.IsZero())
		mockRepo.AssertExpectations(t)
	})

	// Test case: validation error
	t.Run("ValidationError", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		investment := &models.Investment{
			UserID:        0, // Invalid: missing user ID
			Symbol:        "AAPL",
			Quantity:      10,
			PurchasePrice: 150.0,
		}

		// Call service
		err := service.CreateInvestment(ctx, investment)

		// Assert
		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestGetInvestment(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockInvestmentRepository)
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Silence logs

	// Create service
	service := services.NewInvestmentService(mockRepo, logger)

	// Test case: investment found
	t.Run("Found", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		expectedInvestment := &models.Investment{
			ID:            1,
			UserID:        1,
			Symbol:        "AAPL",
			Quantity:      10,
			PurchasePrice: 150.0,
		}

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(1)).Return(expectedInvestment, nil).Once()

		// Call service
		investment, err := service.GetInvestment(ctx, 1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedInvestment, investment)
		mockRepo.AssertExpectations(t)
	})

	// Test case: investment not found
	t.Run("NotFound", func(t *testing.T) {
		// Setup
		ctx := context.Background()

		// Mock expectations
		mockRepo.On("GetByID", ctx, uint(2)).Return(nil, nil).Once()

		// Call service
		investment, err := service.GetInvestment(ctx, 2)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, investment)
		assert.ErrorIs(t, err, services.ErrInvestmentNotFound)
		mockRepo.AssertExpectations(t)
	})
}
