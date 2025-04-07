package service

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/repository"
)

// AIServiceSimple handles AI-related operations with simulated AI functionality
type AIServiceSimple struct {
	repo             *repository.AIRepository
	documentRepo     *repository.DocumentRepository
	verificationRepo *repository.VerificationRepository
	logger           *logrus.Logger
}

// NewAIServiceSimple creates a new AI service with simulated functionality
func NewAIServiceSimple(
	repo *repository.AIRepository,
	documentRepo *repository.DocumentRepository,
	verificationRepo *repository.VerificationRepository,
	logger *logrus.Logger,
) *AIServiceSimple {
	return &AIServiceSimple{
		repo:             repo,
		documentRepo:     documentRepo,
		verificationRepo: verificationRepo,
		logger:           logger,
	}
}

// ListAIModels lists all AI models
func (s *AIServiceSimple) ListAIModels(ctx context.Context) ([]*model.AIModelInfo, error) {
	// Check if we have models in the database
	models, err := s.repo.ListAIModels(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list AI models")
		return nil, err
	}

	// If no models found, create some default ones
	if len(models) == 0 {
		// Create default models
		documentModel := &model.AIModelInfo{
			ID:            uuid.New(),
			Name:          "Document Verification Model",
			Version:       "1.0.0",
			Type:          "DOCUMENT",
			Accuracy:      0.98,
			LastTrainedAt: time.Now().Add(-24 * time.Hour),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		faceModel := &model.AIModelInfo{
			ID:            uuid.New(),
			Name:          "Face Recognition Model",
			Version:       "1.0.0",
			Type:          "FACE",
			Accuracy:      0.95,
			LastTrainedAt: time.Now().Add(-48 * time.Hour),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		riskModel := &model.AIModelInfo{
			ID:            uuid.New(),
			Name:          "Risk Analysis Model",
			Version:       "1.0.0",
			Type:          "RISK",
			Accuracy:      0.92,
			LastTrainedAt: time.Now().Add(-72 * time.Hour),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		anomalyModel := &model.AIModelInfo{
			ID:            uuid.New(),
			Name:          "Anomaly Detection Model",
			Version:       "1.0.0",
			Type:          "ANOMALY",
			Accuracy:      0.90,
			LastTrainedAt: time.Now().Add(-96 * time.Hour),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Save models to database
		if err := s.repo.SaveAIModelInfo(ctx, documentModel); err != nil {
			s.logger.WithError(err).Error("Failed to save document model")
			return nil, err
		}

		if err := s.repo.SaveAIModelInfo(ctx, faceModel); err != nil {
			s.logger.WithError(err).Error("Failed to save face model")
			return nil, err
		}

		if err := s.repo.SaveAIModelInfo(ctx, riskModel); err != nil {
			s.logger.WithError(err).Error("Failed to save risk model")
			return nil, err
		}

		if err := s.repo.SaveAIModelInfo(ctx, anomalyModel); err != nil {
			s.logger.WithError(err).Error("Failed to save anomaly model")
			return nil, err
		}

		// Return the newly created models
		return []*model.AIModelInfo{documentModel, faceModel, riskModel, anomalyModel}, nil
	}

	return models, nil
}

// AnalyzeDocument analyzes a document using simulated AI
func (s *AIServiceSimple) AnalyzeDocument(ctx context.Context, documentID uuid.UUID, verificationID uuid.UUID) (*model.DocumentAnalysisResult, error) {
	// Get document from repository
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		s.logger.WithError(err).WithField("document_id", documentID).Error("Failed to get document for AI analysis")
		return nil, err
	}

	// Simulate AI analysis
	extractedData := map[string]string{
		"full_name":       "John Smith",
		"document_number": "X123456789",
		"date_of_birth":   "1990-01-01",
		"expiry_date":     "2030-01-01",
		"issuing_country": "United States",
	}

	// Convert map to JSON
	extractedDataJSON, err := json.Marshal(extractedData)
	if err != nil {
		s.logger.WithError(err).Error("Failed to marshal extracted data")
		return nil, err
	}

	// Simulate authenticity check (90% chance of being authentic)
	isAuthentic := rand.Float64() > 0.1
	confidence := 70.0 + rand.Float64()*25.0

	// Create issues array
	var issues []string
	if !isAuthentic {
		issues = append(issues, "Document appears to be manipulated")
		issues = append(issues, "Security features missing")
	}

	// Convert issues to JSON
	issuesJSON, err := json.Marshal(issues)
	if err != nil {
		s.logger.WithError(err).Error("Failed to marshal issues")
		return nil, err
	}

	// Create result
	result := &model.DocumentAnalysisResult{
		ID:             uuid.New(),
		VerificationID: verificationID,
		DocumentID:     documentID,
		DocumentType:   string(document.Type),
		IsAuthentic:    isAuthentic,
		Confidence:     confidence,
		ExtractedData:  extractedDataJSON,
		Issues:         issuesJSON,
		CreatedAt:      time.Now(),
	}

	// Save the result
	err = s.repo.SaveDocumentAnalysis(ctx, result)
	if err != nil {
		s.logger.WithError(err).WithField("document_id", documentID).Error("Failed to save document analysis result")
		return nil, err
	}

	return result, nil
}

// MatchFaces matches a selfie with a document photo using simulated AI
func (s *AIServiceSimple) MatchFaces(ctx context.Context, documentID uuid.UUID, selfieID uuid.UUID, verificationID uuid.UUID) (*model.FaceMatchResult, error) {
	// Get document and selfie from repository
	_, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		s.logger.WithError(err).WithField("document_id", documentID).Error("Failed to get document for face matching")
		return nil, err
	}

	_, err = s.documentRepo.GetByID(ctx, selfieID)
	if err != nil {
		s.logger.WithError(err).WithField("selfie_id", selfieID).Error("Failed to get selfie for face matching")
		return nil, err
	}

	// Simulate face matching (85% chance of matching)
	isMatch := rand.Float64() > 0.15
	confidence := 0.0
	if isMatch {
		confidence = 75.0 + rand.Float64()*20.0 // 75-95% confidence for matches
	} else {
		confidence = 30.0 + rand.Float64()*40.0 // 30-70% confidence for non-matches
	}

	// Create result
	result := &model.FaceMatchResult{
		ID:             uuid.New(),
		VerificationID: verificationID,
		DocumentID:     documentID,
		SelfieID:       selfieID,
		IsMatch:        isMatch,
		Confidence:     confidence,
		CreatedAt:      time.Now(),
	}

	// Save the result
	err = s.repo.SaveFaceMatchResult(ctx, result)
	if err != nil {
		s.logger.WithError(err).WithField("document_id", documentID).WithField("selfie_id", selfieID).Error("Failed to save face match result")
		return nil, err
	}

	return result, nil
}

// AnalyzeRisk analyzes risk based on user data and device information
func (s *AIServiceSimple) AnalyzeRisk(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID, deviceInfo model.DeviceInfo) (*model.RiskAnalysisResult, error) {
	// Simulate risk analysis
	riskScore := 5.0 + rand.Float64()*20.0 // 5-25% risk score
	riskLevel := "LOW"
	if riskScore > 15.0 {
		riskLevel = "MEDIUM"
	}
	if riskScore > 20.0 {
		riskLevel = "HIGH"
	}

	// Create risk factors
	riskFactors := []string{}
	if riskScore > 15.0 {
		riskFactors = append(riskFactors, "Unusual location")
	}
	if riskScore > 20.0 {
		riskFactors = append(riskFactors, "Multiple failed attempts")
		riskFactors = append(riskFactors, "IP address associated with fraud")
	}

	// Convert risk factors to JSON
	riskFactorsJSON, err := json.Marshal(riskFactors)
	if err != nil {
		s.logger.WithError(err).Error("Failed to marshal risk factors")
		return nil, err
	}

	// Convert device info to JSON
	deviceInfoMap := map[string]interface{}{
		"ip_address":    deviceInfo.IPAddress,
		"user_agent":    deviceInfo.UserAgent,
		"device_type":   deviceInfo.DeviceType,
		"os":            deviceInfo.OS,
		"browser":       deviceInfo.Browser,
		"location":      deviceInfo.Location,
		"captured_time": deviceInfo.CapturedTime,
	}
	deviceInfoJSON, err := json.Marshal(deviceInfoMap)
	if err != nil {
		s.logger.WithError(err).Error("Failed to marshal device info")
		return nil, err
	}

	// Create result
	result := &model.RiskAnalysisResult{
		ID:             uuid.New(),
		VerificationID: verificationID,
		UserID:         userID,
		RiskScore:      riskScore,
		RiskLevel:      riskLevel,
		RiskFactors:    riskFactorsJSON,
		DeviceInfo:     deviceInfoJSON,
		IPAddress:      deviceInfo.IPAddress,
		Location:       deviceInfo.Location,
		CreatedAt:      time.Now(),
	}

	// Save the result
	err = s.repo.SaveRiskAnalysisResult(ctx, result)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to save risk analysis result")
		return nil, err
	}

	return result, nil
}

// DetectAnomalies detects anomalies in user behavior
func (s *AIServiceSimple) DetectAnomalies(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID, deviceInfo model.DeviceInfo) (*model.AnomalyDetectionResult, error) {
	// Simulate anomaly detection (10% chance of anomaly)
	isAnomaly := rand.Float64() < 0.1
	anomalyScore := 0.0
	anomalyType := ""
	reasons := []string{}

	if isAnomaly {
		anomalyScore = 70.0 + rand.Float64()*30.0 // 70-100% anomaly score
		anomalyType = "SUSPICIOUS_BEHAVIOR"
		reasons = append(reasons, "Multiple verification attempts in short time")
		reasons = append(reasons, "Different device than usual")
		reasons = append(reasons, "Unusual time of day")
	} else {
		anomalyScore = rand.Float64() * 30.0 // 0-30% anomaly score
	}

	// Convert reasons to JSON
	reasonsJSON, err := json.Marshal(reasons)
	if err != nil {
		s.logger.WithError(err).Error("Failed to marshal reasons")
		return nil, err
	}

	// Convert device info to JSON
	deviceInfoMap := map[string]interface{}{
		"ip_address":    deviceInfo.IPAddress,
		"user_agent":    deviceInfo.UserAgent,
		"device_type":   deviceInfo.DeviceType,
		"os":            deviceInfo.OS,
		"browser":       deviceInfo.Browser,
		"location":      deviceInfo.Location,
		"captured_time": deviceInfo.CapturedTime,
	}
	deviceInfoJSON, err := json.Marshal(deviceInfoMap)
	if err != nil {
		s.logger.WithError(err).Error("Failed to marshal device info")
		return nil, err
	}

	// Create result
	result := &model.AnomalyDetectionResult{
		ID:             uuid.New(),
		VerificationID: verificationID,
		UserID:         userID,
		IsAnomaly:      isAnomaly,
		AnomalyScore:   anomalyScore,
		AnomalyType:    anomalyType,
		Reasons:        reasonsJSON,
		DeviceInfo:     deviceInfoJSON,
		CreatedAt:      time.Now(),
	}

	// Save the result
	err = s.repo.SaveAnomalyDetectionResult(ctx, result)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to save anomaly detection result")
		return nil, err
	}

	return result, nil
}

// ProcessDocument processes a document through all AI checks
func (s *AIServiceSimple) ProcessDocument(ctx context.Context, documentID uuid.UUID, selfieID uuid.UUID, verificationID uuid.UUID, deviceInfo model.DeviceInfo) (*model.Verification, error) {
	// Get verification
	verification, err := s.verificationRepo.GetByID(ctx, verificationID)
	if err != nil {
		s.logger.WithError(err).WithField("verification_id", verificationID).Error("Failed to get verification")
		return nil, err
	}

	// Update verification status
	verification.Status = model.VerificationStatusInProcess
	verification.UpdatedAt = time.Now()

	err = s.verificationRepo.Update(ctx, verification)
	if err != nil {
		s.logger.WithError(err).WithField("verification_id", verificationID).Error("Failed to update verification status")
		return nil, err
	}

	// Step 1: Analyze document
	documentAnalysis, err := s.AnalyzeDocument(ctx, documentID, verificationID)
	if err != nil {
		s.logger.WithError(err).WithField("document_id", documentID).Error("Failed to analyze document")
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Failed to analyze document: " + err.Error()
		verification.UpdatedAt = time.Now()
		s.verificationRepo.Update(ctx, verification)
		return verification, err
	}

	// Step 2: Match faces if selfie is provided
	var faceMatch *model.FaceMatchResult
	if selfieID != uuid.Nil {
		faceMatch, err = s.MatchFaces(ctx, documentID, selfieID, verificationID)
		if err != nil {
			s.logger.WithError(err).WithField("document_id", documentID).WithField("selfie_id", selfieID).Error("Failed to match faces")
			verification.Status = model.VerificationStatusFailed
			verification.Notes = "Failed to match faces: " + err.Error()
			verification.UpdatedAt = time.Now()
			s.verificationRepo.Update(ctx, verification)
			return verification, err
		}
	}

	// Step 3: Analyze risk
	riskAnalysis, err := s.AnalyzeRisk(ctx, verification.UserID, verificationID, deviceInfo)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", verification.UserID).Error("Failed to analyze risk")
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Failed to analyze risk: " + err.Error()
		verification.UpdatedAt = time.Now()
		s.verificationRepo.Update(ctx, verification)
		return verification, err
	}

	// Step 4: Detect anomalies
	anomalyDetection, err := s.DetectAnomalies(ctx, verification.UserID, verificationID, deviceInfo)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", verification.UserID).Error("Failed to detect anomalies")
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Failed to detect anomalies: " + err.Error()
		verification.UpdatedAt = time.Now()
		s.verificationRepo.Update(ctx, verification)
		return verification, err
	}

	// Step 5: Determine final verification status
	verification.Status = model.VerificationStatusCompleted
	verification.Notes = "All verification checks passed"
	verification.UpdatedAt = time.Now()
	now := time.Now()
	verification.CompletedAt = &now

	// Check if any checks failed
	if !documentAnalysis.IsAuthentic {
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Document verification failed: Document is not authentic"
	}

	if faceMatch != nil && !faceMatch.IsMatch {
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Face verification failed: Faces do not match"
	}

	if riskAnalysis.RiskLevel == "HIGH" {
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Risk analysis failed: High risk level detected"
	}

	if anomalyDetection.IsAnomaly {
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Anomaly detection failed: Suspicious behavior detected"
	}

	// Update verification
	err = s.verificationRepo.Update(ctx, verification)
	if err != nil {
		s.logger.WithError(err).WithField("verification_id", verificationID).Error("Failed to update verification")
		return nil, err
	}

	return verification, nil
}
