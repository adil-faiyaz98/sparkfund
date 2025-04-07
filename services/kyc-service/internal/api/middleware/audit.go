package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/adil-faiyaz98/sparkfund/services/kyc-service/internal/audit"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditMiddleware creates a middleware for audit logging
func AuditMiddleware(auditLogger audit.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip audit logging for certain paths
		if shouldSkipAudit(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Start time
		startTime := time.Now()

		// Get request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Request.Header.Set("X-Request-ID", requestID)
			c.Writer.Header().Set("X-Request-ID", requestID)
		}

		// Get user ID
		userID, _ := c.Get("userID")
		userIDStr := ""
		if userID != nil {
			userIDStr = userID.(string)
		}

		// Read request body
		var requestBody []byte
		var err error
		if c.Request.Body != nil && c.Request.Method != http.MethodGet {
			requestBody, err = io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
				return
			}
			// Restore request body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create response body buffer
		responseBodyWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseBodyWriter

		// Process request
		c.Next()

		// Get response status
		status := c.Writer.Status()

		// Determine action and resource
		action, resource, resourceID := determineActionAndResource(c)

		// Determine status type
		statusType := determineStatusType(status)

		// Parse request body
		var requestParams interface{}
		if len(requestBody) > 0 && isJSONContentType(c.ContentType()) {
			if err := json.Unmarshal(requestBody, &requestParams); err != nil {
				// If we can't parse as JSON, use as string
				requestParams = string(requestBody)
			}
		}

		// Parse response body
		var responseBody interface{}
		if responseBodyWriter.body.Len() > 0 && isJSONContentType(responseBodyWriter.Header().Get("Content-Type")) {
			if err := json.Unmarshal(responseBodyWriter.body.Bytes(), &responseBody); err != nil {
				// If we can't parse as JSON, use as string
				responseBody = responseBodyWriter.body.String()
			}
		}

		// Get error message
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.Last().Error()
		}

		// Create audit event
		event := &audit.Event{
			Timestamp:     startTime,
			UserID:        userIDStr,
			Action:        action,
			Resource:      resource,
			ResourceID:    resourceID,
			Status:        statusType,
			ClientIP:      c.ClientIP(),
			UserAgent:     c.Request.UserAgent(),
			RequestID:     requestID,
			RequestMethod: c.Request.Method,
			RequestPath:   c.Request.URL.Path,
			RequestParams: requestParams,
			ResponseCode:  status,
			ErrorMessage:  errorMessage,
			Metadata: map[string]interface{}{
				"duration_ms": time.Since(startTime).Milliseconds(),
				"query_params": c.Request.URL.Query(),
				"headers": filterHeaders(c.Request.Header),
			},
		}

		// Log audit event
		if err := auditLogger.Log(context.Background(), event); err != nil {
			c.Error(err)
		}
	}
}

// responseBodyWriter is a custom response writer that captures the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// shouldSkipAudit determines if audit logging should be skipped for a path
func shouldSkipAudit(path string) bool {
	// Skip health checks, metrics, and static files
	return strings.HasPrefix(path, "/health") ||
		strings.HasPrefix(path, "/ready") ||
		strings.HasPrefix(path, "/live") ||
		strings.HasPrefix(path, "/metrics") ||
		strings.HasPrefix(path, "/swagger") ||
		strings.HasSuffix(path, ".js") ||
		strings.HasSuffix(path, ".css") ||
		strings.HasSuffix(path, ".png") ||
		strings.HasSuffix(path, ".jpg") ||
		strings.HasSuffix(path, ".ico")
}

// determineActionAndResource determines the action and resource from the request
func determineActionAndResource(c *gin.Context) (audit.ActionType, audit.ResourceType, string) {
	path := c.Request.URL.Path
	method := c.Request.Method
	
	// Extract resource ID from path
	pathParts := strings.Split(path, "/")
	resourceID := ""
	if len(pathParts) > 3 {
		resourceID = pathParts[3]
	}
	
	// Determine resource type
	var resource audit.ResourceType
	if strings.Contains(path, "/users") || strings.Contains(path, "/auth") {
		resource = audit.ResourceUser
	} else if strings.Contains(path, "/documents") {
		resource = audit.ResourceDocument
	} else if strings.Contains(path, "/verifications") {
		resource = audit.ResourceVerification
	} else if strings.Contains(path, "/kyc") {
		resource = audit.ResourceKYC
	} else if strings.Contains(path, "/analyze") {
		resource = audit.ResourceAnalysis
	} else if strings.Contains(path, "/face-match") {
		resource = audit.ResourceFaceMatch
	} else if strings.Contains(path, "/risk") {
		resource = audit.ResourceRiskAnalysis
	} else if strings.Contains(path, "/anomaly") {
		resource = audit.ResourceAnomalyDetection
	} else {
		resource = audit.ResourceType(pathParts[2])
	}
	
	// Determine action type
	var action audit.ActionType
	switch method {
	case http.MethodGet:
		action = audit.ActionRead
	case http.MethodPost:
		if strings.Contains(path, "/login") {
			action = audit.ActionLogin
		} else if strings.Contains(path, "/logout") {
			action = audit.ActionLogout
		} else if strings.Contains(path, "/verify") {
			action = audit.ActionVerify
		} else if strings.Contains(path, "/approve") {
			action = audit.ActionApprove
		} else if strings.Contains(path, "/reject") {
			action = audit.ActionReject
		} else if strings.Contains(path, "/upload") {
			action = audit.ActionUpload
		} else if strings.Contains(path, "/download") {
			action = audit.ActionDownload
		} else if strings.Contains(path, "/analyze") {
			action = audit.ActionAnalyze
		} else {
			action = audit.ActionCreate
		}
	case http.MethodPut:
		action = audit.ActionUpdate
	case http.MethodDelete:
		action = audit.ActionDelete
	default:
		action = audit.ActionType(method)
	}
	
	return action, resource, resourceID
}

// determineStatusType determines the status type from the HTTP status code
func determineStatusType(status int) audit.StatusType {
	switch {
	case status >= 200 && status < 300:
		return audit.StatusSuccess
	case status >= 400 && status < 500:
		if status == http.StatusUnauthorized || status == http.StatusForbidden {
			return audit.StatusDenied
		}
		return audit.StatusFailure
	case status >= 500:
		return audit.StatusError
	default:
		return audit.StatusSuccess
	}
}

// isJSONContentType checks if the content type is JSON
func isJSONContentType(contentType string) bool {
	return strings.Contains(contentType, "application/json")
}

// filterHeaders filters sensitive headers
func filterHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		// Skip sensitive headers
		if strings.ToLower(key) == "authorization" ||
			strings.ToLower(key) == "cookie" ||
			strings.ToLower(key) == "set-cookie" ||
			strings.ToLower(key) == "x-api-key" {
			continue
		}
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}
