package recommendation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// RecommendationEngine implements the Engine interface
type RecommendationEngine struct {
	repository Repository
	hybridModel *HybridRecommendationModel
}

// NewRecommendationEngine creates a new recommendation engine
func NewRecommendationEngine(repository Repository) *RecommendationEngine {
	return &RecommendationEngine{
		repository:  repository,
		hybridModel: NewHybridRecommendationModel(0.6), // 60% weight to collaborative filtering
	}
}

// GetRecommendation generates investment recommendations for a user
func (e *RecommendationEngine) GetRecommendation(ctx context.Context, request RecommendationRequest) (*PortfolioRecommendation, error) {
	// Get user profile
	userProfile, err := e.repository.GetUserProfile(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Get market data
	marketData, err := e.repository.GetLatestMarketData(ctx)
	if err != nil {
		log.Printf("Warning: failed to get market data: %v", err)
		// Continue without market data
	}

	// Get user interactions
	userInteractions, err := e.repository.GetUserInteractions(ctx, request.UserID, 100)
	if err != nil {
		log.Printf("Warning: failed to get user interactions: %v", err)
		// Continue without user interactions
	}

	// Build user ratings from interactions
	userRatings := make(map[string]float64)
	for _, interaction := range userInteractions {
		// Convert interactions to ratings
		var rating float64
		switch interaction.InteractType {
		case "VIEW":
			rating = 1.0 + float64(interaction.Duration)/60.0 // Duration in minutes adds to rating
		case "SAVE":
			rating = 3.0
		case "PURCHASE":
			rating = 5.0
		case "SELL":
			rating = 2.0
		}

		// Add explicit feedback if available
		if interaction.Feedback > 0 {
			rating += float64(interaction.Feedback)
		}

		// Cap rating at 5.0
		if rating > 5.0 {
			rating = 5.0
		}

		userRatings[interaction.AssetID] = rating
	}

	// Get all investment assets
	allAssets, err := e.repository.GetInvestmentAssets(ctx, nil, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get investment assets: %w", err)
	}

	// Create asset map for quick lookup
	assetMap := make(map[string]InvestmentAsset)
	for _, asset := range allAssets {
		assetMap[asset.ID] = asset

		// Add asset to content-based model
		features := e.extractAssetFeatures(asset, marketData)
		e.hybridModel.contentBasedModel.AddAssetFeatures(asset.ID, features)

		// Add user interactions to collaborative model
		if rating, ok := userRatings[asset.ID]; ok {
			e.hybridModel.collaborativeModel.AddUserInteraction(request.UserID, asset.ID, rating)
		}
	}

	// Get similar users
	similarUsers, err := e.repository.GetSimilarUsers(ctx, request.UserID, 10)
	if err == nil && len(similarUsers) > 0 {
		// Add similar users' interactions
		for _, similarUserID := range similarUsers {
			similarUserInteractions, err := e.repository.GetUserInteractions(ctx, similarUserID, 50)
			if err != nil {
				continue
			}

			for _, interaction := range similarUserInteractions {
				var rating float64
				switch interaction.InteractType {
				case "PURCHASE":
					rating = 5.0
				case "SAVE":
					rating = 3.0
				case "VIEW":
					rating = 1.0
				}

				e.hybridModel.collaborativeModel.AddUserInteraction(similarUserID, interaction.AssetID, rating)
			}
		}
	}

	// Build user profile for content-based filtering
	e.hybridModel.contentBasedModel.BuildUserProfile(request.UserID, userProfile, userRatings)

	// Calculate similarities
	e.hybridModel.collaborativeModel.CalculateUserSimilarities()
	e.hybridModel.collaborativeModel.CalculateItemSimilarities()

	// Get recommended asset IDs
	recommendedAssetIDs := e.hybridModel.GetRecommendedItems(request.UserID, 20)

	// Filter assets based on user preferences
	var filteredAssets []InvestmentAsset
	excludedSectors := make(map[string]bool)
	for _, sector := range request.ExcludedSectors {
		excludedSectors[sector] = true
	}

	preferredSectors := make(map[string]bool)
	for _, sector := range request.PreferredSectors {
		preferredSectors[sector] = true
	}

	for _, assetID := range recommendedAssetIDs {
		asset, exists := assetMap[assetID]
		if !exists {
			continue
		}

		// Skip excluded sectors
		if excludedSectors[asset.Sector] {
			continue
		}

		// Prioritize preferred sectors
		if len(preferredSectors) > 0 && !preferredSectors[asset.Sector] {
			// Include some non-preferred sectors for diversification
			if len(filteredAssets) > 10 {
				continue
			}
		}

		filteredAssets = append(filteredAssets, asset)
	}

	// If we don't have enough assets, add popular ones
	if len(filteredAssets) < 10 {
		popularAssets, err := e.repository.GetPopularAssets(ctx, 20-len(filteredAssets))
		if err == nil {
			for _, asset := range popularAssets {
				// Skip excluded sectors
				if excludedSectors[asset.Sector] {
					continue
				}

				// Check if asset is already in filtered assets
				alreadyIncluded := false
				for _, included := range filteredAssets {
					if included.ID == asset.ID {
						alreadyIncluded = true
						break
					}
				}

				if !alreadyIncluded {
					filteredAssets = append(filteredAssets, asset)
				}
			}
		}
	}

	// Create portfolio optimizer
	optimizer := NewPortfolioOptimizer(filteredAssets)

	// Optimize portfolio
	recommendedAssets := optimizer.OptimizePortfolioWithConstraints(
		request.RiskTolerance,
		request.TimeHorizon,
		request.PreferredSectors,
		request.ExcludedSectors,
		request.Goals,
	)

	// Calculate portfolio metrics
	expectedReturn, riskLevel, diversificationScore := optimizer.CalculatePortfolioMetrics(recommendedAssets)

	// Create portfolio recommendation
	recommendation := &PortfolioRecommendation{
		ID:                   uuid.New().String(),
		UserID:               request.UserID,
		RecommendedAssets:    recommendedAssets,
		TotalExpectedReturn:  expectedReturn,
		PortfolioRiskLevel:   riskLevel,
		DiversificationScore: diversificationScore,
		RebalancingFrequency: getRebalancingFrequency(request.TimeHorizon),
		TimeHorizon:          request.TimeHorizon,
		CreatedAt:            time.Now(),
	}

	// Save recommendation
	err = e.repository.SavePortfolioRecommendation(ctx, recommendation)
	if err != nil {
		log.Printf("Warning: failed to save recommendation: %v", err)
		// Continue without saving
	}

	return recommendation, nil
}

// GetPersonalizedAssets returns personalized investment assets for a user
func (e *RecommendationEngine) GetPersonalizedAssets(ctx context.Context, userID string, count int) ([]RecommendedAsset, error) {
	// Get user profile
	userProfile, err := e.repository.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Create a recommendation request
	request := RecommendationRequest{
		UserID:           userID,
		Amount:           userProfile.InvestmentAmount,
		TimeHorizon:      userProfile.TimeHorizon,
		RiskTolerance:    userProfile.RiskTolerance,
		PreferredSectors: userProfile.PreferredSectors,
		ExcludedSectors:  userProfile.ExcludedSectors,
		Goals:            userProfile.InvestmentGoals,
	}

	// Get full recommendation
	recommendation, err := e.GetRecommendation(ctx, request)
	if err != nil {
		return nil, err
	}

	// Return top N assets
	if len(recommendation.RecommendedAssets) <= count {
		return recommendation.RecommendedAssets, nil
	}

	return recommendation.RecommendedAssets[:count], nil
}

// GetSimilarAssets returns assets similar to the given asset
func (e *RecommendationEngine) GetSimilarAssets(ctx context.Context, assetID string, count int) ([]RecommendedAsset, error) {
	// Get the asset
	asset, err := e.repository.GetInvestmentAsset(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Get similar asset IDs from collaborative model
	similarAssetIDs := e.hybridModel.collaborativeModel.GetSimilarItems(assetID, count*2)

	// Get asset details
	var similarAssets []InvestmentAsset
	for _, id := range similarAssetIDs {
		similarAsset, err := e.repository.GetInvestmentAsset(ctx, id)
		if err != nil {
			continue
		}
		similarAssets = append(similarAssets, *similarAsset)
	}

	// If we don't have enough assets, find more based on features
	if len(similarAssets) < count {
		// Get all assets
		allAssets, err := e.repository.GetInvestmentAssets(ctx, map[string]interface{}{
			"sector": asset.Sector,
		}, count*2)
		if err == nil {
			// Extract features
			assetFeatures := e.extractAssetFeatures(*asset, nil)

			// Calculate similarity scores
			type AssetSimilarity struct {
				Asset      InvestmentAsset
				Similarity float64
			}

			var assetSimilarities []AssetSimilarity
			for _, otherAsset := range allAssets {
				if otherAsset.ID == assetID {
					continue
				}

				// Check if already included
				alreadyIncluded := false
				for _, included := range similarAssets {
					if included.ID == otherAsset.ID {
						alreadyIncluded = true
						break
					}
				}
				if alreadyIncluded {
					continue
				}

				otherFeatures := e.extractAssetFeatures(otherAsset, nil)
				similarity := calculateCosineSimilarity(assetFeatures, otherFeatures)

				assetSimilarities = append(assetSimilarities, AssetSimilarity{
					Asset:      otherAsset,
					Similarity: similarity,
				})
			}

			// Sort by similarity
			sort.Slice(assetSimilarities, func(i, j int) bool {
				return assetSimilarities[i].Similarity > assetSimilarities[j].Similarity
			})

			// Add top assets
			for i := 0; i < len(assetSimilarities) && len(similarAssets) < count; i++ {
				similarAssets = append(similarAssets, assetSimilarities[i].Asset)
			}
		}
	}

	// Convert to recommended assets
	recommendedAssets := make([]RecommendedAsset, len(similarAssets))
	for i, similarAsset := range similarAssets {
		recommendedAssets[i] = RecommendedAsset{
			Asset:             similarAsset,
			Score:             0.9 - float64(i)*0.05, // Decreasing score
			AllocationPercent: 0,                     // No allocation for similar assets
			ExpectedReturn:    calculateExpectedReturn(similarAsset),
			RiskContribution:  similarAsset.Volatility,
			Reasoning:         fmt.Sprintf("Similar to %s in terms of %s and risk profile.", asset.Name, similarAsset.Sector),
		}
	}

	return recommendedAssets, nil
}

// GetDiversificationSuggestions returns suggestions to diversify a portfolio
func (e *RecommendationEngine) GetDiversificationSuggestions(ctx context.Context, userID string) ([]RecommendedAsset, error) {
	// Get user's current portfolio
	portfolio, err := e.repository.GetUserPortfolio(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolio: %w", err)
	}

	// If portfolio is empty, return personalized assets
	if len(portfolio) == 0 {
		return e.GetPersonalizedAssets(ctx, userID, 5)
	}

	// Analyze current portfolio
	sectorAllocation := make(map[string]float64)
	assetTypeAllocation := make(map[string]float64)
	totalValue := 0.0

	for _, asset := range portfolio {
		totalValue += asset.CurrentPrice
	}

	if totalValue > 0 {
		for _, asset := range portfolio {
			weight := asset.CurrentPrice / totalValue
			sectorAllocation[asset.Sector] += weight
			assetTypeAllocation[asset.AssetType] += weight
		}
	}

	// Find underrepresented sectors and asset types
	underrepresentedSectors := make([]string, 0)
	for _, sector := range []string{"TECHNOLOGY", "HEALTHCARE", "FINANCE", "CONSUMER", "ENERGY", "UTILITIES", "REAL_ESTATE"} {
		if sectorAllocation[sector] < 0.05 { // Less than 5%
			underrepresentedSectors = append(underrepresentedSectors, sector)
		}
	}

	underrepresentedTypes := make([]string, 0)
	for _, assetType := range []string{"STOCK", "BOND", "ETF", "REIT"} {
		if assetTypeAllocation[assetType] < 0.1 { // Less than 10%
			underrepresentedTypes = append(underrepresentedTypes, assetType)
		}
	}

	// Get assets in underrepresented sectors and types
	var diversificationAssets []InvestmentAsset
	
	// First try to get assets that match both criteria
	if len(underrepresentedSectors) > 0 && len(underrepresentedTypes) > 0 {
		for _, sector := range underrepresentedSectors {
			for _, assetType := range underrepresentedTypes {
				assets, err := e.repository.GetInvestmentAssets(ctx, map[string]interface{}{
					"sector":     sector,
					"asset_type": assetType,
				}, 2)
				if err == nil && len(assets) > 0 {
					diversificationAssets = append(diversificationAssets, assets...)
				}
			}
		}
	}

	// If we don't have enough, get assets from underrepresented sectors
	if len(diversificationAssets) < 5 && len(underrepresentedSectors) > 0 {
		for _, sector := range underrepresentedSectors {
			assets, err := e.repository.GetInvestmentAssets(ctx, map[string]interface{}{
				"sector": sector,
			}, 2)
			if err == nil && len(assets) > 0 {
				diversificationAssets = append(diversificationAssets, assets...)
			}
		}
	}

	// If we still don't have enough, get assets of underrepresented types
	if len(diversificationAssets) < 5 && len(underrepresentedTypes) > 0 {
		for _, assetType := range underrepresentedTypes {
			assets, err := e.repository.GetInvestmentAssets(ctx, map[string]interface{}{
				"asset_type": assetType,
			}, 2)
			if err == nil && len(assets) > 0 {
				diversificationAssets = append(diversificationAssets, assets...)
			}
		}
	}

	// If we still don't have enough, get popular assets
	if len(diversificationAssets) < 5 {
		assets, err := e.repository.GetPopularAssets(ctx, 5-len(diversificationAssets))
		if err == nil {
			diversificationAssets = append(diversificationAssets, assets...)
		}
	}

	// Convert to recommended assets
	recommendedAssets := make([]RecommendedAsset, len(diversificationAssets))
	for i, asset := range diversificationAssets {
		var reasoning string
		if contains(underrepresentedSectors, asset.Sector) {
			reasoning = fmt.Sprintf("Adds exposure to the underrepresented %s sector in your portfolio.", asset.Sector)
		} else if contains(underrepresentedTypes, asset.AssetType) {
			reasoning = fmt.Sprintf("Adds %s exposure which is underrepresented in your portfolio.", asset.AssetType)
		} else {
			reasoning = "Adds diversification to your portfolio."
		}

		recommendedAssets[i] = RecommendedAsset{
			Asset:             asset,
			Score:             0.8,
			AllocationPercent: 5.0, // Suggest 5% allocation
			ExpectedReturn:    calculateExpectedReturn(asset),
			RiskContribution:  asset.Volatility * 0.05, // 5% allocation
			Reasoning:         reasoning,
		}
	}

	return recommendedAssets, nil
}

// GetRebalancingSuggestions returns suggestions to rebalance a portfolio
func (e *RecommendationEngine) GetRebalancingSuggestions(ctx context.Context, userID string) ([]RecommendedAsset, error) {
	// Get user's current portfolio
	currentPortfolio, err := e.repository.GetUserPortfolio(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolio: %w", err)
	}

	// If portfolio is empty, return personalized assets
	if len(currentPortfolio) == 0 {
		return e.GetPersonalizedAssets(ctx, userID, 5)
	}

	// Get user profile
	userProfile, err := e.repository.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Get latest recommendation
	recommendations, err := e.repository.GetUserRecommendations(ctx, userID, 1)
	if err != nil || len(recommendations) == 0 {
		// If no recommendation exists, create a new one
		request := RecommendationRequest{
			UserID:           userID,
			Amount:           userProfile.InvestmentAmount,
			TimeHorizon:      userProfile.TimeHorizon,
			RiskTolerance:    userProfile.RiskTolerance,
			PreferredSectors: userProfile.PreferredSectors,
			ExcludedSectors:  userProfile.ExcludedSectors,
			Goals:            userProfile.InvestmentGoals,
		}

		recommendation, err := e.GetRecommendation(ctx, request)
		if err != nil {
			return nil, err
		}
		recommendations = []PortfolioRecommendation{*recommendation}
	}

	// Create portfolio optimizer
	optimizer := NewPortfolioOptimizer(currentPortfolio)

	// Get rebalancing suggestions
	suggestions := optimizer.GetRebalancingSuggestions(currentPortfolio, recommendations[0].RecommendedAssets)

	// Convert to recommended assets
	recommendedAssets := make([]RecommendedAsset, 0)
	for _, suggestion := range suggestions {
		if suggestion.Action == "BUY" {
			recommendedAssets = append(recommendedAssets, RecommendedAsset{
				Asset:             suggestion.Asset,
				Score:             0.9,
				AllocationPercent: suggestion.TargetAllocation,
				ExpectedReturn:    calculateExpectedReturn(suggestion.Asset),
				RiskContribution:  suggestion.Asset.Volatility * (suggestion.TargetAllocation / 100.0),
				Reasoning:         suggestion.Reasoning,
			})
		}
	}

	return recommendedAssets, nil
}

// TrainModel trains the recommendation model with new data
func (e *RecommendationEngine) TrainModel(ctx context.Context) error {
	// Get all user interactions
	// This is a simplified implementation - in a real system, you would use a batch process
	// to train the model with all historical data
	
	// Get all users
	userProfiles, err := e.repository.GetAllUserProfiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user profiles: %w", err)
	}

	// Get all assets
	allAssets, err := e.repository.GetInvestmentAssets(ctx, nil, 10000)
	if err != nil {
		return fmt.Errorf("failed to get investment assets: %w", err)
	}

	// Create asset map for quick lookup
	assetMap := make(map[string]InvestmentAsset)
	for _, asset := range allAssets {
		assetMap[asset.ID] = asset
	}

	// Reset the model
	e.hybridModel = NewHybridRecommendationModel(0.6)

	// Add all assets to content-based model
	for _, asset := range allAssets {
		features := e.extractAssetFeatures(asset, nil)
		e.hybridModel.contentBasedModel.AddAssetFeatures(asset.ID, features)
	}

	// Process each user
	for _, userProfile := range userProfiles {
		// Get user interactions
		userInteractions, err := e.repository.GetUserInteractions(ctx, userProfile.UserID, 1000)
		if err != nil {
			continue
		}

		// Build user ratings from interactions
		userRatings := make(map[string]float64)
		for _, interaction := range userInteractions {
			// Convert interactions to ratings
			var rating float64
			switch interaction.InteractType {
			case "VIEW":
				rating = 1.0 + float64(interaction.Duration)/60.0
			case "SAVE":
				rating = 3.0
			case "PURCHASE":
				rating = 5.0
			case "SELL":
				rating = 2.0
			}

			// Add explicit feedback if available
			if interaction.Feedback > 0 {
				rating += float64(interaction.Feedback)
			}

			// Cap rating at 5.0
			if rating > 5.0 {
				rating = 5.0
			}

			userRatings[interaction.AssetID] = rating
			e.hybridModel.collaborativeModel.AddUserInteraction(userProfile.UserID, interaction.AssetID, rating)
		}

		// Build user profile for content-based filtering
		e.hybridModel.contentBasedModel.BuildUserProfile(userProfile.UserID, userProfile, userRatings)
	}

	// Calculate similarities
	e.hybridModel.collaborativeModel.CalculateUserSimilarities()
	e.hybridModel.collaborativeModel.CalculateItemSimilarities()

	// Calculate model metrics
	metrics, err := e.evaluateModel(ctx)
	if err != nil {
		return fmt.Errorf("failed to evaluate model: %w", err)
	}

	// Save model metrics
	err = e.repository.SaveModelMetrics(ctx, metrics)
	if err != nil {
		return fmt.Errorf("failed to save model metrics: %w", err)
	}

	return nil
}

// GetModelMetrics returns the performance metrics of the recommendation model
func (e *RecommendationEngine) GetModelMetrics(ctx context.Context) (*ModelPerformanceMetrics, error) {
	return e.repository.GetLatestModelMetrics(ctx)
}

// Helper methods

// extractAssetFeatures extracts features from an investment asset
func (e *RecommendationEngine) extractAssetFeatures(asset InvestmentAsset, marketData *MarketData) map[string]float64 {
	features := make(map[string]float64)
	
	// Asset type features
	features["type_"+asset.AssetType] = 1.0
	
	// Sector features
	features["sector_"+asset.Sector] = 1.0
	
	// Risk level
	features["risk_level"] = asset.RiskLevel
	
	// Returns and volatility
	if len(asset.HistoricalReturns) > 0 {
		var sum float64
		for _, r := range asset.HistoricalReturns {
			sum += r
		}
		features["avg_return"] = sum / float64(len(asset.HistoricalReturns))
	}
	features["volatility"] = asset.Volatility
	
	// Dividend yield
	features["dividend_yield"] = asset.DividendYield
	
	// ESG score
	features["esg_score"] = asset.ESGScore
	
	// Market cap (normalized to 0-1 range)
	// Assuming market cap is in billions
	marketCapNormalized := math.Min(1.0, asset.MarketCap/1000000000000)
	features["market_cap"] = marketCapNormalized
	
	// Add market trend if available
	if marketData != nil {
		if trend, ok := marketData.MarketTrends[asset.Sector]; ok {
			features["market_trend"] = trend
		}
		
		if performance, ok := marketData.SectorPerformance[asset.Sector]; ok {
			features["sector_performance"] = performance
		}
	}
	
	return features
}

// evaluateModel evaluates the recommendation model
func (e *RecommendationEngine) evaluateModel(ctx context.Context) (*ModelPerformanceMetrics, error) {
	// This is a simplified implementation - in a real system, you would use cross-validation
	// and more sophisticated evaluation metrics
	
	// Get a sample of user interactions for testing
	userInteractions, err := e.repository.GetRecentUserInteractions(ctx, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get user interactions: %w", err)
	}
	
	// Split into training and test sets (80/20)
	trainSize := int(float64(len(userInteractions)) * 0.8)
	trainInteractions := userInteractions[:trainSize]
	testInteractions := userInteractions[trainSize:]
	
	// Build training model
	trainingModel := NewHybridRecommendationModel(0.6)
	
	// Add training interactions
	for _, interaction := range trainInteractions {
		var rating float64
		switch interaction.InteractType {
		case "PURCHASE":
			rating = 5.0
		case "SAVE":
			rating = 3.0
		case "VIEW":
			rating = 1.0
		}
		
		trainingModel.collaborativeModel.AddUserInteraction(interaction.UserID, interaction.AssetID, rating)
	}
	
	// Calculate similarities
	trainingModel.collaborativeModel.CalculateUserSimilarities()
	
	// Evaluate on test set
	var totalError, totalSquaredError float64
	var totalPredictions, correctPredictions int
	
	for _, interaction := range testInteractions {
		var actualRating float64
		switch interaction.InteractType {
		case "PURCHASE":
			actualRating = 5.0
		case "SAVE":
			actualRating = 3.0
		case "VIEW":
			actualRating = 1.0
		}
		
		// Predict rating
		predictedRating := trainingModel.collaborativeModel.PredictUserItemRating(interaction.UserID, interaction.AssetID)
		
		// Skip if no prediction
		if predictedRating == 0 {
			continue
		}
		
		// Calculate error
		error := math.Abs(predictedRating - actualRating)
		totalError += error
		totalSquaredError += error * error
		
		// Count correct predictions (within 1.0)
		if error <= 1.0 {
			correctPredictions++
		}
		
		totalPredictions++
	}
	
	// Calculate metrics
	var accuracy, mae, rmse float64
	
	if totalPredictions > 0 {
		accuracy = float64(correctPredictions) / float64(totalPredictions)
		mae = totalError / float64(totalPredictions)
		rmse = math.Sqrt(totalSquaredError / float64(totalPredictions))
	}
	
	// Create metrics
	metrics := &ModelPerformanceMetrics{
		ModelVersion:        "1.0.0",
		Accuracy:            accuracy,
		Precision:           accuracy, // Simplified
		Recall:              accuracy, // Simplified
		F1Score:             accuracy, // Simplified
		MeanAbsoluteError:   mae,
		RootMeanSquareError: rmse,
		UserSatisfaction:    0.8, // Placeholder
		LastEvaluatedAt:     time.Now(),
	}
	
	return metrics, nil
}

// calculateExpectedReturn calculates the expected return for an asset
func calculateExpectedReturn(asset InvestmentAsset) float64 {
	if len(asset.HistoricalReturns) > 0 {
		var sum float64
		for _, r := range asset.HistoricalReturns {
			sum += r
		}
		return sum / float64(len(asset.HistoricalReturns))
	}
	
	// Default expected returns by asset type
	switch asset.AssetType {
	case "STOCK":
		return 0.08 // 8%
	case "BOND":
		return 0.04 // 4%
	case "ETF":
		return 0.06 // 6%
	default:
		return 0.05 // 5%
	}
}

// getRebalancingFrequency returns the recommended rebalancing frequency
func getRebalancingFrequency(timeHorizon int) string {
	if timeHorizon < 3 {
		return "MONTHLY"
	} else if timeHorizon < 10 {
		return "QUARTERLY"
	} else {
		return "ANNUALLY"
	}
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
