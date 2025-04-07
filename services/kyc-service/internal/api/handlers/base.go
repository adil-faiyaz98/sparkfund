package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/sparkfund/kyc-service/internal/security"
    "github.com/sparkfund/kyc-service/internal/ai/validation"
)

// Handler base struct to manage common dependencies
type BaseHandler struct {
    authConfig       security.AuthConfig
    validationConfig validation.Config
    encryptionConfig security.EncryptionConfig
    auditLogger      *security.AuditLogger
    validationRules  *validation.ValidationRules
}

// NewBaseHandler creates a new base handler with common dependencies
func NewBaseHandler(
    authConfig security.AuthConfig,
    validationConfig validation.Config,
    encryptionConfig security.EncryptionConfig,
    auditLogger *security.AuditLogger,
    validationRules *validation.ValidationRules,
) *BaseHandler {
    return &BaseHandler{
        authConfig:       authConfig,
        validationConfig: validationConfig,
        encryptionConfig: encryptionConfig,
        auditLogger:      auditLogger,
        validationRules:  validationRules,
    }
}

// Common error responses
func respondWithError(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{
        "error": message,
    })
}

// Common success response
func respondWithSuccess(c *gin.Context, data interface{}) {
    c.JSON(200, gin.H{
        "data": data,
    })
}