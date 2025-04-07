package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sparkfund/kyc-service/internal/domain"
	"github.com/sparkfund/kyc-service/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type VerificationHandler struct {
	verificationService service.VerificationService
	logger              *zap.Logger
	auditLogger         service.AuditLogger
}

func NewVerificationHandler(
	verificationService service.VerificationService,
	logger *zap.Logger,
	auditLogger service.AuditLogger,
) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
		logger:              logger,
		auditLogger:         auditLogger,
	}
}

func (h *VerificationHandler) RegisterRoutes(router *gin.Engine) {
	verification := router.Group("/api/v1/verifications")
	{
		verification.POST("/document", h.VerifyDocument)
		verification.POST("/biometric", h.VerifyBiometric)
		verification.GET("/:id", h.GetVerificationStatus)
	}
}

func (h *VerificationHandler) VerifyDocument(c *gin.Context) {
	var req struct {
		DocumentID uuid.UUID `json:"document_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewError("INVALID_REQUEST", "Invalid request format", err.Error()))
		return
	}

	verification, err := h.verificationService.VerifyDocument(c.Request.Context(), req.DocumentID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.auditLogger.Log(c.Request.Context(), "document_verification_initiated", verification.ID, map[string]interface{}{
		"document_id": req.DocumentID,
	})

	c.JSON(http.StatusAccepted, verification)
}

func (h *VerificationHandler) GetVerificationStatus(c *gin.Context) {
	verificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewError("INVALID_ID", "Invalid verification ID", ""))
		return
	}

	verification, err := h.verificationService.GetVerificationStatus(c.Request.Context(), verificationID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, verification)
}

func (h *VerificationHandler) handleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrNotFound:
		c.JSON(http.StatusNotFound, domain.NewError("NOT_FOUND", err.Error(), ""))
	case domain.ErrInvalidInput:
		c.JSON(http.StatusBadRequest, domain.NewError("INVALID_INPUT", err.Error(), ""))
	case domain.ErrVerificationFailed:
		c.JSON(http.StatusUnprocessableEntity, domain.NewError("VERIFICATION_FAILED", err.Error(), ""))
	default:
		h.logger.Error("internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.NewError("INTERNAL_ERROR", "An internal error occurred", ""))
	}
}
