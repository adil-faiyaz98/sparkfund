package nlp

import (
	"math"
	"strings"
	"sync"
	"time"
)

// SentimentAnalyzer analyzes market news for sentiment
type SentimentAnalyzer struct {
	// Lexicon-based sentiment analysis
	positiveLexicon map[string]float64
	negativeLexicon map[string]float64
	
	// Entity recognition
	companyNames    map[string]string // maps company name variations to canonical names
	sectorKeywords  map[string]string // maps sector-related keywords to sectors
	
	// Cached results
	cache      map[string]SentimentResult
	cacheMutex sync.RWMutex
	cacheTime  time.Duration
}

// SentimentResult represents the result of sentiment analysis
type SentimentResult struct {
	Text            string               `json:"text"`
	OverallScore    float64              `json:"overall_score"`    // -1.0 (very negative) to 1.0 (very positive)
	Magnitude       float64              `json:"magnitude"`        // 0.0 (neutral) to 1.0 (strong)
	CompanySentiment map[string]float64  `json:"company_sentiment"` // Company name -> sentiment score
	SectorSentiment  map[string]float64  `json:"sector_sentiment"`  // Sector -> sentiment score
	Keywords        []string             `json:"keywords"`
	Timestamp       time.Time            `json:"timestamp"`
}

// MarketImpact represents the predicted market impact of news
type MarketImpact struct {
	CompanyImpacts map[string]Impact `json:"company_impacts"` // Company symbol -> impact
	SectorImpacts  map[string]Impact `json:"sector_impacts"`  // Sector -> impact
	MarketImpact   Impact            `json:"market_impact"`   // Overall market impact
	Timestamp      time.Time         `json:"timestamp"`
}

// Impact represents the predicted impact on price and volume
type Impact struct {
	Symbol        string  `json:"symbol,omitempty"`
	Sector        string  `json:"sector,omitempty"`
	PriceImpact   float64 `json:"price_impact"`   // Predicted percentage change
	VolumeImpact  float64 `json:"volume_impact"`  // Predicted percentage change
	Confidence    float64 `json:"confidence"`     // 0.0 to 1.0
	TimePeriod    string  `json:"time_period"`    // "SHORT_TERM", "MEDIUM_TERM", "LONG_TERM"
}

// NewSentimentAnalyzer creates a new sentiment analyzer
func NewSentimentAnalyzer() *SentimentAnalyzer {
	analyzer := &SentimentAnalyzer{
		positiveLexicon: loadPositiveLexicon(),
		negativeLexicon: loadNegativeLexicon(),
		companyNames:    loadCompanyNames(),
		sectorKeywords:  loadSectorKeywords(),
		cache:           make(map[string]SentimentResult),
		cacheTime:       15 * time.Minute,
	}
	return analyzer
}

// AnalyzeText analyzes text for sentiment
func (a *SentimentAnalyzer) AnalyzeText(text string) SentimentResult {
	// Check cache first
	a.cacheMutex.RLock()
	if result, ok := a.cache[text]; ok {
		if time.Since(result.Timestamp) < a.cacheTime {
			a.cacheMutex.RUnlock()
			return result
		}
	}
	a.cacheMutex.RUnlock()
	
	// Preprocess text
	processedText := preprocess(text)
	words := strings.Fields(processedText)
	
	// Calculate overall sentiment
	var positiveScore, negativeScore float64
	for _, word := range words {
		if score, ok := a.positiveLexicon[word]; ok {
			positiveScore += score
		}
		if score, ok := a.negativeLexicon[word]; ok {
			negativeScore += score
		}
	}
	
	// Calculate overall score and magnitude
	overallScore := (positiveScore - negativeScore) / math.Max(1.0, positiveScore+negativeScore)
	magnitude := math.Min(1.0, (positiveScore+negativeScore)/float64(len(words)))
	
	// Extract company and sector sentiment
	companySentiment := make(map[string]float64)
	sectorSentiment := make(map[string]float64)
	
	// Extract company sentiment
	for company, variations := range a.companyNames {
		for _, variation := range strings.Split(variations, ",") {
			if strings.Contains(processedText, strings.ToLower(variation)) {
				// Calculate sentiment in the vicinity of the company mention
				companySentiment[company] = a.calculateLocalSentiment(words, variation)
				break
			}
		}
	}
	
	// Extract sector sentiment
	for sector, keywords := range a.sectorKeywords {
		for _, keyword := range strings.Split(keywords, ",") {
			if strings.Contains(processedText, strings.ToLower(keyword)) {
				// Calculate sentiment in the vicinity of the sector mention
				sectorSentiment[sector] = a.calculateLocalSentiment(words, keyword)
				break
			}
		}
	}
	
	// Extract keywords
	keywords := extractKeywords(words, 10)
	
	// Create result
	result := SentimentResult{
		Text:             text,
		OverallScore:     overallScore,
		Magnitude:        magnitude,
		CompanySentiment: companySentiment,
		SectorSentiment:  sectorSentiment,
		Keywords:         keywords,
		Timestamp:        time.Now(),
	}
	
	// Cache result
	a.cacheMutex.Lock()
	a.cache[text] = result
	a.cacheMutex.Unlock()
	
	return result
}

// PredictMarketImpact predicts the market impact of news
func (a *SentimentAnalyzer) PredictMarketImpact(sentiment SentimentResult) MarketImpact {
	companyImpacts := make(map[string]Impact)
	sectorImpacts := make(map[string]Impact)
	
	// Calculate company impacts
	for company, sentiment := range sentiment.CompanySentiment {
		// Simple model: sentiment directly affects price
		priceImpact := sentiment * 0.02 // 2% max impact
		volumeImpact := math.Abs(sentiment) * 0.05 // 5% max impact
		
		// Determine time period based on magnitude
		timePeriod := "SHORT_TERM"
		if sentiment.Magnitude > 0.7 {
			timePeriod = "MEDIUM_TERM"
		}
		
		companyImpacts[company] = Impact{
			Symbol:       company,
			PriceImpact:  priceImpact,
			VolumeImpact: volumeImpact,
			Confidence:   sentiment.Magnitude,
			TimePeriod:   timePeriod,
		}
	}
	
	// Calculate sector impacts
	for sector, sentiment := range sentiment.SectorSentiment {
		// Sectors typically have less volatility than individual stocks
		priceImpact := sentiment * 0.01 // 1% max impact
		volumeImpact := math.Abs(sentiment) * 0.03 // 3% max impact
		
		// Determine time period based on magnitude
		timePeriod := "SHORT_TERM"
		if sentiment.Magnitude > 0.7 {
			timePeriod = "MEDIUM_TERM"
		}
		
		sectorImpacts[sector] = Impact{
			Sector:       sector,
			PriceImpact:  priceImpact,
			VolumeImpact: volumeImpact,
			Confidence:   sentiment.Magnitude,
			TimePeriod:   timePeriod,
		}
	}
	
	// Calculate overall market impact
	// Market impact is typically less than individual sectors
	marketImpact := Impact{
		PriceImpact:  sentiment.OverallScore * 0.005, // 0.5% max impact
		VolumeImpact: math.Abs(sentiment.OverallScore) * 0.02, // 2% max impact
		Confidence:   sentiment.Magnitude,
		TimePeriod:   "SHORT_TERM",
	}
	
	return MarketImpact{
		CompanyImpacts: companyImpacts,
		SectorImpacts:  sectorImpacts,
		MarketImpact:   marketImpact,
		Timestamp:      time.Now(),
	}
}

// calculateLocalSentiment calculates sentiment in the vicinity of a term
func (a *SentimentAnalyzer) calculateLocalSentiment(words []string, term string) float64 {
	term = strings.ToLower(term)
	
	// Find the term in the words
	var termIndex int
	for i, word := range words {
		if strings.ToLower(word) == term {
			termIndex = i
			break
		}
	}
	
	// Calculate sentiment in a window around the term
	windowSize := 5
	startIndex := math.Max(0, float64(termIndex-windowSize))
	endIndex := math.Min(float64(len(words)-1), float64(termIndex+windowSize))
	
	var positiveScore, negativeScore float64
	for i := int(startIndex); i <= int(endIndex); i++ {
		word := words[i]
		if score, ok := a.positiveLexicon[word]; ok {
			positiveScore += score
		}
		if score, ok := a.negativeLexicon[word]; ok {
			negativeScore += score
		}
	}
	
	// Calculate overall score
	windowLength := endIndex - startIndex + 1
	if windowLength == 0 {
		return 0
	}
	
	overallScore := (positiveScore - negativeScore) / math.Max(1.0, positiveScore+negativeScore)
	return overallScore
}

// preprocess preprocesses text for sentiment analysis
func preprocess(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)
	
	// Remove punctuation
	text = strings.Map(func(r rune) rune {
		if strings.ContainsRune(".,!?;:()\"-", r) {
			return ' '
		}
		return r
	}, text)
	
	// Replace multiple spaces with a single space
	text = strings.Join(strings.Fields(text), " ")
	
	return text
}

// extractKeywords extracts the most important keywords from text
func extractKeywords(words []string, maxKeywords int) []string {
	// Simple implementation: remove stopwords and return most frequent words
	stopwords := getStopwords()
	wordFreq := make(map[string]int)
	
	for _, word := range words {
		if !stopwords[word] && len(word) > 2 {
			wordFreq[word]++
		}
	}
	
	// Sort words by frequency
	type wordCount struct {
		word  string
		count int
	}
	
	var wordCounts []wordCount
	for word, count := range wordFreq {
		wordCounts = append(wordCounts, wordCount{word, count})
	}
	
	// Sort by count (descending)
	sort.Slice(wordCounts, func(i, j int) bool {
		return wordCounts[i].count > wordCounts[j].count
	})
	
	// Extract top keywords
	var keywords []string
	for i := 0; i < len(wordCounts) && i < maxKeywords; i++ {
		keywords = append(keywords, wordCounts[i].word)
	}
	
	return keywords
}

// loadPositiveLexicon loads the positive sentiment lexicon
func loadPositiveLexicon() map[string]float64 {
	// In a real implementation, this would load from a file or database
	// This is a simplified version with a few examples
	return map[string]float64{
		"good":      0.8,
		"great":     0.9,
		"excellent": 1.0,
		"positive":  0.8,
		"profit":    0.7,
		"growth":    0.7,
		"increase":  0.6,
		"up":        0.5,
		"gain":      0.7,
		"success":   0.8,
		"bullish":   0.9,
		"strong":    0.7,
		"improve":   0.6,
		"beat":      0.7,
		"exceed":    0.8,
		"outperform": 0.9,
		"rally":     0.8,
		"recovery":  0.7,
		"opportunity": 0.6,
		"innovative": 0.7,
		"launch":    0.6,
		"partnership": 0.7,
		"expansion": 0.7,
		"dividend":  0.6,
		"upgrade":   0.8,
	}
}

// loadNegativeLexicon loads the negative sentiment lexicon
func loadNegativeLexicon() map[string]float64 {
	// In a real implementation, this would load from a file or database
	// This is a simplified version with a few examples
	return map[string]float64{
		"bad":       0.8,
		"poor":      0.7,
		"negative":  0.8,
		"loss":      0.9,
		"decline":   0.7,
		"decrease":  0.6,
		"down":      0.5,
		"fall":      0.6,
		"drop":      0.6,
		"bearish":   0.9,
		"weak":      0.7,
		"worsen":    0.8,
		"miss":      0.7,
		"below":     0.6,
		"underperform": 0.8,
		"sell":      0.6,
		"concern":   0.7,
		"risk":      0.7,
		"volatile":  0.6,
		"uncertainty": 0.7,
		"layoff":    0.9,
		"cut":       0.7,
		"debt":      0.6,
		"litigation": 0.8,
		"downgrade": 0.8,
		"recall":    0.8,
		"investigation": 0.8,
		"delay":     0.6,
		"warning":   0.7,
	}
}

// loadCompanyNames loads company names and their variations
func loadCompanyNames() map[string]string {
	// In a real implementation, this would load from a file or database
	// This is a simplified version with a few examples
	return map[string]string{
		"AAPL": "apple,apple inc",
		"MSFT": "microsoft,microsoft corporation",
		"GOOGL": "google,alphabet,alphabet inc",
		"AMZN": "amazon,amazon.com",
		"META": "meta,facebook,meta platforms",
		"TSLA": "tesla,tesla inc",
		"NVDA": "nvidia,nvidia corporation",
		"JPM": "jpmorgan,jp morgan,jpmorgan chase",
		"BAC": "bank of america,bofa",
		"WMT": "walmart,walmart inc",
	}
}

// loadSectorKeywords loads sector keywords
func loadSectorKeywords() map[string]string {
	// In a real implementation, this would load from a file or database
	// This is a simplified version with a few examples
	return map[string]string{
		"TECHNOLOGY": "tech,technology,software,hardware,semiconductor,ai,artificial intelligence,cloud",
		"HEALTHCARE": "healthcare,health,medical,pharma,pharmaceutical,biotech,biotechnology",
		"FINANCE": "finance,financial,bank,banking,investment,insurance",
		"CONSUMER": "consumer,retail,e-commerce,ecommerce",
		"ENERGY": "energy,oil,gas,renewable,solar,wind",
		"INDUSTRIAL": "industrial,manufacturing,construction,aerospace",
		"UTILITIES": "utilities,utility,electric,water,gas",
		"REAL_ESTATE": "real estate,property,reit",
		"COMMUNICATION": "communication,telecom,media,entertainment",
		"MATERIALS": "materials,chemical,mining,metal",
	}
}

// getStopwords returns a set of common stopwords
func getStopwords() map[string]bool {
	// In a real implementation, this would be more comprehensive
	stopwords := []string{
		"a", "an", "the", "and", "or", "but", "if", "then", "else", "when",
		"at", "by", "for", "with", "about", "against", "between", "into",
		"through", "during", "before", "after", "above", "below", "to", "from",
		"up", "down", "in", "out", "on", "off", "over", "under", "again",
		"further", "then", "once", "here", "there", "when", "where", "why",
		"how", "all", "any", "both", "each", "few", "more", "most", "other",
		"some", "such", "no", "nor", "not", "only", "own", "same", "so",
		"than", "too", "very", "s", "t", "can", "will", "just", "don",
		"should", "now", "d", "ll", "m", "o", "re", "ve", "y", "ain", "aren",
		"couldn", "didn", "doesn", "hadn", "hasn", "haven", "isn", "ma",
		"mightn", "mustn", "needn", "shan", "shouldn", "wasn", "weren", "won",
		"wouldn", "i", "me", "my", "myself", "we", "our", "ours", "ourselves",
		"you", "your", "yours", "yourself", "yourselves", "he", "him", "his",
		"himself", "she", "her", "hers", "herself", "it", "its", "itself",
		"they", "them", "their", "theirs", "themselves", "what", "which",
		"who", "whom", "this", "that", "these", "those", "am", "is", "are",
		"was", "were", "be", "been", "being", "have", "has", "had", "having",
		"do", "does", "did", "doing", "would", "should", "could", "ought",
		"i'm", "you're", "he's", "she's", "it's", "we're", "they're", "i've",
		"you've", "we've", "they've", "i'd", "you'd", "he'd", "she'd", "we'd",
		"they'd", "i'll", "you'll", "he'll", "she'll", "we'll", "they'll",
		"isn't", "aren't", "wasn't", "weren't", "hasn't", "haven't", "hadn't",
		"doesn't", "don't", "didn't", "won't", "wouldn't", "shan't", "shouldn't",
		"can't", "cannot", "couldn't", "mustn't", "let's", "that's", "who's",
		"what's", "here's", "there's", "when's", "where's", "why's", "how's",
	}
	
	stopwordMap := make(map[string]bool)
	for _, word := range stopwords {
		stopwordMap[word] = true
	}
	
	return stopwordMap
}
