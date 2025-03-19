package service

import (
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/adil-faiyaz98/money-pulse/internal/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Setup test database
	testDB := testutil.NewTestDB(t)
	defer testDB.Close(t)

	// Create repository
	repo := NewPostgresUserRepository(testDB.DB)
	service := NewUserService(repo)

	// Create test context
	ctx := testutil.CreateTestContext(t)

	// Test data
	testUser := &users.User{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password123",
	}

	t.Run("Create and Retrieve User", func(t *testing.T) {
		// Create user
		err := service.CreateUser(ctx, testUser)
		require.NoError(t, err)
		require.NotEmpty(t, testUser.ID)

		// Retrieve user
		user, err := service.GetUser(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.FirstName, user.FirstName)
		assert.Equal(t, testUser.LastName, user.LastName)
	})

	t.Run("Get User by Email", func(t *testing.T) {
		// Get user by email
		user, err := service.GetUserByEmail(ctx, testUser.Email)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.FirstName, user.FirstName)
		assert.Equal(t, testUser.LastName, user.LastName)
	})

	t.Run("Update User", func(t *testing.T) {
		// Update user
		testUser.FirstName = "Jane"
		err := service.UpdateUser(ctx, testUser)
		require.NoError(t, err)

		// Verify update
		updated, err := service.GetUser(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, "Jane", updated.FirstName)
	})

	t.Run("Delete User", func(t *testing.T) {
		// Delete user
		err := service.DeleteUser(ctx, testUser.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = service.GetUser(ctx, testUser.ID)
		assert.Error(t, err)
	})
}
