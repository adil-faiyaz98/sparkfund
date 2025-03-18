package transactions

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "DEPOSIT"
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"
	TransactionTypeTransfer   TransactionType = "TRANSFER"
	TransactionTypePayment    TransactionType = "PAYMENT"
	TransactionTypeInterest   TransactionType = "INTEREST"
	TransactionTypeFee        TransactionType = "FEE"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
)

type Transaction struct {
	ID                 uuid.UUID         `json:"id" gorm:"type:uuid;primary_key"`
	UserID             uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	AccountID          uuid.UUID         `json:"account_id" gorm:"type:uuid;not null"`
	Type               TransactionType   `json:"type" gorm:"type:varchar(20);not null"`
	Status             TransactionStatus `json:"status" gorm:"type:varchar(20);not null;default:'PENDING'"`
	Amount             float64           `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency           string            `json:"currency" gorm:"type:varchar(3);not null"`
	Description        string            `json:"description" gorm:"type:text"`
	Category           string            `json:"category" gorm:"type:varchar(50)"`
	Tags               []string          `json:"tags" gorm:"type:text[]"`
	Metadata           string            `json:"metadata" gorm:"type:jsonb"`
	SourceAccount      *uuid.UUID        `json:"source_account" gorm:"type:uuid"`
	DestinationAccount *uuid.UUID        `json:"destination_account" gorm:"type:uuid"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
	CompletedAt        *time.Time        `json:"completed_at"`
	FailedAt           *time.Time        `json:"failed_at"`
	CancelledAt        *time.Time        `json:"cancelled_at"`
}

type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByID(id uuid.UUID) (*Transaction, error)
	GetByUserID(userID uuid.UUID) ([]*Transaction, error)
	GetByAccountID(accountID uuid.UUID) ([]*Transaction, error)
	Update(transaction *Transaction) error
	Delete(id uuid.UUID) error
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	GetTransaction(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]*Transaction, error)
	GetAccountTransactions(ctx context.Context, accountID uuid.UUID) ([]*Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status TransactionStatus, err error) error
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
}
