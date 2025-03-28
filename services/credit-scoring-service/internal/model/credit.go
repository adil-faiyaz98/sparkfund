package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditScoreRange string

const (
	CreditScoreRangeExcellent CreditScoreRange = "excellent" // 800-850
	CreditScoreRangeVeryGood  CreditScoreRange = "very_good" // 740-799
	CreditScoreRangeGood      CreditScoreRange = "good"      // 670-739
	CreditScoreRangeFair      CreditScoreRange = "fair"      // 580-669
	CreditScoreRangePoor      CreditScoreRange = "poor"      // 300-579
)

type CreditHistoryStatus string

const (
	CreditHistoryStatusActive    CreditHistoryStatus = "active"
	CreditHistoryStatusClosed    CreditHistoryStatus = "closed"
	CreditHistoryStatusDefaulted CreditHistoryStatus = "defaulted"
)

type CreditScore struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null"`
	Score           int              `gorm:"not null"`
	ScoreRange      CreditScoreRange `gorm:"type:varchar(20);not null"`
	LastUpdated     time.Time        `gorm:"not null"`
	Factors         string           `gorm:"type:jsonb"`
	Recommendations string           `gorm:"type:jsonb"`
	CreatedAt       time.Time        `gorm:"not null"`
	UpdatedAt       time.Time        `gorm:"not null"`
	DeletedAt       gorm.DeletedAt   `gorm:"index"`
}

type CreditHistory struct {
	ID                uuid.UUID           `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID            uuid.UUID           `gorm:"type:uuid;not null"`
	AccountType       string              `gorm:"type:varchar(50);not null"`
	Institution       string              `gorm:"type:varchar(100);not null"`
	AccountNumber     string              `gorm:"type:varchar(50)"`
	Status            CreditHistoryStatus `gorm:"type:varchar(20);not null"`
	CreditLimit       float64             `gorm:"not null"`
	CurrentBalance    float64             `gorm:"not null"`
	PaymentHistory    string              `gorm:"type:jsonb"`
	OpenDate          time.Time           `gorm:"not null"`
	CloseDate         *time.Time          `gorm:"default:null"`
	LastPaymentDate   *time.Time          `gorm:"default:null"`
	LastPaymentAmount *float64            `gorm:"default:null"`
	CreatedAt         time.Time           `gorm:"not null"`
	UpdatedAt         time.Time           `gorm:"not null"`
	DeletedAt         gorm.DeletedAt      `gorm:"index"`
}

type CreditCheckRequest struct {
	UserID            uuid.UUID  `json:"userId" binding:"required"`
	AccountType       string     `json:"accountType" binding:"required"`
	Institution       string     `json:"institution" binding:"required"`
	AccountNumber     string     `json:"accountNumber"`
	CreditLimit       float64    `json:"creditLimit" binding:"required,gt=0"`
	CurrentBalance    float64    `json:"currentBalance" binding:"required,gte=0"`
	OpenDate          time.Time  `json:"openDate" binding:"required"`
	CloseDate         *time.Time `json:"closeDate,omitempty"`
	LastPaymentDate   *time.Time `json:"lastPaymentDate,omitempty"`
	LastPaymentAmount *float64   `json:"lastPaymentAmount,omitempty"`
}

type CreditScoreResponse struct {
	ID              uuid.UUID        `json:"id"`
	UserID          uuid.UUID        `json:"userId"`
	Score           int              `json:"score"`
	ScoreRange      CreditScoreRange `json:"scoreRange"`
	LastUpdated     time.Time        `json:"lastUpdated"`
	Factors         []string         `json:"factors"`
	Recommendations []string         `json:"recommendations"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

type CreditHistoryResponse struct {
	ID                uuid.UUID           `json:"id"`
	UserID            uuid.UUID           `json:"userId"`
	AccountType       string              `json:"accountType"`
	Institution       string              `json:"institution"`
	AccountNumber     string              `json:"accountNumber"`
	Status            CreditHistoryStatus `json:"status"`
	CreditLimit       float64             `json:"creditLimit"`
	CurrentBalance    float64             `json:"currentBalance"`
	PaymentHistory    []PaymentRecord     `json:"paymentHistory"`
	OpenDate          time.Time           `json:"openDate"`
	CloseDate         *time.Time          `json:"closeDate,omitempty"`
	LastPaymentDate   *time.Time          `json:"lastPaymentDate,omitempty"`
	LastPaymentAmount *float64            `json:"lastPaymentAmount,omitempty"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
}

type PaymentRecord struct {
	Date     time.Time `json:"date"`
	Amount   float64   `json:"amount"`
	Status   string    `json:"status"`
	LateDays int       `json:"lateDays,omitempty"`
}

type CreditScoreFactors struct {
	PaymentHistoryWeight    float64 `json:"paymentHistoryWeight"`
	CreditUtilizationWeight float64 `json:"creditUtilizationWeight"`
	CreditHistoryWeight     float64 `json:"creditHistoryWeight"`
	AccountMixWeight        float64 `json:"accountMixWeight"`
	NewCreditWeight         float64 `json:"newCreditWeight"`
}

type CreditScoreRecommendations struct {
	ImprovePaymentHistory   []string `json:"improvePaymentHistory"`
	ReduceCreditUtilization []string `json:"reduceCreditUtilization"`
	BuildCreditHistory      []string `json:"buildCreditHistory"`
	DiversifyAccountMix     []string `json:"diversifyAccountMix"`
	ManageNewCredit         []string `json:"manageNewCredit"`
}
