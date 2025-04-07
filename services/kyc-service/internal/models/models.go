package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	FirstName string    `gorm:"type:varchar(100);not null" json:"first_name" validate:"required"`
	LastName  string    `gorm:"type:varchar(100);not null" json:"last_name" validate:"required"`
	Role      string    `gorm:"type:varchar(50);not null;default:'user'" json:"role" validate:"required,oneof=user admin"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Document represents a KYC document
type Document struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type" validate:"required,oneof=PASSPORT DRIVERS_LICENSE ID_CARD RESIDENCE_PERMIT UTILITY_BILL BANK_STATEMENT"`
	FileURL     string    `gorm:"type:varchar(255);not null" json:"file_url" validate:"required,url"`
	FileName    string    `gorm:"type:varchar(255);not null" json:"file_name" validate:"required"`
	FileSize    int64     `gorm:"not null" json:"file_size" validate:"required,gt=0"`
	ContentType string    `gorm:"type:varchar(100);not null" json:"content_type" validate:"required"`
	Status      string    `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status" validate:"required,oneof=PENDING APPROVED REJECTED"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new document
func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// Selfie represents a user's selfie for face matching
type Selfie struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	FileURL     string    `gorm:"type:varchar(255);not null" json:"file_url" validate:"required,url"`
	FileName    string    `gorm:"type:varchar(255);not null" json:"file_name" validate:"required"`
	FileSize    int64     `gorm:"not null" json:"file_size" validate:"required,gt=0"`
	ContentType string    `gorm:"type:varchar(100);not null" json:"content_type" validate:"required"`
	Status      string    `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status" validate:"required,oneof=PENDING APPROVED REJECTED"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new selfie
func (s *Selfie) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// Verification represents a KYC verification process
type Verification struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	KycID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"kyc_id" validate:"required"`
	DocumentID  uuid.UUID `gorm:"type:uuid;not null;index" json:"document_id" validate:"required"`
	SelfieID    uuid.UUID `gorm:"type:uuid;index" json:"selfie_id"`
	Method      string    `gorm:"type:varchar(50);not null" json:"method" validate:"required,oneof=AI MANUAL"`
	Status      string    `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status" validate:"required,oneof=PENDING IN_PROGRESS APPROVED REJECTED"`
	Notes       string    `gorm:"type:text" json:"notes"`
	CompletedAt *time.Time `gorm:"index" json:"completed_at"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new verification
func (v *Verification) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// DocumentAnalysis represents the result of document analysis
type DocumentAnalysis struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VerificationID uuid.UUID `gorm:"type:uuid;not null;index" json:"verification_id" validate:"required"`
	DocumentID     uuid.UUID `gorm:"type:uuid;not null;index" json:"document_id" validate:"required"`
	DocumentType   string    `gorm:"type:varchar(50);not null" json:"document_type" validate:"required"`
	IsAuthentic    bool      `gorm:"not null" json:"is_authentic"`
	Confidence     float64   `gorm:"not null" json:"confidence" validate:"required,min=0,max=1"`
	ExtractedData  JSON      `gorm:"type:jsonb;not null" json:"extracted_data"`
	Issues         []string  `gorm:"type:text[]" json:"issues"`
	CreatedAt      time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new document analysis
func (da *DocumentAnalysis) BeforeCreate(tx *gorm.DB) error {
	if da.ID == uuid.Nil {
		da.ID = uuid.New()
	}
	return nil
}

// FaceMatch represents the result of face matching
type FaceMatch struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VerificationID uuid.UUID `gorm:"type:uuid;not null;index" json:"verification_id" validate:"required"`
	DocumentID     uuid.UUID `gorm:"type:uuid;not null;index" json:"document_id" validate:"required"`
	SelfieID       uuid.UUID `gorm:"type:uuid;not null;index" json:"selfie_id" validate:"required"`
	IsMatch        bool      `gorm:"not null" json:"is_match"`
	Confidence     float64   `gorm:"not null" json:"confidence" validate:"required,min=0,max=1"`
	CreatedAt      time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new face match
func (fm *FaceMatch) BeforeCreate(tx *gorm.DB) error {
	if fm.ID == uuid.Nil {
		fm.ID = uuid.New()
	}
	return nil
}

// RiskAnalysis represents the result of risk analysis
type RiskAnalysis struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VerificationID uuid.UUID `gorm:"type:uuid;not null;index" json:"verification_id" validate:"required"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	RiskScore      float64   `gorm:"not null" json:"risk_score" validate:"required,min=0,max=1"`
	RiskLevel      string    `gorm:"type:varchar(50);not null" json:"risk_level" validate:"required,oneof=LOW MEDIUM HIGH"`
	RiskFactors    []string  `gorm:"type:text[]" json:"risk_factors"`
	DeviceInfo     JSON      `gorm:"type:jsonb;not null" json:"device_info"`
	IPAddress      string    `gorm:"type:varchar(50);not null;index" json:"ip_address" validate:"required,ip"`
	Location       string    `gorm:"type:varchar(255)" json:"location"`
	CreatedAt      time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new risk analysis
func (ra *RiskAnalysis) BeforeCreate(tx *gorm.DB) error {
	if ra.ID == uuid.Nil {
		ra.ID = uuid.New()
	}
	return nil
}

// AnomalyDetection represents the result of anomaly detection
type AnomalyDetection struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VerificationID uuid.UUID `gorm:"type:uuid;not null;index" json:"verification_id" validate:"required"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	IsAnomaly      bool      `gorm:"not null" json:"is_anomaly"`
	AnomalyScore   float64   `gorm:"not null" json:"anomaly_score" validate:"required,min=0,max=1"`
	AnomalyType    string    `gorm:"type:varchar(100)" json:"anomaly_type"`
	Reasons        []string  `gorm:"type:text[]" json:"reasons"`
	DeviceInfo     JSON      `gorm:"type:jsonb;not null" json:"device_info"`
	CreatedAt      time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new anomaly detection
func (ad *AnomalyDetection) BeforeCreate(tx *gorm.DB) error {
	if ad.ID == uuid.Nil {
		ad.ID = uuid.New()
	}
	return nil
}

// JSON is a custom type for JSON data
type JSON map[string]interface{}

// Tables returns all models that should be migrated
func Tables() []interface{} {
	return []interface{}{
		&User{},
		&Document{},
		&Selfie{},
		&Verification{},
		&DocumentAnalysis{},
		&FaceMatch{},
		&RiskAnalysis{},
		&AnomalyDetection{},
	}
}
