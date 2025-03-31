package security

import (
	"context"
	"time"
)

// TransactionStore defines the interface for transaction storage
type TransactionStore interface {
	// GetTransactionsByTimeRange retrieves transactions within a time range
	GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*Transaction, error)

	// GetTransaction retrieves a specific transaction
	GetTransaction(ctx context.Context, transactionID string) (*Transaction, error)

	// SaveTransaction saves a new transaction
	SaveTransaction(ctx context.Context, transaction *Transaction) error

	// UpdateTransaction updates an existing transaction
	UpdateTransaction(ctx context.Context, transaction *Transaction) error

	// GetUserTransactions retrieves all transactions for a user
	GetUserTransactions(ctx context.Context, userID string) ([]*Transaction, error)
}

// MFAStore defines the interface for MFA storage
type MFAStore interface {
	SaveFactor(ctx context.Context, userID string, factor *MFAFactor) error
	GetFactors(ctx context.Context, userID string) ([]*MFAFactor, error)
	UpdateFactor(ctx context.Context, userID string, factor *MFAFactor) error
	DeleteFactor(ctx context.Context, userID string, factorID string) error
}

// RBACStore defines the interface for RBAC storage
type RBACStore interface {
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	GetRolePermissions(ctx context.Context, role string) ([]string, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	UpdateUserRoles(ctx context.Context, userID string, roles []string) error
}

// Transaction represents a financial transaction
type Transaction struct {
	ID            string
	UserID        string
	Amount        float64
	Currency      string
	Timestamp     time.Time
	Status        string
	RiskLevel     string
	RiskScore     float64
	Location      Location
	Device        DeviceInfo
	RecipientID   string
	RecipientName string
	Description   string
	Category      string
	Flags         []string
}

// Location represents transaction location information
type Location struct {
	Country    string
	City       string
	IP         string
	Latitude   float64
	Longitude  float64
	IsVPN      bool
	IsProxy    bool
	Confidence float64
}

// DeviceInfo represents device information
type DeviceInfo struct {
	DeviceID      string
	DeviceType    string
	OS            string
	Browser       string
	UserAgent     string
	IsKnownDevice bool
	FirstSeen     time.Time
	LastSeen      time.Time
}
