package ml

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
)

// Client represents a SageMaker runtime client
type Client struct {
	client *sagemakerruntime.Client
}

// Config holds SageMaker configuration
type Config struct {
	Region    string
	Endpoint  string
	ModelName string
}

// NewClient creates a new SageMaker runtime client
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := sagemakerruntime.NewFromConfig(awsCfg)
	return &Client{client: client}, nil
}

// TransactionRisk represents the risk assessment of a transaction
type TransactionRisk struct {
	UserID          string  `json:"user_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Category        string  `json:"category"`
	MerchantID      string  `json:"merchant_id"`
	MerchantName    string  `json:"merchant_name"`
	MerchantCountry string  `json:"merchant_country"`
	CardPresent     bool    `json:"card_present"`
	OnlinePayment   bool    `json:"online_payment"`
	DeviceID        string  `json:"device_id"`
	IPAddress       string  `json:"ip_address"`
	TransactionTime string  `json:"transaction_time"`
}

// RiskScore represents the fraud risk assessment result
type RiskScore struct {
	Score       float64            `json:"score"`
	Risk        string             `json:"risk"`
	Factors     []string           `json:"factors"`
	Confidence  float64            `json:"confidence"`
	Explanation map[string]float64 `json:"explanation"`
}

// AssessTransactionRisk assesses the fraud risk of a transaction
func (c *Client) AssessTransactionRisk(ctx context.Context, transaction TransactionRisk) (*RiskScore, error) {
	payload, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	input := &sagemakerruntime.InvokeEndpointInput{
		EndpointName: aws.String("fraud-detection-endpoint"),
		ContentType:  aws.String("application/json"),
		Body:         payload,
	}

	output, err := c.client.InvokeEndpoint(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke endpoint: %w", err)
	}

	var score RiskScore
	if err := json.Unmarshal(output.Body, &score); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &score, nil
}

// InvestmentRecommendation represents an investment recommendation
type InvestmentRecommendation struct {
	UserID           string    `json:"user_id"`
	RiskTolerance    string    `json:"risk_tolerance"`
	InvestmentGoal   string    `json:"investment_goal"`
	TimeHorizon      int       `json:"time_horizon"`
	CurrentPortfolio Portfolio `json:"current_portfolio"`
}

// Portfolio represents a user's investment portfolio
type Portfolio struct {
	Stocks      float64 `json:"stocks"`
	Bonds       float64 `json:"bonds"`
	Cash        float64 `json:"cash"`
	RealEstate  float64 `json:"real_estate"`
	Commodities float64 `json:"commodities"`
	Crypto      float64 `json:"crypto"`
}

// PortfolioRecommendation represents the recommended portfolio allocation
type PortfolioRecommendation struct {
	Allocation      Portfolio          `json:"allocation"`
	ExpectedReturn  float64            `json:"expected_return"`
	Risk            float64            `json:"risk"`
	Rebalancing     []Rebalancing      `json:"rebalancing"`
	Recommendations []Recommendation   `json:"recommendations"`
	Analysis        map[string]float64 `json:"analysis"`
}

// Rebalancing represents a portfolio rebalancing action
type Rebalancing struct {
	AssetClass string  `json:"asset_class"`
	Action     string  `json:"action"`
	Amount     float64 `json:"amount"`
	Reason     string  `json:"reason"`
}

// Recommendation represents an investment recommendation
type Recommendation struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Impact      float64 `json:"impact"`
	Confidence  float64 `json:"confidence"`
}

// GetInvestmentRecommendations gets investment recommendations for a user
func (c *Client) GetInvestmentRecommendations(ctx context.Context, req InvestmentRecommendation) (*PortfolioRecommendation, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	input := &sagemakerruntime.InvokeEndpointInput{
		EndpointName: aws.String("investment-advisor-endpoint"),
		ContentType:  aws.String("application/json"),
		Body:         payload,
	}

	output, err := c.client.InvokeEndpoint(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke endpoint: %w", err)
	}

	var recommendation PortfolioRecommendation
	if err := json.Unmarshal(output.Body, &recommendation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &recommendation, nil
}

// Example usage:
// client, err := ml.NewClient(ctx, ml.Config{
//     Region: "us-west-2",
// })
// if err != nil {
//     log.Fatal(err)
// }
//
// risk, err := client.AssessTransactionRisk(ctx, transaction)
// recommendation, err := client.GetInvestmentRecommendations(ctx, request)
