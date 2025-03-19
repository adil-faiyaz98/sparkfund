package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeSavings    AccountType = "SAVINGS"
	AccountTypeChecking   AccountType = "CHECKING"
	AccountTypeInvestment AccountType = "INVESTMENT"
	AccountTypeCredit     AccountType = "CREDIT"
)

type Account struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Type          AccountType
	Name          string
	AccountNumber string
	Balance       float64
	Currency      string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewAccount(userID uuid.UUID, accountType AccountType, name string) (*Account, error) {
	if name == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}

	if !isValidAccountType(accountType) {
		return nil, fmt.Errorf("invalid account type: %s", accountType)
	}

	now := time.Now()
	return &Account{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      accountType,
		Name:      name,
		Balance:   0,
		Currency:  "USD",
		Status:    "ACTIVE",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func isValidAccountType(accountType AccountType) bool {
	switch accountType {
	case AccountTypeSavings,
		AccountTypeChecking,
		AccountTypeInvestment,
		AccountTypeCredit:
		return true
	default:
		return false
	}
}
