package models

import (
	"github.com/go-playground/validator/v10"
)

type Investment struct {
	InvestmentId string  `json:"investmentId" db:"investment_id"`
	ClientId     string  `json:"clientId" db:"client_id"`
	PortfolioId  string  `json:"portfolioId" db:"portfolio_id"`
	Type         string  `json:"type" db:"type"`
	Amount       float64 `json:"amount" db:"amount"`
	PurchaseDate string  `json:"purchaseDate" db:"purchase_date"`
}

type InvestmentRecommendation struct {
	InvestmentId   string  `json:"investmentId"`
	Recommendation string  `json:"recommendation"` // E.g., "Hold", "Buy", "Sell"
	Confidence     float64 `json:"confidence"`     // Confidence level (0-1)
	Rationale      string  `json:"rationale"`      // Explanation for the recommendation
}

type Portfolio struct {
	PortfolioId int    `json:"portfolioId" db:"portfolio_id"`
	ClientId    int    `json:"clientId" db:"client_id" validate:"required,gt=0"`
	Name        string `json:"name" db:"name" validate:"required"`
	Description string `json:"description" db:"description"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Will struct {
	WillId       int    `json:"willId"`
	ClientId     int64  `json:"clientId"`
	Beneficiary  string `json:"beneficiary"`
	Distribution string `json:"distribution"`
}

type WithdrawalThreshold struct {
	ThresholdAmount float64 `json:"thresholdAmount"`
}

type InvestmentForecast struct {
	InvestmentId int64  `json:"investmentId"`
	Forecast     string `json:"forecast"`
}

func ValidateStruct(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}
