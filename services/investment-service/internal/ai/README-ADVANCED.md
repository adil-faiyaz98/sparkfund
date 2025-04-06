# Advanced AI Models for Investment Service

This document describes the advanced AI models implemented for the Investment Service, including Natural Language Processing (NLP) for market news analysis, Time Series Forecasting for price prediction, and Reinforcement Learning (RL) for portfolio optimization.

## Overview

The Investment Service now includes three advanced AI components:

1. **Natural Language Processing (NLP)**: Analyzes market news to extract sentiment and generate investment signals
2. **Time Series Forecasting**: Predicts future asset prices and market movements
3. **Reinforcement Learning (RL)**: Optimizes portfolio allocations and learns investment strategies over time

These components work together with the existing recommendation and security systems to provide a comprehensive AI-powered investment platform.

## Natural Language Processing (NLP)

The NLP system analyzes market news to extract sentiment and generate investment signals.

### Key Components

- **SentimentAnalyzer**: Analyzes text for sentiment, extracting company and sector-specific sentiment
- **NewsAnalyzer**: Analyzes news articles and generates investment signals

### Features

- Lexicon-based sentiment analysis with company and sector recognition
- Market impact prediction based on news sentiment
- Investment signal generation with confidence scores and reasoning
- Caching for efficient processing of repeated content

### API Endpoints

- `POST /api/v1/ai/advanced/news/analyze`: Analyze a news article and generate investment signals

### Example Usage

```json
// Request
POST /api/v1/ai/advanced/news/analyze
{
  "article": {
    "id": "news-123",
    "title": "Apple Reports Record Quarterly Revenue",
    "content": "Apple Inc. today announced financial results for its fiscal 2023 third quarter ended June 24, 2023. The Company posted quarterly revenue of $81.8 billion, down 1 percent year over year, and quarterly earnings per diluted share of $1.26, up 5 percent year over year.",
    "source_id": "bloomberg",
    "published_at": "2023-08-03T16:30:00Z"
  }
}

// Response
{
  "article_id": "news-123",
  "sentiment": {
    "overall_score": 0.65,
    "magnitude": 0.8,
    "company_sentiment": {
      "AAPL": 0.72
    },
    "sector_sentiment": {
      "TECHNOLOGY": 0.5
    },
    "keywords": ["apple", "revenue", "quarterly", "record", "financial"]
  },
  "market_impact": {
    "company_impacts": {
      "AAPL": {
        "price_impact": 0.014,
        "volume_impact": 0.036,
        "confidence": 0.8,
        "time_period": "SHORT_TERM"
      }
    },
    "sector_impacts": {
      "TECHNOLOGY": {
        "price_impact": 0.005,
        "volume_impact": 0.015,
        "confidence": 0.6,
        "time_period": "SHORT_TERM"
      }
    }
  },
  "investment_signals": [
    {
      "symbol": "AAPL",
      "action": "BUY",
      "strength": 0.75,
      "confidence": 0.68,
      "reasoning": "Positive sentiment detected for AAPL with expected price increase of 1.40%. Strong sentiment magnitude indicates high confidence. Impact expected in the short term."
    }
  ]
}
```

## Time Series Forecasting

The Time Series Forecasting system predicts future asset prices and market movements.

### Key Components

- **PriceForecaster**: Forecasts future prices using multiple time series models
- **MarketPredictor**: Predicts market movements using technical and fundamental indicators

### Features

- Price forecasting using ensemble of ARIMA and EMA models
- Confidence intervals for price predictions
- Market direction prediction with probability scores
- Technical indicator analysis (SMA, RSI, MACD)
- Fundamental and sentiment indicator integration

### API Endpoints

- `POST /api/v1/ai/advanced/price/forecast`: Forecast future prices for a symbol
- `POST /api/v1/ai/advanced/market/predict`: Predict market movements for a symbol

### Example Usage

```json
// Request
POST /api/v1/ai/advanced/price/forecast
{
  "symbol": "AAPL",
  "historical_prices": [
    {
      "symbol": "AAPL",
      "timestamp": "2023-07-01T00:00:00Z",
      "open": 190.5,
      "high": 192.7,
      "low": 189.8,
      "close": 191.2,
      "volume": 75000000
    },
    // ... more price points
  ],
  "days": 30
}

// Response
{
  "symbol": "AAPL",
  "current_price": 191.2,
  "forecast_points": [
    {
      "timestamp": "2023-08-01T00:00:00Z",
      "price": 193.5
    },
    // ... more forecast points
  ],
  "confidence_interval": {
    "lower": [190.1, 189.5, ...],
    "upper": [196.9, 198.2, ...],
    "level": 0.95
  },
  "model_accuracy": 0.82,
  "created_at": "2023-07-31T12:00:00Z"
}
```

## Reinforcement Learning (RL)

The Reinforcement Learning system optimizes portfolio allocations and learns investment strategies over time.

### Key Components

- **PortfolioOptimizer**: Optimizes portfolios using reinforcement learning
- **RLAgent**: Learns investment strategies through experience

### Features

- Portfolio optimization based on risk tolerance and time horizon
- Trade recommendations with reasoning
- Learning from past investment decisions
- Adaptive exploration/exploitation balance
- Experience replay for efficient learning

### API Endpoints

- `POST /api/v1/ai/advanced/portfolio/optimize`: Optimize a portfolio using reinforcement learning
- `POST /api/v1/ai/advanced/portfolio/action`: Get the best action for a given portfolio state

### Example Usage

```json
// Request
POST /api/v1/ai/advanced/portfolio/optimize
{
  "portfolio": {
    "user_id": "user-123",
    "risk_tolerance": 0.7,
    "time_horizon": 10,
    "allocations": {
      "AAPL": 0.25,
      "MSFT": 0.25,
      "GOOGL": 0.2,
      "AMZN": 0.15,
      "BND": 0.15
    }
  },
  "available_assets": [
    {
      "symbol": "AAPL",
      "asset_type": "STOCK",
      "sector": "TECHNOLOGY",
      "price": 191.2,
      "volatility": 0.25,
      "returns": [0.12, 0.08, 0.15, -0.05, 0.1]
    },
    // ... more assets
  ]
}

// Response
{
  "initial_portfolio": {
    "user_id": "user-123",
    "risk_tolerance": 0.7,
    "time_horizon": 10,
    "allocations": {
      "AAPL": 0.25,
      "MSFT": 0.25,
      "GOOGL": 0.2,
      "AMZN": 0.15,
      "BND": 0.15
    },
    "expected_return": 0.09,
    "risk": 0.18,
    "sharpe_ratio": 0.39
  },
  "optimized_portfolio": {
    "user_id": "user-123",
    "risk_tolerance": 0.7,
    "time_horizon": 10,
    "allocations": {
      "AAPL": 0.3,
      "MSFT": 0.25,
      "GOOGL": 0.2,
      "AMZN": 0.15,
      "BND": 0.1
    },
    "expected_return": 0.095,
    "risk": 0.19,
    "sharpe_ratio": 0.42
  },
  "improvement": 0.077,
  "trade_recommendations": [
    {
      "symbol": "AAPL",
      "action": "BUY",
      "current_allocation": 25.0,
      "target_allocation": 30.0,
      "change_amount": 5.0,
      "reasoning": "Increase allocation to AAPL by 5.0% (from 25.0% to 30.0%). This increases the portfolio's expected return. This improves the portfolio's risk-adjusted return (Sharpe ratio)."
    },
    {
      "symbol": "BND",
      "action": "SELL",
      "current_allocation": 15.0,
      "target_allocation": 10.0,
      "change_amount": -5.0,
      "reasoning": "Decrease allocation to BND by 5.0% (from 15.0% to 10.0%). This better aligns with your aggressive risk profile."
    }
  ]
}
```

## Automated Investments

The system can generate automated investment decisions by combining signals from multiple AI models.

### Features

- Integration of news analysis, price forecasting, and market prediction
- Confidence-weighted decision making
- Multi-source reasoning for investment decisions
- Adaptive investment amounts based on confidence

### API Endpoints

- `POST /api/v1/ai/advanced/investment/automated`: Generate an automated investment decision

### Example Usage

```json
// Request
POST /api/v1/ai/advanced/investment/automated
{
  "user_id": "user-123",
  "news_article": {
    "id": "news-123",
    "title": "Apple Reports Record Quarterly Revenue",
    "content": "Apple Inc. today announced financial results...",
    "source_id": "bloomberg",
    "published_at": "2023-08-03T16:30:00Z"
  },
  "symbol": "AAPL",
  "historical_prices": [
    // ... historical price data
  ],
  "indicators": [
    {
      "name": "PE_RATIO",
      "value": 28.5,
      "type": "FUNDAMENTAL"
    },
    // ... more indicators
  ]
}

// Response
{
  "user_id": "user-123",
  "symbol": "AAPL",
  "action": "BUY",
  "amount": 7500.0,
  "confidence": 0.72,
  "reasoning": "Based on news analysis: Positive sentiment detected for AAPL with expected price increase of 1.40%. Strong sentiment magnitude indicates high confidence. Impact expected in the short term. Price forecast confirms with +3.50% expected return over forecast period. Market prediction confirms with UP direction and +2.80% expected return.",
  "data_sources": ["NEWS", "PRICE_FORECAST", "MARKET_PREDICTION"],
  "timestamp": "2023-08-03T17:15:00Z"
}
```

## Integration with Existing AI Models

The advanced AI models integrate with the existing recommendation and security systems:

1. **Recommendation System**: The RL portfolio optimizer can use the recommendation system's asset suggestions as input for optimization.

2. **Security System**: The automated investment system checks transactions with the security system before execution.

## Implementation Details

### NLP Implementation

The NLP system uses a lexicon-based approach for sentiment analysis, with specialized dictionaries for financial terminology. Entity recognition is performed using predefined company and sector dictionaries.

### Time Series Implementation

The time series forecasting system uses an ensemble approach combining:
- ARIMA (AutoRegressive Integrated Moving Average) for trend and seasonality
- EMA (Exponential Moving Average) for smoothing and trend following

### RL Implementation

The RL system uses:
- Q-learning with function approximation
- Experience replay for efficient learning
- Exploration/exploitation balance with epsilon-greedy policy
- Modern Portfolio Theory for baseline optimization

## Getting Started

To use the advanced AI models in your code:

```go
// Initialize AI service
aiService := ai.NewAIService(recommendationService, securityService)

// Register API endpoints
advancedAIController := api.NewAdvancedAIController(aiService)
advancedAIController.RegisterRoutes(router)
```

## Future Enhancements

Planned enhancements for the advanced AI models:

1. **NLP**:
   - Deep learning-based sentiment analysis
   - Named entity recognition with transformer models
   - Event extraction and impact prediction

2. **Time Series**:
   - LSTM and transformer-based forecasting
   - Multi-variate time series analysis
   - Volatility forecasting with GARCH models

3. **RL**:
   - Deep Q-Networks (DQN) for portfolio optimization
   - Policy gradient methods for continuous action spaces
   - Multi-agent reinforcement learning for market simulation
