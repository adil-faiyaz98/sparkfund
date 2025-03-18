package investments

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// InvestmentType represents the type of investment
type InvestmentType string

const (
	InvestmentTypeStock      InvestmentType = "STOCK"
	InvestmentTypeBond       InvestmentType = "BOND"
	InvestmentTypeMutualFund InvestmentType = "MUTUAL_FUND"
	InvestmentTypeETF        InvestmentType = "ETF"
	InvestmentTypeCrypto     InvestmentType = "CRYPTO"
)

// InvestmentStatus represents the current status of an investment
type InvestmentStatus string

const (
	InvestmentStatusActive  InvestmentStatus = "ACTIVE"
	InvestmentStatusSold    InvestmentStatus = "SOLD"
	InvestmentStatusPending InvestmentStatus = "PENDING"
	InvestmentStatusFailed  InvestmentStatus = "FAILED"
)

// Investment represents a financial investment
type Investment struct {
	ID            uuid.UUID        `json:"id" gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID        `json:"user_id" gorm:"type:uuid;not null"`
	AccountID     uuid.UUID        `json:"account_id" gorm:"type:uuid;not null"`
	Type          InvestmentType   `json:"type" gorm:"type:varchar(20);not null"`
	Status        InvestmentStatus `json:"status" gorm:"type:varchar(20);not null;default:'PENDING'"`
	Symbol        string           `json:"symbol" gorm:"type:varchar(10);not null"`
	Quantity      float64          `json:"quantity" gorm:"type:decimal(15,6);not null"`
	PurchasePrice float64          `json:"purchase_price" gorm:"type:decimal(15,2);not null"`
	CurrentPrice  float64          `json:"current_price" gorm:"type:decimal(15,2);not null"`
	Currency      string           `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`
	PurchaseDate  time.Time        `json:"purchase_date" gorm:"not null"`
	LastUpdated   time.Time        `json:"last_updated" gorm:"not null"`
	CreatedAt     time.Time        `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"not null"`
}

// InvestmentRepository defines the interface for investment data operations
type InvestmentRepository interface {
	Create(investment *Investment) error
	GetByID(id uuid.UUID) (*Investment, error)
	GetByUserID(userID uuid.UUID) ([]*Investment, error)
	GetByAccountID(accountID uuid.UUID) ([]*Investment, error)
	Update(investment *Investment) error
	Delete(id uuid.UUID) error
	GetBySymbol(symbol string) ([]*Investment, error)
}

// InvestmentService defines the interface for investment business logic operations
type InvestmentService interface {
	CreateInvestment(ctx context.Context, investment *Investment) error
	GetInvestment(ctx context.Context, id uuid.UUID) (*Investment, error)
	GetUserInvestments(ctx context.Context, userID uuid.UUID) ([]*Investment, error)
	GetAccountInvestments(ctx context.Context, accountID uuid.UUID) ([]*Investment, error)
	UpdateInvestment(ctx context.Context, investment *Investment) error
	DeleteInvestment(ctx context.Context, id uuid.UUID) error
	GetInvestmentsBySymbol(ctx context.Context, symbol string) ([]*Investment, error)
	UpdateInvestmentPrice(ctx context.Context, id uuid.UUID, newPrice float64) error
}
