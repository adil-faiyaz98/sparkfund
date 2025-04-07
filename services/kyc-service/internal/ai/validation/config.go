package validation

import "time"

// Config contains all validation-related configuration
type Config struct {
	// General settings
	MinValidationScore float64
	RetryDelay         time.Duration
	MaxRetries         int

	// Document validation
	MinDocumentScore         float64
	DocumentQualityThreshold float64
	DocumentExpiryThreshold  time.Duration
	AllowedDocumentTypes     []string

	// Biometric validation
	MinBiometricScore           float64
	ManipulationThreshold       float64
	PresentationAttackThreshold float64
	BiometricWeights            BiometricWeights

	// Threat validation
	MinThreatScore             float64
	SyntheticIdentityThreshold float64
	IdentityTheftThreshold     float64
	FraudRiskThreshold         float64
	AMLRiskThreshold           float64
	SanctionsThreshold         float64
	PEPThreshold               float64
	ThreatWeights              ThreatWeights

	// Behavior validation
	MinBehaviorScore           float64
	UnusualPatternThreshold    float64
	VelocityThreshold          float64
	AnomalyThreshold           float64
	FraudPatternThreshold      float64
	CompliancePatternThreshold float64
	BehaviorWeights            BehaviorWeights

	// AI configuration
	AIConfig *AIConfig
}

// BiometricWeights defines weights for biometric validation
type BiometricWeights struct {
	FaceMatch  float64
	Liveness   float64
	FraudCheck float64
}

// ThreatWeights defines weights for threat validation
type ThreatWeights struct {
	Identity   float64
	Financial  float64
	Compliance float64
}

// BehaviorWeights defines weights for behavior validation components
type BehaviorWeights struct {
	Pattern float64 // Weight for pattern analysis
	Anomaly float64 // Weight for anomaly detection
	Risk    float64 // Weight for risk pattern assessment
}

// AIConfig defines the AI-specific configuration settings
type AIConfig struct {
	// Model configuration
	ModelPath      string        // Path to AI models
	BatchSize      int           // Batch size for model inference
	UpdateInterval time.Duration // Interval for model updates
	MinAccuracy    float64       // Minimum required model accuracy

	// Feature configuration
	EnableGPU  bool // Whether to use GPU acceleration
	MaxWorkers int  // Maximum number of concurrent workers
	CacheSize  int  // Size of the model cache

	// Security settings
	EncryptionEnabled bool          // Whether to encrypt model data
	AccessControl     AccessControl // Access control settings
}

// AccessControl defines access control settings for AI operations
type AccessControl struct {
	Enabled           bool
	AllowedRoles      []string
	MaxFailedAttempts int
	LockoutDuration   time.Duration
}
