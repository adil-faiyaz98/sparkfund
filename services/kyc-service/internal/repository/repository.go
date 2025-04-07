package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"sparkfund/services/kyc-service/internal/model"
)

// Repository is the base repository interface
type Repository interface {
	GetDB() *gorm.DB
}

// BaseRepository is the base repository implementation
type BaseRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *gorm.DB, logger *logrus.Logger) *BaseRepository {
	return &BaseRepository{
		db:     db,
		logger: logger,
	}
}

// GetDB returns the database connection
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}

// UserRepository handles user-related database operations
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// List lists all users with pagination
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// SessionRepository handles session-related database operations
type SessionRepository struct {
	*BaseRepository
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *gorm.DB, logger *logrus.Logger) *SessionRepository {
	return &SessionRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Create creates a new session
func (r *SessionRepository) Create(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID gets a session by ID
func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Session, error) {
	var session model.Session
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &session, nil
}

// GetByRefreshToken gets a session by refresh token
func (r *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	var session model.Session
	err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &session, nil
}

// GetByUserID gets sessions by user ID
func (r *SessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Session, error) {
	var sessions []*model.Session
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// Update updates a session
func (r *SessionRepository) Update(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// Delete deletes a session
func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Session{}, id).Error
}

// DeleteByUserID deletes all sessions for a user
func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.Session{}).Error
}

// DeleteExpired deletes all expired sessions
func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&model.Session{}).Error
}

// DocumentRepository handles document-related database operations
type DocumentRepository struct {
	*BaseRepository
}

// NewDocumentRepository creates a new document repository
func NewDocumentRepository(db *gorm.DB, logger *logrus.Logger) *DocumentRepository {
	return &DocumentRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Create creates a new document
func (r *DocumentRepository) Create(ctx context.Context, document *model.Document) error {
	return r.db.WithContext(ctx).Create(document).Error
}

// GetByID gets a document by ID
func (r *DocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Document, error) {
	var document model.Document
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&document).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("document not found")
		}
		return nil, err
	}
	return &document, nil
}

// GetByUserID gets documents by user ID
func (r *DocumentRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Document{}).Where("user_id = ?", userID)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documents).Error
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// GetByType gets documents by type
func (r *DocumentRepository) GetByType(ctx context.Context, documentType string, page, pageSize int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Document{}).Where("type = ?", documentType)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documents).Error
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// Update updates a document
func (r *DocumentRepository) Update(ctx context.Context, document *model.Document) error {
	return r.db.WithContext(ctx).Save(document).Error
}

// Delete deletes a document
func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Document{}, id).Error
}

// List lists all documents with pagination
func (r *DocumentRepository) List(ctx context.Context, page, pageSize int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Document{})

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documents).Error
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// DeleteExpired deletes all expired documents
func (r *DocumentRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ? AND expires_at IS NOT NULL", time.Now()).Delete(&model.Document{}).Error
}

// VerificationRepository handles verification-related database operations
type VerificationRepository struct {
	*BaseRepository
}

// NewVerificationRepository creates a new verification repository
func NewVerificationRepository(db *gorm.DB, logger *logrus.Logger) *VerificationRepository {
	return &VerificationRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Create creates a new verification
func (r *VerificationRepository) Create(ctx context.Context, verification *model.Verification) error {
	return r.db.WithContext(ctx).Create(verification).Error
}

// GetByID gets a verification by ID
func (r *VerificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Verification, error) {
	var verification model.Verification
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&verification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("verification not found")
		}
		return nil, err
	}
	return &verification, nil
}

// GetByUserID gets verifications by user ID
func (r *VerificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*model.Verification, int64, error) {
	var verifications []*model.Verification
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Verification{}).Where("user_id = ?", userID)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&verifications).Error
	if err != nil {
		return nil, 0, err
	}

	return verifications, total, nil
}

// GetByStatus gets verifications by status
func (r *VerificationRepository) GetByStatus(ctx context.Context, status model.VerificationStatus, page, pageSize int) ([]*model.Verification, int64, error) {
	var verifications []*model.Verification
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Verification{}).Where("status = ?", status)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&verifications).Error
	if err != nil {
		return nil, 0, err
	}

	return verifications, total, nil
}

// Update updates a verification
func (r *VerificationRepository) Update(ctx context.Context, verification *model.Verification) error {
	return r.db.WithContext(ctx).Save(verification).Error
}

// Delete deletes a verification
func (r *VerificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Verification{}, id).Error
}

// List lists all verifications with pagination
func (r *VerificationRepository) List(ctx context.Context, page, pageSize int) ([]*model.Verification, int64, error) {
	var verifications []*model.Verification
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Verification{})

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&verifications).Error
	if err != nil {
		return nil, 0, err
	}

	return verifications, total, nil
}

// AIRepository handles AI-related database operations
type AIRepository struct {
	*BaseRepository
}

// NewAIRepository creates a new AI repository
func NewAIRepository(db *gorm.DB, logger *logrus.Logger) *AIRepository {
	return &AIRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// SaveDocumentAnalysis saves a document analysis result
func (r *AIRepository) SaveDocumentAnalysis(ctx context.Context, analysis *model.DocumentAnalysisResult) error {
	// Convert maps to JSON
	if analysis.ExtractedData == nil {
		extractedData, err := json.Marshal(map[string]string{})
		if err != nil {
			return err
		}
		analysis.ExtractedData = extractedData
	}

	if analysis.Issues == nil {
		issues, err := json.Marshal([]string{})
		if err != nil {
			return err
		}
		analysis.Issues = issues
	}

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
	// Convert maps to JSON
	if result.RiskFactors == nil {
		riskFactors, err := json.Marshal([]string{})
		if err != nil {
			return err
		}
		result.RiskFactors = riskFactors
	}

	if result.DeviceInfo == nil {
		deviceInfo, err := json.Marshal(map[string]string{})
		if err != nil {
			return err
		}
		result.DeviceInfo = deviceInfo
	}

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
	// Convert maps to JSON
	if result.Reasons == nil {
		reasons, err := json.Marshal([]string{})
		if err != nil {
			return err
		}
		result.Reasons = reasons
	}

	if result.DeviceInfo == nil {
		deviceInfo, err := json.Marshal(map[string]string{})
		if err != nil {
			return err
		}
		result.DeviceInfo = deviceInfo
	}

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
