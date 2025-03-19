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
	Name          string
	Type          AccountType
	Balance       float64
	Currency      string
	AccountNumber string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewAccount(userID uuid.UUID, name string, accountType AccountType, currency string) (*Account, error) {
	if name == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}

	if !isValidAccountType(accountType) {
		return nil, fmt.Errorf("invalid account type: %s", accountType)
	}

	if currency == "" {
		return nil, fmt.Errorf("currency cannot be empty")
	}

	now := time.Now()
	return &Account{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          name,
		Type:          accountType,
		Balance:       0,
		Currency:      currency,
		AccountNumber: generateAccountNumber(),
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (a *Account) Update(name string, accountType AccountType) error {
	if name == "" {
		return fmt.Errorf("account name cannot be empty")
	}

	if !isValidAccountType(accountType) {
		return fmt.Errorf("invalid account type: %s", accountType)
	}

	a.Name = name
	a.Type = accountType
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Account) UpdateBalance(amount float64) error {
	if amount < 0 && a.Balance+amount < 0 {
		return fmt.Errorf("insufficient funds")
	}
	a.Balance += amount
	a.UpdatedAt = time.Now()
	return nil
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

func generateAccountNumber() string {
	return fmt.Sprintf("ACC%s", uuid.New().String()[:8])
}
