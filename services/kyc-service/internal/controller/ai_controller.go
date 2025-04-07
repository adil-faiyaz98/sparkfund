package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/service"
)

// AIController handles AI-related endpoints
type AIController struct {
	aiService *service.AIService
	logger    *logrus.Logger
}

// NewAIController creates a new AI controller
func NewAIController(aiService *service.AIService, logger *logrus.Logger) *AIController {
	return &AIController{
		aiService: aiService,
		logger:    logger,
	}
}

// AnalyzeDocument godoc
// @Summary Analyze a document using AI
// @Description Analyzes a document for authenticity and extracts information
// @Tags ai
// @Accept json
// @Produce json
// @Param request body AnalyzeDocumentRequest true "Document analysis request"
// @Success 200 {object} AnalyzeDocumentResponse "Document analysis result"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Document not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /ai/analyze-document [post]
func (c *AIController) AnalyzeDocument(ctx *gin.Context) {
	var req AnalyzeDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: err.Error(),
		})
		return
	}

	// Parse UUIDs
	documentID, err := uuid.Parse(req.DocumentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
			Details: err.Error(),
		})
		return
	}

	verificationID, err := uuid.Parse(req.VerificationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid verification ID",
			Details: err.Error(),
		})
		return
	}

	// Analyze document
	result, err := c.aiService.AnalyzeDocument(ctx, documentID, verificationID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to analyze document")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to analyze document",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, AnalyzeDocumentResponse{
		ID:             result.ID.String(),
		VerificationID: result.VerificationID.String(),
		DocumentID:     result.DocumentID.String(),
		DocumentType:   result.DocumentType,
		IsAuthentic:    result.IsAuthentic,
		Confidence:     result.Confidence,
		ExtractedData:  result.ExtractedData,
		Issues:         result.Issues,
		CreatedAt:      result.CreatedAt,
	})
}

// MatchFaces godoc
// @Summary Match faces between a selfie and a document
// @Description Compares a selfie with a document photo to verify identity
// @Tags ai
// @Accept json
// @Produce json
// @Param request body MatchFacesRequest true "Face matching request"
// @Success 200 {object} MatchFacesResponse "Face matching result"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Document or selfie not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /ai/match-faces [post]
func (c *AIController) MatchFaces(ctx *gin.Context) {
	var req MatchFacesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: err.Error(),
		})
		return
	}

	// Parse UUIDs
	documentID, err := uuid.Parse(req.DocumentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
			Details: err.Error(),
		})
		return
	}

	selfieID, err := uuid.Parse(req.SelfieID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid selfie ID",
			Details: err.Error(),
		})
		return
	}

	verificationID, err := uuid.Parse(req.VerificationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid verification ID",
			Details: err.Error(),
		})
		return
	}

	// Match faces
	result, err := c.aiService.MatchFaces(ctx, documentID, selfieID, verificationID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to match faces")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to match faces",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, MatchFacesResponse{
		ID:             result.ID.String(),
		VerificationID: result.VerificationID.String(),
		DocumentID:     result.DocumentID.String(),
		SelfieID:       result.SelfieID.String(),
		IsMatch:        result.IsMatch,
		Confidence:     result.Confidence,
		CreatedAt:      result.CreatedAt,
	})
}

// AnalyzeRisk godoc
// @Summary Analyze risk for a user
// @Description Analyzes risk based on user data and device information
// @Tags ai
// @Accept json
// @Produce json
// @Param request body AnalyzeRiskRequest true "Risk analysis request"
// @Success 200 {object} AnalyzeRiskResponse "Risk analysis result"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /ai/analyze-risk [post]
func (c *AIController) AnalyzeRisk(ctx *gin.Context) {
	var req AnalyzeRiskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: err.Error(),
		})
		return
	}

	// Parse UUIDs
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	verificationID, err := uuid.Parse(req.VerificationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid verification ID",
			Details: err.Error(),
		})
		return
	}

	// Create device info
	deviceInfo := model.DeviceInfo{
		IPAddress:    req.DeviceInfo.IPAddress,
		UserAgent:    req.DeviceInfo.UserAgent,
		DeviceType:   req.DeviceInfo.DeviceType,
		OS:           req.DeviceInfo.OS,
		Browser:      req.DeviceInfo.Browser,
		MacAddress:   req.DeviceInfo.MacAddress,
		Location:     req.DeviceInfo.Location,
		Coordinates:  req.DeviceInfo.Coordinates,
		ISP:          req.DeviceInfo.ISP,
		CountryCode:  req.DeviceInfo.CountryCode,
		CapturedTime: req.DeviceInfo.CapturedTime,
	}

	// Analyze risk
	result, err := c.aiService.AnalyzeRisk(ctx, userID, verificationID, deviceInfo)
	if err != nil {
		c.logger.WithError(err).Error("Failed to analyze risk")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to analyze risk",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, AnalyzeRiskResponse{
		ID:             result.ID.String(),
		VerificationID: result.VerificationID.String(),
		UserID:         result.UserID.String(),
		RiskScore:      result.RiskScore,
		RiskLevel:      result.RiskLevel,
		RiskFactors:    result.RiskFactors,
		DeviceInfo:     result.DeviceInfo,
		IPAddress:      result.IPAddress,
		Location:       result.Location,
		CreatedAt:      result.CreatedAt,
	})
}

// DetectAnomalies godoc
// @Summary Detect anomalies in user behavior
// @Description Detects anomalies in user behavior based on historical patterns
// @Tags ai
// @Accept json
// @Produce json
// @Param request body DetectAnomaliesRequest true "Anomaly detection request"
// @Success 200 {object} DetectAnomaliesResponse "Anomaly detection result"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /ai/detect-anomalies [post]
func (c *AIController) DetectAnomalies(ctx *gin.Context) {
	var req DetectAnomaliesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: err.Error(),
		})
		return
	}

	// Parse UUIDs
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
			Details: err.Error(),
		})
		return
	}

	verificationID, err := uuid.Parse(req.VerificationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid verification ID",
			Details: err.Error(),
		})
		return
	}

	// Create device info
	deviceInfo := model.DeviceInfo{
		IPAddress:    req.DeviceInfo.IPAddress,
		UserAgent:    req.DeviceInfo.UserAgent,
		DeviceType:   req.DeviceInfo.DeviceType,
		OS:           req.DeviceInfo.OS,
		Browser:      req.DeviceInfo.Browser,
		MacAddress:   req.DeviceInfo.MacAddress,
		Location:     req.DeviceInfo.Location,
		Coordinates:  req.DeviceInfo.Coordinates,
		ISP:          req.DeviceInfo.ISP,
		CountryCode:  req.DeviceInfo.CountryCode,
		CapturedTime: req.DeviceInfo.CapturedTime,
	}

	// Detect anomalies
	result, err := c.aiService.DetectAnomalies(ctx, userID, verificationID, deviceInfo)
	if err != nil {
		c.logger.WithError(err).Error("Failed to detect anomalies")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to detect anomalies",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, DetectAnomaliesResponse{
		ID:             result.ID.String(),
		VerificationID: result.VerificationID.String(),
		UserID:         result.UserID.String(),
		IsAnomaly:      result.IsAnomaly,
		AnomalyScore:   result.AnomalyScore,
		AnomalyType:    result.AnomalyType,
		Reasons:        result.Reasons,
		DeviceInfo:     result.DeviceInfo,
		CreatedAt:      result.CreatedAt,
	})
}

// ProcessDocument godoc
// @Summary Process a document through all AI checks
// @Description Processes a document through document analysis, face matching, risk analysis, and anomaly detection
// @Tags ai
// @Accept json
// @Produce json
// @Param request body ProcessDocumentRequest true "Document processing request"
// @Success 200 {object} ProcessDocumentResponse "Document processing result"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Document or verification not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /ai/process-document [post]
func (c *AIController) ProcessDocument(ctx *gin.Context) {
	var req ProcessDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: err.Error(),
		})
		return
	}

	// Parse UUIDs
	documentID, err := uuid.Parse(req.DocumentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
			Details: err.Error(),
		})
		return
	}

	var selfieID uuid.UUID
	if req.SelfieID != "" {
		selfieID, err = uuid.Parse(req.SelfieID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid selfie ID",
				Details: err.Error(),
			})
			return
		}
	}

	verificationID, err := uuid.Parse(req.VerificationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid verification ID",
			Details: err.Error(),
		})
		return
	}

	// Create device info
	deviceInfo := model.DeviceInfo{
		IPAddress:    req.DeviceInfo.IPAddress,
		UserAgent:    req.DeviceInfo.UserAgent,
		DeviceType:   req.DeviceInfo.DeviceType,
		OS:           req.DeviceInfo.OS,
		Browser:      req.DeviceInfo.Browser,
		MacAddress:   req.DeviceInfo.MacAddress,
		Location:     req.DeviceInfo.Location,
		Coordinates:  req.DeviceInfo.Coordinates,
		ISP:          req.DeviceInfo.ISP,
		CountryCode:  req.DeviceInfo.CountryCode,
		CapturedTime: req.DeviceInfo.CapturedTime,
	}

	// Process document
	verification, err := c.aiService.ProcessDocument(ctx, documentID, selfieID, verificationID, deviceInfo)
	if err != nil {
		c.logger.WithError(err).Error("Failed to process document")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to process document",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, ProcessDocumentResponse{
		VerificationID: verification.ID.String(),
		UserID:         verification.UserID.String(),
		DocumentID:     documentID.String(),
		Status:         string(verification.Status),
		Notes:          verification.Notes,
		CompletedAt:    verification.CompletedAt,
	})
}

// GetAIModels godoc
// @Summary Get AI models
// @Description Get a list of AI models
// @Tags ai
// @Produce json
// @Success 200 {object} GetAIModelsResponse "AI models"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /ai/models [get]
func (c *AIController) GetAIModels(ctx *gin.Context) {
	// Get AI models
	models, err := c.aiService.ListAIModels(ctx)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get AI models")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get AI models",
			Details: err.Error(),
		})
		return
	}

	// Convert models to response format
	var modelResponses []AIModelResponse
	for _, model := range models {
		modelResponses = append(modelResponses, AIModelResponse{
			ID:            model.ID.String(),
			Name:          model.Name,
			Version:       model.Version,
			Type:          model.Type,
			Accuracy:      model.Accuracy,
			LastTrainedAt: model.LastTrainedAt,
		})
	}

	// Return response
	ctx.JSON(http.StatusOK, GetAIModelsResponse{
		Models: modelResponses,
	})
}

// RegisterRoutes registers the AI controller routes
func (c *AIController) RegisterRoutes(router *gin.Engine) {
	ai := router.Group("/api/v1/ai")
	{
		ai.POST("/analyze-document", c.AnalyzeDocument)
		ai.POST("/match-faces", c.MatchFaces)
		ai.POST("/analyze-risk", c.AnalyzeRisk)
		ai.POST("/detect-anomalies", c.DetectAnomalies)
		ai.POST("/process-document", c.ProcessDocument)
		ai.GET("/models", c.GetAIModels)
	}
}

// Request and response types

// AnalyzeDocumentRequest represents a request to analyze a document
type AnalyzeDocumentRequest struct {
	DocumentID     string `json:"document_id" binding:"required"`
	VerificationID string `json:"verification_id" binding:"required"`
}

// AnalyzeDocumentResponse represents a response from document analysis
type AnalyzeDocumentResponse struct {
	ID             string            `json:"id"`
	VerificationID string            `json:"verification_id"`
	DocumentID     string            `json:"document_id"`
	DocumentType   string            `json:"document_type"`
	IsAuthentic    bool              `json:"is_authentic"`
	Confidence     float64           `json:"confidence"`
	ExtractedData  map[string]string `json:"extracted_data"`
	Issues         []string          `json:"issues"`
	CreatedAt      time.Time         `json:"created_at"`
}

// MatchFacesRequest represents a request to match faces
type MatchFacesRequest struct {
	DocumentID     string `json:"document_id" binding:"required"`
	SelfieID       string `json:"selfie_id" binding:"required"`
	VerificationID string `json:"verification_id" binding:"required"`
}

// MatchFacesResponse represents a response from face matching
type MatchFacesResponse struct {
	ID             string    `json:"id"`
	VerificationID string    `json:"verification_id"`
	DocumentID     string    `json:"document_id"`
	SelfieID       string    `json:"selfie_id"`
	IsMatch        bool      `json:"is_match"`
	Confidence     float64   `json:"confidence"`
	CreatedAt      time.Time `json:"created_at"`
}

// DeviceInfoRequest represents device information in a request
type DeviceInfoRequest struct {
	IPAddress    string    `json:"ip_address" binding:"required"`
	UserAgent    string    `json:"user_agent" binding:"required"`
	DeviceType   string    `json:"device_type"`
	OS           string    `json:"os"`
	Browser      string    `json:"browser"`
	MacAddress   string    `json:"mac_address,omitempty"`
	Location     string    `json:"location,omitempty"`
	Coordinates  string    `json:"coordinates,omitempty"`
	ISP          string    `json:"isp,omitempty"`
	CountryCode  string    `json:"country_code,omitempty"`
	CapturedTime time.Time `json:"captured_time"`
}

// AnalyzeRiskRequest represents a request to analyze risk
type AnalyzeRiskRequest struct {
	UserID         string           `json:"user_id" binding:"required"`
	VerificationID string           `json:"verification_id" binding:"required"`
	DeviceInfo     DeviceInfoRequest `json:"device_info" binding:"required"`
}

// AnalyzeRiskResponse represents a response from risk analysis
type AnalyzeRiskResponse struct {
	ID             string            `json:"id"`
	VerificationID string            `json:"verification_id"`
	UserID         string            `json:"user_id"`
	RiskScore      float64           `json:"risk_score"`
	RiskLevel      string            `json:"risk_level"`
	RiskFactors    []string          `json:"risk_factors"`
	DeviceInfo     map[string]string `json:"device_info"`
	IPAddress      string            `json:"ip_address"`
	Location       string            `json:"location"`
	CreatedAt      time.Time         `json:"created_at"`
}

// DetectAnomaliesRequest represents a request to detect anomalies
type DetectAnomaliesRequest struct {
	UserID         string           `json:"user_id" binding:"required"`
	VerificationID string           `json:"verification_id" binding:"required"`
	DeviceInfo     DeviceInfoRequest `json:"device_info" binding:"required"`
}

// DetectAnomaliesResponse represents a response from anomaly detection
type DetectAnomaliesResponse struct {
	ID             string            `json:"id"`
	VerificationID string            `json:"verification_id"`
	UserID         string            `json:"user_id"`
	IsAnomaly      bool              `json:"is_anomaly"`
	AnomalyScore   float64           `json:"anomaly_score"`
	AnomalyType    string            `json:"anomaly_type"`
	Reasons        []string          `json:"reasons"`
	DeviceInfo     map[string]string `json:"device_info"`
	CreatedAt      time.Time         `json:"created_at"`
}

// ProcessDocumentRequest represents a request to process a document
type ProcessDocumentRequest struct {
	DocumentID     string           `json:"document_id" binding:"required"`
	SelfieID       string           `json:"selfie_id,omitempty"`
	VerificationID string           `json:"verification_id" binding:"required"`
	DeviceInfo     DeviceInfoRequest `json:"device_info" binding:"required"`
}

// ProcessDocumentResponse represents a response from document processing
type ProcessDocumentResponse struct {
	VerificationID string     `json:"verification_id"`
	UserID         string     `json:"user_id"`
	DocumentID     string     `json:"document_id"`
	Status         string     `json:"status"`
	Notes          string     `json:"notes"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// AIModelResponse represents an AI model in a response
type AIModelResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	Type          string    `json:"type"`
	Accuracy      float64   `json:"accuracy"`
	LastTrainedAt time.Time `json:"last_trained_at"`
}

// GetAIModelsResponse represents a response containing AI models
type GetAIModelsResponse struct {
	Models []AIModelResponse `json:"models"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
