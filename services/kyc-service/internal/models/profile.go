package models

import (
	"time"

	"github.com/google/uuid"
)

// KYCStatus represents the overall KYC status of a user
type KYCStatus string

const (
	KYCStatusPending   KYCStatus = "pending"
	KYCStatusInReview  KYCStatus = "in_review"
	KYCStatusApproved  KYCStatus = "approved"
	KYCStatusRejected  KYCStatus = "rejected"
	KYCStatusSuspended KYCStatus = "suspended"
)

// RiskLevel represents the risk level assessment of a user
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

// EmploymentStatus represents the employment status of a user
type EmploymentStatus string

const (
	EmploymentStatusEmployed     EmploymentStatus = "employed"
	EmploymentStatusSelfEmployed EmploymentStatus = "self_employed"
	EmploymentStatusRetired      EmploymentStatus = "retired"
	EmploymentStatusUnemployed   EmploymentStatus = "unemployed"
)

// TransactionFrequency represents the expected transaction frequency
type TransactionFrequency string

const (
	TransactionFrequencyLow    TransactionFrequency = "low"
	TransactionFrequencyMedium TransactionFrequency = "medium"
	TransactionFrequencyHigh   TransactionFrequency = "high"
)

// InvestmentExperience represents the user's investment experience level
type InvestmentExperience string

const (
	InvestmentExperienceBeginner     InvestmentExperience = "beginner"
	InvestmentExperienceIntermediate InvestmentExperience = "intermediate"
	InvestmentExperienceAdvanced     InvestmentExperience = "advanced"
)

// Address represents a physical address
type Address struct {
	Street     string `json:"street" gorm:"type:varchar(255);not null"`
	City       string `json:"city" gorm:"type:varchar(100);not null"`
	State      string `json:"state" gorm:"type:varchar(100);not null"`
	Country    string `json:"country" gorm:"type:varchar(100);not null"`
	PostalCode string `json:"postal_code" gorm:"type:varchar(20);not null"`
}

// PersonalInfo represents personal information of a user
type PersonalInfo struct {
	FullName    string    `json:"full_name" gorm:"type:varchar(255);not null"`
	DateOfBirth time.Time `json:"date_of_birth" gorm:"type:date;not null"`
	Nationality string    `json:"nationality" gorm:"type:varchar(100);not null"`
	TaxID       string    `json:"tax_id" gorm:"type:varchar(50);not null"`
	Address     Address   `json:"address" gorm:"embedded"`
}

// EmploymentInfo represents employment information of a user
type EmploymentInfo struct {
	Occupation       string           `json:"occupation" gorm:"type:varchar(100);not null"`
	Employer         string           `json:"employer" gorm:"type:varchar(255);not null"`
	EmploymentStatus EmploymentStatus `json:"employment_status" gorm:"type:varchar(20);not null"`
	AnnualIncome     float64          `json:"annual_income" gorm:"type:decimal(15,2);not null"`
	SourceOfFunds    string           `json:"source_of_funds" gorm:"type:text;not null"`
}

// FinancialInfo represents financial information of a user
type FinancialInfo struct {
	ExpectedTransactionVolume    float64              `json:"expected_transaction_volume" gorm:"type:decimal(15,2);not null"`
	ExpectedTransactionFrequency TransactionFrequency `json:"expected_transaction_frequency" gorm:"type:varchar(20);not null"`
	InvestmentExperience         InvestmentExperience `json:"investment_experience" gorm:"type:varchar(20);not null"`
	InvestmentGoals              []string             `json:"investment_goals" gorm:"type:text[];not null"`
}

// KYCProfile represents the complete KYC profile of a user
type KYCProfile struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;unique"`
	Status         KYCStatus      `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	RiskLevel      RiskLevel      `json:"risk_level" gorm:"type:varchar(20);not null;default:'high'"`
	RiskScore      float64        `json:"risk_score" gorm:"type:float;not null;default:100"`
	PersonalInfo   PersonalInfo   `json:"personal_info" gorm:"embedded"`
	EmploymentInfo EmploymentInfo `json:"employment_info" gorm:"embedded"`
	FinancialInfo  FinancialInfo  `json:"financial_info" gorm:"embedded"`
	CreatedAt      time.Time      `json:"created_at" gorm:"type:timestamp with time zone;not null"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"type:timestamp with time zone;not null"`
	LastReviewDate time.Time      `json:"last_review_date" gorm:"type:timestamp with time zone"`
	NextReviewDate time.Time      `json:"next_review_date" gorm:"type:timestamp with time zone"`
	DeletedAt      *time.Time     `json:"-" gorm:"type:timestamp with time zone"`
}

// TableName specifies the table name for the KYCProfile model
func (KYCProfile) TableName() string {
	return "kyc_profiles"
}
