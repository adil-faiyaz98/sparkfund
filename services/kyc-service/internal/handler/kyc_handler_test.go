package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sparkfund/kyc-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKYCService struct {
	mock.Mock
}

func (m *MockKYCService) SubmitKYC(userID uuid.UUID, req *model.KYCRequest) (*model.KYCResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.KYCResponse), args.Error(1)
}

func (m *MockKYCService) GetKYCStatus(id uuid.UUID) (*model.KYCResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.KYCResponse), args.Error(1)
}

func (m *MockKYCService) VerifyKYC(id uuid.UUID, verifiedBy uuid.UUID) error {
	args := m.Called(id, verifiedBy)
	return args.Error(0)
}

func (m *MockKYCService) RejectKYC(id uuid.UUID, reason string) error {
	args := m.Called(id, reason)
	return args.Error(0)
}

func (m *MockKYCService) ListPendingKYC() ([]model.KYC, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.KYC), args.Error(1)
}

func setupTest() (*gin.Engine, *MockKYCService) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockKYCService)
	handler := NewKYCHandler(mockService)
	router := gin.New()
	handler.RegisterRoutes(router)
	return router, mockService
}

func TestSubmitKYC(t *testing.T) {
	router, mockService := setupTest()

	t.Run("successful submission", func(t *testing.T) {
		req := &model.KYCRequest{
			FirstName:      "John",
			LastName:       "Doe",
			DateOfBirth:    "1990-01-01",
			Address:        "123 Main St",
			City:           "New York",
			Country:        "US",
			PostalCode:     "10001",
			DocumentType:   "passport",
			DocumentNumber: "AB123456",
			DocumentFront:  "base64_front",
			DocumentBack:   "base64_back",
			SelfieImage:    "base64_selfie",
		}

		resp := &model.KYCResponse{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Status: model.KYCStatusPending,
		}

		mockService.On("SubmitKYC", mock.AnythingOfType("uuid.UUID"), req).Return(resp, nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/kyc", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.KYCResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, resp.ID, response.ID)
		assert.Equal(t, resp.Status, response.Status)
	})

	t.Run("invalid request", func(t *testing.T) {
		req := &model.KYCRequest{} // Empty request

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/kyc", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetKYCStatus(t *testing.T) {
	router, mockService := setupTest()

	t.Run("successful status retrieval", func(t *testing.T) {
		id := uuid.New()
		resp := &model.KYCResponse{
			ID:     id,
			UserID: uuid.New(),
			Status: model.KYCStatusPending,
		}

		mockService.On("GetKYCStatus", id).Return(resp, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/kyc/"+id.String(), nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.KYCResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, resp.ID, response.ID)
		assert.Equal(t, resp.Status, response.Status)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/kyc/invalid-uuid", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestVerifyKYC(t *testing.T) {
	router, mockService := setupTest()

	t.Run("successful verification", func(t *testing.T) {
		id := uuid.New()
		mockService.On("VerifyKYC", id, mock.AnythingOfType("uuid.UUID")).Return(nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/kyc/"+id.String()+"/verify", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/kyc/invalid-uuid/verify", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRejectKYC(t *testing.T) {
	router, mockService := setupTest()

	t.Run("successful rejection", func(t *testing.T) {
		id := uuid.New()
		reason := "Document quality is poor"
		mockService.On("RejectKYC", id, reason).Return(nil)

		body := map[string]string{"reason": reason}
		bodyBytes, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/kyc/"+id.String()+"/reject", bytes.NewBuffer(bodyBytes))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing reason", func(t *testing.T) {
		id := uuid.New()
		body := map[string]string{} // Empty body

		bodyBytes, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/kyc/"+id.String()+"/reject", bytes.NewBuffer(bodyBytes))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListPendingKYC(t *testing.T) {
	router, mockService := setupTest()

	t.Run("successful list retrieval", func(t *testing.T) {
		kycs := []model.KYC{
			{
				ID:     uuid.New(),
				UserID: uuid.New(),
				Status: model.KYCStatusPending,
			},
			{
				ID:     uuid.New(),
				UserID: uuid.New(),
				Status: model.KYCStatusPending,
			},
		}

		mockService.On("ListPendingKYC").Return(kycs, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/kyc/pending", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []model.KYC
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
	})
}
