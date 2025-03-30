package handlers

import (
	"net/http"
	"strconv"

	"sparkfund/services/kyc-service/internal/models"
	"sparkfund/services/kyc-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProfileHandlers handles HTTP requests for KYC profile operations
type ProfileHandlers struct {
	profileService *service.ProfileService
}

// NewProfileHandlers creates new profile handlers
func NewProfileHandlers(profileService *service.ProfileService) *ProfileHandlers {
	return &ProfileHandlers{
		profileService: profileService,
	}
}

// CreateProfile handles profile creation requests
func (h *ProfileHandlers) CreateProfile(c *gin.Context) {
	// Get user ID from context
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse request body
	var profile models.KYCProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate profile
	if err := h.profileService.ValidateProfile(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create profile
	created, err := h.profileService.CreateProfile(c.Request.Context(), userID, &profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetProfile handles profile retrieval requests
func (h *ProfileHandlers) GetProfile(c *gin.Context) {
	// Get user ID from context
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get profile
	profile, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile handles profile update requests
func (h *ProfileHandlers) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse request body
	var profile models.KYCProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate profile
	if err := h.profileService.ValidateProfile(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update profile
	err = h.profileService.UpdateProfile(c.Request.Context(), userID, &profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

// UpdateProfileStatus handles profile status update requests
func (h *ProfileHandlers) UpdateProfileStatus(c *gin.Context) {
	// Parse user ID
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse request body
	var req struct {
		Status models.ProfileStatus `json:"status" binding:"required"`
		Notes  string               `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update profile status
	err = h.profileService.UpdateProfileStatus(c.Request.Context(), userID, req.Status, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile status updated successfully"})
}

// UpdateRiskLevel handles risk level update requests
func (h *ProfileHandlers) UpdateRiskLevel(c *gin.Context) {
	// Parse user ID
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse request body
	var req struct {
		RiskLevel models.RiskLevel `json:"risk_level" binding:"required"`
		Notes     string           `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update risk level
	err = h.profileService.UpdateRiskLevel(c.Request.Context(), userID, req.RiskLevel, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "risk level updated successfully"})
}

// ListProfiles handles profile listing requests
func (h *ProfileHandlers) ListProfiles(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *models.ProfileStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := models.ProfileStatus(statusStr)
		if s.IsValid() {
			status = &s
		}
	}

	var riskLevel *models.RiskLevel
	if riskStr := c.Query("risk_level"); riskStr != "" {
		r := models.RiskLevel(riskStr)
		if r.IsValid() {
			riskLevel = &r
		}
	}

	// List profiles
	profiles, total, err := h.profileService.ListProfiles(c.Request.Context(), status, riskLevel, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetProfileStats handles profile statistics requests
func (h *ProfileHandlers) GetProfileStats(c *gin.Context) {
	// Get profile statistics
	stats, err := h.profileService.GetProfileStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DeleteProfile handles profile deletion requests
func (h *ProfileHandlers) DeleteProfile(c *gin.Context) {
	// Parse user ID
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Delete profile
	err = h.profileService.DeleteProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile deleted successfully"})
}
