package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adil-faiyaz98/sparkfund/services/kyc-service/internal/cache"
	"github.com/adil-faiyaz98/sparkfund/services/kyc-service/internal/model"
)

// CacheService handles caching operations
type CacheService struct {
	cache cache.Cache
}

// NewCacheService creates a new cache service
func NewCacheService(cache cache.Cache) *CacheService {
	return &CacheService{
		cache: cache,
	}
}

// GetDocument retrieves a document from the cache
func (s *CacheService) GetDocument(ctx context.Context, id string) (*model.Document, error) {
	var document model.Document
	err := s.cache.Get(ctx, fmt.Sprintf("document:%s", id), &document)
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// SetDocument stores a document in the cache
func (s *CacheService) SetDocument(ctx context.Context, document *model.Document, ttl time.Duration) error {
	return s.cache.Set(ctx, fmt.Sprintf("document:%s", document.ID.String()), document, ttl)
}

// DeleteDocument removes a document from the cache
func (s *CacheService) DeleteDocument(ctx context.Context, id string) error {
	return s.cache.Delete(ctx, fmt.Sprintf("document:%s", id))
}

// GetVerification retrieves a verification from the cache
func (s *CacheService) GetVerification(ctx context.Context, id string) (*model.Verification, error) {
	var verification model.Verification
	err := s.cache.Get(ctx, fmt.Sprintf("verification:%s", id), &verification)
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// SetVerification stores a verification in the cache
func (s *CacheService) SetVerification(ctx context.Context, verification *model.Verification, ttl time.Duration) error {
	return s.cache.Set(ctx, fmt.Sprintf("verification:%s", verification.ID.String()), verification, ttl)
}

// DeleteVerification removes a verification from the cache
func (s *CacheService) DeleteVerification(ctx context.Context, id string) error {
	return s.cache.Delete(ctx, fmt.Sprintf("verification:%s", id))
}

// GetProfile retrieves a profile from the cache
func (s *CacheService) GetProfile(ctx context.Context, userID string) (*model.KYC, error) {
	var profile model.KYC
	err := s.cache.Get(ctx, fmt.Sprintf("profile:%s", userID), &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// SetProfile stores a profile in the cache
func (s *CacheService) SetProfile(ctx context.Context, profile *model.KYC, ttl time.Duration) error {
	return s.cache.Set(ctx, fmt.Sprintf("profile:%s", profile.UserID.String()), profile, ttl)
}

// DeleteProfile removes a profile from the cache
func (s *CacheService) DeleteProfile(ctx context.Context, userID string) error {
	return s.cache.Delete(ctx, fmt.Sprintf("profile:%s", userID))
}

// GetDocumentAnalysis retrieves a document analysis from the cache
func (s *CacheService) GetDocumentAnalysis(ctx context.Context, id string) (*model.DocumentAnalysisResult, error) {
	var analysis model.DocumentAnalysisResult
	err := s.cache.Get(ctx, fmt.Sprintf("analysis:%s", id), &analysis)
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

// SetDocumentAnalysis stores a document analysis in the cache
func (s *CacheService) SetDocumentAnalysis(ctx context.Context, analysis *model.DocumentAnalysisResult, ttl time.Duration) error {
	return s.cache.Set(ctx, fmt.Sprintf("analysis:%s", analysis.ID.String()), analysis, ttl)
}

// DeleteDocumentAnalysis removes a document analysis from the cache
func (s *CacheService) DeleteDocumentAnalysis(ctx context.Context, id string) error {
	return s.cache.Delete(ctx, fmt.Sprintf("analysis:%s", id))
}

// GetFaceMatch retrieves a face match from the cache
func (s *CacheService) GetFaceMatch(ctx context.Context, id string) (*model.FaceMatchResult, error) {
	var match model.FaceMatchResult
	err := s.cache.Get(ctx, fmt.Sprintf("face_match:%s", id), &match)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

// SetFaceMatch stores a face match in the cache
func (s *CacheService) SetFaceMatch(ctx context.Context, match *model.FaceMatchResult, ttl time.Duration) error {
	return s.cache.Set(ctx, fmt.Sprintf("face_match:%s", match.ID.String()), match, ttl)
}

// DeleteFaceMatch removes a face match from the cache
func (s *CacheService) DeleteFaceMatch(ctx context.Context, id string) error {
	return s.cache.Delete(ctx, fmt.Sprintf("face_match:%s", id))
}

// GetRiskAnalysis retrieves a risk analysis from the cache
func (s *CacheService) GetRiskAnalysis(ctx context.Context, id string) (*model.RiskAnalysisResult, error) {
	var analysis model.RiskAnalysisResult
	err := s.cache.Get(ctx, fmt.Sprintf("risk_analysis:%s", id), &analysis)
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

// SetRiskAnalysis stores a risk analysis in the cache
func (s *CacheService) SetRiskAnalysis(ctx context.Context, analysis *model.RiskAnalysisResult, ttl time.Duration) error {
	return s.cache.Set(ctx, fmt.Sprintf("risk_analysis:%s", analysis.ID.String()), analysis, ttl)
}

// DeleteRiskAnalysis removes a risk analysis from the cache
func (s *CacheService) DeleteRiskAnalysis(ctx context.Context, id string) error {
	return s.cache.Delete(ctx, fmt.Sprintf("risk_analysis:%s", id))
}

// GetAnomalyDetection retrieves an anomaly detection from the cache
func (s *CacheService) GetAnomalyDetection(ctx context.Context, id string) (*model.AnomalyDetectionResult, error) {
	var detection model.AnomalyDetectionResult
	err := s.cache.Get(ctx, fmt.Sprintf("anomaly_detection:%s", id), &detection)
	if err != nil {
		return nil, err
	}
	return &detection, nil
}

// SetAnomalyDetection stores an anomaly detection in the cache
func (s *CacheService) SetAnomalyDetection(ctx context.Context, detection *model.AnomalyDetectionResult, ttl time.Duration) error {
	return s.cache.Set(ctx, fmt.Sprintf("anomaly_detection:%s", detection.ID.String()), detection, ttl)
}

// DeleteAnomalyDetection removes an anomaly detection from the cache
func (s *CacheService) DeleteAnomalyDetection(ctx context.Context, id string) error {
	return s.cache.Delete(ctx, fmt.Sprintf("anomaly_detection:%s", id))
}

// IncrementCounter increments a counter in the cache
func (s *CacheService) IncrementCounter(ctx context.Context, key string, value int64) (int64, error) {
	return s.cache.Increment(ctx, fmt.Sprintf("counter:%s", key), value)
}

// GetCounter retrieves a counter from the cache
func (s *CacheService) GetCounter(ctx context.Context, key string) (int64, error) {
	var counter int64
	err := s.cache.Get(ctx, fmt.Sprintf("counter:%s", key), &counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// SetCounter sets a counter in the cache
func (s *CacheService) SetCounter(ctx context.Context, key string, value int64, ttl time.Duration) error {
	return s.cache.Set(ctx, fmt.Sprintf("counter:%s", key), value, ttl)
}

// Close closes the cache connection
func (s *CacheService) Close() error {
	return s.cache.Close()
}
