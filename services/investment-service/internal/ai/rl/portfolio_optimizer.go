package rl

import (
	"context"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Asset represents an investment asset
type Asset struct {
	Symbol     string    `json:"symbol"`
	AssetType  string    `json:"asset_type"`
	Sector     string    `json:"sector"`
	Price      float64   `json:"price"`
	Volatility float64   `json:"volatility"`
	Returns    []float64 `json:"returns"`
}

// Portfolio represents an investment portfolio
type Portfolio struct {
	UserID         string                `json:"user_id"`
	RiskTolerance  float64               `json:"risk_tolerance"`
	TimeHorizon    int                   `json:"time_horizon"`
	Allocations    map[string]float64    `json:"allocations"`    // Symbol -> allocation percentage
	ExpectedReturn float64               `json:"expected_return"`
	Risk           float64               `json:"risk"`
	SharpeRatio    float64               `json:"sharpe_ratio"`
	Timestamp      time.Time             `json:"timestamp"`
}

// OptimizationResult represents the result of portfolio optimization
type OptimizationResult struct {
	InitialPortfolio  Portfolio            `json:"initial_portfolio"`
	OptimizedPortfolio Portfolio            `json:"optimized_portfolio"`
	Improvement       float64              `json:"improvement"`        // Percentage improvement in Sharpe ratio
	TradeRecommendations []TradeRecommendation `json:"trade_recommendations"`
	Timestamp         time.Time            `json:"timestamp"`
}

// TradeRecommendation represents a recommended trade
type TradeRecommendation struct {
	Symbol        string  `json:"symbol"`
	Action        string  `json:"action"`        // "BUY", "SELL", "HOLD"
	CurrentAllocation float64 `json:"current_allocation"`
	TargetAllocation  float64 `json:"target_allocation"`
	ChangeAmount   float64 `json:"change_amount"`
	Reasoning      string  `json:"reasoning"`
}

// PortfolioOptimizer optimizes portfolios using reinforcement learning
type PortfolioOptimizer struct {
	// RL model parameters
	learningRate     float64
	discountFactor   float64
	explorationRate  float64
	
	// Q-table: state -> action -> value
	qTable           map[string]map[string]float64
	qTableMutex      sync.RWMutex
	
	// Cache of optimization results
	resultCache      map[string]OptimizationResult
	cacheMutex       sync.RWMutex
	cacheTime        time.Duration
	
	// Risk-free rate for Sharpe ratio calculation
	riskFreeRate     float64
}

// NewPortfolioOptimizer creates a new portfolio optimizer
func NewPortfolioOptimizer() *PortfolioOptimizer {
	return &PortfolioOptimizer{
		learningRate:    0.1,
		discountFactor:  0.9,
		explorationRate: 0.1,
		qTable:          make(map[string]map[string]float64),
		resultCache:     make(map[string]OptimizationResult),
		cacheTime:       24 * time.Hour,
		riskFreeRate:    0.02, // 2% risk-free rate
	}
}

// OptimizePortfolio optimizes a portfolio using reinforcement learning
func (o *PortfolioOptimizer) OptimizePortfolio(ctx context.Context, portfolio Portfolio, availableAssets []Asset) (*OptimizationResult, error) {
	// Check cache first
	cacheKey := portfolio.UserID
	o.cacheMutex.RLock()
	if result, ok := o.resultCache[cacheKey]; ok {
		if time.Since(result.Timestamp) < o.cacheTime {
			o.cacheMutex.RUnlock()
			return &result, nil
		}
	}
	o.cacheMutex.RUnlock()
	
	// Calculate initial portfolio metrics
	initialPortfolio := portfolio
	initialPortfolio.ExpectedReturn = calculateExpectedReturn(portfolio, availableAssets)
	initialPortfolio.Risk = calculatePortfolioRisk(portfolio, availableAssets)
	initialPortfolio.SharpeRatio = calculateSharpeRatio(initialPortfolio.ExpectedReturn, initialPortfolio.Risk, o.riskFreeRate)
	
	// Create a copy of the portfolio for optimization
	optimizedPortfolio := Portfolio{
		UserID:        portfolio.UserID,
		RiskTolerance: portfolio.RiskTolerance,
		TimeHorizon:   portfolio.TimeHorizon,
		Allocations:   make(map[string]float64),
		Timestamp:     time.Now(),
	}
	
	// Copy allocations
	for symbol, allocation := range portfolio.Allocations {
		optimizedPortfolio.Allocations[symbol] = allocation
	}
	
	// Run reinforcement learning optimization
	optimizedPortfolio = o.runRLOptimization(optimizedPortfolio, availableAssets)
	
	// Calculate optimized portfolio metrics
	optimizedPortfolio.ExpectedReturn = calculateExpectedReturn(optimizedPortfolio, availableAssets)
	optimizedPortfolio.Risk = calculatePortfolioRisk(optimizedPortfolio, availableAssets)
	optimizedPortfolio.SharpeRatio = calculateSharpeRatio(optimizedPortfolio.ExpectedReturn, optimizedPortfolio.Risk, o.riskFreeRate)
	
	// Calculate improvement
	improvement := 0.0
	if initialPortfolio.SharpeRatio > 0 {
		improvement = (optimizedPortfolio.SharpeRatio - initialPortfolio.SharpeRatio) / initialPortfolio.SharpeRatio
	} else if optimizedPortfolio.SharpeRatio > 0 {
		improvement = 1.0 // 100% improvement from zero or negative
	}
	
	// Generate trade recommendations
	tradeRecommendations := generateTradeRecommendations(initialPortfolio, optimizedPortfolio, availableAssets)
	
	// Create result
	result := OptimizationResult{
		InitialPortfolio:     initialPortfolio,
		OptimizedPortfolio:   optimizedPortfolio,
		Improvement:          improvement,
		TradeRecommendations: tradeRecommendations,
		Timestamp:            time.Now(),
	}
	
	// Cache result
	o.cacheMutex.Lock()
	o.resultCache[cacheKey] = result
	o.cacheMutex.Unlock()
	
	return &result, nil
}

// runRLOptimization runs the reinforcement learning optimization algorithm
func (o *PortfolioOptimizer) runRLOptimization(portfolio Portfolio, availableAssets []Asset) Portfolio {
	// This is a simplified implementation of RL for portfolio optimization
	// In a real implementation, this would be more sophisticated
	
	// Number of iterations
	iterations := 100
	
	// Best portfolio found so far
	bestPortfolio := portfolio
	bestSharpeRatio := calculateSharpeRatio(
		calculateExpectedReturn(portfolio, availableAssets),
		calculatePortfolioRisk(portfolio, availableAssets),
		o.riskFreeRate,
	)
	
	// Run iterations
	for i := 0; i < iterations; i++ {
		// Create a copy of the portfolio
		currentPortfolio := Portfolio{
			UserID:        portfolio.UserID,
			RiskTolerance: portfolio.RiskTolerance,
			TimeHorizon:   portfolio.TimeHorizon,
			Allocations:   make(map[string]float64),
		}
		
		// Copy allocations
		for symbol, allocation := range portfolio.Allocations {
			currentPortfolio.Allocations[symbol] = allocation
		}
		
		// Apply random adjustments (exploration)
		if rand.Float64() < o.explorationRate {
			currentPortfolio = randomlyAdjustPortfolio(currentPortfolio, availableAssets)
		} else {
			// Apply learned policy (exploitation)
			currentPortfolio = o.applyLearnedPolicy(currentPortfolio, availableAssets)
		}
		
		// Normalize allocations to sum to 100%
		normalizeAllocations(currentPortfolio.Allocations)
		
		// Calculate portfolio metrics
		expectedReturn := calculateExpectedReturn(currentPortfolio, availableAssets)
		risk := calculatePortfolioRisk(currentPortfolio, availableAssets)
		sharpeRatio := calculateSharpeRatio(expectedReturn, risk, o.riskFreeRate)
		
		// Update Q-table
		o.updateQTable(portfolio, currentPortfolio, sharpeRatio)
		
		// Update best portfolio if better
		if sharpeRatio > bestSharpeRatio {
			bestPortfolio = currentPortfolio
			bestSharpeRatio = sharpeRatio
		}
	}
	
	// Apply risk tolerance adjustment
	bestPortfolio = adjustForRiskTolerance(bestPortfolio, availableAssets)
	
	return bestPortfolio
}

// applyLearnedPolicy applies the learned policy to a portfolio
func (o *PortfolioOptimizer) applyLearnedPolicy(portfolio Portfolio, availableAssets []Asset) Portfolio {
	// Get state representation
	state := getStateRepresentation(portfolio)
	
	// Get Q-values for this state
	o.qTableMutex.RLock()
	qValues, exists := o.qTable[state]
	o.qTableMutex.RUnlock()
	
	if !exists {
		// No learned policy yet, use random adjustments
		return randomlyAdjustPortfolio(portfolio, availableAssets)
	}
	
	// Find best action
	var bestAction string
	bestValue := -math.MaxFloat64
	
	for action, value := range qValues {
		if value > bestValue {
			bestValue = value
			bestAction = action
		}
	}
	
	if bestAction == "" {
		// No best action found, use random adjustments
		return randomlyAdjustPortfolio(portfolio, availableAssets)
	}
	
	// Apply the action
	return applyAction(portfolio, bestAction, availableAssets)
}

// updateQTable updates the Q-table based on the observed reward
func (o *PortfolioOptimizer) updateQTable(oldPortfolio, newPortfolio Portfolio, reward float64) {
	// Get state and action representations
	oldState := getStateRepresentation(oldPortfolio)
	newState := getStateRepresentation(newPortfolio)
	action := getActionRepresentation(oldPortfolio, newPortfolio)
	
	// Get current Q-value
	o.qTableMutex.RLock()
	qValues, exists := o.qTable[oldState]
	if !exists {
		qValues = make(map[string]float64)
	}
	currentQValue := qValues[action]
	
	// Get max Q-value for next state
	nextQValues, nextExists := o.qTable[newState]
	o.qTableMutex.RUnlock()
	
	maxNextQValue := 0.0
	if nextExists {
		for _, value := range nextQValues {
			if value > maxNextQValue {
				maxNextQValue = value
			}
		}
	}
	
	// Calculate new Q-value
	newQValue := currentQValue + o.learningRate*(reward + o.discountFactor*maxNextQValue - currentQValue)
	
	// Update Q-table
	o.qTableMutex.Lock()
	if _, exists := o.qTable[oldState]; !exists {
		o.qTable[oldState] = make(map[string]float64)
	}
	o.qTable[oldState][action] = newQValue
	o.qTableMutex.Unlock()
}

// randomlyAdjustPortfolio makes random adjustments to a portfolio
func randomlyAdjustPortfolio(portfolio Portfolio, availableAssets []Asset) Portfolio {
	// Create a copy of the portfolio
	adjustedPortfolio := Portfolio{
		UserID:        portfolio.UserID,
		RiskTolerance: portfolio.RiskTolerance,
		TimeHorizon:   portfolio.TimeHorizon,
		Allocations:   make(map[string]float64),
	}
	
	// Copy allocations
	for symbol, allocation := range portfolio.Allocations {
		adjustedPortfolio.Allocations[symbol] = allocation
	}
	
	// Randomly select assets to adjust
	numAdjustments := 2 + rand.Intn(3) // 2-4 adjustments
	
	for i := 0; i < numAdjustments; i++ {
		// Randomly select an asset
		assetIndex := rand.Intn(len(availableAssets))
		asset := availableAssets[assetIndex]
		
		// Randomly adjust allocation
		adjustment := (rand.Float64() - 0.5) * 0.1 // -5% to +5%
		
		currentAllocation := adjustedPortfolio.Allocations[asset.Symbol]
		newAllocation := currentAllocation + adjustment
		
		// Ensure allocation is non-negative
		if newAllocation < 0 {
			newAllocation = 0
		}
		
		adjustedPortfolio.Allocations[asset.Symbol] = newAllocation
	}
	
	// Normalize allocations to sum to 100%
	normalizeAllocations(adjustedPortfolio.Allocations)
	
	return adjustedPortfolio
}

// applyAction applies an action to a portfolio
func applyAction(portfolio Portfolio, action string, availableAssets []Asset) Portfolio {
	// In a real implementation, this would parse the action and apply it
	// This is a simplified implementation
	
	// Create a copy of the portfolio
	adjustedPortfolio := Portfolio{
		UserID:        portfolio.UserID,
		RiskTolerance: portfolio.RiskTolerance,
		TimeHorizon:   portfolio.TimeHorizon,
		Allocations:   make(map[string]float64),
	}
	
	// Copy allocations
	for symbol, allocation := range portfolio.Allocations {
		adjustedPortfolio.Allocations[symbol] = allocation
	}
	
	// Parse action (simplified)
	// Format: "ADJUST_SYMBOL_AMOUNT"
	parts := strings.Split(action, "_")
	if len(parts) != 3 {
		return portfolio // Invalid action
	}
	
	symbol := parts[1]
	amountStr := parts[2]
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return portfolio // Invalid action
	}
	
	// Apply adjustment
	currentAllocation := adjustedPortfolio.Allocations[symbol]
	newAllocation := currentAllocation + amount
	
	// Ensure allocation is non-negative
	if newAllocation < 0 {
		newAllocation = 0
	}
	
	adjustedPortfolio.Allocations[symbol] = newAllocation
	
	// Normalize allocations to sum to 100%
	normalizeAllocations(adjustedPortfolio.Allocations)
	
	return adjustedPortfolio
}

// adjustForRiskTolerance adjusts a portfolio based on risk tolerance
func adjustForRiskTolerance(portfolio Portfolio, availableAssets []Asset) Portfolio {
	// Create a copy of the portfolio
	adjustedPortfolio := Portfolio{
		UserID:        portfolio.UserID,
		RiskTolerance: portfolio.RiskTolerance,
		TimeHorizon:   portfolio.TimeHorizon,
		Allocations:   make(map[string]float64),
	}
	
	// Copy allocations
	for symbol, allocation := range portfolio.Allocations {
		adjustedPortfolio.Allocations[symbol] = allocation
	}
	
	// Calculate current portfolio risk
	currentRisk := calculatePortfolioRisk(portfolio, availableAssets)
	
	// Calculate target risk based on risk tolerance
	// Risk tolerance is 0.0 (very conservative) to 1.0 (very aggressive)
	targetRisk := 0.05 + portfolio.RiskTolerance * 0.15 // 5% to 20% volatility
	
	// If current risk is close to target risk, no adjustment needed
	if math.Abs(currentRisk - targetRisk) < 0.02 {
		return adjustedPortfolio
	}
	
	// Sort assets by volatility
	type AssetVolatility struct {
		Symbol     string
		Volatility float64
	}
	
	var assetVolatilities []AssetVolatility
	for _, asset := range availableAssets {
		assetVolatilities = append(assetVolatilities, AssetVolatility{
			Symbol:     asset.Symbol,
			Volatility: asset.Volatility,
		})
	}
	
	sort.Slice(assetVolatilities, func(i, j int) bool {
		return assetVolatilities[i].Volatility < assetVolatilities[j].Volatility
	})
	
	// Adjust allocations based on risk tolerance
	if currentRisk > targetRisk {
		// Current portfolio is too risky, shift towards less volatile assets
		for i := 0; i < len(assetVolatilities)/2; i++ {
			// Increase allocation to less volatile assets
			lowVolAsset := assetVolatilities[i]
			highVolAsset := assetVolatilities[len(assetVolatilities)-i-1]
			
			// Shift 5% from high volatility to low volatility
			adjustment := math.Min(0.05, adjustedPortfolio.Allocations[highVolAsset.Symbol])
			adjustedPortfolio.Allocations[highVolAsset.Symbol] -= adjustment
			adjustedPortfolio.Allocations[lowVolAsset.Symbol] += adjustment
		}
	} else {
		// Current portfolio is too conservative, shift towards more volatile assets
		for i := 0; i < len(assetVolatilities)/2; i++ {
			// Increase allocation to more volatile assets
			lowVolAsset := assetVolatilities[i]
			highVolAsset := assetVolatilities[len(assetVolatilities)-i-1]
			
			// Shift 5% from low volatility to high volatility
			adjustment := math.Min(0.05, adjustedPortfolio.Allocations[lowVolAsset.Symbol])
			adjustedPortfolio.Allocations[lowVolAsset.Symbol] -= adjustment
			adjustedPortfolio.Allocations[highVolAsset.Symbol] += adjustment
		}
	}
	
	// Normalize allocations to sum to 100%
	normalizeAllocations(adjustedPortfolio.Allocations)
	
	return adjustedPortfolio
}

// normalizeAllocations normalizes allocations to sum to 100%
func normalizeAllocations(allocations map[string]float64) {
	// Calculate sum
	var sum float64
	for _, allocation := range allocations {
		sum += allocation
	}
	
	// Normalize
	if sum > 0 {
		for symbol, allocation := range allocations {
			allocations[symbol] = allocation / sum
		}
	}
}

// getStateRepresentation gets a string representation of a portfolio state
func getStateRepresentation(portfolio Portfolio) string {
	// In a real implementation, this would create a more sophisticated state representation
	// This is a simplified implementation
	
	// Sort symbols for consistent representation
	var symbols []string
	for symbol := range portfolio.Allocations {
		symbols = append(symbols, symbol)
	}
	sort.Strings(symbols)
	
	// Create state representation
	var state strings.Builder
	state.WriteString(fmt.Sprintf("RT%.2f_TH%d", portfolio.RiskTolerance, portfolio.TimeHorizon))
	
	for _, symbol := range symbols {
		allocation := portfolio.Allocations[symbol]
		state.WriteString(fmt.Sprintf("_%s%.2f", symbol, allocation))
	}
	
	return state.String()
}

// getActionRepresentation gets a string representation of an action
func getActionRepresentation(oldPortfolio, newPortfolio Portfolio) string {
	// In a real implementation, this would create a more sophisticated action representation
	// This is a simplified implementation
	
	// Find the symbol with the largest change
	var maxChangeSymbol string
	var maxChangeAmount float64
	
	for symbol, newAllocation := range newPortfolio.Allocations {
		oldAllocation := oldPortfolio.Allocations[symbol]
		change := newAllocation - oldAllocation
		
		if math.Abs(change) > math.Abs(maxChangeAmount) {
			maxChangeSymbol = symbol
			maxChangeAmount = change
		}
	}
	
	// Create action representation
	return fmt.Sprintf("ADJUST_%s_%.4f", maxChangeSymbol, maxChangeAmount)
}

// calculateExpectedReturn calculates the expected return of a portfolio
func calculateExpectedReturn(portfolio Portfolio, assets []Asset) float64 {
	// Create a map of assets for quick lookup
	assetMap := make(map[string]Asset)
	for _, asset := range assets {
		assetMap[asset.Symbol] = asset
	}
	
	// Calculate weighted average of expected returns
	var expectedReturn float64
	
	for symbol, allocation := range portfolio.Allocations {
		asset, exists := assetMap[symbol]
		if !exists {
			continue
		}
		
		// Calculate expected return for this asset
		assetReturn := calculateAssetExpectedReturn(asset)
		
		// Add to weighted average
		expectedReturn += allocation * assetReturn
	}
	
	return expectedReturn
}

// calculateAssetExpectedReturn calculates the expected return of an asset
func calculateAssetExpectedReturn(asset Asset) float64 {
	// If historical returns are available, use their average
	if len(asset.Returns) > 0 {
		var sum float64
		for _, r := range asset.Returns {
			sum += r
		}
		return sum / float64(len(asset.Returns))
	}
	
	// Otherwise, use a default based on asset type
	switch asset.AssetType {
	case "STOCK":
		return 0.08 // 8% expected return
	case "BOND":
		return 0.04 // 4% expected return
	case "ETF":
		return 0.06 // 6% expected return
	default:
		return 0.05 // 5% expected return
	}
}

// calculatePortfolioRisk calculates the risk (volatility) of a portfolio
func calculatePortfolioRisk(portfolio Portfolio, assets []Asset) float64 {
	// Create a map of assets for quick lookup
	assetMap := make(map[string]Asset)
	for _, asset := range assets {
		assetMap[asset.Symbol] = asset
	}
	
	// In a real implementation, this would calculate the portfolio variance using a covariance matrix
	// This is a simplified implementation that ignores correlations
	
	var portfolioVariance float64
	
	for symbol, allocation := range portfolio.Allocations {
		asset, exists := assetMap[symbol]
		if !exists {
			continue
		}
		
		// Add to portfolio variance (ignoring correlations)
		portfolioVariance += allocation * allocation * asset.Volatility * asset.Volatility
	}
	
	// Return portfolio volatility (standard deviation)
	return math.Sqrt(portfolioVariance)
}

// calculateSharpeRatio calculates the Sharpe ratio of a portfolio
func calculateSharpeRatio(expectedReturn, risk, riskFreeRate float64) float64 {
	if risk == 0 {
		return 0 // Avoid division by zero
	}
	
	return (expectedReturn - riskFreeRate) / risk
}

// generateTradeRecommendations generates trade recommendations
func generateTradeRecommendations(initialPortfolio, optimizedPortfolio Portfolio, assets []Asset) []TradeRecommendation {
	var recommendations []TradeRecommendation
	
	// Create a map of assets for quick lookup
	assetMap := make(map[string]Asset)
	for _, asset := range assets {
		assetMap[asset.Symbol] = asset
	}
	
	// Find all symbols in either portfolio
	allSymbols := make(map[string]bool)
	for symbol := range initialPortfolio.Allocations {
		allSymbols[symbol] = true
	}
	for symbol := range optimizedPortfolio.Allocations {
		allSymbols[symbol] = true
	}
	
	// Generate recommendations for each symbol
	for symbol := range allSymbols {
		initialAllocation := initialPortfolio.Allocations[symbol]
		optimizedAllocation := optimizedPortfolio.Allocations[symbol]
		
		// Calculate change
		change := optimizedAllocation - initialAllocation
		
		// Determine action
		action := "HOLD"
		if change > 0.01 { // 1% threshold for buy
			action = "BUY"
		} else if change < -0.01 { // -1% threshold for sell
			action = "SELL"
		}
		
		// Skip if no significant change
		if action == "HOLD" {
			continue
		}
		
		// Generate reasoning
		reasoning := generateTradeReasoning(symbol, action, change, assetMap[symbol], initialPortfolio, optimizedPortfolio)
		
		// Create recommendation
		recommendation := TradeRecommendation{
			Symbol:            symbol,
			Action:            action,
			CurrentAllocation: initialAllocation * 100, // Convert to percentage
			TargetAllocation:  optimizedAllocation * 100, // Convert to percentage
			ChangeAmount:      change * 100, // Convert to percentage
			Reasoning:         reasoning,
		}
		
		recommendations = append(recommendations, recommendation)
	}
	
	// Sort recommendations by absolute change amount (descending)
	sort.Slice(recommendations, func(i, j int) bool {
		return math.Abs(recommendations[i].ChangeAmount) > math.Abs(recommendations[j].ChangeAmount)
	})
	
	return recommendations
}

// generateTradeReasoning generates reasoning for a trade recommendation
func generateTradeReasoning(symbol, action string, change float64, asset Asset, initialPortfolio, optimizedPortfolio Portfolio) string {
	var reasoning string
	
	switch action {
	case "BUY":
		reasoning = fmt.Sprintf("Increase allocation to %s by %.1f%% (from %.1f%% to %.1f%%). ", 
			symbol, change*100, initialPortfolio.Allocations[symbol]*100, optimizedPortfolio.Allocations[symbol]*100)
		
		// Add reasoning based on portfolio characteristics
		if optimizedPortfolio.ExpectedReturn > initialPortfolio.ExpectedReturn {
			reasoning += "This increases the portfolio's expected return. "
		}
		
		if optimizedPortfolio.SharpeRatio > initialPortfolio.SharpeRatio {
			reasoning += "This improves the portfolio's risk-adjusted return (Sharpe ratio). "
		}
		
		// Add reasoning based on asset characteristics
		if asset.AssetType == "BOND" && initialPortfolio.RiskTolerance < 0.3 {
			reasoning += "This aligns with your conservative risk profile. "
		} else if asset.AssetType == "STOCK" && initialPortfolio.RiskTolerance > 0.7 {
			reasoning += "This aligns with your aggressive risk profile. "
		}
		
	case "SELL":
		reasoning = fmt.Sprintf("Decrease allocation to %s by %.1f%% (from %.1f%% to %.1f%%). ", 
			symbol, -change*100, initialPortfolio.Allocations[symbol]*100, optimizedPortfolio.Allocations[symbol]*100)
		
		// Add reasoning based on portfolio characteristics
		if optimizedPortfolio.Risk < initialPortfolio.Risk {
			reasoning += "This reduces the portfolio's overall risk. "
		}
		
		if optimizedPortfolio.SharpeRatio > initialPortfolio.SharpeRatio {
			reasoning += "This improves the portfolio's risk-adjusted return (Sharpe ratio). "
		}
		
		// Add reasoning based on asset characteristics
		if asset.AssetType == "STOCK" && initialPortfolio.RiskTolerance < 0.3 {
			reasoning += "This better aligns with your conservative risk profile. "
		} else if asset.AssetType == "BOND" && initialPortfolio.RiskTolerance > 0.7 {
			reasoning += "This better aligns with your aggressive risk profile. "
		}
	}
	
	return reasoning
}
