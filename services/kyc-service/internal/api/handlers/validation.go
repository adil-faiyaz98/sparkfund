package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/sparkfund/kyc-service/internal/ai/validation"
)

type ValidationHandler struct {
    *BaseHandler
    validator *validation.Validator
}

func NewValidationHandler(base *BaseHandler, validator *validation.Validator) *ValidationHandler {
    return &ValidationHandler{
        BaseHandler: base,
        validator:   validator,
    }
}

// HandleDocumentValidation validates KYC documents
func (h *ValidationHandler) HandleDocumentValidation(c *gin.Context) {
    var req struct {
        DocumentType string `json:"document_type" binding:"required"`
        DocumentData string `json:"document_data" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        respondWithError(c, http.StatusBadRequest, "Invalid request format")
        return
    }

    result, err := h.validator.ValidateDocument(c.Request.Context(), req.DocumentType, req.DocumentData)
    if err != nil {
        respondWithError(c, http.StatusInternalServerError, "Validation failed")
        return
    }

    respondWithSuccess(c, result)
}

// HandleBiometricValidation validates biometric data
func (h *ValidationHandler) HandleBiometricValidation(c *gin.Context) {
    var req struct {
        BiometricType string `json:"biometric_type" binding:"required"`
        BiometricData string `json:"biometric_data" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        respondWithError(c, http.StatusBadRequest, "Invalid request format")
        return
    }

    result, err := h.validator.ValidateBiometric(c.Request.Context(), req.BiometricType, req.BiometricData)
    if err != nil {
        respondWithError(c, http.StatusInternalServerError, "Biometric validation failed")
        return
    }

    respondWithSuccess(c, result)
}