package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sparkfund/kyc-service/internal/model"
	"net/http"
)

type KYCHandler struct {
	*BaseHandler
}

func NewKYCHandler(base *BaseHandler) *KYCHandler {
	return &KYCHandler{BaseHandler: base}
}

func (h *KYCHandler) RegisterRoutes(r *gin.RouterGroup) {
	kyc := r.Group("/kyc")
	{
		// Document operations
		kyc.POST("/documents", h.UploadDocument)
		kyc.GET("/documents/:id", h.GetDocument)
		kyc.DELETE("/documents/:id", h.DeleteDocument)
		kyc.GET("/documents/pending", h.ListPendingDocuments)

		// Verification operations
		kyc.POST("/verify", h.VerifyDocument)
		kyc.GET("/verify/:id", h.GetVerificationStatus)
		kyc.POST("/validate", h.ValidateIdentity)

		// KYC status operations
		kyc.GET("/status/:userId", h.GetKYCStatus)
		kyc.PUT("/status/:userId", h.UpdateKYCStatus)
	}
}

func (h *KYCHandler) VerifyDocument(c *gin.Context) {
	var req model.DocumentVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.services.KYC.VerifyDocument(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *KYCHandler) UploadDocument(c *gin.Context) {
	kycID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewError("INVALID_ID", "Invalid KYC ID", ""))
		return
	}

	file, header, err := c.Request.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewError("INVALID_FILE", "Invalid file upload", ""))
		return
	}
	defer file.Close()

	docType := domain.DocumentType(c.PostForm("type"))
	if !isValidDocumentType(docType) {
		c.JSON(http.StatusBadRequest, domain.NewError("INVALID_DOC_TYPE", "Invalid document type", ""))
		return
	}

	doc, err := h.kycService.UploadDocument(c.Request.Context(), kycID, docType, file, header)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.auditLogger.Log(c.Request.Context(), "document_uploaded", doc.ID, map[string]interface{}{
		"kyc_id": kycID,
		"type":   docType,
		"size":   header.Size,
	})

	c.JSON(http.StatusCreated, doc)
}

func (h *KYCHandler) GetKYCStatus(c *gin.Context) {
	kycID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewError("INVALID_ID", "Invalid KYC ID", ""))
		return
	}

	kyc, err := h.kycService.GetKYCStatus(c.Request.Context(), kycID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, kyc)
}

func (h *KYCHandler) handleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrNotFound:
		c.JSON(http.StatusNotFound, domain.NewError("NOT_FOUND", err.Error(), ""))
	case domain.ErrInvalidInput:
		c.JSON(http.StatusBadRequest, domain.NewError("INVALID_INPUT", err.Error(), ""))
	case domain.ErrUnauthorized:
		c.JSON(http.StatusUnauthorized, domain.NewError("UNAUTHORIZED", err.Error(), ""))
	case domain.ErrForbidden:
		c.JSON(http.StatusForbidden, domain.NewError("FORBIDDEN", err.Error(), ""))
	default:
		h.logger.Error("internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.NewError("INTERNAL_ERROR", "An internal error occurred", ""))
	}
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	id, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, domain.ErrUnauthorized
	}
	return id.(uuid.UUID), nil
}

func isValidDocumentType(docType domain.DocumentType) bool {
	validTypes := map[domain.DocumentType]bool{
		domain.DocTypePassport:      true,
		domain.DocTypeDriverLicense: true,
		domain.DocTypeIDCard:        true,
		domain.DocTypeUtilityBill:   true,
		domain.DocTypeBankStatement: true,
	}
	return validTypes[docType]
}
