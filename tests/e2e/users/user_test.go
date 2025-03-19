package users

import (
	"encoding/json"
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/adil-faiyaz98/money-pulse/internal/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test")
	}

	server := testutil.NewTestServer(t)
	ctx := testutil.CreateTestContext(t)

	t.Run("Create User", func(t *testing.T) {
		user := &users.User{
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password123",
		}

		w := server.SendRequest("POST", "/api/v1/users", user)
		assert.Equal(t, 201, w.Code)

		var response struct {
			Data users.User `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.Data.ID)
		assert.Equal(t, user.Email, response.Data.Email)
		assert.Equal(t, user.FirstName, response.Data.FirstName)
		assert.Equal(t, user.LastName, response.Data.LastName)
	})

	t.Run("Get User", func(t *testing.T) {
		// First create a user
		user := &users.User{
			Email:     "test2@example.com",
			FirstName: "Jane",
			LastName:  "Doe",
			Password:  "password123",
		}

		w := server.SendRequest("POST", "/api/v1/users", user)
		assert.Equal(t, 201, w.Code)

		var createResponse struct {
			Data users.User `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&createResponse)
		require.NoError(t, err)

		// Then get the user
		w = server.SendRequest("GET", "/api/v1/users/"+createResponse.Data.ID.String(), nil)
		assert.Equal(t, 200, w.Code)

		var getResponse struct {
			Data users.User `json:"data"`
		}
		err = json.NewDecoder(w.Body).Decode(&getResponse)
		require.NoError(t, err)

		assert.Equal(t, createResponse.Data.ID, getResponse.Data.ID)
		assert.Equal(t, createResponse.Data.Email, getResponse.Data.Email)
		assert.Equal(t, createResponse.Data.FirstName, getResponse.Data.FirstName)
		assert.Equal(t, createResponse.Data.LastName, getResponse.Data.LastName)
	})

	t.Run("Get User by Email", func(t *testing.T) {
		// Get user by email
		w := server.SendRequest("GET", "/api/v1/users/by-email/test@example.com", nil)
		assert.Equal(t, 200, w.Code)

		var response struct {
			Data users.User `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "test@example.com", response.Data.Email)
		assert.Equal(t, "John", response.Data.FirstName)
		assert.Equal(t, "Doe", response.Data.LastName)
	})

	t.Run("Update User", func(t *testing.T) {
		// First create a user
		user := &users.User{
			Email:     "test3@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password123",
		}

		w := server.SendRequest("POST", "/api/v1/users", user)
		assert.Equal(t, 201, w.Code)

		var createResponse struct {
			Data users.User `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&createResponse)
		require.NoError(t, err)

		// Update the user
		updateUser := &users.User{
			ID:        createResponse.Data.ID,
			Email:     createResponse.Data.Email,
			FirstName: "Jane",
			LastName:  "Doe",
		}

		w = server.SendRequest("PUT", "/api/v1/users/"+createResponse.Data.ID.String(), updateUser)
		assert.Equal(t, 200, w.Code)

		var updateResponse struct {
			Data users.User `json:"data"`
		}
		err = json.NewDecoder(w.Body).Decode(&updateResponse)
		require.NoError(t, err)

		assert.Equal(t, "Jane", updateResponse.Data.FirstName)
	})

	t.Run("Delete User", func(t *testing.T) {
		// First create a user
		user := &users.User{
			Email:     "test4@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password123",
		}

		w := server.SendRequest("POST", "/api/v1/users", user)
		assert.Equal(t, 201, w.Code)

		var createResponse struct {
			Data users.User `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&createResponse)
		require.NoError(t, err)

		// Delete the user
		w = server.SendRequest("DELETE", "/api/v1/users/"+createResponse.Data.ID.String(), nil)
		assert.Equal(t, 204, w.Code)

		// Try to get the deleted user
		w = server.SendRequest("GET", "/api/v1/users/"+createResponse.Data.ID.String(), nil)
		assert.Equal(t, 404, w.Code)
	})
}
