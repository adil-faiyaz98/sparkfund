package service

import (
	"context"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/repository"
)

// AIService handles AI-related operations
type AIService struct {
	repo             *repository.AIRepository
	docRepo          *repository.DocumentRepository
	verificationRepo *repository.VerificationRepository
	logger           *logrus.Logger
	awsSession       *session.Session
	rekognition      *rekognition.Rekognition
	textract         *textract.Textract
	modelPath        string
}

// NewAIService creates a new AI service
func NewAIService(
	repo *repository.AIRepository,
	docRepo *repository.DocumentRepository,
	verificationRepo *repository.VerificationRepository,
	logger *logrus.Logger,
	awsSession *session.Session,
	modelPath string,
) *AIService {
	rekognitionClient := rekognition.New(awsSession)
	textractClient := textract.New(awsSession)

	return &AIService{
		repo:             repo,
		docRepo:          docRepo,
		verificationRepo: verificationRepo,
		logger:           logger,
		awsSession:       awsSession,
		rekognition:      rekognitionClient,
		textract:         textractClient,
		modelPath:        modelPath,
	}
}

// AnalyzeDocument analyzes a document using AI
func (s *AIService) AnalyzeDocument(ctx context.Context, documentID uuid.UUID, verificationID uuid.UUID) (*model.DocumentAnalysisResult, error) {
	// Get document from repository
	document, err := s.docRepo.GetByID(ctx, documentID)
	if err != nil {
		s.logger.WithError(err).WithField("document_id", documentID).Error("Failed to get document for AI analysis")
		return nil, err
	}

	// Read document file
	file, err := os.Open(document.Path)
	if err != nil {
		s.logger.WithError(err).WithField("document_path", document.Path).Error("Failed to open document file")
		return nil, err
	}
	defer file.Close()

	// Read file content
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		s.logger.WithError(err).WithField("document_path", document.Path).Error("Failed to read document file")
		return nil, err
	}

	// Use AWS Textract to extract text and analyze document
	extractedData := make(map[string]string)
	issues := make([]string, 0)
	isAuthentic := true
	confidence := 95.0

	// Check if we're in production or test mode
	if os.Getenv("APP_ENV") == "production" {
		// Use AWS Textract in production
		textractInput := &textract.DetectDocumentTextInput{
			Document: &textract.Document{
				Bytes: fileContent,
			},
		}

		textractOutput, err := s.textract.DetectDocumentText(textractInput)
		if err != nil {
			s.logger.WithError(err).WithField("document_id", documentID).Error("Failed to extract text from document")
			return nil, err
		}

		// Process extracted text
		extractedText := ""
		for _, block := range textractOutput.Blocks {
			if *block.BlockType == "LINE" {
				extractedText += *block.Text + "\n"
			}
		}

		// Parse extracted text based on document type
		switch document.Type {
		case model.DocumentTypePassport:
			extractedData, issues, isAuthentic, confidence = s.parsePassport(extractedText)
		case model.DocumentTypeDriversLicense:
			extractedData, issues, isAuthentic, confidence = s.parseDriversLicense(extractedText)
		case model.DocumentTypeIDCard:
			extractedData, issues, isAuthentic, confidence = s.parseIDCard(extractedText)
		default:
			extractedData, issues, isAuthentic, confidence = s.parseGenericDocument(extractedText)
		}
	} else {
		// Use mock data for testing
		extractedData = map[string]string{
			"full_name":       "John Smith",
			"document_number": "X123456789",
			"date_of_birth":   "1990-01-01",
			"expiry_date":     "2030-01-01",
			"issuing_country": "United States",
		}

		// Simulate some issues in test mode (10% chance)
		if rand.Float64() < 0.1 {
			isAuthentic = false
			confidence = 30.0 + rand.Float64()*40.0
			issues = append(issues, "Document appears to be manipulated")
			issues = append(issues, "Security features missing")
		}
	}

	// Create result
	result := &model.DocumentAnalysisResult{
		ID:             uuid.New(),
		VerificationID: verificationID,
		DocumentID:     documentID,
		DocumentType:   document.Type,
		IsAuthentic:    isAuthentic,
		Confidence:     confidence,
		ExtractedData:  extractedData,
		Issues:         issues,
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

// parsePassport parses passport text and validates it
func (s *AIService) parsePassport(text string) (map[string]string, []string, bool, float64) {
	// In a real implementation, we would use regex and other techniques to extract structured data
	// For this demo, we'll return mock data with some validation
	extractedData := map[string]string{
		"full_name":       "John Smith",
		"document_number": "X123456789",
		"date_of_birth":   "1990-01-01",
		"expiry_date":     "2030-01-01",
		"issuing_country": "United States",
	}

	issues := []string{}
	isAuthentic := true
	confidence := 95.0

	// Validate expiry date
	expiryDate, err := time.Parse("2006-01-02", extractedData["expiry_date"])
	if err != nil {
		issues = append(issues, "Invalid expiry date format")
		isAuthentic = false
		confidence -= 20.0
	} else if expiryDate.Before(time.Now()) {
		issues = append(issues, "Document is expired")
		isAuthentic = false
		confidence -= 10.0
	}

	// Validate document number format (should be alphanumeric)
	if !isAlphanumeric(extractedData["document_number"]) {
		issues = append(issues, "Invalid document number format")
		isAuthentic = false
		confidence -= 15.0
	}

	return extractedData, issues, isAuthentic, confidence
}

// parseDriversLicense parses driver's license text and validates it
func (s *AIService) parseDriversLicense(text string) (map[string]string, []string, bool, float64) {
	// Similar implementation to parsePassport but with driver's license specific logic
	extractedData := map[string]string{
		"full_name":      "John Smith",
		"license_number": "DL123456789",
		"date_of_birth":  "1990-01-01",
		"expiry_date":    "2025-01-01",
		"issuing_state":  "California",
		"address":        "123 Main St, San Francisco, CA 94105",
	}

	issues := []string{}
	isAuthentic := true
	confidence := 92.0

	// Add validation logic here

	return extractedData, issues, isAuthentic, confidence
}

// parseIDCard parses ID card text and validates it
func (s *AIService) parseIDCard(text string) (map[string]string, []string, bool, float64) {
	// Similar implementation to parsePassport but with ID card specific logic
	extractedData := map[string]string{
		"full_name":         "John Smith",
		"id_number":         "ID123456789",
		"date_of_birth":     "1990-01-01",
		"expiry_date":       "2028-01-01",
		"issuing_authority": "Department of Home Affairs",
	}

	issues := []string{}
	isAuthentic := true
	confidence := 90.0

	// Add validation logic here

	return extractedData, issues, isAuthentic, confidence
}

// parseGenericDocument parses generic document text
func (s *AIService) parseGenericDocument(text string) (map[string]string, []string, bool, float64) {
	// For generic documents, we extract what we can but with lower confidence
	extractedData := map[string]string{
		"text": text,
	}

	issues := []string{"Unknown document type, limited validation possible"}
	return extractedData, issues, true, 70.0
}

// isAlphanumeric checks if a string contains only alphanumeric characters
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// MatchFaces matches a selfie with a document photo
func (s *AIService) MatchFaces(ctx context.Context, documentID uuid.UUID, selfieID uuid.UUID, verificationID uuid.UUID) (*model.FaceMatchResult, error) {
	// Get document and selfie from repository
	document, err := s.docRepo.GetByID(ctx, documentID)
	if err != nil {
		s.logger.WithError(err).WithField("document_id", documentID).Error("Failed to get document for face matching")
		return nil, err
	}

	selfie, err := s.docRepo.GetByID(ctx, selfieID)
	if err != nil {
		s.logger.WithError(err).WithField("selfie_id", selfieID).Error("Failed to get selfie for face matching")
		return nil, err
	}

	// Variables to store result
	isMatch := false
	confidence := 0.0

	// Check if we're in production or test mode
	if os.Getenv("APP_ENV") == "production" {
		// Load document image
		docFile, err := os.Open(document.Path)
		if err != nil {
			s.logger.WithError(err).WithField("document_path", document.Path).Error("Failed to open document file")
			return nil, err
		}
		defer docFile.Close()

		docContent, err := ioutil.ReadAll(docFile)
		if err != nil {
			s.logger.WithError(err).WithField("document_path", document.Path).Error("Failed to read document file")
			return nil, err
		}

		// Load selfie image
		selfieFile, err := os.Open(selfie.Path)
		if err != nil {
			s.logger.WithError(err).WithField("selfie_path", selfie.Path).Error("Failed to open selfie file")
			return nil, err
		}
		defer selfieFile.Close()

		selfieContent, err := ioutil.ReadAll(selfieFile)
		if err != nil {
			s.logger.WithError(err).WithField("selfie_path", selfie.Path).Error("Failed to read selfie file")
			return nil, err
		}

		// Use AWS Rekognition to compare faces
		compareFacesInput := &rekognition.CompareFacesInput{
			SourceImage: &rekognition.Image{
				Bytes: selfieContent,
			},
			TargetImage: &rekognition.Image{
				Bytes: docContent,
			},
			SimilarityThreshold: aws.Float64(70.0), // Minimum similarity threshold
		}

		compareFacesOutput, err := s.rekognition.CompareFaces(compareFacesInput)
		if err != nil {
			s.logger.WithError(err).WithField("document_id", documentID).WithField("selfie_id", selfieID).Error("Failed to compare faces")
			return nil, err
		}

		// Process comparison result
		if len(compareFacesOutput.FaceMatches) > 0 {
			// Faces match
			isMatch = true
			confidence = *compareFacesOutput.FaceMatches[0].Similarity
		} else {
			// No match found
			isMatch = false
			confidence = 0.0
		}
	} else {
		// Use mock data for testing
		isMatch = rand.Float64() > 0.15 // 85% chance of matching
		if isMatch {
			confidence = 75.0 + rand.Float64()*20.0 // 75-95% confidence for matches
		} else {
			confidence = 30.0 + rand.Float64()*40.0 // 30-70% confidence for non-matches
		}
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

// AnalyzeRisk analyzes risk for a user
func (s *AIService) AnalyzeRisk(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID, deviceInfo model.DeviceInfo) (*model.RiskAnalysisResult, error) {
	// In a real implementation, we would:
	// 1. Analyze user data, verification history, and device info
	// 2. Check against known fraud patterns
	// 3. Calculate a risk score
	// 4. Return the risk analysis result

	// For this demo, we'll simulate the risk analysis
	riskScore := 10.0 + rand.Float64()*40.0
	riskLevel := "LOW"
	if riskScore > 30.0 {
		riskLevel = "MEDIUM"
	}
	if riskScore > 70.0 {
		riskLevel = "HIGH"
	}

	riskFactors := []string{}
	if rand.Float64() > 0.7 {
		riskFactors = append(riskFactors, "Unusual location")
	}
	if rand.Float64() > 0.8 {
		riskFactors = append(riskFactors, "IP address associated with previous fraud")
	}
	if rand.Float64() > 0.9 {
		riskFactors = append(riskFactors, "Multiple verification attempts")
	}

	deviceInfoMap := map[string]string{
		"ip_address":   deviceInfo.IPAddress,
		"user_agent":   deviceInfo.UserAgent,
		"device_type":  deviceInfo.DeviceType,
		"os":           deviceInfo.OS,
		"browser":      deviceInfo.Browser,
		"mac_address":  deviceInfo.MacAddress,
		"location":     deviceInfo.Location,
		"coordinates":  deviceInfo.Coordinates,
		"isp":          deviceInfo.ISP,
		"country_code": deviceInfo.CountryCode,
	}

	result := &model.RiskAnalysisResult{
		ID:             uuid.New(),
		VerificationID: verificationID,
		UserID:         userID,
		RiskScore:      riskScore,
		RiskLevel:      riskLevel,
		RiskFactors:    riskFactors,
		DeviceInfo:     deviceInfoMap,
		IPAddress:      deviceInfo.IPAddress,
		Location:       deviceInfo.Location,
		CreatedAt:      time.Now(),
	}

	// Save the result
	err := s.repo.SaveRiskAnalysisResult(ctx, result)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to save risk analysis result")
		return nil, err
	}

	return result, nil
}

// DetectAnomalies detects anomalies in user behavior
func (s *AIService) DetectAnomalies(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID, deviceInfo model.DeviceInfo) (*model.AnomalyDetectionResult, error) {
	// In a real implementation, we would:
	// 1. Analyze user's historical behavior
	// 2. Compare current behavior with historical patterns
	// 3. Detect anomalies using ML models
	// 4. Return the anomaly detection result

	// For this demo, we'll simulate the anomaly detection
	isAnomaly := rand.Float64() > 0.8 // 20% chance of anomaly
	anomalyScore := 10.0 + rand.Float64()*30.0
	if isAnomaly {
		anomalyScore = 70.0 + rand.Float64()*30.0
	}

	anomalyType := ""
	reasons := []string{}
	if isAnomaly {
		anomalyTypes := []string{"LOCATION", "DEVICE", "BEHAVIOR", "TIME"}
		anomalyType = anomalyTypes[rand.Intn(len(anomalyTypes))]

		if anomalyType == "LOCATION" {
			reasons = append(reasons, "Login from unusual location")
		} else if anomalyType == "DEVICE" {
			reasons = append(reasons, "New device used for verification")
		} else if anomalyType == "BEHAVIOR" {
			reasons = append(reasons, "Unusual verification pattern")
		} else if anomalyType == "TIME" {
			reasons = append(reasons, "Verification attempt at unusual time")
		}
	}

	deviceInfoMap := map[string]string{
		"ip_address":   deviceInfo.IPAddress,
		"user_agent":   deviceInfo.UserAgent,
		"device_type":  deviceInfo.DeviceType,
		"os":           deviceInfo.OS,
		"browser":      deviceInfo.Browser,
		"mac_address":  deviceInfo.MacAddress,
		"location":     deviceInfo.Location,
		"coordinates":  deviceInfo.Coordinates,
		"isp":          deviceInfo.ISP,
		"country_code": deviceInfo.CountryCode,
	}

	result := &model.AnomalyDetectionResult{
		ID:             uuid.New(),
		VerificationID: verificationID,
		UserID:         userID,
		IsAnomaly:      isAnomaly,
		AnomalyScore:   anomalyScore,
		AnomalyType:    anomalyType,
		Reasons:        reasons,
		DeviceInfo:     deviceInfoMap,
		CreatedAt:      time.Now(),
	}

	// Save the result
	err := s.repo.SaveAnomalyDetectionResult(ctx, result)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to save anomaly detection result")
		return nil, err
	}

	return result, nil
}

// ProcessDocument processes a document through all AI checks
func (s *AIService) ProcessDocument(ctx context.Context, documentID uuid.UUID, selfieID uuid.UUID, verificationID uuid.UUID, deviceInfo model.DeviceInfo) (*model.Verification, error) {
	// Get verification
	verification, err := s.verificationRepo.GetByID(ctx, verificationID)
	if err != nil {
		s.logger.WithError(err).WithField("verification_id", verificationID).Error("Failed to get verification for AI processing")
		return nil, err
	}

	// 1. Analyze document
	docResult, err := s.AnalyzeDocument(ctx, documentID, verificationID)
	if err != nil {
		s.logger.WithError(err).WithField("document_id", documentID).Error("Document analysis failed")
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Document analysis failed: " + err.Error()
		s.verificationRepo.Update(ctx, verification)
		return verification, err
	}

	// 2. Match faces if selfie provided
	var faceResult *model.FaceMatchResult
	if selfieID != uuid.Nil {
		faceResult, err = s.MatchFaces(ctx, documentID, selfieID, verificationID)
		if err != nil {
			s.logger.WithError(err).WithField("document_id", documentID).WithField("selfie_id", selfieID).Error("Face matching failed")
			verification.Status = model.VerificationStatusFailed
			verification.Notes = "Face matching failed: " + err.Error()
			s.verificationRepo.Update(ctx, verification)
			return verification, err
		}
	}

	// 3. Analyze risk
	riskResult, err := s.AnalyzeRisk(ctx, verification.UserID, verificationID, deviceInfo)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", verification.UserID).Error("Risk analysis failed")
		// Continue despite risk analysis failure
	}

	// 4. Detect anomalies
	anomalyResult, err := s.DetectAnomalies(ctx, verification.UserID, verificationID, deviceInfo)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", verification.UserID).Error("Anomaly detection failed")
		// Continue despite anomaly detection failure
	}

	// Determine verification result based on all checks
	if !docResult.IsAuthentic {
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Document appears to be inauthentic"
	} else if selfieID != uuid.Nil && !faceResult.IsMatch {
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Face in selfie does not match document"
	} else if riskResult != nil && riskResult.RiskLevel == "HIGH" {
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "High risk verification"
	} else if anomalyResult != nil && anomalyResult.IsAnomaly {
		verification.Status = model.VerificationStatusFailed
		verification.Notes = "Anomalous verification detected"
	} else {
		verification.Status = model.VerificationStatusCompleted
		verification.Notes = "All verification checks passed"
		now := time.Now()
		verification.CompletedAt = &now
	}

	// Update verification
	err = s.verificationRepo.Update(ctx, verification)
	if err != nil {
		s.logger.WithError(err).WithField("verification_id", verificationID).Error("Failed to update verification after AI processing")
		return nil, err
	}

	return verification, nil
}

// GetDocumentAnalysis gets document analysis result
func (s *AIService) GetDocumentAnalysis(ctx context.Context, analysisID uuid.UUID) (*model.DocumentAnalysisResult, error) {
	return s.repo.GetDocumentAnalysis(ctx, analysisID)
}

// GetFaceMatchResult gets face match result
func (s *AIService) GetFaceMatchResult(ctx context.Context, matchID uuid.UUID) (*model.FaceMatchResult, error) {
	return s.repo.GetFaceMatchResult(ctx, matchID)
}

// GetRiskAnalysisResult gets risk analysis result
func (s *AIService) GetRiskAnalysisResult(ctx context.Context, analysisID uuid.UUID) (*model.RiskAnalysisResult, error) {
	return s.repo.GetRiskAnalysisResult(ctx, analysisID)
}

// GetAnomalyDetectionResult gets anomaly detection result
func (s *AIService) GetAnomalyDetectionResult(ctx context.Context, detectionID uuid.UUID) (*model.AnomalyDetectionResult, error) {
	return s.repo.GetAnomalyDetectionResult(ctx, detectionID)
}

// GetAIModelInfo gets AI model information
func (s *AIService) GetAIModelInfo(ctx context.Context, modelID uuid.UUID) (*model.AIModelInfo, error) {
	return s.repo.GetAIModelInfo(ctx, modelID)
}

// ListAIModels lists all AI models
func (s *AIService) ListAIModels(ctx context.Context) ([]*model.AIModelInfo, error) {
	return s.repo.ListAIModels(ctx)
}

// ExtractTextFromImage extracts text from an image
func (s *AIService) ExtractTextFromImage(ctx context.Context, imageReader io.Reader) (map[string]string, error) {
	// In a real implementation, we would:
	// 1. Use AWS Textract to extract text from the image
	// 2. Parse the extracted text into structured data
	// 3. Return the structured data

	// For this demo, we'll return mock data
	return map[string]string{
		"full_name":       "John Smith",
		"document_number": "X123456789",
		"date_of_birth":   "1990-01-01",
		"expiry_date":     "2030-01-01",
		"issuing_country": "United States",
	}, nil
}
