package ai

import (
	"context"
	"math"
	"time"

	"github.com/sparkfund/security-monitoring/internal/models"
)

// AnomalyDetector handles AI-based anomaly detection
type AnomalyDetector struct {
	config *AnomalyConfig
	model  AnomalyModel
}

// AnomalyConfig defines configuration for anomaly detection
type AnomalyConfig struct {
	ModelPath        string
	BatchSize        int
	UpdateInterval   time.Duration
	Threshold        float64
	Features         []string
	TrainingWindow   time.Duration
	MinSamples       int
	MaxAnomalies     int
	RetrainingPeriod time.Duration
}

// AnomalyModel defines the interface for anomaly detection models
type AnomalyModel interface {
	Train(ctx context.Context, data []models.Transaction) error
	Predict(ctx context.Context, tx *models.Transaction) (float64, error)
	Update(ctx context.Context, data []models.Transaction) error
	Save(path string) error
	Load(path string) error
}

// NewAnomalyDetector creates a new anomaly detector
func NewAnomalyDetector(config AnomalyConfig) (*AnomalyDetector, error) {
	model, err := loadModel(config.ModelPath)
	if err != nil {
		return nil, err
	}

	return &AnomalyDetector{
		config: &config,
		model:  model,
	}, nil
}

// DetectAnomalies analyzes transactions for anomalies
func (d *AnomalyDetector) DetectAnomalies(ctx context.Context, tx *models.Transaction) (*AnomalyResult, error) {
	// Extract features from transaction
	features := d.extractFeatures(tx)

	// Get anomaly score from model
	score, err := d.model.Predict(ctx, tx)
	if err != nil {
		return nil, err
	}

	// Determine if transaction is anomalous
	isAnomaly := score > d.config.Threshold

	// Generate explanation if anomalous
	var explanation string
	if isAnomaly {
		explanation = d.generateExplanation(tx, score)
	}

	return &AnomalyResult{
		IsAnomaly:   isAnomaly,
		Score:       score,
		Explanation: explanation,
		Features:    features,
		Timestamp:   time.Now(),
	}, nil
}

// TrainModel trains the anomaly detection model
func (d *AnomalyDetector) TrainModel(ctx context.Context, data []models.Transaction) error {
	// Validate training data
	if len(data) < d.config.MinSamples {
		return ErrInsufficientData
	}

	// Train model
	if err := d.model.Train(ctx, data); err != nil {
		return err
	}

	// Save updated model
	return d.model.Save(d.config.ModelPath)
}

// UpdateModel updates the model with new data
func (d *AnomalyDetector) UpdateModel(ctx context.Context, data []models.Transaction) error {
	return d.model.Update(ctx, data)
}

// Helper functions

func (d *AnomalyDetector) extractFeatures(tx *models.Transaction) map[string]float64 {
	features := make(map[string]float64)

	// Extract time-based features
	hour := float64(tx.Timestamp.Hour())
	dayOfWeek := float64(tx.Timestamp.Weekday())
	features["hour"] = hour
	features["day_of_week"] = dayOfWeek

	// Extract amount-based features
	features["amount"] = tx.Amount
	features["amount_log"] = log10(tx.Amount)

	// Extract location-based features
	features["location_risk"] = d.calculateLocationRisk(tx.Location)
	features["is_vpn"] = boolToFloat(tx.Location.IsVPN)
	features["is_proxy"] = boolToFloat(tx.Location.IsProxy)

	// Extract device-based features
	features["device_risk"] = d.calculateDeviceRisk(tx.Device)
	features["is_known_device"] = boolToFloat(tx.Device.IsKnownDevice)

	// Extract recipient-based features
	features["recipient_risk"] = d.calculateRecipientRisk(tx.RecipientID)

	return features
}

func (d *AnomalyDetector) calculateLocationRisk(loc models.Location) float64 {
	// Implement location risk calculation
	// This could consider:
	// - Country risk score
	// - Distance from usual locations
	// - VPN/proxy usage
	// - Time zone consistency
	return 0.0
}

func (d *AnomalyDetector) calculateDeviceRisk(device models.DeviceInfo) float64 {
	// Implement device risk calculation
	// This could consider:
	// - Device type
	// - OS/browser combination
	// - Known device status
	// - First/last seen timing
	return 0.0
}

func (d *AnomalyDetector) calculateRecipientRisk(recipientID string) float64 {
	// Implement recipient risk calculation
	// This could consider:
	// - Transaction history
	// - Relationship with user
	// - Recipient reputation
	return 0.0
}

func (d *AnomalyDetector) generateExplanation(tx *models.Transaction, score float64) string {
	// Generate human-readable explanation of why the transaction is anomalous
	// This could include:
	// - Unusual amount
	// - New location
	// - New device
	// - Unusual time
	// - New recipient
	return "Transaction shows unusual patterns"
}

// AnomalyResult represents the result of anomaly detection
type AnomalyResult struct {
	IsAnomaly   bool
	Score       float64
	Explanation string
	Features    map[string]float64
	Timestamp   time.Time
}

// Helper functions for feature extraction
func log10(x float64) float64 {
	return math.Log10(x)
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
