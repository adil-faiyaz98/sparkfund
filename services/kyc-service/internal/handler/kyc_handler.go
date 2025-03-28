package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sparkfund/kyc-service/internal/model"
	"github.com/sparkfund/kyc-service/internal/service"
)

type KYCHandler struct {
	kycService *service.KYCService
}

func NewKYCHandler(kycService *service.KYCService) *KYCHandler {
	return &KYCHandler{kycService: kycService}
}

func (h *KYCHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/kyc", h.SubmitKYC)
		v1.GET("/kyc/:id", h.GetKYCStatus)
		v1.POST("/kyc/:id/verify", h.VerifyKYC)
		v1.POST("/kyc/:id/reject", h.RejectKYC)
		v1.GET("/kyc/pending", h.ListPendingKYC)
	}
}

func (h *KYCHandler) SubmitKYC(c *gin.Context) {
	var req model.KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// In a real application, you would get the user ID from the authentication context
	userID := uuid.New() // This is just for demonstration

	resp, err := h.kycService.SubmitKYC(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *KYCHandler) GetKYCStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KYC ID"})
		return
	}

	resp, err := h.kycService.GetKYCStatus(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *KYCHandler) VerifyKYC(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KYC ID"})
		return
	}

	// In a real application, you would get the verifier ID from the authentication context
	verifiedBy := uuid.New() // This is just for demonstration

	if err := h.kycService.VerifyKYC(id, verifiedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *KYCHandler) RejectKYC(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KYC ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.kycService.RejectKYC(id, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *KYCHandler) ListPendingKYC(c *gin.Context) {
	kycs, err := h.kycService.ListPendingKYC()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kycs)
}
