-- Drop tables in reverse order of creation to respect foreign key constraints
DROP TABLE IF EXISTS anomaly_detection_results;
DROP TABLE IF EXISTS fraud_detection_results;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS user_security_profiles;
DROP TABLE IF EXISTS user_portfolio;
DROP TABLE IF EXISTS model_performance_metrics;
DROP TABLE IF EXISTS portfolio_recommendations;
DROP TABLE IF EXISTS user_interactions;
DROP TABLE IF EXISTS market_data;
DROP TABLE IF EXISTS investment_assets;
DROP TABLE IF EXISTS user_profiles;
