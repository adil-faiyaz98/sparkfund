package ai

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis"
	"github.com/sirupsen/logrus"
	"github.com/sparkfund/services/kyc-service/internal/models"
)

// DocumentVerificationService uses AI for document verification
type DocumentVerificationService struct {
	documentClient DocumentAIClient
	cache          *redis.Client
	logger         *logrus.Logger
	config         *config.AIConfig
}

// NewDocumentVerificationService creates a new document verification service
func NewDocumentVerificationService(
	documentClient DocumentAIClient,
	cache *redis.Client,
	logger *logrus.Logger,
	config *config.AIConfig,
) *DocumentVerificationService {
	return &DocumentVerificationService{
		documentClient: documentClient,
		cache:          cache,
		logger:         logger,
		config:         config,
	}
}

// VerifyIdentityDocument performs AI verification of identity documents
func (s *DocumentVerificationService) VerifyIdentityDocument(ctx context.Context, document *models.IdentityDocument) (*models.DocumentVerificationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DocumentVerificationService.VerifyIdentityDocument")
	defer span.Finish()

	// Check if we have cached results for this document
	cacheKey := s.generateCacheKey(document.DocumentID, document.Type)
	cachedResult, err := s.getFromCache(ctx, cacheKey)
	if err == nil && cachedResult != nil {
		s.logger.WithFields(logrus.Fields{
			"document_id":   document.DocumentID,
			"document_type": document.Type,
			"cache_hit":     true,
		}).Info("Using cached document verification result")
		return cachedResult, nil
	}

	// Use circuit breaker for external AI service call
	var result *models.DocumentVerificationResult
	err = s.withCircuitBreaker(ctx, "document_verification", func() error {
		// Process document front side
		frontAnalysis, err := s.documentClient.AnalyzeDocument(ctx, &DocumentAnalysisRequest{
			Image:        document.FrontImage,
			DocumentType: document.Type,
			Side:         "front",
			Options:      s.config.DocumentVerificationOptions,
		})
		if err != nil {
			return fmt.Errorf("failed to analyze document front: %w", err)
		}

		// Process document back side (if available)
		var backAnalysis *DocumentAnalysisResponse
		if document.BackImage != "" {
			backAnalysis, err = s.documentClient.AnalyzeDocument(ctx, &DocumentAnalysisRequest{
				Image:        document.BackImage,
				DocumentType: document.Type,
				Side:         "back",
				Options:      s.config.DocumentVerificationOptions,
			})
			if err != nil {
				return fmt.Errorf("failed to analyze document back: %w", err)
			}
		}

		// Validate the document based on AI analysis
		result = s.validateDocument(frontAnalysis, backAnalysis, document)
		return nil
	})

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"document_id":   document.DocumentID,
			"document_type": document.Type,
		}).Error("Document verification failed")
		return nil, err
	}

	// Cache successful results
	if result.OverallVerificationStatus == models.VerificationStatusApproved ||
		result.OverallVerificationStatus == models.VerificationStatusNeedsReview {
		if err := s.saveToCache(ctx, cacheKey, result, s.config.DocumentCacheTTL); err != nil {
			s.logger.WithError(err).Warn("Failed to cache document verification result")
		}
	}

	// Log results for audit
	s.logger.WithFields(logrus.Fields{
		"document_id":            document.DocumentID,
		"document_type":          document.Type,
		"verification_status":    result.OverallVerificationStatus,
		"confidence_score":       result.ConfidenceScore,
		"extracted_fields_count": len(result.ExtractedFields),
	}).Info("Document verification completed")

	return result, nil
}

// BiometricVerification compares a selfie with the document photo
func (s *DocumentVerificationService) BiometricVerification(ctx context.Context, selfieImage string, documentImage string, userID string) (*models.BiometricVerificationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DocumentVerificationService.BiometricVerification")
	defer span.Finish()

	// Generate cache key for biometric verification
	cacheKey := s.generateBiometricCacheKey(userID, documentImage, selfieImage)

	// Check cache first
	cachedResult, err := s.getBiometricFromCache(ctx, cacheKey)
	if err == nil && cachedResult != nil {
		s.logger.WithField("user_id", userID).Info("Using cached biometric verification result")
		return cachedResult, nil
	}

	// Use circuit breaker for AI service
	var result *models.BiometricVerificationResult
	err = s.withCircuitBreaker(ctx, "biometric_verification", func() error {
		// Call face comparison API
		resp, err := s.documentClient.CompareFaces(ctx, &FaceComparisonRequest{
			FaceImage1: documentImage,
			FaceImage2: selfieImage,
			Options:    s.config.BiometricVerificationOptions,
		})
		if err != nil {
			return fmt.Errorf("face comparison failed: %w", err)
		}

		// Process the comparison results
		result = &models.BiometricVerificationResult{
			MatchScore:          resp.MatchScore,
			LivenessScore:       resp.LivenessScore,
			MatchStatus:         s.determineMatchStatus(resp.MatchScore),
			VerificationTime:    time.Now().UTC(),
			PotentialFraudFlags: s.detectFraudIndicators(resp),
		}

		return nil
	})

	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Biometric verification failed")
		return nil, err
	}

	// Cache successful verifications
	if result.MatchStatus == models.BiometricMatchStatusMatch ||
		result.MatchStatus == models.BiometricMatchStatusPossibleMatch {
		if err := s.saveBiometricToCache(ctx, cacheKey, result, s.config.BiometricCacheTTL); err != nil {
			s.logger.WithError(err).Warn("Failed to cache biometric verification result")
		}
	}

	// Log results
	s.logger.WithFields(logrus.Fields{
		"user_id":        userID,
		"match_score":    result.MatchScore,
		"match_status":   result.MatchStatus,
		"liveness_score": result.LivenessScore,
	}).Info("Biometric verification completed")

	return result, nil
}

// Helper functions for document verification
func (s *DocumentVerificationService) prepareImageForAI(imageData string) (string, error) {
	// If the image is already a URL, return as is
	if strings.HasPrefix(imageData, "http") {
		return imageData, nil
	}

	// If base64, validate and clean up the data
	if !strings.Contains(imageData, ";base64,") {
		return "", errors.New("unsupported image format")
	}

	// Extract the base64 content
	parts := strings.Split(imageData, ";base64,")
	if len(parts) != 2 {
		return "", errors.New("invalid base64 image format")
	}

	// Validate base64 content
	_, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid base64 encoding: %w", err)
	}

	return imageData, nil
}

func (s *DocumentVerificationService) processAIResult(aiResult *DocumentAIResult, document *models.IdentityDocument) *models.DocumentVerificationResult {
	result := &models.DocumentVerificationResult{
		DocumentID:       document.ID,
		UserID:           document.UserID,
		VerificationType: "ai",
		AuthenticitySore: aiResult.AuthenticitySore,
		DataMatchScore:   aiResult.DataMatchScore,
		VerificationTime: time.Now().UTC(),
		ModelVersion:     aiResult.ModelVersion,
		ExtractedData:    aiResult.ExtractedData,
		Flags:            make([]string, 0),
	}

	// Add flags for potential issues
	if aiResult.IsSuspectedFraud {
		result.Flags = append(result.Flags, "SUSPECTED_FRAUD")
	}

	if aiResult.IsExpired {
		result.Flags = append(result.Flags, "DOCUMENT_EXPIRED")
	}

	if aiResult.HasManipulationMarks {
		result.Flags = append(result.Flags, "MANIPULATION_DETECTED")
	}

	// Validate data consistency
	if document.DateOfBirth.Format("2006-01-02") != aiResult.ExtractedData.DateOfBirth {
		result.Flags = append(result.Flags, "DOB_MISMATCH")
	}

	if !strings.EqualFold(document.FirstName, aiResult.ExtractedData.FirstName) {
		result.Flags = append(result.Flags, "FIRSTNAME_MISMATCH")
	}

	if !strings.EqualFold(document.LastName, aiResult.ExtractedData.LastName) {
		result.Flags = append(result.Flags, "LASTNAME_MISMATCH")
	}

	// Determine validity based on scores and flags
	result.IsValid = aiResult.AuthenticitySore >= s.config.MinAuthenticityScore &&
		aiResult.DataMatchScore >= s.config.MinDataMatchScore &&
		!aiResult.IsSuspectedFraud &&
		!aiResult.IsExpired

	return result
}

// Cache functions
func (s *DocumentVerificationService) getFromCache(ctx context.Context, key string) (*models.DocumentVerificationResult, error) {
	data, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var result models.DocumentVerificationResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *DocumentVerificationService) saveToCache(ctx context.Context, key string, result *models.DocumentVerificationResult, ttl time.Duration) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, key, string(data), ttl).Err()
}

// Helper circuit breaker function
func (s *DocumentVerificationService) executeWithCircuitBreaker(ctx context.Context, fn func() error) error {
	circuitBreaker := s.documentClient.GetCircuitBreaker()
	if circuitBreaker == nil {
		return fn()
	}

	return circuitBreaker.Execute(fn)
}

// Utility function to hash strings
func hashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Get/save biometric results from/to cache
func (s *DocumentVerificationService) getBiometricFromCache(ctx context.Context, key string) (*models.BiometricVerificationResult, error) {
	data, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var result models.BiometricVerificationResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *DocumentVerificationService) saveBiometricToCache(ctx context.Context, key string, result *models.BiometricVerificationResult, ttl time.Duration) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, key, string(data), ttl).Err()
}

func (s *DocumentVerificationService) validateDocument(frontAnalysis *DocumentAnalysisResponse, backAnalysis *DocumentAnalysisResponse, document *models.IdentityDocument) *models.DocumentVerificationResult {
	result := &models.DocumentVerificationResult{
		DocumentID:       document.DocumentID,
		VerificationTime: time.Now().UTC(),
		ExtractedFields:  make(map[string]string),
		ValidationChecks: make([]models.ValidationCheck, 0),
		SecurityFeatures: make([]models.SecurityFeatureCheck, 0),
		ConfidenceScore:  frontAnalysis.OverallConfidence,
	}

	// Extract fields from document
	for k, v := range frontAnalysis.ExtractedFields {
		result.ExtractedFields[k] = v
	}
	if backAnalysis != nil {
		for k, v := range backAnalysis.ExtractedFields {
			result.ExtractedFields[k] = v
		}
	}

	// Perform validation checks
	s.performBasicValidationChecks(result, document)
	s.performSecurityFeatureChecks(result, frontAnalysis, backAnalysis)
	s.performConsistencyChecks(result, document)

	// Determine overall verification status
	result.OverallVerificationStatus = s.determineOverallStatus(result)

	return result
}

// Helper circuit breaker function
func (s *DocumentVerificationService) withCircuitBreaker(ctx context.Context, operation string, fn func() error) error {
	circuitBreaker := NewCircuitBreaker(fmt.Sprintf("document_ai_%s", operation))
	return circuitBreaker.Execute(fn)
}

// Utility function to hash strings
func (s *DocumentVerificationService) generateCacheKey(documentID, documentType string) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s:%s", documentID, documentType)))
	return fmt.Sprintf("doc_verify:%s", base64.URLEncoding.EncodeToString(hash.Sum(nil)))
}

func (s *DocumentVerificationService) generateBiometricCacheKey(userID, documentImage, selfieImage string) string {
	docHash := sha256.New()
	docHash.Write([]byte(documentImage))
	docHashStr := base64.URLEncoding.EncodeToString(docHash.Sum(nil)[:8]) // Use first 8 bytes

	selfieHash := sha256.New()
	selfieHash.Write([]byte(selfieImage))
	selfieHashStr := base64.URLEncoding.EncodeToString(selfieHash.Sum(nil)[:8])

	return fmt.Sprintf("bio_verify:%s:%s:%s", userID, docHashStr, selfieHashStr)
}

func (s *DocumentVerificationService) determineMatchStatus(score float64) models.BiometricMatchStatus {
	if score >= s.config.BiometricHighConfidenceThreshold {
		return models.BiometricMatchStatusMatch
	} else if score >= s.config.BiometricMediumConfidenceThreshold {
		return models.BiometricMatchStatusPossibleMatch
	}
	return models.BiometricMatchStatusNoMatch
}

func (s *DocumentVerificationService) detectFraudIndicators(resp *FaceComparisonResponse) []string {
	var flags []string

	if resp.SpoofingDetected {
		flags = append(flags, "POTENTIAL_SPOOFING_DETECTED")
	}

	if resp.MaskDetected {
		flags = append(flags, "FACE_MASK_DETECTED")
	}

	if resp.DeepfakeScore > s.config.DeepfakeThreshold {
		flags = append(flags, "POTENTIAL_DEEPFAKE")
	}

	return flags
}

// AnalyzeCustomerRisk performs AI-based AML risk analysis on a customer
func (s *AMLRiskAnalysisService) AnalyzeCustomerRisk(ctx context.Context, customer *models.Customer) (*models.CustomerRiskAssessment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AMLRiskAnalysisService.AnalyzeCustomerRisk")
	defer span.Finish()

	s.logger.WithField("customer_id", customer.ID).Info("Starting customer risk analysis")

	// Collect customer data for analysis
	customerData, err := s.prepareCustomerData(ctx, customer)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare customer data: %w", err)
	}

	// Use AI model to calculate base risk score
	riskScore, riskFactors, err := s.mlClient.PredictCustomerRisk(ctx, customerData)
	if err != nil {
		return nil, fmt.Errorf("failed to predict customer risk: %w", err)
	}

	// Perform graph analysis to detect complex relationships
	graphAnalysisResults, err := s.performGraphAnalysis(ctx, customer.ID)
	if err != nil {
		s.logger.WithError(err).WithField("customer_id", customer.ID).Warn("Graph analysis failed, continuing with basic risk assessment")
		// Continue with basic assessment even if graph analysis fails
	} else {
		// Adjust risk score based on graph analysis results
		riskScore = s.adjustRiskScoreFromGraphAnalysis(riskScore, graphAnalysisResults)
		for k, v := range graphAnalysisResults.RiskFactors {
			riskFactors[k] = v
		}
	}

	// Check watchlists and sanctions
	watchlistResults, err := s.checkWatchlists(ctx, customer)
	if err != nil {
		s.logger.WithError(err).WithField("customer_id", customer.ID).Warn("Watchlist check failed, continuing with available data")
	} else if watchlistResults.Matches {
		// Significant risk increase if on watchlists
		riskScore = math.Min(100, riskScore*1.5)
		riskFactors["watchlist_match"] = watchlistResults.Score
	}

	// Determine risk level based on score
	riskLevel := s.determineRiskLevel(riskScore)

	// Generate recommendations based on risk level and factors
	recommendations := s.generateRiskMitigationRecommendations(riskLevel, riskFactors)

	// Create the final assessment
	assessment := &models.CustomerRiskAssessment{
		CustomerID:       customer.ID,
		RiskScore:        riskScore,
		RiskLevel:        riskLevel,
		RiskFactors:      riskFactors,
		Recommendations:  recommendations,
		AssessmentDate:   time.Now().UTC(),
		WatchlistMatches: watchlistResults.Details,
		DataSources:      []string{"customer_data", "transaction_history", "watchlists"},
	}

	if graphAnalysisResults != nil {
		assessment.NetworkRiskInfo = &models.NetworkRiskInfo{
			HighRiskConnections:   graphAnalysisResults.HighRiskConnections,
			AnomalousTransactions: graphAnalysisResults.AnomalousTransactionCount,
			ClusterRisk:           graphAnalysisResults.ClusterRiskScore,
		}
	}

	s.logger.WithFields(logrus.Fields{
		"customer_id":        customer.ID,
		"risk_score":         riskScore,
		"risk_level":         riskLevel,
		"risk_factors_count": len(riskFactors),
	}).Info("Customer risk analysis completed")

	return assessment, nil
}

// AnalyzeTransactionRisk evaluates AML risk for a specific transaction
func (s *AMLRiskAnalysisService) AnalyzeTransactionRisk(ctx context.Context, transaction *models.Transaction) (*models.TransactionRiskAssessment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AMLRiskAnalysisService.AnalyzeTransactionRisk")
	defer span.Finish()

	s.logger.WithField("transaction_id", transaction.ID).Info("Starting transaction risk analysis")

	// Get customer information
	customer, err := s.transactionAPI.GetCustomerInfo(ctx, transaction.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer info: %w", err)
	}

	// Get transaction history for pattern analysis
	transactionHistory, err := s.transactionAPI.GetTransactionHistory(ctx, transaction.CustomerID, s.config.TransactionLookbackDays)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to get complete transaction history, using available data")
		// Continue with potentially incomplete history
	}

	// Prepare features for ML model
	features, err := s.prepareTransactionFeatures(transaction, customer, transactionHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transaction features: %w", err)
	}

	// Use AI model for risk prediction
	riskScore, anomalyScore, patterns, err := s.mlClient.PredictTransactionRisk(ctx, features)
	if err != nil {
		return nil, fmt.Errorf("failed to predict transaction risk: %w", err)
	}

	// Get network context for the transaction
	networkContext, err := s.getTransactionNetworkContext(ctx, transaction)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to get network context, continuing with available data")
	} else {
		// Adjust risk score based on network insights
		riskScore = s.adjustRiskScoreFromNetwork(riskScore, networkContext)
	}

	// Extract risk patterns and signals
	riskSignals := s.extractRiskSignals(patterns, anomalyScore, networkContext)

	// Determine if transaction should be allowed, need review, or be blocked
	decision, rationale := s.determineTransactionDecision(riskScore, riskSignals, transaction)

	// Create assessment result
	assessment := &models.TransactionRiskAssessment{
		TransactionID:       transaction.ID,
		CustomerID:          transaction.CustomerID,
		RiskScore:           riskScore,
		AnomalyScore:        anomalyScore,
		Decision:            decision,
		DecisionRationale:   rationale,
		RiskSignals:         riskSignals,
		AssessmentTimestamp: time.Now().UTC(),
		RequiredReviews:     s.determineRequiredReviews(riskScore, riskSignals),
	}

	s.logger.WithFields(logrus.Fields{
		"transaction_id": transaction.ID,
		"risk_score":     riskScore,
		"anomaly_score":  anomalyScore,
		"decision":       decision,
	}).Info("Transaction risk analysis completed")

	return assessment, nil
}

// Helper functions
func (s *AMLRiskAnalysisService) prepareCustomerData(ctx context.Context, customer *models.Customer) (map[string]interface{}, error) {
	// Collect and format customer data for ML model
	data := map[string]interface{}{
		"customer_id":            customer.ID,
		"risk_category":          customer.RiskCategory,
		"account_age_days":       time.Since(customer.AccountCreatedAt).Hours() / 24,
		"nationality":            customer.Nationality,
		"residence_country":      customer.ResidenceCountry,
		"occupation":             customer.Occupation,
		"business_type":          customer.BusinessType,
		"politically_exposed":    customer.PoliticallyExposed,
		"high_risk_jurisdiction": s.isHighRiskJurisdiction(customer.ResidenceCountry),
	}

	// Get additional data from transaction API
	return data, nil
}
