package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/sparkfund/kyc-service/internal/ai/validation"
)

type ValidationMiddleware struct {
    config validation.Config
}

func NewValidationMiddleware(config validation.Config) *ValidationMiddleware {
    return &ValidationMiddleware{
        config: config,
    }
}

// ValidateRequest validates incoming requests against configured rules
func (m *ValidationMiddleware) ValidateRequest() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Validate request size
        if c.Request.ContentLength > m.config.MaxRequestSize {
            c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
                "error": "Request too large",
            })
            return
        }

        // Validate content type
        contentType := c.GetHeader("Content-Type")
        if !isValidContentType(contentType) {
            c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
                "error": "Unsupported media type",
            })
            return
        }

        c.Next()
    }
}

// ValidateDocumentUpload validates document uploads
func (m *ValidationMiddleware) ValidateDocumentUpload() gin.HandlerFunc {
    return func(c *gin.Context) {
        file, err := c.FormFile("document")
        if err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error": "Invalid document upload",
            })
            return
        }

        // Validate file size
        if file.Size > m.config.MaxDocumentSize {
            c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
                "error": "Document too large",
            })
            return
        }

        // Validate file type
        if !isAllowedDocumentType(file.Header.Get("Content-Type"), m.config.AllowedDocumentTypes) {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error": "Document type not allowed",
            })
            return
        }

        c.Next()
    }
}

// Helper functions
func isValidContentType(contentType string) bool {
    validTypes := []string{
        "application/json",
        "multipart/form-data",
        "application/x-www-form-urlencoded",
    }

    for _, valid := range validTypes {
        if contentType == valid {
            return true
        }
    }
    return false
}

func isAllowedDocumentType(fileType string, allowedTypes []string) bool {
    for _, allowed := range allowedTypes {
        if fileType == allowed {
            return true
        }
    }
    return false
}