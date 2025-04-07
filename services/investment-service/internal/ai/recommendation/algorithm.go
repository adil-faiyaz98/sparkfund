package recommendation

import (
	"math"
	"sort"
)

// CollaborativeFilteringModel implements collaborative filtering for investment recommendations
type CollaborativeFilteringModel struct {
	// User-item matrix: maps user IDs to maps of asset IDs to ratings
	userItemMatrix map[string]map[string]float64
	
	// User similarity matrix: maps user IDs to maps of other user IDs to similarity scores
	userSimilarityMatrix map[string]map[string]float64
	
	// Item similarity matrix: maps asset IDs to maps of other asset IDs to similarity scores
	itemSimilarityMatrix map[string]map[string]float64
}

// NewCollaborativeFilteringModel creates a new collaborative filtering model
func NewCollaborativeFilteringModel() *CollaborativeFilteringModel {
	return &CollaborativeFilteringModel{
		userItemMatrix:       make(map[string]map[string]float64),
		userSimilarityMatrix: make(map[string]map[string]float64),
		itemSimilarityMatrix: make(map[string]map[string]float64),
	}
}

// AddUserInteraction adds a user interaction to the model
func (m *CollaborativeFilteringModel) AddUserInteraction(userID, assetID string, rating float64) {
	// Initialize user's ratings if not exists
	if _, exists := m.userItemMatrix[userID]; !exists {
		m.userItemMatrix[userID] = make(map[string]float64)
	}
	
	// Add or update the rating
	m.userItemMatrix[userID][assetID] = rating
}

// CalculateUserSimilarities calculates similarities between users using cosine similarity
func (m *CollaborativeFilteringModel) CalculateUserSimilarities() {
	for userA := range m.userItemMatrix {
		if _, exists := m.userSimilarityMatrix[userA]; !exists {
			m.userSimilarityMatrix[userA] = make(map[string]float64)
		}
		
		for userB := range m.userItemMatrix {
			if userA == userB {
				m.userSimilarityMatrix[userA][userB] = 1.0 // Self-similarity is 1.0
				continue
			}
			
			// Calculate cosine similarity between userA and userB
			similarity := m.calculateCosineSimilarity(m.userItemMatrix[userA], m.userItemMatrix[userB])
			m.userSimilarityMatrix[userA][userB] = similarity
		}
	}
}

// CalculateItemSimilarities calculates similarities between items using cosine similarity
func (m *CollaborativeFilteringModel) CalculateItemSimilarities() {
	// Create item-user matrix (transpose of user-item matrix)
	itemUserMatrix := make(map[string]map[string]float64)
	
	for userID, userRatings := range m.userItemMatrix {
		for assetID, rating := range userRatings {
			if _, exists := itemUserMatrix[assetID]; !exists {
				itemUserMatrix[assetID] = make(map[string]float64)
			}
			itemUserMatrix[assetID][userID] = rating
		}
	}
	
	// Calculate similarities between items
	for itemA := range itemUserMatrix {
		if _, exists := m.itemSimilarityMatrix[itemA]; !exists {
			m.itemSimilarityMatrix[itemA] = make(map[string]float64)
		}
		
		for itemB := range itemUserMatrix {
			if itemA == itemB {
				m.itemSimilarityMatrix[itemA][itemB] = 1.0 // Self-similarity is 1.0
				continue
			}
			
			// Calculate cosine similarity between itemA and itemB
			similarity := m.calculateCosineSimilarity(itemUserMatrix[itemA], itemUserMatrix[itemB])
			m.itemSimilarityMatrix[itemA][itemB] = similarity
		}
	}
}

// PredictUserItemRating predicts a user's rating for an item using user-based collaborative filtering
func (m *CollaborativeFilteringModel) PredictUserItemRating(userID, assetID string) float64 {
	// Check if the user has already rated this item
	if userRatings, exists := m.userItemMatrix[userID]; exists {
		if rating, rated := userRatings[assetID]; rated {
			return rating
		}
	}
	
	// Check if we have user similarities
	if _, exists := m.userSimilarityMatrix[userID]; !exists {
		return 0.0
	}
	
	// Find similar users who have rated this item
	type UserSimilarity struct {
		UserID     string
		Similarity float64
		Rating     float64
	}
	
	var similarUsers []UserSimilarity
	
	for otherUserID, similarity := range m.userSimilarityMatrix[userID] {
		if otherUserID == userID {
			continue
		}
		
		if otherUserRatings, exists := m.userItemMatrix[otherUserID]; exists {
			if rating, rated := otherUserRatings[assetID]; rated {
				similarUsers = append(similarUsers, UserSimilarity{
					UserID:     otherUserID,
					Similarity: similarity,
					Rating:     rating,
				})
			}
		}
	}
	
	// Sort similar users by similarity (descending)
	sort.Slice(similarUsers, func(i, j int) bool {
		return similarUsers[i].Similarity > similarUsers[j].Similarity
	})
	
	// Use top K similar users for prediction
	k := 5
	if len(similarUsers) > k {
		similarUsers = similarUsers[:k]
	}
	
	// Calculate weighted average of ratings
	var sumSimilarityRating, sumSimilarity float64
	
	for _, user := range similarUsers {
		sumSimilarityRating += user.Similarity * user.Rating
		sumSimilarity += math.Abs(user.Similarity)
	}
	
	if sumSimilarity == 0 {
		return 0.0
	}
	
	return sumSimilarityRating / sumSimilarity
}

// GetSimilarItems returns items similar to the given item
func (m *CollaborativeFilteringModel) GetSimilarItems(assetID string, count int) []string {
	if _, exists := m.itemSimilarityMatrix[assetID]; !exists {
		return nil
	}
	
	// Create a slice of item similarities
	type ItemSimilarity struct {
		ItemID     string
		Similarity float64
	}
	
	var similarItems []ItemSimilarity
	
	for otherItemID, similarity := range m.itemSimilarityMatrix[assetID] {
		if otherItemID == assetID {
			continue
		}
		
		similarItems = append(similarItems, ItemSimilarity{
			ItemID:     otherItemID,
			Similarity: similarity,
		})
	}
	
	// Sort similar items by similarity (descending)
	sort.Slice(similarItems, func(i, j int) bool {
		return similarItems[i].Similarity > similarItems[j].Similarity
	})
	
	// Return top N similar items
	if len(similarItems) > count {
		similarItems = similarItems[:count]
	}
	
	// Extract item IDs
	result := make([]string, len(similarItems))
	for i, item := range similarItems {
		result[i] = item.ItemID
	}
	
	return result
}

// GetRecommendedItems returns recommended items for a user
func (m *CollaborativeFilteringModel) GetRecommendedItems(userID string, count int) []string {
	// Get all items the user hasn't rated yet
	ratedItems := make(map[string]bool)
	if userRatings, exists := m.userItemMatrix[userID]; exists {
		for itemID := range userRatings {
			ratedItems[itemID] = true
		}
	}
	
	// Calculate predicted ratings for unrated items
	type ItemPrediction struct {
		ItemID     string
		Prediction float64
	}
	
	var predictions []ItemPrediction
	
	// Get all unique item IDs
	allItems := make(map[string]bool)
	for _, userRatings := range m.userItemMatrix {
		for itemID := range userRatings {
			allItems[itemID] = true
		}
	}
	
	// Predict ratings for unrated items
	for itemID := range allItems {
		if !ratedItems[itemID] {
			prediction := m.PredictUserItemRating(userID, itemID)
			if prediction > 0 {
				predictions = append(predictions, ItemPrediction{
					ItemID:     itemID,
					Prediction: prediction,
				})
			}
		}
	}
	
	// Sort predictions by predicted rating (descending)
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].Prediction > predictions[j].Prediction
	})
	
	// Return top N recommended items
	if len(predictions) > count {
		predictions = predictions[:count]
	}
	
	// Extract item IDs
	result := make([]string, len(predictions))
	for i, prediction := range predictions {
		result[i] = prediction.ItemID
	}
	
	return result
}

// calculateCosineSimilarity calculates the cosine similarity between two vectors
func (m *CollaborativeFilteringModel) calculateCosineSimilarity(vectorA, vectorB map[string]float64) float64 {
	// Find common keys
	var dotProduct, magnitudeA, magnitudeB float64
	
	// Calculate dot product and magnitudes
	for key, valueA := range vectorA {
		if valueB, exists := vectorB[key]; exists {
			dotProduct += valueA * valueB
		}
		magnitudeA += valueA * valueA
	}
	
	for _, valueB := range vectorB {
		magnitudeB += valueB * valueB
	}
	
	// Calculate cosine similarity
	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)
	
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0.0
	}
	
	return dotProduct / (magnitudeA * magnitudeB)
}

// ContentBasedFilteringModel implements content-based filtering for investment recommendations
type ContentBasedFilteringModel struct {
	// Asset features: maps asset IDs to feature vectors
	assetFeatures map[string]map[string]float64
	
	// User profiles: maps user IDs to feature preference vectors
	userProfiles map[string]map[string]float64
}

// NewContentBasedFilteringModel creates a new content-based filtering model
func NewContentBasedFilteringModel() *ContentBasedFilteringModel {
	return &ContentBasedFilteringModel{
		assetFeatures: make(map[string]map[string]float64),
		userProfiles:  make(map[string]map[string]float64),
	}
}

// AddAssetFeatures adds features for an asset
func (m *ContentBasedFilteringModel) AddAssetFeatures(assetID string, features map[string]float64) {
	m.assetFeatures[assetID] = features
}

// BuildUserProfile builds a user profile based on their rated items
func (m *ContentBasedFilteringModel) BuildUserProfile(userID string, userRatings map[string]float64) {
	// Initialize user profile
	m.userProfiles[userID] = make(map[string]float64)
	
	// Calculate weighted average of features from rated items
	var totalWeight float64
	
	for assetID, rating := range userRatings {
		if features, exists := m.assetFeatures[assetID]; exists {
			weight := rating // Use rating as weight
			totalWeight += weight
			
			for feature, value := range features {
				if _, exists := m.userProfiles[userID][feature]; !exists {
					m.userProfiles[userID][feature] = 0
				}
				m.userProfiles[userID][feature] += value * weight
			}
		}
	}
	
	// Normalize user profile
	if totalWeight > 0 {
		for feature := range m.userProfiles[userID] {
			m.userProfiles[userID][feature] /= totalWeight
		}
	}
}

// PredictUserItemRating predicts a user's rating for an item using content-based filtering
func (m *ContentBasedFilteringModel) PredictUserItemRating(userID, assetID string) float64 {
	// Check if we have user profile and asset features
	userProfile, userExists := m.userProfiles[userID]
	assetFeatures, assetExists := m.assetFeatures[assetID]
	
	if !userExists || !assetExists {
		return 0.0
	}
	
	// Calculate similarity between user profile and asset features
	similarity := m.calculateCosineSimilarity(userProfile, assetFeatures)
	
	// Scale similarity to rating scale (assuming 0-5 rating scale)
	return similarity * 5.0
}

// GetRecommendedItems returns recommended items for a user
func (m *ContentBasedFilteringModel) GetRecommendedItems(userID string, count int, excludeItems map[string]bool) []string {
	// Check if we have user profile
	if _, exists := m.userProfiles[userID]; !exists {
		return nil
	}
	
	// Calculate predicted ratings for all items
	type ItemPrediction struct {
		ItemID     string
		Prediction float64
	}
	
	var predictions []ItemPrediction
	
	for assetID := range m.assetFeatures {
		if excludeItems != nil && excludeItems[assetID] {
			continue
		}
		
		prediction := m.PredictUserItemRating(userID, assetID)
		if prediction > 0 {
			predictions = append(predictions, ItemPrediction{
				ItemID:     assetID,
				Prediction: prediction,
			})
		}
	}
	
	// Sort predictions by predicted rating (descending)
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].Prediction > predictions[j].Prediction
	})
	
	// Return top N recommended items
	if len(predictions) > count {
		predictions = predictions[:count]
	}
	
	// Extract item IDs
	result := make([]string, len(predictions))
	for i, prediction := range predictions {
		result[i] = prediction.ItemID
	}
	
	return result
}

// calculateCosineSimilarity calculates the cosine similarity between two feature vectors
func (m *ContentBasedFilteringModel) calculateCosineSimilarity(vectorA, vectorB map[string]float64) float64 {
	// Calculate dot product and magnitudes
	var dotProduct, magnitudeA, magnitudeB float64
	
	for feature, valueA := range vectorA {
		if valueB, exists := vectorB[feature]; exists {
			dotProduct += valueA * valueB
		}
		magnitudeA += valueA * valueA
	}
	
	for _, valueB := range vectorB {
		magnitudeB += valueB * valueB
	}
	
	// Calculate cosine similarity
	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)
	
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0.0
	}
	
	return dotProduct / (magnitudeA * magnitudeB)
}

// HybridRecommendationModel combines collaborative filtering and content-based filtering
type HybridRecommendationModel struct {
	collaborativeModel *CollaborativeFilteringModel
	contentBasedModel  *ContentBasedFilteringModel
	
	// Weight for blending predictions (0.0 to 1.0)
	// Higher values give more weight to collaborative filtering
	collaborativeWeight float64
}

// NewHybridRecommendationModel creates a new hybrid recommendation model
func NewHybridRecommendationModel(collaborativeWeight float64) *HybridRecommendationModel {
	return &HybridRecommendationModel{
		collaborativeModel:  NewCollaborativeFilteringModel(),
		contentBasedModel:   NewContentBasedFilteringModel(),
		collaborativeWeight: collaborativeWeight,
	}
}

// PredictUserItemRating predicts a user's rating for an item using hybrid filtering
func (m *HybridRecommendationModel) PredictUserItemRating(userID, assetID string) float64 {
	// Get predictions from both models
	collaborativePrediction := m.collaborativeModel.PredictUserItemRating(userID, assetID)
	contentBasedPrediction := m.contentBasedModel.PredictUserItemRating(userID, assetID)
	
	// If one model couldn't make a prediction, use the other
	if collaborativePrediction == 0 {
		return contentBasedPrediction
	}
	if contentBasedPrediction == 0 {
		return collaborativePrediction
	}
	
	// Blend predictions using weighted average
	return m.collaborativeWeight*collaborativePrediction + (1-m.collaborativeWeight)*contentBasedPrediction
}

// GetRecommendedItems returns recommended items for a user
func (m *HybridRecommendationModel) GetRecommendedItems(userID string, count int) []string {
	// Get recommendations from both models
	excludeItems := make(map[string]bool)
	if userRatings, exists := m.collaborativeModel.userItemMatrix[userID]; exists {
		for itemID := range userRatings {
			excludeItems[itemID] = true
		}
	}
	
	collaborativeRecommendations := m.collaborativeModel.GetRecommendedItems(userID, count*2)
	contentBasedRecommendations := m.contentBasedModel.GetRecommendedItems(userID, count*2, excludeItems)
	
	// Combine and deduplicate recommendations
	recommendedItems := make(map[string]float64)
	
	// Add collaborative recommendations with weight
	for i, itemID := range collaborativeRecommendations {
		// Use position as a proxy for score (higher position = higher score)
		score := float64(len(collaborativeRecommendations) - i)
		recommendedItems[itemID] = m.collaborativeWeight * score
	}
	
	// Add content-based recommendations with weight
	for i, itemID := range contentBasedRecommendations {
		// Use position as a proxy for score (higher position = higher score)
		score := float64(len(contentBasedRecommendations) - i)
		if existingScore, exists := recommendedItems[itemID]; exists {
			recommendedItems[itemID] = existingScore + (1-m.collaborativeWeight)*score
		} else {
			recommendedItems[itemID] = (1 - m.collaborativeWeight) * score
		}
	}
	
	// Convert map to slice for sorting
	type ItemScore struct {
		ItemID string
		Score  float64
	}
	
	var items []ItemScore
	for itemID, score := range recommendedItems {
		items = append(items, ItemScore{
			ItemID: itemID,
			Score:  score,
		})
	}
	
	// Sort by score (descending)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Score > items[j].Score
	})
	
	// Return top N items
	if len(items) > count {
		items = items[:count]
	}
	
	// Extract item IDs
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.ItemID
	}
	
	return result
}
