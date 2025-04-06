package rl

import (
	"context"
	"math"
	"sync"
	"time"
)

// State represents the state of the environment
type State struct {
	Portfolio      Portfolio            `json:"portfolio"`
	MarketState    map[string]float64   `json:"market_state"`    // Market indicators
	AssetPrices    map[string]float64   `json:"asset_prices"`    // Current asset prices
	AssetTrends    map[string]float64   `json:"asset_trends"`    // Price trends (-1.0 to 1.0)
	Timestamp      time.Time            `json:"timestamp"`
}

// Action represents an action to take
type Action struct {
	Type           string              `json:"type"`           // "BUY", "SELL", "HOLD", "REBALANCE"
	Symbol         string              `json:"symbol,omitempty"`
	Amount         float64             `json:"amount,omitempty"`
	Allocations    map[string]float64  `json:"allocations,omitempty"` // For rebalance actions
}

// Experience represents a learning experience
type Experience struct {
	State          State               `json:"state"`
	Action         Action              `json:"action"`
	Reward         float64             `json:"reward"`
	NextState      State               `json:"next_state"`
	Done           bool                `json:"done"`
}

// RLAgent is a reinforcement learning agent for portfolio management
type RLAgent struct {
	// Learning parameters
	learningRate     float64
	discountFactor   float64
	explorationRate  float64
	
	// Experience replay buffer
	experienceBuffer []Experience
	bufferSize       int
	bufferMutex      sync.RWMutex
	
	// Neural network weights (simplified)
	weights          map[string]float64
	weightsMutex     sync.RWMutex
	
	// Training state
	trainingEpisodes int
	lastTrainingTime time.Time
}

// NewRLAgent creates a new reinforcement learning agent
func NewRLAgent() *RLAgent {
	return &RLAgent{
		learningRate:     0.01,
		discountFactor:   0.95,
		explorationRate:  0.1,
		experienceBuffer: make([]Experience, 0),
		bufferSize:       1000,
		weights:          make(map[string]float64),
		trainingEpisodes: 0,
		lastTrainingTime: time.Now(),
	}
}

// GetAction gets the best action for a given state
func (a *RLAgent) GetAction(ctx context.Context, state State) Action {
	// Exploration: randomly select an action
	if rand.Float64() < a.explorationRate {
		return a.getRandomAction(state)
	}
	
	// Exploitation: select the best action according to the policy
	return a.getBestAction(state)
}

// AddExperience adds an experience to the replay buffer
func (a *RLAgent) AddExperience(experience Experience) {
	a.bufferMutex.Lock()
	defer a.bufferMutex.Unlock()
	
	// Add experience to buffer
	a.experienceBuffer = append(a.experienceBuffer, experience)
	
	// Trim buffer if it exceeds the maximum size
	if len(a.experienceBuffer) > a.bufferSize {
		a.experienceBuffer = a.experienceBuffer[1:]
	}
}

// Train trains the agent using experiences from the replay buffer
func (a *RLAgent) Train(ctx context.Context) error {
	a.bufferMutex.RLock()
	bufferSize := len(a.experienceBuffer)
	a.bufferMutex.RUnlock()
	
	// Skip training if buffer is too small
	if bufferSize < 100 {
		return nil
	}
	
	// Number of training iterations
	iterations := 10
	
	// Batch size
	batchSize := 32
	
	// Train for multiple iterations
	for i := 0; i < iterations; i++ {
		// Sample a batch of experiences
		batch := a.sampleExperiences(batchSize)
		
		// Update weights using the batch
		a.updateWeights(batch)
	}
	
	// Update training state
	a.trainingEpisodes++
	a.lastTrainingTime = time.Now()
	
	// Decay exploration rate
	a.explorationRate = math.Max(0.01, a.explorationRate*0.995)
	
	return nil
}

// getRandomAction returns a random action
func (a *RLAgent) getRandomAction(state State) Action {
	// Get all available symbols
	var symbols []string
	for symbol := range state.AssetPrices {
		symbols = append(symbols, symbol)
	}
	
	// No symbols available
	if len(symbols) == 0 {
		return Action{
			Type: "HOLD",
		}
	}
	
	// Randomly select an action type
	actionTypes := []string{"BUY", "SELL", "HOLD", "REBALANCE"}
	actionType := actionTypes[rand.Intn(len(actionTypes))]
	
	switch actionType {
	case "BUY", "SELL":
		// Randomly select a symbol
		symbol := symbols[rand.Intn(len(symbols))]
		
		// Randomly select an amount (1-10% of portfolio)
		amount := (1 + rand.Float64()*9) / 100.0
		
		return Action{
			Type:   actionType,
			Symbol: symbol,
			Amount: amount,
		}
		
	case "REBALANCE":
		// Generate random allocations
		allocations := make(map[string]float64)
		remainingAllocation := 1.0
		
		for i, symbol := range symbols {
			if i == len(symbols)-1 {
				// Last symbol gets the remaining allocation
				allocations[symbol] = remainingAllocation
			} else {
				// Random allocation between 0% and remaining allocation
				allocation := rand.Float64() * remainingAllocation
				allocations[symbol] = allocation
				remainingAllocation -= allocation
			}
		}
		
		return Action{
			Type:        "REBALANCE",
			Allocations: allocations,
		}
		
	default: // HOLD
		return Action{
			Type: "HOLD",
		}
	}
}

// getBestAction returns the best action according to the policy
func (a *RLAgent) getBestAction(state State) Action {
	// Get all available symbols
	var symbols []string
	for symbol := range state.AssetPrices {
		symbols = append(symbols, symbol)
	}
	
	// No symbols available
	if len(symbols) == 0 {
		return Action{
			Type: "HOLD",
		}
	}
	
	// Calculate Q-values for all possible actions
	bestAction := Action{
		Type: "HOLD",
	}
	bestQValue := a.calculateQValue(state, bestAction)
	
	// Check BUY and SELL actions for each symbol
	for _, symbol := range symbols {
		// Check BUY action
		buyAction := Action{
			Type:   "BUY",
			Symbol: symbol,
			Amount: 0.05, // 5% of portfolio
		}
		buyQValue := a.calculateQValue(state, buyAction)
		
		if buyQValue > bestQValue {
			bestAction = buyAction
			bestQValue = buyQValue
		}
		
		// Check SELL action
		sellAction := Action{
			Type:   "SELL",
			Symbol: symbol,
			Amount: 0.05, // 5% of portfolio
		}
		sellQValue := a.calculateQValue(state, sellAction)
		
		if sellQValue > bestQValue {
			bestAction = sellAction
			bestQValue = sellQValue
		}
	}
	
	// Check REBALANCE action
	rebalanceAction := a.generateRebalanceAction(state)
	rebalanceQValue := a.calculateQValue(state, rebalanceAction)
	
	if rebalanceQValue > bestQValue {
		bestAction = rebalanceAction
		bestQValue = rebalanceQValue
	}
	
	return bestAction
}

// generateRebalanceAction generates a rebalance action based on the current state
func (a *RLAgent) generateRebalanceAction(state State) Action {
	// Get all available symbols
	var symbols []string
	for symbol := range state.AssetPrices {
		symbols = append(symbols, symbol)
	}
	
	// Generate allocations based on asset trends and market state
	allocations := make(map[string]float64)
	totalScore := 0.0
	
	for _, symbol := range symbols {
		// Calculate a score for each asset based on its trend and market state
		trend := state.AssetTrends[symbol]
		
		// Adjust trend based on market state
		marketState := 0.0
		for indicator, value := range state.MarketState {
			if indicator == "MARKET_TREND" {
				marketState = value
				break
			}
		}
		
		// Combine trend and market state
		score := trend + marketState*0.5
		
		// Ensure score is positive
		score = math.Max(0.1, score+1.0)
		
		allocations[symbol] = score
		totalScore += score
	}
	
	// Normalize allocations
	for symbol := range allocations {
		allocations[symbol] /= totalScore
	}
	
	return Action{
		Type:        "REBALANCE",
		Allocations: allocations,
	}
}

// calculateQValue calculates the Q-value for a state-action pair
func (a *RLAgent) calculateQValue(state State, action Action) float64 {
	// Extract features from state and action
	features := a.extractFeatures(state, action)
	
	// Calculate Q-value as a linear combination of features and weights
	qValue := 0.0
	
	a.weightsMutex.RLock()
	for feature, value := range features {
		weight, exists := a.weights[feature]
		if !exists {
			weight = 0.0
		}
		qValue += weight * value
	}
	a.weightsMutex.RUnlock()
	
	return qValue
}

// extractFeatures extracts features from a state-action pair
func (a *RLAgent) extractFeatures(state State, action Action) map[string]float64 {
	features := make(map[string]float64)
	
	// Basic features
	features["bias"] = 1.0
	
	// Portfolio features
	features["portfolio_expected_return"] = state.Portfolio.ExpectedReturn
	features["portfolio_risk"] = state.Portfolio.Risk
	features["portfolio_sharpe_ratio"] = state.Portfolio.SharpeRatio
	
	// Market state features
	for indicator, value := range state.MarketState {
		features["market_"+indicator] = value
	}
	
	// Action-specific features
	features["action_type_"+action.Type] = 1.0
	
	switch action.Type {
	case "BUY", "SELL":
		// Asset-specific features
		if price, exists := state.AssetPrices[action.Symbol]; exists {
			features["asset_price_"+action.Symbol] = price
		}
		
		if trend, exists := state.AssetTrends[action.Symbol]; exists {
			features["asset_trend_"+action.Symbol] = trend
		}
		
		// Action amount
		features["action_amount"] = action.Amount
		
		// Current allocation
		if allocation, exists := state.Portfolio.Allocations[action.Symbol]; exists {
			features["current_allocation_"+action.Symbol] = allocation
		}
		
	case "REBALANCE":
		// Calculate allocation changes
		for symbol, targetAllocation := range action.Allocations {
			currentAllocation := 0.0
			if allocation, exists := state.Portfolio.Allocations[symbol]; exists {
				currentAllocation = allocation
			}
			
			features["allocation_change_"+symbol] = targetAllocation - currentAllocation
		}
	}
	
	return features
}

// sampleExperiences samples a batch of experiences from the replay buffer
func (a *RLAgent) sampleExperiences(batchSize int) []Experience {
	a.bufferMutex.RLock()
	defer a.bufferMutex.RUnlock()
	
	// Ensure batch size is not larger than buffer size
	if batchSize > len(a.experienceBuffer) {
		batchSize = len(a.experienceBuffer)
	}
	
	// Sample experiences without replacement
	indices := rand.Perm(len(a.experienceBuffer))[:batchSize]
	
	batch := make([]Experience, batchSize)
	for i, idx := range indices {
		batch[i] = a.experienceBuffer[idx]
	}
	
	return batch
}

// updateWeights updates the weights using a batch of experiences
func (a *RLAgent) updateWeights(batch []Experience) {
	a.weightsMutex.Lock()
	defer a.weightsMutex.Unlock()
	
	// For each experience in the batch
	for _, experience := range batch {
		// Calculate target Q-value
		targetQ := experience.Reward
		
		if !experience.Done {
			// Get best action for next state
			bestNextAction := a.getBestAction(experience.NextState)
			
			// Calculate Q-value for next state and best action
			nextQ := a.calculateQValue(experience.NextState, bestNextAction)
			
			// Update target Q-value
			targetQ += a.discountFactor * nextQ
		}
		
		// Calculate current Q-value
		currentQ := a.calculateQValue(experience.State, experience.Action)
		
		// Calculate TD error
		tdError := targetQ - currentQ
		
		// Extract features
		features := a.extractFeatures(experience.State, experience.Action)
		
		// Update weights
		for feature, value := range features {
			// Get current weight
			weight, exists := a.weights[feature]
			if !exists {
				weight = 0.0
			}
			
			// Update weight
			a.weights[feature] = weight + a.learningRate * tdError * value
		}
	}
}
