package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSecret = "test-secret-key"
)

type JWTTestSuite struct {
	secret []byte
}

func NewJWTTestSuite() *JWTTestSuite {
	return &JWTTestSuite{
		secret: []byte(testSecret),
	}
}

func (s *JWTTestSuite) TestValidToken(t *testing.T) {
	token := s.generateValidToken(t)
	req := httptest.NewRequest("GET", "/api/v1/investments", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	w := httptest.NewRecorder()
	handler := JWTAuthMiddleware(s.secret)
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func (s *JWTTestSuite) TestExpiredToken(t *testing.T) {
	token := s.generateExpiredToken(t)
	req := httptest.NewRequest("GET", "/api/v1/investments", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	w := httptest.NewRecorder()
	handler := JWTAuthMiddleware(s.secret)
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func (s *JWTTestSuite) TestInvalidSignature(t *testing.T) {
	token := s.generateTokenWithInvalidSignature(t)
	req := httptest.NewRequest("GET", "/api/v1/investments", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	w := httptest.NewRecorder()
	handler := JWTAuthMiddleware(s.secret)
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func (s *JWTTestSuite) TestMissingToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/investments", nil)
	
	w := httptest.NewRecorder()
	handler := JWTAuthMiddleware(s.secret)
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func (s *JWTTestSuite) TestTokenTampering(t *testing.T) {
	token := s.generateValidToken(t)
	tamperedToken := s.tamperWithToken(token)
	
	req := httptest.NewRequest("GET", "/api/v1/investments", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tamperedToken))
	
	w := httptest.NewRecorder()
	handler := JWTAuthMiddleware(s.secret)
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func (s *JWTTestSuite) TestAlgorithmNone(t *testing.T) {
	token := s.generateTokenWithAlgorithmNone(t)
	req := httptest.NewRequest("GET", "/api/v1/investments", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	w := httptest.NewRecorder()
	handler := JWTAuthMiddleware(s.secret)
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func (s *JWTTestSuite) TestWeakSecret(t *testing.T) {
	weakSecret := "weak-secret"
	token := s.generateTokenWithWeakSecret(t, weakSecret)
	req := httptest.NewRequest("GET", "/api/v1/investments", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	w := httptest.NewRecorder()
	handler := JWTAuthMiddleware(s.secret)
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Helper functions
func (s *JWTTestSuite) generateValidToken(t *testing.T) string {
	claims := jwt.MapClaims{
		"sub":   "test-user",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"roles": []string{"user"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	require.NoError(t, err)
	return tokenString
}

func (s *JWTTestSuite) generateExpiredToken(t *testing.T) string {
	claims := jwt.MapClaims{
		"sub":   "test-user",
		"exp":   time.Now().Add(-time.Hour).Unix(),
		"iat":   time.Now().Add(-2 * time.Hour).Unix(),
		"roles": []string{"user"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	require.NoError(t, err)
	return tokenString
}

func (s *JWTTestSuite) generateTokenWithInvalidSignature(t *testing.T) string {
	claims := jwt.MapClaims{
		"sub":   "test-user",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"roles": []string{"user"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("wrong-secret"))
	require.NoError(t, err)
	return tokenString
}

func (s *JWTTestSuite) tamperWithToken(token string) string {
	parts := splitToken(token)
	payload := parts[1]
	
	// Decode and modify the payload
	decodedPayload, _ := base64.RawURLEncoding.DecodeString(payload)
	claims := make(map[string]interface{})
	jwt.Unmarshal(decodedPayload, &claims)
	
	// Modify the claims
	claims["roles"] = []string{"admin"}
	
	// Re-encode the payload
	modifiedPayload, _ := jwt.Marshal(claims)
	parts[1] = base64.RawURLEncoding.EncodeToString(modifiedPayload)
	
	// Reconstruct the token with the original signature
	return joinToken(parts)
}

func (s *JWTTestSuite) generateTokenWithAlgorithmNone(t *testing.T) string {
	claims := jwt.MapClaims{
		"sub":   "test-user",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"roles": []string{"user"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(nil)
	require.NoError(t, err)
	return tokenString
}

func (s *JWTTestSuite) generateTokenWithWeakSecret(t *testing.T, weakSecret string) string {
	claims := jwt.MapClaims{
		"sub":   "test-user",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"roles": []string{"user"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(weakSecret))
	require.NoError(t, err)
	return tokenString
}

func splitToken(token string) []string {
	parts := make([]string, 3)
	copy(parts, strings.Split(token, "."))
	return parts
}

func joinToken(parts []string) string {
	return strings.Join(parts, ".")
} 