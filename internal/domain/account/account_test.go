package account

import (
	"testing"
	"time"
)

// Unit test example
func TestAccount_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	tests := []struct {
		name        string
		accountName string
		balance     float64
		wantErr     bool
	}{
		{
			name:        "valid account",
			accountName: "Savings",
			balance:     1000.0,
			wantErr:     false,
		},
		{
			name:        "empty name",
			accountName: "",
			balance:     500.0,
			wantErr:     true,
		},
		{
			name:        "negative balance",
			accountName: "Checking",
			balance:     -100.0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccount(tt.accountName, tt.balance)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if account.Name != tt.accountName {
				t.Errorf("account.Name = %v, want %v", account.Name, tt.accountName)
			}

			if account.Balance != tt.balance {
				t.Errorf("account.Balance = %v, want %v", account.Balance, tt.balance)
			}
		})
	}
}

// Integration test example
func TestIntegration_AccountRepository_Save(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Setup test database connection
	// This would connect to a test database
	// db := setupTestDatabase()
	// defer cleanupTestDatabase()

	// Create a repository
	// repo := NewAccountRepository(db)

	// Option 2: Use blank identifier
	_ = &Account(
		ID:        "test-id",
		Name:      "Test Account",
		Balance:   1000.0,
		CreatedAt: time.Now(),
	}

	// Test saving the account
	// err := repo.Save(account)
	// if err != nil {
	//     t.Fatalf("Failed to save account: %v", err)
	// }

	// Verify the account was saved correctly
	// savedAccount, err := repo.FindByID(account.ID)
	// if err != nil {
	//     t.Fatalf("Failed to retrieve saved account: %v", err)
	// }

	// Compare fields
	// if savedAccount.Name != account.Name {
	//     t.Errorf("account.Name = %v, want %v", savedAccount.Name, account.Name)
	// }
}
