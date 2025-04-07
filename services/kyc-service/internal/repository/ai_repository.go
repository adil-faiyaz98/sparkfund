package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sparkfund/services/kyc-service/internal/model"
)

// AIRepository handles AI-related database operations
type AIRepository struct {
	db *gorm.DB
}

// NewAIRepository creates a new AI repository
func NewAIRepository(db *gorm.DB) *AIRepository {
	return &AIRepository{
		db: db,
	}
}

// SaveDocumentAnalysis saves a document analysis result
func (r *AIRepository) SaveDocumentAnalysis(ctx context.Context, analysis *model.DocumentAnalysisResult) error {
	return r.db.WithContext(ctx).Create(analysis).Error
}

// GetDocumentAnalysis gets a document analysis result by ID
func (r *AIRepository) GetDocumentAnalysis(ctx context.Context, id uuid.UUID) (*model.DocumentAnalysisResult, error) {
	var analysis model.DocumentAnalysisResult
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&analysis).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("document analysis not found")
		}
		return nil, err
	}
	return &analysis, nil
}

// GetDocumentAnalysisByVerificationID gets document analysis results by verification ID
func (r *AIRepository) GetDocumentAnalysisByVerificationID(ctx context.Context, verificationID uuid.UUID) ([]*model.DocumentAnalysisResult, error) {
	var analyses []*model.DocumentAnalysisResult
	err := r.db.WithContext(ctx).Where("verification_id = ?", verificationID).Find(&analyses).Error
	if err != nil {
		return nil, err
	}
	return analyses, nil
}

// SaveFaceMatchResult saves a face match result
func (r *AIRepository) SaveFaceMatchResult(ctx context.Context, result *model.FaceMatchResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

// GetFaceMatchResult gets a face match result by ID
func (r *AIRepository) GetFaceMatchResult(ctx context.Context, id uuid.UUID) (*model.FaceMatchResult, error) {
	var result model.FaceMatchResult
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("face match result not found")
		}
		return nil, err
	}
	return &result, nil
}

// GetFaceMatchResultByVerificationID gets face match results by verification ID
func (r *AIRepository) GetFaceMatchResultByVerificationID(ctx context.Context, verificationID uuid.UUID) ([]*model.FaceMatchResult, error) {
	var results []*model.FaceMatchResult
	err := r.db.WithContext(ctx).Where("verification_id = ?", verificationID).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SaveRiskAnalysisResult saves a risk analysis result
func (r *AIRepository) SaveRiskAnalysisResult(ctx context.Context, result *model.RiskAnalysisResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

// GetRiskAnalysisResult gets a risk analysis result by ID
func (r *AIRepository) GetRiskAnalysisResult(ctx context.Context, id uuid.UUID) (*model.RiskAnalysisResult, error) {
	var result model.RiskAnalysisResult
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("risk analysis result not found")
		}
		return nil, err
	}
	return &result, nil
}

// GetRiskAnalysisResultByVerificationID gets risk analysis results by verification ID
func (r *AIRepository) GetRiskAnalysisResultByVerificationID(ctx context.Context, verificationID uuid.UUID) ([]*model.RiskAnalysisResult, error) {
	var results []*model.RiskAnalysisResult
	err := r.db.WithContext(ctx).Where("verification_id = ?", verificationID).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SaveAnomalyDetectionResult saves an anomaly detection result
func (r *AIRepository) SaveAnomalyDetectionResult(ctx context.Context, result *model.AnomalyDetectionResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

// GetAnomalyDetectionResult gets an anomaly detection result by ID
func (r *AIRepository) GetAnomalyDetectionResult(ctx context.Context, id uuid.UUID) (*model.AnomalyDetectionResult, error) {
	var result model.AnomalyDetectionResult
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("anomaly detection result not found")
		}
		return nil, err
	}
	return &result, nil
}

// GetAnomalyDetectionResultByVerificationID gets anomaly detection results by verification ID
func (r *AIRepository) GetAnomalyDetectionResultByVerificationID(ctx context.Context, verificationID uuid.UUID) ([]*model.AnomalyDetectionResult, error) {
	var results []*model.AnomalyDetectionResult
	err := r.db.WithContext(ctx).Where("verification_id = ?", verificationID).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SaveAIModelInfo saves AI model information
func (r *AIRepository) SaveAIModelInfo(ctx context.Context, info *model.AIModelInfo) error {
	return r.db.WithContext(ctx).Create(info).Error
}

// GetAIModelInfo gets AI model information by ID
func (r *AIRepository) GetAIModelInfo(ctx context.Context, id uuid.UUID) (*model.AIModelInfo, error) {
	var info model.AIModelInfo
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&info).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("AI model info not found")
		}
		return nil, err
	}
	return &info, nil
}

// ListAIModels lists all AI models
func (r *AIRepository) ListAIModels(ctx context.Context) ([]*model.AIModelInfo, error) {
	var models []*model.AIModelInfo
	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// GetAIModelByType gets AI model information by type
func (r *AIRepository) GetAIModelByType(ctx context.Context, modelType string) (*model.AIModelInfo, error) {
	var info model.AIModelInfo
	err := r.db.WithContext(ctx).Where("type = ?", modelType).Order("version DESC").First(&info).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("AI model info not found")
		}
		return nil, err
	}
	return &info, nil
}
