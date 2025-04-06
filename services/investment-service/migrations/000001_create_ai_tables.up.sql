-- User profiles for recommendation system
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id VARCHAR(36) PRIMARY KEY,
    risk_tolerance FLOAT NOT NULL,
    investment_goals JSONB NOT NULL DEFAULT '[]',
    time_horizon INTEGER NOT NULL,
    age INTEGER,
    income FLOAT,
    liquid_net_worth FLOAT,
    investment_amount FLOAT,
    preferred_sectors JSONB NOT NULL DEFAULT '[]',
    excluded_sectors JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Investment assets
CREATE TABLE IF NOT EXISTS investment_assets (
    id VARCHAR(36) PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    asset_type VARCHAR(20) NOT NULL,
    sector VARCHAR(50) NOT NULL,
    risk_level FLOAT NOT NULL,
    historical_returns JSONB NOT NULL DEFAULT '[]',
    volatility FLOAT NOT NULL,
    current_price FLOAT NOT NULL,
    one_year_target FLOAT,
    five_year_target FLOAT,
    dividend_yield FLOAT,
    market_cap FLOAT,
    esg_score FLOAT,
    features JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Market data
CREATE TABLE IF NOT EXISTS market_data (
    date TIMESTAMP WITH TIME ZONE PRIMARY KEY,
    market_trends JSONB NOT NULL DEFAULT '{}',
    economic_indicators JSONB NOT NULL DEFAULT '{}',
    sector_performance JSONB NOT NULL DEFAULT '{}'
);

-- User interactions
CREATE TABLE IF NOT EXISTS user_interactions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES user_profiles(user_id),
    asset_id VARCHAR(36) NOT NULL REFERENCES investment_assets(id),
    interact_type VARCHAR(20) NOT NULL,
    rating FLOAT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    amount FLOAT,
    duration INTEGER,
    feedback INTEGER
);

-- Portfolio recommendations
CREATE TABLE IF NOT EXISTS portfolio_recommendations (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES user_profiles(user_id),
    recommended_assets JSONB NOT NULL,
    total_expected_return FLOAT NOT NULL,
    portfolio_risk_level FLOAT NOT NULL,
    diversification_score FLOAT NOT NULL,
    rebalancing_frequency VARCHAR(20) NOT NULL,
    time_horizon INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Model performance metrics
CREATE TABLE IF NOT EXISTS model_performance_metrics (
    id SERIAL PRIMARY KEY,
    model_version VARCHAR(20) NOT NULL,
    accuracy FLOAT NOT NULL,
    precision FLOAT NOT NULL,
    recall FLOAT NOT NULL,
    f1_score FLOAT NOT NULL,
    mean_absolute_error FLOAT NOT NULL,
    root_mean_square_error FLOAT NOT NULL,
    user_satisfaction FLOAT NOT NULL,
    last_evaluated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- User portfolios
CREATE TABLE IF NOT EXISTS user_portfolio (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES user_profiles(user_id),
    asset_id VARCHAR(36) NOT NULL REFERENCES investment_assets(id),
    quantity FLOAT NOT NULL,
    purchase_price FLOAT NOT NULL,
    purchase_date TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(user_id, asset_id)
);

-- Security profiles for fraud detection
CREATE TABLE IF NOT EXISTS user_security_profiles (
    user_id VARCHAR(36) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_login TIMESTAMP WITH TIME ZONE NOT NULL,
    usual_ip_addresses JSONB NOT NULL DEFAULT '[]',
    usual_device_ids JSONB NOT NULL DEFAULT '[]',
    usual_locations JSONB NOT NULL DEFAULT '[]',
    average_transaction_amount FLOAT NOT NULL DEFAULT 0,
    transaction_frequency FLOAT NOT NULL DEFAULT 0,
    risk_score FLOAT NOT NULL DEFAULT 0.5
);

-- Transactions for fraud and anomaly detection
CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    amount FLOAT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    asset_id VARCHAR(36),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    ip_address VARCHAR(45),
    device_id VARCHAR(100),
    location JSONB,
    user_agent TEXT
);

-- Fraud detection results
CREATE TABLE IF NOT EXISTS fraud_detection_results (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL REFERENCES transactions(id),
    user_id VARCHAR(36) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    fraud_score FLOAT NOT NULL,
    fraud_level VARCHAR(20) NOT NULL,
    fraud_indicators JSONB NOT NULL,
    action VARCHAR(20) NOT NULL
);

-- Anomaly detection results
CREATE TABLE IF NOT EXISTS anomaly_detection_results (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL REFERENCES transactions(id),
    user_id VARCHAR(36) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    anomaly_score FLOAT NOT NULL,
    anomaly_level VARCHAR(20) NOT NULL,
    anomaly_indicators JSONB NOT NULL,
    action VARCHAR(20) NOT NULL
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_interactions_user_id ON user_interactions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_interactions_asset_id ON user_interactions(asset_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_recommendations_user_id ON portfolio_recommendations(user_id);
CREATE INDEX IF NOT EXISTS idx_user_portfolio_user_id ON user_portfolio(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_fraud_detection_results_user_id ON fraud_detection_results(user_id);
CREATE INDEX IF NOT EXISTS idx_anomaly_detection_results_user_id ON anomaly_detection_results(user_id);
