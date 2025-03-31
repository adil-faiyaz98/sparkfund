package ai

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/sparkfund/services/kyc-service/internal/models"
)

// AMLRiskAnalysisService uses AI for AML risk analysis
type AMLRiskAnalysisService struct {
	mlClient       AMLModelClient
	graphDB        GraphDatabaseClient
	transactionAPI TransactionAPIClient
	logger         *logrus.Logger
	config         *config.AIConfig
}

// NewAMLRiskAnalysisService creates a new AML risk analysis service
func NewAMLRiskAnalysisService(
	mlClient AMLModelClient,
	graphDB GraphDatabaseClient,
	transactionAPI TransactionAPIClient,
	logger *logrus.Logger,
	config *config.AIConfig,
) *AMLRiskAnalysisService {
	return &AMLRiskAnalysisService{
		mlClient:       mlClient,
		graphDB:        graphDB,
		transactionAPI: transactionAPI,
		logger:         logger,
		config:         config,
	}
}

// AnalyzeCustomerRisk performs AI-based AML risk analysis on a customer
func (s *AMLRiskAnalysisService) AnalyzeCustomerRisk(ctx context.Context, customerData *models.CustomerData) (*models.AMLRiskAssessment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AMLRiskAnalysisService.AnalyzeCustomerRisk")
	defer span.Finish()

	logger := s.logger.WithFields(logrus.Fields{
		"customer_id":    customerData.ID,
		"correlation_id": ctx.Value("correlation_id"),
	})

	logger.Info("Starting AML risk analysis for customer")

	// Collect features for risk analysis
	features, err := s.collectCustomerFeatures(ctx, customerData)
	if err != nil {
		logger.WithError(err).Error("Failed to collect customer features")
		return nil, fmt.Errorf("feature collection failed: %w", err)
	}

	// Get network connections from graph database
	networkConnections, err := s.getNetworkConnections(ctx, customerData.ID)
	if err != nil {
		logger.WithError(err).Warn("Failed to retrieve network connections")
		// Continue with analysis even if network connections can't be retrieved
	}

	// Add network analysis features
	networkFeatures := s.analyzeNetwork(networkConnections)
	for k, v := range networkFeatures {
		features[k] = v
	}

	// Call ML model for risk prediction
	prediction, err := s.mlClient.PredictCustomerRisk(ctx, features)
	if err != nil {
		logger.WithError(err).Error("Failed to predict customer risk")
		return nil, fmt.Errorf("risk prediction failed: %w", err)
	}

	// Get feature importance
	featureImportance, err := s.mlClient.GetFeatureImportance(ctx, features)
	if err != nil {
		logger.WithError(err).Warn("Failed to get feature importance")
		// Continue without feature importance
	}

	// Determine risk factors based on feature importance
	riskFactors := s.determineRiskFactors(featureImportance, prediction)

	// Apply compliance rules
	complianceFlags := s.applyComplianceRules(ctx, customerData, prediction)

	// Generate recommendations based on risk factors
	recommendations := s.generateRecommendations(riskFactors, complianceFlags)

	// Prepare final assessment
	assessment := &models.AMLRiskAssessment{
		CustomerID:      customerData.ID,
		RiskScore:       prediction.Score,
		RiskLevel:       s.getRiskLevel(prediction.Score),
		RiskFactors:     riskFactors,
		Recommendations: recommendations,
		ComplianceFlags: complianceFlags,
		Timestamp:       time.Now(),
		ModelVersion:    prediction.ModelVersion,
		Explainability:  featureImportance,
	}

	logger.WithFields(logrus.Fields{
		"risk_score":   assessment.RiskScore,
		"risk_level":   assessment.RiskLevel,
		"factor_count": len(assessment.RiskFactors),
		"flags_count":  len(assessment.ComplianceFlags),
	}).Info("AML risk analysis completed")

	return assessment, nil
}

// AnalyzeTransactionRisk evaluates AML risk for a specific transaction
func (s *AMLRiskAnalysisService) AnalyzeTransactionRisk(ctx context.Context, transaction *models.Transaction) (*models.TransactionRiskAssessment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AMLRiskAnalysisService.AnalyzeTransactionRisk")
	defer span.Finish()

	logger := s.logger.WithFields(logrus.Fields{
		"transaction_id": transaction.ID,
		"customer_id":    transaction.CustomerID,
		"amount":         transaction.Amount,
		"correlation_id": ctx.Value("correlation_id"),
	})

	logger.Info("Starting transaction risk analysis")

	// Get customer risk profile
	customerRisk, err := s.getStoredCustomerRisk(ctx, transaction.CustomerID)
	if err != nil {
		logger.WithError(err).Warn("Failed to retrieve customer risk profile")
		// Continue with default risk profile
	}

	// Collect transaction features
	features, err := s.collectTransactionFeatures(ctx, transaction, customerRisk)
	if err != nil {
		logger.WithError(err).Error("Failed to collect transaction features")
		return nil, fmt.Errorf("feature collection failed: %w", err)
	}

	// Enrich with historical pattern data
	historyFeatures, err := s.enrichWithHistoricalPatterns(ctx, transaction)
	if err != nil {
		logger.WithError(err).Warn("Failed to enrich with historical patterns")
		// Continue without historical patterns
	} else {
		for k, v := range historyFeatures {
			features[k] = v
		}
	}

	// Call ML model for transaction risk prediction
	prediction, err := s.mlClient.PredictTransactionRisk(ctx, features)
	if err != nil {
		logger.WithError(err).Error("Failed to predict transaction risk")
		return nil, fmt.Errorf("risk prediction failed: %w", err)
	}

	// Get feature importance for explainability
	featureImportance, err := s.mlClient.GetTransactionFeatureImportance(ctx, features)
	if err != nil {
		logger.WithError(err).Warn("Failed to get feature importance")
		// Continue without feature importance
	}

	// Detect anomalies in transaction
	anomalies := s.detectAnomalies(ctx, transaction, features)

	// Apply transaction compliance rules
	complianceFlags, reportingRequired := s.applyTransactionComplianceRules(ctx, transaction, prediction, anomalies)

	// Generate risk mitigation actions
	mitigationActions := s.generateMitigationActions(prediction.Score, complianceFlags, reportingRequired)

	// Prepare final assessment
	assessment := &models.TransactionRiskAssessment{
		TransactionID:     transaction.ID,
		CustomerID:        transaction.CustomerID,
		RiskScore:         prediction.Score,
		RiskLevel:         s.getRiskLevel(prediction.Score),
		Anomalies:         anomalies,
		ComplianceFlags:   complianceFlags,
		ReportingRequired: reportingRequired,
		MitigationActions: mitigationActions,
		Timestamp:         time.Now(),
		ModelVersion:      prediction.ModelVersion,
		Explainability:    featureImportance,
	}

	logger.WithFields(logrus.Fields{
		"risk_score":         assessment.RiskScore,
		"risk_level":         assessment.RiskLevel,
		"anomalies_detected": len(assessment.Anomalies),
		"reporting_required": assessment.ReportingRequired,
	}).Info("Transaction risk analysis completed")

	return assessment, nil
}

// Helper functions
func (s *AMLRiskAnalysisService) collectCustomerFeatures(ctx context.Context, customerData *models.CustomerData) (map[string]interface{}, error) {
	features := make(map[string]interface{})

	// Basic customer attributes
	features["age"] = s.calculateAge(customerData.DateOfBirth)
	features["country_of_residence"] = customerData.Address.Country
	features["nationality"] = customerData.Nationality
	features["politically_exposed"] = customerData.PoliticallyExposed

	// Risk factors based on location
	countryRisk, err := s.getCountryRiskScore(customerData.Address.Country)
	if err != nil {
		return nil, err
	}
	features["country_risk_score"] = countryRisk

	// Customer history and activity patterns
	customerActivity, err := s.transactionAPI.GetCustomerActivitySummary(ctx, customerData.ID)
	if err == nil {
		features["account_age_days"] = customerActivity.AccountAgeDays
		features["avg_monthly_transaction_count"] = customerActivity.AvgMonthlyTransactionCount
		features["avg_transaction_amount"] = customerActivity.AvgTransactionAmount
		features["largest_transaction_amount"] = customerActivity.LargestTransactionAmount
		features["international_transaction_ratio"] = customerActivity.InternationalTransactionRatio
		features["transaction_countries_count"] = float64(len(customerActivity.TransactionCountries))
	}

	// Additional occupation and business type risk factors
	features["occupation"] = customerData.Occupation
	features["business_type"] = customerData.BusinessType
	features["income_source"] = customerData.IncomeSource

	return features, nil
}

func (s *AMLRiskAnalysisService) getNetworkConnections(ctx context.Context, customerID string) (*models.NetworkConnections, error) {
	return s.graphDB.GetEntityConnections(ctx, customerID, s.config.NetworkDepth)
}

func (s *AMLRiskAnalysisService) analyzeNetwork(connections *models.NetworkConnections) map[string]interface{} {
	if connections == nil {
		return make(map[string]interface{})
	}

	features := make(map[string]interface{})

	// Network size and characteristics
	features["network_size"] = len(connections.Entities)

	// Count high-risk connections
	highRiskCount := 0
	for _, entity := range connections.Entities {
		if entity.RiskScore > s.config.HighRiskThreshold {
			highRiskCount++
		}
	}
	features["high_risk_connections"] = highRiskCount

	// Network density (if available)
	if connections.Stats != nil {
		features["network_density"] = connections.Stats.Density
		features["avg_path_length"] = connections.Stats.AvgPathLength
	}

	return features
}

// Determine risk level from score
func (s *AMLRiskAnalysisService) getRiskLevel(score float64) string {
	switch {
	case score >= s.config.HighRiskThreshold:
		return "HIGH"
	case score >= s.config.MediumRiskThreshold:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

// Calculate age from date of birth
func (s *AMLRiskAnalysisService) calculateAge(dateOfBirth time.Time) int {
	now := time.Now()
	age := now.Year() - dateOfBirth.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if now.YearDay() < dateOfBirth.YearDay() {
		age--
	}

	return age
}

// Get country risk score from configuration or API
func (s *AMLRiskAnalysisService) getCountryRiskScore(countryCode string) (float64, error) {
	// Implementation would fetch from config or external API
	// For now, return a default based on high-risk countries list
	highRiskCountries := map[string]bool{
		"AF": true, "BY": true, "MM": true, "KP": true, "RU": true,
		"SY": true, "VE": true, "YE": true, "ZW": true,
	}

	if highRiskCountries[strings.ToUpper(countryCode)] {
		return 0.9, nil // High risk score
	}

	return 0.3, nil // Default moderate-low risk
}

// Helper functions

func (s *AMLRiskAnalysisService) gatherCustomerFeatures(ctx context.Context, customer models.Customer) (map[string]interface{}, error) {
	// Gather various features about the customer for ML analysis
	features := make(map[string]interface{})

	// Basic customer data
	features["customer_age"] = calculateAge(customer.DateOfBirth)
	features["customer_country"] = customer.Country
	features["customer_occupation"] = customer.Occupation
	features["account_age_days"] = time.Since(customer.CreatedAt).Hours() / 24

	// PEP status
	pepCheck, err := s.performPEPCheck(ctx, customer)
	if err != nil {
		return nil, err
	}
	features["is_pep"] = pepCheck.IsPEP
	features["pep_score"] = pepCheck.Score

	// Sanctions check
	sanctionsCheck, err := s.performSanctionsCheck(ctx, customer)
	if err != nil {
		return nil, err
	}
	features["sanctions_match"] = sanctionsCheck.IsMatch
	features["sanctions_score"] = sanctionsCheck.Score

	// Add geographic risk
	features["country_risk_score"] = getCountryRiskScore(customer.Country)

	return features, nil
}

func (s *AMLRiskAnalysisService) extractTransactionFeatures(ctx context.Context, transactions []models.Transaction) (map[string]interface{}, error) {
	// Extract features from transaction history
	features := make(map[string]interface{})

	if len(transactions) == 0 {
		return features, nil
	}

	// Calculate transaction statistics
	var amounts []float64
	var countries = make(map[string]int)
	var categories = make(map[string]int)
	var counterparties = make(map[string]int)

	for _, tx := range transactions {
		amounts = append(amounts, tx.Amount)
		countries[tx.CountryCode]++
		categories[tx.Category]++
		counterparties[tx.CounterpartyID]++
	}

	// Add statistical features
	features["tx_count"] = len(transactions)
	features["tx_avg_amount"] = calculateAverage(amounts)
	features["tx_max_amount"] = calculateMax(amounts)
	features["tx_volume_30d"] = calculateVolumeInDays(transactions, 30)
	features["tx_unique_countries"] = len(countries)
	features["tx_unique_categories"] = len(categories)
	features["tx_unique_counterparties"] = len(counterparties)
	features["tx_high_risk_country_ratio"] = calculateHighRiskCountryRatio(countries)

	return features, nil
}

// Completing the cut-off function
func (s *AMLRiskAnalysisService) performNetworkAnalysis(ctx context.Context, customerID string) (map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ai.performNetworkAnalysis")
	defer span.Finish()

	features := make(map[string]interface{})

	// Query the graph database for network patterns
	networkQuery := &NetworkQuery{
		CustomerID:          customerID,
		Depth:               3,                   // Up to 3 degrees of connection
		IncludeEntities:     true,                // Include companies, trusts, etc.
		IncludeTransactions: true,                // Include transaction connections
		TimeWindow:          90 * 24 * time.Hour, // Last 90 days
	}

	networkResult, err := s.graphDB.QueryNetwork(ctx, networkQuery)
	if err != nil {
		return features, fmt.Errorf("graph database query failed: %w", err)
	}

	// Extract network metrics
	features["network_size"] = networkResult.NetworkSize
	features["direct_connections"] = networkResult.DirectConnections
	features["high_risk_connections"] = networkResult.HighRiskConnections
	features["pep_connections"] = networkResult.PEPConnections
	features["sanctioned_entity_connections"] = networkResult.SanctionedEntityConnections
	features["circular_transactions"] = networkResult.CircularTransactionPatterns
	features["shared_addresses"] = networkResult.SharedAddressCount
	features["shared_contact_info"] = networkResult.SharedContactInfoCount
	features["shell_company_connections"] = networkResult.ShellCompanyConnections
	features["network_density"] = networkResult.NetworkDensity

	// Calculate network risk score
	networkRiskScore := s.calculateNetworkRiskScore(networkResult)
	features["network_risk_score"] = networkRiskScore

	// If risk score is high, extract more detailed features
	if networkRiskScore > 70 {
		detailedFeatures, err := s.extractHighRiskNetworkFeatures(ctx, networkResult)
		if err != nil {
			s.logger.WithError(err).Warn("Failed to extract detailed network features")
		} else {
			for k, v := range detailedFeatures {
				features[k] = v
			}
		}
	}

	return features, nil
}

func (s *AMLRiskAnalysisService) extractSingleTransactionFeatures(ctx context.Context, tx models.Transaction) (map[string]interface{}, error) {
	features := make(map[string]interface{})

	// Basic transaction data
	features["tx_amount"] = tx.Amount
	features["tx_currency"] = tx.Currency
	features["tx_type"] = tx.Type
	features["tx_category"] = tx.Category
	features["tx_country_code"] = tx.CountryCode
	features["tx_has_description"] = tx.Description != ""
	features["tx_description_length"] = len(tx.Description)
	features["tx_hour_of_day"] = tx.Timestamp.Hour()
	features["tx_day_of_week"] = tx.Timestamp.Weekday()
	features["tx_weekend"] = tx.Timestamp.Weekday() == time.Saturday || tx.Timestamp.Weekday() == time.Sunday

	// Country risk
	features["tx_country_risk"] = getCountryRiskScore(tx.CountryCode)

	// Transaction amount features
	features["tx_amount_usd"] = convertToUSD(tx.Amount, tx.Currency)

	// Round number detection (e.g., exactly 10000.00)
	amountStr := fmt.Sprintf("%.2f", tx.Amount)
	features["tx_round_amount"] = strings.HasSuffix(amountStr, "000.00") ||
		strings.HasSuffix(amountStr, "500.00") ||
		strings.HasSuffix(amountStr, "999.99")

	// Just below threshold detection (e.g., 9999.99 - potentially structuring)
	for _, threshold := range s.config.StructuringThresholds {
		thresholdUSD := convertToUSD(threshold, "USD")
		amountUSD := features["tx_amount_usd"].(float64)

		if amountUSD >= thresholdUSD*0.95 && amountUSD < thresholdUSD {
			features["tx_just_below_threshold"] = true
			features["tx_threshold_proximity"] = (thresholdUSD - amountUSD) / thresholdUSD
			break
		}
	}

	// Counterparty features
	if tx.CounterpartyID != "" {
		counterpartyRisk, err := s.getCounterpartyRisk(ctx, tx.CounterpartyID)
		if err != nil {
			s.logger.WithError(err).Warn("Failed to get counterparty risk")
		} else {
			features["counterparty_risk_score"] = counterpartyRisk
			features["high_risk_counterparty"] = counterpartyRisk > s.config.HighRiskCounterpartyThreshold
		}
	}

	return features, nil
}

func (s *AMLRiskAnalysisService) addHistoricalContext(ctx context.Context, tx models.Transaction) (map[string]interface{}, error) {
	features := make(map[string]interface{})

	// Get historical transactions (last 90 days)
	endDate := tx.Timestamp
	startDate := endDate.AddDate(0, 0, -90)

	historicalTxs, err := s.transactionAPI.GetCustomerTransactionsInRange(ctx, tx.CustomerID, startDate, endDate)
	if err != nil {
		return features, fmt.Errorf("failed to retrieve historical transactions: %w", err)
	}

	// Skip the current transaction
	var filteredTxs []models.Transaction
	for _, htx := range historicalTxs {
		if htx.ID != tx.ID {
			filteredTxs = append(filteredTxs, htx)
		}
	}

	if len(filteredTxs) == 0 {
		// No history available
		features["no_history_available"] = true
		return features, nil
	}

	// Calculate historical statistics
	amounts := make([]float64, 0, len(filteredTxs))
	sameCounterpartyTxs := 0
	sameCountryTxs := 0
	sameTypeTxs := 0

	for _, htx := range filteredTxs {
		amounts = append(amounts, htx.Amount)

		if htx.CounterpartyID == tx.CounterpartyID && tx.CounterpartyID != "" {
			sameCounterpartyTxs++
		}

		if htx.CountryCode == tx.CountryCode {
			sameCountryTxs++
		}

		if htx.Type == tx.Type {
			sameTypeTxs++
		}
	}

	// Calculate z-score for the current transaction amount
	mean := calculateAverage(amounts)
	stdDev := calculateStdDev(amounts, mean)

	var zScore float64
	if stdDev > 0 {
		zScore = (tx.Amount - mean) / stdDev
	} else {
		zScore = 0
	}

	features["historical_tx_count"] = len(filteredTxs)
	features["historical_tx_avg_amount"] = mean
	features["historical_tx_max_amount"] = calculateMax(amounts)
	features["tx_amount_zscore"] = zScore
	features["tx_amount_unusual"] = math.Abs(zScore) > 2.5 // More than 2.5 std devs from mean

	totalTxs := float64(len(filteredTxs))
	features["historical_same_counterparty_ratio"] = float64(sameCounterpartyTxs) / totalTxs
	features["historical_same_country_ratio"] = float64(sameCountryTxs) / totalTxs
	features["historical_same_type_ratio"] = float64(sameTypeTxs) / totalTxs

	// Check for velocity anomalies (sudden increase in transaction frequency)
	velocityFeatures := s.detectVelocityAnomalies(ctx, filteredTxs, tx)
	for k, v := range velocityFeatures {
		features[k] = v
	}

	return features, nil
}

func (s *AMLRiskAnalysisService) detectVelocityAnomalies(ctx context.Context, historicalTxs []models.Transaction, currentTx models.Transaction) map[string]interface{} {
	features := make(map[string]interface{})

	// Sort transactions by timestamp
	sort.Slice(historicalTxs, func(i, j int) bool {
		return historicalTxs[i].Timestamp.Before(historicalTxs[j].Timestamp)
	})

	// Calculate transactions per day for the past periods
	periods := map[string]time.Duration{
		"7d":  7 * 24 * time.Hour,
		"30d": 30 * 24 * time.Hour,
		"90d": 90 * 24 * time.Hour,
	}

	for periodName, duration := range periods {
		cutoffTime := currentTx.Timestamp.Add(-duration)
		var periodTxs []models.Transaction

		for _, tx := range historicalTxs {
			if tx.Timestamp.After(cutoffTime) {
				periodTxs = append(periodTxs, tx)
			}
		}

		txCount := len(periodTxs)
		daysInPeriod := duration.Hours() / 24
		txPerDay := float64(txCount) / daysInPeriod

		features[fmt.Sprintf("tx_per_day_%s", periodName)] = txPerDay

		// Calculate total volume for the period
		var totalVolume float64
		for _, tx := range periodTxs {
			totalVolume += convertToUSD(tx.Amount, tx.Currency)
		}

		features[fmt.Sprintf("volume_per_day_%s", periodName)] = totalVolume / daysInPeriod
	}

	// Detect velocity anomaly by comparing 7d rate to 30d rate
	if features["tx_per_day_7d"].(float64) > 0 && features["tx_per_day_30d"].(float64) > 0 {
		velocityRatio := features["tx_per_day_7d"].(float64) / features["tx_per_day_30d"].(float64)
		features["velocity_ratio_7d_30d"] = velocityRatio
		features["velocity_anomaly"] = velocityRatio > s.config.VelocityAnomalyThreshold
	}

	return features
}

func (s *AMLRiskAnalysisService) combineFeatures(customerFeatures, transactionFeatures, networkFeatures map[string]interface{}) map[string]interface{} {
	combined := make(map[string]interface{})

	// Add customer features
	for k, v := range customerFeatures {
		combined[k] = v
	}

	// Add transaction features
	for k, v := range transactionFeatures {
		combined[k] = v
	}

	// Add network features
	for k, v := range networkFeatures {
		combined[k] = v
	}

	return combined
}

func (s *AMLRiskAnalysisService) generateRiskFactors(prediction *AMLRiskPrediction) []string {
	var factors []string

	for factor, score := range prediction.FactorScores {
		if score >= s.config.RiskFactorThreshold {
			factors = append(factors, factor)
		}
	}

	return factors
}

func (s *AMLRiskAnalysisService) generateTransactionRiskFactors(prediction *TransactionRiskPrediction) []string {
	var factors []string

	for factor, score := range prediction.FactorScores {
		if score >= s.config.RiskFactorThreshold {
			factors = append(factors, factor)
		}
	}

	// Add triggered rules
	for _, rule := range prediction.TriggerRules {
		factors = append(factors, fmt.Sprintf("Rule triggered: %s", rule))
	}

	return factors
}

func (s *AMLRiskAnalysisService) generateRecommendations(prediction *AMLRiskPrediction) []string {
	var recommendations []string

	// Generate recommendations based on risk level
	switch determineRiskLevel(prediction.RiskScore) {
	case "extreme", "high":
		recommendations = append(recommendations, "Conduct enhanced due diligence")
		recommendations = append(recommendations, "Consider filing Suspicious Activity Report")
		recommendations = append(recommendations, "Review and restrict transaction limits")
	case "medium":
		recommendations = append(recommendations, "Monitor account activity closely")
		recommendations = append(recommendations, "Verify source of funds for large transactions")
	case "low":
		recommendations = append(recommendations, "Apply standard monitoring procedures")
	}

	// Add specific recommendations based on triggered factors
	for factor, score := range prediction.FactorScores {
		if score >= s.config.RiskFactorThreshold {
			switch factor {
			case "high_value_transactions":
				recommendations = append(recommendations, "Request source of funds documentation")
			case "unusual_cross_border_activity":
				recommendations = append(recommendations, "Review international transaction patterns")
			case "pep_connection":
				recommendations = append(recommendations, "Apply PEP-specific enhanced due diligence")
			case "complex_ownership_structure":
				recommendations = append(recommendations, "Conduct beneficial ownership investigation")
			case "suspicious_network_patterns":
				recommendations = append(recommendations, "Map and analyze customer's transaction network")
			}
		}
	}

	return recommendations
}

func determineRiskLevel(score float64) string {
	switch {
	case score >= 90:
		return "extreme"
	case score >= 70:
		return "high"
	case score >= 40:
		return "medium"
	default:
		return "low"
	}
}

func (s *AMLRiskAnalysisService) triggerSARConsideration(customerID string, assessment *models.AMLRiskAssessment) {
	// This would typically integrate with a workflow or case management system
	s.logger.WithFields(logrus.Fields{
		"customer_id":   customerID,
		"risk_score":    assessment.OverallRiskScore,
		"assessment_id": assessment.AssessmentID,
	}).Info("Triggered SAR consideration workflow")

	// Create SAR consideration case (implementation would depend on case management system)
	// ...
}

func (s *AMLRiskAnalysisService) triggerHighRiskTransactionAlert(tx models.Transaction, assessment *models.TransactionRiskAssessment) {
	// This would integrate with an alerting system
	s.logger.WithFields(logrus.Fields{
		"transaction_id": tx.ID,
		"customer_id":    tx.CustomerID,
		"amount":         tx.Amount,
		"currency":       tx.Currency,
		"risk_score":     assessment.RiskScore,
		"risk_level":     assessment.RiskLevel,
	}).Warn("High risk transaction detected")

	// Create alert (implementation would depend on alerting system)
	// ...
}

// Utility functions
func calculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var sum float64
	for _, v := range values {
		sum += v
	}

	return sum / float64(len(values))
}

func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) < 2 {
		return 0
	}

	var sumSquaredDiff float64
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}

	variance := sumSquaredDiff / float64(len(values)-1)
	return math.Sqrt(variance)
}

func calculateMax(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}

	return max
}
