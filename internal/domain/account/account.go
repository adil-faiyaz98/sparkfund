package account

import (
	"fmt"
	"time"
)

type Account struct {
	ID        string
	Name      string
	Balance   float64
	CreatedAt time.Time
}

func NewAccount(name string, balance float64) (*Account, error) {
	if name == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}
	if balance < 0 {
		return nil, fmt.Errorf("balance cannot be negative")
	}

	return &Account{
		Name:      name,
		Balance:   balance,
		CreatedAt: time.Now(),
	}, nil
}
