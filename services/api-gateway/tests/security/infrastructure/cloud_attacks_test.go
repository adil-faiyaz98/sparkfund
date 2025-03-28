package infrastructure

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestAWSMetadataSSRF tests for SSRF vulnerabilities to AWS metadata API
func TestAWSMetadataSSRF(t *testing.T) {
	t.Run("AWS Metadata SSRF", func(t *testing.T) {
		payloads := []string{
			"http://169.254.169.254/latest/meta-data/",
			"http://[::1]:169254/latest/meta-data/",
			"http://metadata.google.internal/computeMetadata/v1/",
			"http://169.254.169.254/latest/user-data",
			"http://169.254.169.254/latest/dynamic/instance-identity/document",
			"http://169.254.169.254/latest/meta-data/iam/security-credentials/",
			"http://169.254.169.254/latest/meta-data/iam/security-credentials/ROLE_NAME",
		}

		for _, payload := range payloads {
			// Test via URL parameter
			req := httptest.NewRequest("GET", "/api/proxy?url="+payload, nil)
			w := httptest.NewRecorder()
			handler := setupCloudSecureHandler()
			handler.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "SSRF attempt should be blocked: %s", payload)

			// Test via POST body
			jsonData := map[string]string{"url": payload}
			body, _ := json.Marshal(jsonData)
			req = httptest.NewRequest("POST", "/api/proxy", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w = httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "SSRF attempt in POST body should be blocked: %s", payload)
		}
	})
}

// TestKubernetesRBAC tests for Kubernetes RBAC vulnerabilities
func TestKubernetesRBAC(t *testing.T) {
	t.Run("Kubernetes RBAC Bypass", func(t *testing.T) {
		endpoints := []string{
			"/api/v1/namespaces",
			"/api/v1/pods",
			"/apis/apps/v1/deployments",
			"/api/v1/secrets",
			"/api/v1/configmaps",
			"/apis/rbac.authorization.k8s.io/v1/roles",
			"/apis/rbac.authorization.k8s.io/v1/rolebindings",
			"/api/v1/services",
			"/api/v1/endpoints",
			"/apis/networking.k8s.io/v1/networkpolicies",
		}

		for _, endpoint := range endpoints {
			req := httptest.NewRequest("GET", "/api/k8s"+endpoint, nil)
			w := httptest.NewRecorder()
			handler := setupCloudSecureHandler()
			handler.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "Kubernetes API access attempt should be blocked: %s", endpoint)
		}
	})
}

// TestS3BucketExposure tests for S3 bucket exposure vulnerabilities
func TestS3BucketExposure(t *testing.T) {
	t.Run("S3 Bucket Access", func(t *testing.T) {
		payloads := []string{
			"https://s3.amazonaws.com/company-private-bucket/",
			"https://company-private-bucket.s3.amazonaws.com/",
			"https://s3.amazonaws.com/company-private-bucket/*",
			"https://s3.amazonaws.com/company-private-bucket/backup/",
			"https://s3.amazonaws.com/company-private-bucket/config/",
			"https://s3.amazonaws.com/company-private-bucket/logs/",
			"https://s3.amazonaws.com/company-private-bucket/secrets/",
		}

		for _, payload := range payloads {
			// Test direct access
			req := httptest.NewRequest("GET", "/api/files?path="+payload, nil)
			w := httptest.NewRecorder()
			handler := setupCloudSecureHandler()
			handler.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "S3 bucket access attempt should be blocked: %s", payload)

			// Test via POST
			jsonData := map[string]string{"path": payload}
			body, _ := json.Marshal(jsonData)
			req = httptest.NewRequest("POST", "/api/files", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w = httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "S3 bucket access attempt in POST should be blocked: %s", payload)
		}
	})
}

// TestRateLimiting tests for rate limiting vulnerabilities
func TestRateLimiting(t *testing.T) {
	t.Run("Rate Limiting", func(t *testing.T) {
		handler := setupCloudSecureHandler()
		endpoints := []string{
			"/api/auth/login",
			"/api/auth/register",
			"/api/investments",
			"/api/transactions",
		}

		for _, endpoint := range endpoints {
			// Test burst requests
			for i := 0; i < 100; i++ {
				req := httptest.NewRequest("GET", endpoint, nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)

				if i > 50 { // Assuming rate limit is 50 requests per minute
					assert.Equal(t, http.StatusTooManyRequests, w.Code, "Rate limiting should be enforced for burst requests")
				}
				time.Sleep(10 * time.Millisecond)
			}

			// Test sustained requests
			for i := 0; i < 10; i++ {
				req := httptest.NewRequest("GET", endpoint, nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)
				time.Sleep(100 * time.Millisecond)
			}

			// Verify rate limit headers
			req := httptest.NewRequest("GET", endpoint, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			assert.Contains(t, w.Header().Get("X-RateLimit-Limit"), "50", "Rate limit header should be present")
			assert.Contains(t, w.Header().Get("X-RateLimit-Remaining"), "0", "Rate limit remaining should be 0")
		}
	})
}

// setupCloudSecureHandler implements cloud-specific security checks
func setupCloudSecureHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement cloud-specific security checks here
		// 1. SSRF protection for cloud metadata
		// 2. Kubernetes RBAC validation
		// 3. S3 bucket access control
		// 4. Rate limiting
		w.WriteHeader(http.StatusOK)
	})
} 