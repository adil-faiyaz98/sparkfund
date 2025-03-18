package accounts

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeSavings    AccountType = "SAVINGS"
	AccountTypeChecking   AccountType = "CHECKING"
	AccountTypeInvestment AccountType = "INVESTMENT"
	AccountTypeCredit     AccountType = "CREDIT"
)

// AccountStatus represents the current status of an account
type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "ACTIVE"
	AccountStatusInactive  AccountStatus = "INACTIVE"
	AccountStatusSuspended AccountStatus = "SUSPENDED"
	AccountStatusClosed    AccountStatus = "CLOSED"
)

// Account represents a financial account
type Account struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID     `json:"user_id" gorm:"type:uuid;not null"`
	Type          AccountType   `json:"type" gorm:"type:varchar(20);not null"`
	Status        AccountStatus `json:"status" gorm:"type:varchar(20);not null;default:'ACTIVE'"`
	Balance       float64       `json:"balance" gorm:"type:decimal(15,2);not null;default:0"`
	Currency      string        `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`
	Name          string        `json:"name" gorm:"type:varchar(100);not null"`
	Description   string        `json:"description" gorm:"type:text"`
	AccountNumber string        `json:"account_number" gorm:"type:varchar(20);unique"`
	RoutingNumber string        `json:"routing_number" gorm:"type:varchar(20)"`
	CreatedAt     time.Time     `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"not null"`
}

// AccountRepository defines the interface for account data operations
type AccountRepository interface {
	Create(account *Account) error
	GetByID(id uuid.UUID) (*Account, error)
	GetByUserID(userID uuid.UUID) ([]*Account, error)
	Update(account *Account) error
	Delete(id uuid.UUID) error
	GetByAccountNumber(accountNumber string) (*Account, error)
}

// AccountService defines the interface for account business logic operations
type AccountService interface {
	CreateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, id uuid.UUID) (*Account, error)
	GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	GetAccountByNumber(ctx context.Context, accountNumber string) (*Account, error)
}
