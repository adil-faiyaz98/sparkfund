package application

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SQLInjectionTestSuite struct {
	server *httptest.Server
}

func NewSQLInjectionTestSuite() *SQLInjectionTestSuite {
	return &SQLInjectionTestSuite{}
}

func (s *SQLInjectionTestSuite) TestBasicSQLInjection(t *testing.T) {
	payloads := []string{
		"' OR '1'='1",
		"' OR 1=1--",
		"'; DROP TABLE users; --",
		"' UNION SELECT * FROM users; --",
		"' OR 'x'='x",
		"' OR '1'='1' --",
		"' OR '1'='1' #",
		"' OR '1'='1' /*",
	}

	for _, payload := range payloads {
		t.Run(payload, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/investments?user_id="+payload, nil)
			w := httptest.NewRecorder()
			
			handler := InvestmentHandler{}
			handler.GetInvestment(w, req)
			
			assert.NotEqual(t, http.StatusOK, w.Code, "SQL injection attempt should be blocked")
		})
	}
}

func (s *SQLInjectionTestSuite) TestBlindSQLInjection(t *testing.T) {
	payloads := []string{
		"' AND 1=1 AND '1'='1",
		"' AND 1=2 AND '1'='1",
		"' AND (SELECT COUNT(*) FROM users)>0 AND '1'='1",
		"' AND (SELECT COUNT(*) FROM users)>100 AND '1'='1",
		"' AND (SELECT LENGTH(password) FROM users WHERE id=1)>0 AND '1'='1",
	}

	for _, payload := range payloads {
		t.Run(payload, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/investments?user_id="+payload, nil)
			w := httptest.NewRecorder()
			
			handler := InvestmentHandler{}
			handler.GetInvestment(w, req)
			
			assert.NotEqual(t, http.StatusOK, w.Code, "Blind SQL injection attempt should be blocked")
		})
	}
}

func (s *SQLInjectionTestSuite) TestTimeBasedSQLInjection(t *testing.T) {
	payloads := []string{
		"' AND (SELECT SLEEP(5)) AND '1'='1",
		"' AND (SELECT pg_sleep(5)) AND '1'='1",
		"' AND (SELECT WAITFOR DELAY '0:0:5') AND '1'='1",
	}

	for _, payload := range payloads {
		t.Run(payload, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/investments?user_id="+payload, nil)
			w := httptest.NewRecorder()
			
			handler := InvestmentHandler{}
			handler.GetInvestment(w, req)
			
			assert.NotEqual(t, http.StatusOK, w.Code, "Time-based SQL injection attempt should be blocked")
		})
	}
}

func (s *SQLInjectionTestSuite) TestErrorBasedSQLInjection(t *testing.T) {
	payloads := []string{
		"' AND (SELECT 1 FROM (SELECT COUNT(*),CONCAT(VERSION(),FLOOR(RAND(0)*2))x FROM INFORMATION_SCHEMA.TABLES GROUP BY x)a) AND '1'='1",
		"' AND (SELECT 1 FROM (SELECT COUNT(*),CONCAT(DATABASE(),FLOOR(RAND(0)*2))x FROM INFORMATION_SCHEMA.TABLES GROUP BY x)a) AND '1'='1",
		"' AND (SELECT 1 FROM (SELECT COUNT(*),CONCAT(TABLE_NAME,FLOOR(RAND(0)*2))x FROM INFORMATION_SCHEMA.TABLES GROUP BY x)a) AND '1'='1",
	}

	for _, payload := range payloads {
		t.Run(payload, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/investments?user_id="+payload, nil)
			w := httptest.NewRecorder()
			
			handler := InvestmentHandler{}
			handler.GetInvestment(w, req)
			
			assert.NotEqual(t, http.StatusOK, w.Code, "Error-based SQL injection attempt should be blocked")
		})
	}
}

func (s *SQLInjectionTestSuite) TestStackedQueries(t *testing.T) {
	payloads := []string{
		"'; INSERT INTO users (username, password) VALUES ('hacker', 'password'); --",
		"'; UPDATE users SET password='hacked' WHERE username='admin'; --",
		"'; DELETE FROM users WHERE username='admin'; --",
		"'; CREATE TABLE hacked (id INT); --",
		"'; DROP TABLE users; --",
	}

	for _, payload := range payloads {
		t.Run(payload, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/investments?user_id="+payload, nil)
			w := httptest.NewRecorder()
			
			handler := InvestmentHandler{}
			handler.GetInvestment(w, req)
			
			assert.NotEqual(t, http.StatusOK, w.Code, "Stacked query SQL injection attempt should be blocked")
		})
	}
}

func (s *SQLInjectionTestSuite) TestPOSTSQLInjection(t *testing.T) {
	payloads := []map[string]interface{}{
		{"user_id": "' OR '1'='1"},
		{"user_id": "'; DROP TABLE users; --"},
		{"user_id": "' UNION SELECT * FROM users; --"},
	}

	for _, payload := range payloads {
		t.Run(payload["user_id"].(string), func(t *testing.T) {
			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/investments", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			handler := InvestmentHandler{}
			handler.CreateInvestment(w, req)
			
			assert.NotEqual(t, http.StatusOK, w.Code, "POST SQL injection attempt should be blocked")
		})
	}
}

func (s *SQLInjectionTestSuite) TestNoSQLInjection(t *testing.T) {
	payloads := []map[string]interface{}{
		{"$gt": ""},
		{"$ne": null},
		{"$regex": ".*"},
		{"$exists": true},
		{"$in": []string{"admin", "user"}},
	}

	for _, payload := range payloads {
		t.Run(fmt.Sprintf("%v", payload), func(t *testing.T) {
			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/investments", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			handler := InvestmentHandler{}
			handler.CreateInvestment(w, req)
			
			assert.NotEqual(t, http.StatusOK, w.Code, "NoSQL injection attempt should be blocked")
		})
	}
} 