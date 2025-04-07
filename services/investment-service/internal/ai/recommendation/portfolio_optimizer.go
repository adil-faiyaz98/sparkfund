package recommendation

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// PortfolioOptimizer implements portfolio optimization using Modern Portfolio Theory
type PortfolioOptimizer struct {
	assets []InvestmentAsset
}

// NewPortfolioOptimizer creates a new portfolio optimizer
func NewPortfolioOptimizer(assets []InvestmentAsset) *PortfolioOptimizer {
	return &PortfolioOptimizer{
		assets: assets,
	}
}

// OptimizePortfolio optimizes a portfolio based on risk tolerance and time horizon
func (po *PortfolioOptimizer) OptimizePortfolio(riskTolerance float64, timeHorizon int) []RecommendedAsset {
	// Calculate expected returns and covariance matrix
	expectedReturns := make([]float64, len(po.assets))
	for i, asset := range po.assets {
		// Simple expected return calculation based on historical returns
		if len(asset.HistoricalReturns) > 0 {
			var sum float64
			for _, r := range asset.HistoricalReturns {
				sum += r
			}
			expectedReturns[i] = sum / float64(len(asset.HistoricalReturns))
		} else {
			// Fallback if no historical returns
			expectedReturns[i] = 0.05 // Default 5% return
		}
	}
	
	// Adjust allocation based on risk tolerance and time horizon
	// Higher risk tolerance and longer time horizon allow for more aggressive allocation
	stockAllocation := riskTolerance * (0.7 + 0.3*math.Min(1.0, float64(timeHorizon)/20.0))
	bondAllocation := 1.0 - stockAllocation
	
	// Create recommended assets with allocations
	recommendedAssets := make([]RecommendedAsset, 0, len(po.assets))
	
	// Sort assets by expected return / risk ratio (Sharpe ratio without risk-free rate)
	type AssetRatio struct {
		Index int
		Ratio float64
	}
	
	var assetRatios []AssetRatio
	for i, asset := range po.assets {
		// Avoid division by zero
		risk := math.Max(0.001, asset.Volatility)
		ratio := expectedReturns[i] / risk
		assetRatios = append(assetRatios, AssetRatio{
			Index: i,
			Ratio: ratio,
		})
	}
	
	// Sort by ratio (descending)
	sort.Slice(assetRatios, func(i, j int) bool {
		return assetRatios[i].Ratio > assetRatios[j].Ratio
	})
	
	// Allocate to assets based on their type and ratio
	stocksAllocated := 0.0
	bondsAllocated := 0.0
	
	for _, ar := range assetRatios {
		asset := po.assets[ar.Index]
		var allocation float64
		
		if asset.AssetType == "STOCK" || asset.AssetType == "ETF" {
			// Allocate to stocks/ETFs
			if stocksAllocated < stockAllocation {
				// Allocate more to higher ratio assets
				allocation = math.Min(0.2, stockAllocation-stocksAllocated)
				stocksAllocated += allocation
			}
		} else if asset.AssetType == "BOND" {
			// Allocate to bonds
			if bondsAllocated < bondAllocation {
				// Allocate more to higher ratio assets
				allocation = math.Min(0.2, bondAllocation-bondsAllocated)
				bondsAllocated += allocation
			}
		}
		
		if allocation > 0 {
			recommendedAssets = append(recommendedAssets, RecommendedAsset{
				Asset:             asset,
				AllocationPercent: allocation * 100.0, // Convert to percentage
				ExpectedReturn:    expectedReturns[ar.Index],
				RiskContribution:  asset.Volatility * allocation,
				ConfidenceScore:   0.8, // Default confidence score
				RecommendationTags: getRecommendationTags(asset, allocation),
				Reasoning:         generateReasoning(asset, allocation, expectedReturns[ar.Index], riskTolerance),
			})
		}
	}
	
	return recommendedAssets
}

// OptimizePortfolioWithConstraints optimizes a portfolio with additional constraints
func (po *PortfolioOptimizer) OptimizePortfolioWithConstraints(
	riskTolerance float64, 
	timeHorizon int, 
	preferredSectors []string, 
	excludedSectors []string,
	goals []string,
) []RecommendedAsset {
	// Filter assets based on excluded sectors
	var filteredAssets []InvestmentAsset
	excludedSectorsMap := make(map[string]bool)
	for _, sector := range excludedSectors {
		excludedSectorsMap[sector] = true
	}
	
	for _, asset := range po.assets {
		if !excludedSectorsMap[asset.Sector] {
			filteredAssets = append(filteredAssets, asset)
		}
	}
	
	// Create a new optimizer with filtered assets
	filteredOptimizer := NewPortfolioOptimizer(filteredAssets)
	
	// Get base recommendations
	recommendations := filteredOptimizer.OptimizePortfolio(riskTolerance, timeHorizon)
	
	// Adjust allocations based on preferred sectors
	preferredSectorsMap := make(map[string]bool)
	for _, sector := range preferredSectors {
		preferredSectorsMap[sector] = true
	}
	
	// Boost allocations for preferred sectors
	if len(preferredSectorsMap) > 0 {
		// Calculate total allocation to preferred sectors
		var preferredAllocation float64
		for i, rec := range recommendations {
			if preferredSectorsMap[rec.Asset.Sector] {
				preferredAllocation += rec.AllocationPercent / 100.0
			}
		}
		
		// Target preferred allocation based on risk tolerance
		targetPreferredAllocation := 0.4 + 0.3*riskTolerance // 40-70% to preferred sectors
		
		// Adjust allocations if needed
		if preferredAllocation < targetPreferredAllocation && preferredAllocation > 0 {
			// Boost preferred sectors
			boostFactor := targetPreferredAllocation / preferredAllocation
			
			// Apply boost and normalize
			var totalAllocation float64
			for i := range recommendations {
				if preferredSectorsMap[recommendations[i].Asset.Sector] {
					recommendations[i].AllocationPercent *= boostFactor
				}
				totalAllocation += recommendations[i].AllocationPercent / 100.0
			}
			
			// Normalize to 100%
			for i := range recommendations {
				recommendations[i].AllocationPercent = (recommendations[i].AllocationPercent / 100.0) / totalAllocation * 100.0
			}
		}
	}
	
	// Adjust based on investment goals
	for _, goal := range goals {
		switch goal {
		case "RETIREMENT":
			// For retirement, increase allocation to dividend stocks and bonds for income
			for i := range recommendations {
				if recommendations[i].Asset.DividendYield > 0.03 {
					recommendations[i].RecommendationTags = append(recommendations[i].RecommendationTags, "RETIREMENT_INCOME")
				}
			}
		case "EDUCATION":
			// For education, focus on medium-term growth with moderate risk
			// This is handled by the time horizon parameter
		case "HOUSE":
			// For house purchase, focus on lower risk investments
			// This is handled by the risk tolerance parameter
		}
	}
	
	return recommendations
}

// CalculatePortfolioMetrics calculates metrics for a portfolio
func (po *PortfolioOptimizer) CalculatePortfolioMetrics(recommendations []RecommendedAsset) (float64, float64, float64) {
	// Calculate portfolio expected return, risk, and diversification
	var totalExpectedReturn, totalRisk float64
	sectorCounts := make(map[string]int)
	assetTypeCounts := make(map[string]int)
	
	for _, rec := range recommendations {
		weight := rec.AllocationPercent / 100.0
		totalExpectedReturn += rec.ExpectedReturn * weight
		totalRisk += rec.RiskContribution
		
		sectorCounts[rec.Asset.Sector]++
		assetTypeCounts[rec.Asset.AssetType]++
	}
	
	// Calculate diversification score based on sector and asset type distribution
	sectorDiversification := float64(len(sectorCounts)) / math.Max(1.0, float64(len(recommendations)))
	assetTypeDiversification := float64(len(assetTypeCounts)) / math.Max(1.0, float64(len(recommendations)))
	
	// Combined diversification score
	diversificationScore := (sectorDiversification + assetTypeDiversification) / 2.0
	
	return totalExpectedReturn, totalRisk, diversificationScore
}

// GetRebalancingSuggestions generates suggestions for rebalancing a portfolio
func (po *PortfolioOptimizer) GetRebalancingSuggestions(currentPortfolio []InvestmentAsset, targetPortfolio []RecommendedAsset) []RebalancingSuggestion {
	// Create maps for current and target portfolios
	currentMap := make(map[string]float64)
	for _, asset := range currentPortfolio {
		currentMap[asset.ID] = 0 // Will be filled with actual allocations
	}
	
	targetMap := make(map[string]float64)
	for _, rec := range targetPortfolio {
		targetMap[rec.Asset.ID] = rec.AllocationPercent
	}
	
	// Calculate current allocations
	var totalValue float64
	for _, asset := range currentPortfolio {
		totalValue += asset.CurrentPrice
	}
	
	if totalValue > 0 {
		for i, asset := range currentPortfolio {
			currentMap[asset.ID] = (asset.CurrentPrice / totalValue) * 100.0
		}
	}
	
	// Generate rebalancing suggestions
	var suggestions []RebalancingSuggestion
	
	// Assets to sell (reduce allocation)
	for assetID, currentAllocation := range currentMap {
		targetAllocation, exists := targetMap[assetID]
		if !exists {
			targetAllocation = 0
		}
		
		if currentAllocation > targetAllocation + 1.0 { // 1% threshold
			// Find asset details
			var asset InvestmentAsset
			for _, a := range currentPortfolio {
				if a.ID == assetID {
					asset = a
					break
				}
			}
			
			suggestions = append(suggestions, RebalancingSuggestion{
				Asset:            asset,
				CurrentAllocation: currentAllocation,
				TargetAllocation:  targetAllocation,
				Action:           "SELL",
				ChangeAmount:     currentAllocation - targetAllocation,
				Reasoning:        fmt.Sprintf("Reduce exposure to %s to align with target allocation", asset.Name),
			})
		}
	}
	
	// Assets to buy (increase allocation)
	for assetID, targetAllocation := range targetMap {
		currentAllocation, exists := currentMap[assetID]
		if !exists {
			currentAllocation = 0
		}
		
		if targetAllocation > currentAllocation + 1.0 { // 1% threshold
			// Find asset details
			var asset InvestmentAsset
			for _, rec := range targetPortfolio {
				if rec.Asset.ID == assetID {
					asset = rec.Asset
					break
				}
			}
			
			suggestions = append(suggestions, RebalancingSuggestion{
				Asset:            asset,
				CurrentAllocation: currentAllocation,
				TargetAllocation:  targetAllocation,
				Action:           "BUY",
				ChangeAmount:     targetAllocation - currentAllocation,
				Reasoning:        fmt.Sprintf("Increase exposure to %s to align with target allocation", asset.Name),
			})
		}
	}
	
	// Sort suggestions by change amount (descending)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].ChangeAmount > suggestions[j].ChangeAmount
	})
	
	return suggestions
}

// RebalancingSuggestion represents a suggestion for rebalancing a portfolio
type RebalancingSuggestion struct {
	Asset             InvestmentAsset `json:"asset"`
	CurrentAllocation float64         `json:"current_allocation"` // Current allocation percentage
	TargetAllocation  float64         `json:"target_allocation"`  // Target allocation percentage
	Action            string          `json:"action"`             // "BUY" or "SELL"
	ChangeAmount      float64         `json:"change_amount"`      // Percentage points to change
	Reasoning         string          `json:"reasoning"`          // Explanation for the suggestion
}

// Helper functions

// getRecommendationTags returns recommendation tags for an asset
func getRecommendationTags(asset InvestmentAsset, allocation float64) []string {
	tags := make([]string, 0)
	
	// Add tags based on asset type
	if asset.AssetType == "STOCK" {
		tags = append(tags, "GROWTH")
	} else if asset.AssetType == "BOND" {
		tags = append(tags, "INCOME")
	} else if asset.AssetType == "ETF" {
		tags = append(tags, "DIVERSIFICATION")
	}
	
	// Add tags based on risk level
	if asset.RiskLevel < 0.3 {
		tags = append(tags, "LOW_RISK")
	} else if asset.RiskLevel > 0.7 {
		tags = append(tags, "HIGH_RISK")
	}
	
	// Add tags based on allocation
	if allocation > 15.0 {
		tags = append(tags, "CORE_HOLDING")
	} else {
		tags = append(tags, "SATELLITE")
	}
	
	// Add tags based on dividend yield
	if asset.DividendYield > 0.03 {
		tags = append(tags, "DIVIDEND")
	}
	
	// Add tags based on ESG score
	if asset.ESGScore > 0.7 {
		tags = append(tags, "ESG_FRIENDLY")
	}
	
	return tags
}

// generateReasoning generates reasoning for a recommendation
func generateReasoning(asset InvestmentAsset, allocation float64, expectedReturn float64, riskTolerance float64) string {
	var reasoning string
	
	// Base reasoning on asset type
	if asset.AssetType == "STOCK" {
		reasoning = "This stock offers growth potential "
	} else if asset.AssetType == "BOND" {
		reasoning = "This bond provides income stability "
	} else if asset.AssetType == "ETF" {
		reasoning = "This ETF offers diversified exposure "
	} else {
		reasoning = "This investment "
	}
	
	// Add sector information
	reasoning += "in the " + asset.Sector + " sector. "
	
	// Add return information
	reasoning += "Historical performance suggests a potential return of " + formatPercent(expectedReturn) + ". "
	
	// Add risk information
	if asset.RiskLevel < 0.3 {
		reasoning += "It has a low risk profile "
	} else if asset.RiskLevel > 0.7 {
		reasoning += "It has a higher risk profile "
	} else {
		reasoning += "It has a moderate risk profile "
	}
	
	reasoning += "which aligns with your " + getRiskToleranceDescription(riskTolerance) + " risk tolerance. "
	
	// Add allocation reasoning
	reasoning += "We recommend allocating " + formatPercent(allocation) + " of your portfolio to this investment."
	
	return reasoning
}

// getRiskToleranceDescription returns a description of risk tolerance
func getRiskToleranceDescription(riskTolerance float64) string {
	if riskTolerance < 0.3 {
		return "conservative"
	} else if riskTolerance < 0.7 {
		return "moderate"
	} else {
		return "aggressive"
	}
}

// formatPercent formats a number as a percentage string
func formatPercent(value float64) string {
	return fmt.Sprintf("%.1f%%", value)
}

// getCurrentTime returns the current time
func getCurrentTime() time.Time {
	return time.Now()
}

// generateUUID generates a UUID string
func generateUUID() string {
	return fmt.Sprintf("rec-%d", time.Now().UnixNano())
}
