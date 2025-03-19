package accounts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/accounts"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestServer struct {
	router *gin.Engine
	t      *testing.T
}

func NewTestServer(t *testing.T) *TestServer {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup test database
	testDB := testutil.NewTestDB(t)
	t.Cleanup(func() {
		testDB.Close(t)
	})

	// Create repository and service
	repo := NewPostgresAccountRepository(testDB.DB)
	service := NewAccountService(repo)

	// Create handler
	handler := NewAccountHandler(service)

	// Register routes
	handler.RegisterRoutes(router)

	return &TestServer{
		router: router,
		t:      t,
	}
}

func (ts *TestServer) sendRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(ts.t, err)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)
	return w
}

func TestAccountAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test")
	}

	server := NewTestServer(t)
	ctx := testutil.CreateTestContext(t)

	t.Run("Create Account", func(t *testing.T) {
		account := &accounts.Account{
			UserID:  uuid.New(),
			Name:    "Test Account",
			Type:    accounts.AccountTypeSavings,
			Balance: 1000.0,
		}

		w := server.sendRequest(http.MethodPost, "/api/v1/accounts", account)
		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Data accounts.Account `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.Data.ID)
		assert.NotEmpty(t, response.Data.AccountNumber)
		assert.Equal(t, account.Name, response.Data.Name)
		assert.Equal(t, account.Type, response.Data.Type)
		assert.Equal(t, account.Balance, response.Data.Balance)
	})

	t.Run("Get Account", func(t *testing.T) {
		// First create an account
		account := &accounts.Account{
			UserID:  uuid.New(),
			Name:    "Test Account",
			Type:    accounts.AccountTypeSavings,
			Balance: 1000.0,
		}

		w := server.sendRequest(http.MethodPost, "/api/v1/accounts", account)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse struct {
			Data accounts.Account `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&createResponse)
		require.NoError(t, err)

		// Then get the account
		w = server.sendRequest(http.MethodGet, "/api/v1/accounts/"+createResponse.Data.ID.String(), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var getResponse struct {
			Data accounts.Account `json:"data"`
		}
		err = json.NewDecoder(w.Body).Decode(&getResponse)
		require.NoError(t, err)

		assert.Equal(t, createResponse.Data.ID, getResponse.Data.ID)
		assert.Equal(t, createResponse.Data.Name, getResponse.Data.Name)
		assert.Equal(t, createResponse.Data.Type, getResponse.Data.Type)
		assert.Equal(t, createResponse.Data.Balance, getResponse.Data.Balance)
	})

	t.Run("Get User Accounts", func(t *testing.T) {
		userID := uuid.New()

		// Create multiple accounts for the same user
		accounts := []*accounts.Account{
			{
				UserID:  userID,
				Name:    "Account 1",
				Type:    accounts.AccountTypeSavings,
				Balance: 1000.0,
			},
			{
				UserID:  userID,
				Name:    "Account 2",
				Type:    accounts.AccountTypeChecking,
				Balance: 2000.0,
			},
		}

		for _, account := range accounts {
			w := server.sendRequest(http.MethodPost, "/api/v1/accounts", account)
			assert.Equal(t, http.StatusCreated, w.Code)
		}

		// Get all accounts for the user
		w := server.sendRequest(http.MethodGet, "/api/v1/users/"+userID.String()+"/accounts", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data []accounts.Account `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Data, 2)
		accountMap := make(map[string]accounts.Account)
		for _, acc := range response.Data {
			accountMap[acc.Name] = acc
		}

		assert.Contains(t, accountMap, "Account 1")
		assert.Contains(t, accountMap, "Account 2")
	})

	t.Run("Update Account", func(t *testing.T) {
		// First create an account
		account := &accounts.Account{
			UserID:  uuid.New(),
			Name:    "Test Account",
			Type:    accounts.AccountTypeSavings,
			Balance: 1000.0,
		}

		w := server.sendRequest(http.MethodPost, "/api/v1/accounts", account)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse struct {
			Data accounts.Account `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&createResponse)
		require.NoError(t, err)

		// Update the account
		updateAccount := &accounts.Account{
			ID:      createResponse.Data.ID,
			UserID:  createResponse.Data.UserID,
			Name:    "Updated Account",
			Type:    accounts.AccountTypeChecking,
			Balance: 2000.0,
		}

		w = server.sendRequest(http.MethodPut, "/api/v1/accounts/"+createResponse.Data.ID.String(), updateAccount)
		assert.Equal(t, http.StatusOK, w.Code)

		var updateResponse struct {
			Data accounts.Account `json:"data"`
		}
		err = json.NewDecoder(w.Body).Decode(&updateResponse)
		require.NoError(t, err)

		assert.Equal(t, "Updated Account", updateResponse.Data.Name)
		assert.Equal(t, accounts.AccountTypeChecking, updateResponse.Data.Type)
		assert.Equal(t, 2000.0, updateResponse.Data.Balance)
	})

	t.Run("Delete Account", func(t *testing.T) {
		// First create an account
		account := &accounts.Account{
			UserID:  uuid.New(),
			Name:    "Test Account",
			Type:    accounts.AccountTypeSavings,
			Balance: 1000.0,
		}

		w := server.sendRequest(http.MethodPost, "/api/v1/accounts", account)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse struct {
			Data accounts.Account `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&createResponse)
		require.NoError(t, err)

		// Delete the account
		w = server.sendRequest(http.MethodDelete, "/api/v1/accounts/"+createResponse.Data.ID.String(), nil)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Try to get the deleted account
		w = server.sendRequest(http.MethodGet, "/api/v1/accounts/"+createResponse.Data.ID.String(), nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Get Account by Number", func(t *testing.T) {
		// First create an account
		account := &accounts.Account{
			UserID:  uuid.New(),
			Name:    "Test Account",
			Type:    accounts.AccountTypeSavings,
			Balance: 1000.0,
		}

		w := server.sendRequest(http.MethodPost, "/api/v1/accounts", account)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse struct {
			Data accounts.Account `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&createResponse)
		require.NoError(t, err)

		// Get account by number
		w = server.sendRequest(http.MethodGet, "/api/v1/accounts/by-number/"+createResponse.Data.AccountNumber, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var getResponse struct {
			Data accounts.Account `json:"data"`
		}
		err = json.NewDecoder(w.Body).Decode(&getResponse)
		require.NoError(t, err)

		assert.Equal(t, createResponse.Data.ID, getResponse.Data.ID)
		assert.Equal(t, createResponse.Data.AccountNumber, getResponse.Data.AccountNumber)
		assert.Equal(t, createResponse.Data.Name, getResponse.Data.Name)
		assert.Equal(t, createResponse.Data.Type, getResponse.Data.Type)
		assert.Equal(t, createResponse.Data.Balance, getResponse.Data.Balance)
	})
}
