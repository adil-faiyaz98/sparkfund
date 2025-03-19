package testutil

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// TestServer represents a test HTTP server
type TestServer struct {
	Router *gin.Engine
	T      *testing.T
}

// NewTestServer creates a new test server
func NewTestServer(t *testing.T) *TestServer {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return &TestServer{
		Router: router,
		T:      t,
	}
}

// SendRequest sends an HTTP request to the test server
func (ts *TestServer) SendRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(ts.T, err)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.Router.ServeHTTP(w, req)
	return w
}

// DecodeResponse decodes a JSON response into the given interface
func (ts *TestServer) DecodeResponse(w *httptest.ResponseRecorder, v interface{}) {
	err := json.NewDecoder(w.Body).Decode(v)
	require.NoError(ts.T, err)
}

// AssertResponse asserts that the response has the expected status code and body
func (ts *TestServer) AssertResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	require.Equal(ts.T, expectedStatus, w.Code)
	if expectedBody != nil {
		var actualBody interface{}
		ts.DecodeResponse(w, &actualBody)
		require.Equal(ts.T, expectedBody, actualBody)
	}
}
