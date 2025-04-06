# AI Models for Investment Service

This directory contains the AI models used by the Investment Service for recommendation, fraud detection, and anomaly detection.

## Overview

The Investment Service uses three main AI components:

1. **Recommendation System**: Provides personalized investment recommendations based on user profiles, market data, and collaborative filtering.
2. **Fraud Detection**: Analyzes transactions for potential fraud using various risk indicators.
3. **Anomaly Detection**: Identifies unusual investment patterns that may indicate risky behavior.

## Recommendation System

The recommendation system uses a hybrid approach combining:

- **Collaborative Filtering**: Recommends investments based on similar users' preferences
- **Content-Based Filtering**: Recommends investments based on asset features and user preferences
- **Portfolio Optimization**: Optimizes asset allocation based on Modern Portfolio Theory

### Key Features

- Personalized investment recommendations based on risk tolerance and goals
- Similar asset recommendations for diversification
- Portfolio rebalancing suggestions
- Diversification recommendations to reduce portfolio risk

### API Endpoints

- `POST /api/v1/ai/recommendations`: Get personalized portfolio recommendations
- `GET /api/v1/ai/recommendations/personalized/:userId`: Get personalized asset recommendations
- `GET /api/v1/ai/recommendations/similar/:assetId`: Get similar assets to a given asset
- `GET /api/v1/ai/recommendations/diversification/:userId`: Get diversification suggestions
- `GET /api/v1/ai/recommendations/rebalancing/:userId`: Get portfolio rebalancing suggestions

## Fraud Detection

The fraud detection system analyzes transactions for potential fraud using:

- Unusual transaction amounts
- Unusual locations
- Unusual devices or IP addresses
- Unusual transaction frequency
- Velocity checks (multiple transactions in short time)

### Key Features

- Real-time transaction analysis
- User profile building for baseline behavior
- Risk scoring and categorization
- Configurable action recommendations (approve, review, reject)

## Anomaly Detection

The anomaly detection system identifies unusual investment patterns using:

- Unusual investment amounts
- Unusual asset choices
- Unusual timing
- Market-contrary behavior
- Pattern breaks in investment behavior

### Key Features

- Behavioral analysis of investment patterns
- Market context integration
- Anomaly scoring and categorization
- Configurable alerts and monitoring

### API Endpoints

- `POST /api/v1/ai/security/analyze`: Analyze a transaction for fraud and anomalies

## Data Models

### User Profile

```go
type UserProfile struct {
    UserID           string    // Unique identifier for the user
    RiskTolerance    float64   // 0.0 (conservative) to 1.0 (aggressive)
    InvestmentGoals  []string  // e.g., "RETIREMENT", "EDUCATION", "HOUSE"
    TimeHorizon      int       // Investment horizon in years
    Age              int       // User's age
    Income           float64   // Annual income
    LiquidNetWorth   float64   // Liquid net worth
    PreferredSectors []string  // Preferred sectors
    ExcludedSectors  []string  // Sectors to exclude
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

### Investment Asset

```go
type InvestmentAsset struct {
    ID                string    // Unique identifier for the asset
    Symbol            string    // Trading symbol
    Name              string    // Full name
    AssetType         string    // e.g., "STOCK", "BOND", "ETF", "CRYPTO"
    Sector            string    // e.g., "TECHNOLOGY", "HEALTHCARE", "FINANCE"
    RiskLevel         float64   // 0.0 (low risk) to 1.0 (high risk)
    HistoricalReturns []float64 // Historical annual returns
    Volatility        float64   // Standard deviation of returns
    CurrentPrice      float64
    OneYearTarget     float64
    FiveYearTarget    float64
    DividendYield     float64
    MarketCap         float64
    ESGScore          float64   // Environmental, Social, and Governance score
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

### Transaction

```go
type Transaction struct {
    ID              string    // Unique identifier for the transaction
    UserID          string    // User who made the transaction
    Amount          float64   // Transaction amount
    Currency        string    // Currency code
    TransactionType string    // DEPOSIT, WITHDRAWAL, INVESTMENT, SALE
    AssetID         string    // Asset involved (for investments/sales)
    Timestamp       time.Time // When the transaction occurred
    IPAddress       string    // IP address used
    DeviceID        string    // Device identifier
    Location        Location  // Geographical location
    UserAgent       string    // Browser/app user agent
}
```

## Implementation Details

### Recommendation Algorithm

The recommendation system uses a hybrid approach:

1. **User Profile Analysis**: Analyzes user's risk tolerance, goals, and time horizon
2. **Collaborative Filtering**: Finds similar users and their preferred investments
3. **Content-Based Filtering**: Matches user preferences with asset features
4. **Portfolio Optimization**: Applies Modern Portfolio Theory to optimize asset allocation
5. **Diversification Analysis**: Ensures proper diversification across sectors and asset types

### Fraud Detection Algorithm

The fraud detection system uses a weighted scoring approach:

1. **Amount Analysis**: Compares transaction amount to user's historical patterns
2. **Location Analysis**: Checks if transaction location matches usual locations
3. **Device Analysis**: Verifies if device and IP are recognized
4. **Frequency Analysis**: Checks if transaction timing matches usual patterns
5. **Velocity Analysis**: Detects multiple transactions in short time periods

### Anomaly Detection Algorithm

The anomaly detection system analyzes investment behavior:

1. **Amount Analysis**: Detects unusual investment amounts
2. **Asset Analysis**: Identifies investments in unusual asset types or sectors
3. **Timing Analysis**: Detects unusual investment timing
4. **Market Context**: Identifies market-contrary behavior
5. **Pattern Analysis**: Detects breaks in established investment patterns

## Database Schema

The AI models use the following database tables:

- `user_profiles`: Stores user investment profiles
- `investment_assets`: Stores investment asset details
- `market_data`: Stores market conditions and trends
- `user_interactions`: Stores user interactions with assets
- `portfolio_recommendations`: Stores generated portfolio recommendations
- `model_performance_metrics`: Stores AI model performance metrics
- `user_portfolio`: Stores users' current portfolios
- `user_security_profiles`: Stores user security profiles for fraud detection
- `transactions`: Stores transaction details for security analysis
- `fraud_detection_results`: Stores fraud detection results
- `anomaly_detection_results`: Stores anomaly detection results

## Getting Started

To use the AI models in your code:

```go
// Initialize recommendation service
recommendationRepo := recommendation.NewPostgresRepository(db)
recommendationEngine := recommendation.NewRecommendationEngine(recommendationRepo)
recommendationService := recommendation.NewService(recommendationEngine, recommendationRepo)

// Initialize security service
securityRepo := security.NewPostgresRepository(db)
securityService := security.NewSecurityService(securityRepo)

// Register API endpoints
aiController := api.NewAIController(recommendationService, securityService)
aiController.RegisterRoutes(router)
```
