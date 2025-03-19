package service

import (
	"context"
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/investments"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockInvestmentRepository struct {
	testutil.MockRepository
}

func (m *MockInvestmentRepository) Create(ctx context.Context, investment *investments.Investment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *MockInvestmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*investments.Investment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*investments.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*investments.Investment, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*investments.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) Update(ctx context.Context, investment *investments.Investment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *MockInvestmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestInvestmentService_CreateInvestment(t *testing.T) {
	// Setup
	mockRepo := new(MockInvestmentRepository)
	service := NewInvestmentService(mockRepo)
	ctx := context.Background()

	// Test data
	userID := uuid.New()
	testInvestment := &investments.Investment{
		UserID:      userID,
		Type:        investments.InvestmentTypeStocks,
		Amount:      5000.0,
		RiskLevel:   investments.RiskLevelMedium,
		Duration:    36,
		Description: "Stock portfolio investment",
	}

	t.Run("Create Valid Investment", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("Create", ctx, mock.AnythingOfType("*investments.Investment")).Return(nil)

		// Execute
		err := service.CreateInvestment(ctx, testInvestment)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, testInvestment.ID)
		assert.Equal(t, investments.InvestmentStatusActive, testInvestment.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create Investment with Invalid Type", func(t *testing.T) {
		// Setup test data
		invalidInvestment := *testInvestment
		invalidInvestment.Type = "invalid_type"

		// Execute
		err := service.CreateInvestment(ctx, &invalidInvestment)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, invalidInvestment.ID)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Create Investment with Negative Amount", func(t *testing.T) {
		// Setup test data
		invalidInvestment := *testInvestment
		invalidInvestment.Amount = -1000.0

		// Execute
		err := service.CreateInvestment(ctx, &invalidInvestment)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, invalidInvestment.ID)
		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestInvestmentService_GetInvestment(t *testing.T) {
	// Setup
	mockRepo := new(MockInvestmentRepository)
	service := NewInvestmentService(mockRepo)
	ctx := context.Background()

	// Test data
	investmentID := uuid.New()
	testInvestment := &investments.Investment{
		ID:          investmentID,
		UserID:      uuid.New(),
		Type:        investments.InvestmentTypeStocks,
		Amount:      5000.0,
		RiskLevel:   investments.RiskLevelMedium,
		Duration:    36,
		Description: "Stock portfolio investment",
		Status:      investments.InvestmentStatusActive,
	}

	t.Run("Get Existing Investment", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, investmentID).Return(testInvestment, nil)

		// Execute
		investment, err := service.GetInvestment(ctx, investmentID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, testInvestment.ID, investment.ID)
		assert.Equal(t, testInvestment.UserID, investment.UserID)
		assert.Equal(t, testInvestment.Type, investment.Type)
		assert.Equal(t, testInvestment.Amount, investment.Amount)
		assert.Equal(t, testInvestment.RiskLevel, investment.RiskLevel)
		assert.Equal(t, testInvestment.Duration, investment.Duration)
		assert.Equal(t, testInvestment.Description, investment.Description)
		assert.Equal(t, testInvestment.Status, investment.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get Non-Existing Investment", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, investmentID).Return(nil, investments.ErrInvestmentNotFound)

		// Execute
		investment, err := service.GetInvestment(ctx, investmentID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, investment)
		assert.Equal(t, investments.ErrInvestmentNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestInvestmentService_UpdateInvestmentStatus(t *testing.T) {
	// Setup
	mockRepo := new(MockInvestmentRepository)
	service := NewInvestmentService(mockRepo)
	ctx := context.Background()

	// Test data
	investmentID := uuid.New()
	testInvestment := &investments.Investment{
		ID:          investmentID,
		UserID:      uuid.New(),
		Type:        investments.InvestmentTypeStocks,
		Amount:      5000.0,
		RiskLevel:   investments.RiskLevelMedium,
		Duration:    36,
		Description: "Stock portfolio investment",
		Status:      investments.InvestmentStatusActive,
	}

	t.Run("Update Investment Status Successfully", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, investmentID).Return(testInvestment, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*investments.Investment")).Return(nil)

		// Execute
		err := service.UpdateInvestmentStatus(ctx, investmentID, investments.InvestmentStatusCompleted, "Investment completed successfully")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, investments.InvestmentStatusCompleted, testInvestment.Status)
		assert.Equal(t, "Investment completed successfully", testInvestment.StatusReason)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update Non-Existing Investment Status", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, investmentID).Return(nil, investments.ErrInvestmentNotFound)

		// Execute
		err := service.UpdateInvestmentStatus(ctx, investmentID, investments.InvestmentStatusCompleted, "Investment completed successfully")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, investments.ErrInvestmentNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update Investment Status with Invalid Status", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, investmentID).Return(testInvestment, nil)

		// Execute
		err := service.UpdateInvestmentStatus(ctx, investmentID, "invalid_status", "Invalid status")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, investments.InvestmentStatusActive, testInvestment.Status)
		mockRepo.AssertNotCalled(t, "Update")
	})
}
