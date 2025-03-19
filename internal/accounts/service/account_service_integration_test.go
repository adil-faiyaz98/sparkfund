package service

import (
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/accounts"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Setup test database
	testDB := testutil.NewTestDB(t)
	defer testDB.Close(t)

	// Create repository
	repo := NewPostgresAccountRepository(testDB.DB)
	service := NewAccountService(repo)

	// Create test context
	ctx := testutil.CreateTestContext(t)

	// Test data
	userID := uuid.New()
	testAccount := &accounts.Account{
		UserID:  userID,
		Name:    "Test Account",
		Type:    accounts.AccountTypeSavings,
		Balance: 1000.0,
	}

	t.Run("Create and Retrieve Account", func(t *testing.T) {
		// Create account
		err := service.CreateAccount(ctx, testAccount)
		require.NoError(t, err)
		require.NotEmpty(t, testAccount.ID)
		require.NotEmpty(t, testAccount.AccountNumber)

		// Retrieve account
		account, err := service.GetAccount(ctx, testAccount.ID)
		require.NoError(t, err)
		assert.Equal(t, testAccount.ID, account.ID)
		assert.Equal(t, testAccount.Name, account.Name)
		assert.Equal(t, testAccount.Type, account.Type)
		assert.Equal(t, testAccount.Balance, account.Balance)
	})

	t.Run("Get User Accounts", func(t *testing.T) {
		// Create another account for the same user
		anotherAccount := &accounts.Account{
			UserID:  userID,
			Name:    "Another Account",
			Type:    accounts.AccountTypeChecking,
			Balance: 2000.0,
		}
		err := service.CreateAccount(ctx, anotherAccount)
		require.NoError(t, err)

		// Get all accounts for the user
		accounts, err := service.GetUserAccounts(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, accounts, 2)

		// Verify account details
		accountMap := make(map[uuid.UUID]*accounts.Account)
		for _, acc := range accounts {
			accountMap[acc.ID] = acc
		}

		assert.Contains(t, accountMap, testAccount.ID)
		assert.Contains(t, accountMap, anotherAccount.ID)
	})

	t.Run("Update Account", func(t *testing.T) {
		// Update account
		testAccount.Name = "Updated Account"
		testAccount.Type = accounts.AccountTypeChecking
		err := service.UpdateAccount(ctx, testAccount)
		require.NoError(t, err)

		// Verify update
		updated, err := service.GetAccount(ctx, testAccount.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Account", updated.Name)
		assert.Equal(t, accounts.AccountTypeChecking, updated.Type)
	})

	t.Run("Delete Account", func(t *testing.T) {
		// Delete account
		err := service.DeleteAccount(ctx, testAccount.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = service.GetAccount(ctx, testAccount.ID)
		assert.Error(t, err)
	})

	t.Run("Get Account by Number", func(t *testing.T) {
		// Create a new account
		newAccount := &accounts.Account{
			UserID:  userID,
			Name:    "New Account",
			Type:    accounts.AccountTypeSavings,
			Balance: 3000.0,
		}
		err := service.CreateAccount(ctx, newAccount)
		require.NoError(t, err)

		// Get account by number
		account, err := service.GetAccountByNumber(ctx, newAccount.AccountNumber)
		require.NoError(t, err)
		assert.Equal(t, newAccount.ID, account.ID)
		assert.Equal(t, newAccount.Name, account.Name)
		assert.Equal(t, newAccount.Type, account.Type)
		assert.Equal(t, newAccount.Balance, account.Balance)
	})
}
