package security

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

// MFAManager handles multi-factor authentication
type MFAManager struct {
	config *MFAConfig
}

// MFAConfig defines MFA configuration
type MFAConfig struct {
	RequiredFactors    int
	FactorTypes        []string
	TokenValidity      time.Duration
	MaxRetries         int
	LockoutDuration    time.Duration
	RateLimitPerMinute int
}

// NewMFAManager creates a new MFA manager
func NewMFAManager(config MFAConfig) *MFAManager {
	return &MFAManager{
		config: &config,
	}
}

// VerifyMFA verifies multiple authentication factors
func (m *MFAManager) VerifyMFA(ctx context.Context, factors []MFAFactor) (MFAStatus, error) {
	if len(factors) < m.config.RequiredFactors {
		return MFAStatus{}, fmt.Errorf("insufficient factors provided")
	}

	verifiedFactors := make([]string, 0)
	riskLevel := "low"

	for _, factor := range factors {
		// Check if factor type is allowed
		if !m.isAllowedFactorType(factor.Type) {
			return MFAStatus{}, fmt.Errorf("unsupported factor type: %s", factor.Type)
		}

		// Verify factor
		verified, err := m.verifyFactor(ctx, factor)
		if err != nil {
			return MFAStatus{}, fmt.Errorf("factor verification failed: %v", err)
		}

		if verified {
			verifiedFactors = append(verifiedFactors, factor.Type)
			// Update risk level based on factor type
			riskLevel = m.updateRiskLevel(riskLevel, factor.Type)
		}
	}

	if len(verifiedFactors) < m.config.RequiredFactors {
		return MFAStatus{}, fmt.Errorf("insufficient verified factors")
	}

	return MFAStatus{
		Verified:     true,
		FactorsUsed:  verifiedFactors,
		LastVerified: time.Now(),
		RiskLevel:    riskLevel,
	}, nil
}

// verifyFactor verifies a single authentication factor
func (m *MFAManager) verifyFactor(ctx context.Context, factor MFAFactor) (bool, error) {
	switch factor.Type {
	case "totp":
		return m.verifyTOTP(factor)
	case "sms":
		return m.verifySMS(factor)
	case "email":
		return m.verifyEmail(factor)
	case "hardware":
		return m.verifyHardwareToken(factor)
	case "biometric":
		return m.verifyBiometric(factor)
	default:
		return false, fmt.Errorf("unsupported factor type: %s", factor.Type)
	}
}

// verifyTOTP verifies a Time-based One-Time Password
func (m *MFAManager) verifyTOTP(factor MFAFactor) (bool, error) {
	// Decode secret
	secret, err := base32.StdEncoding.DecodeString(factor.Secret)
	if err != nil {
		return false, fmt.Errorf("invalid secret: %v", err)
	}

	// Verify TOTP
	valid := totp.Validate(factor.Secret, factor.Metadata["code"].(string))
	if !valid {
		return false, fmt.Errorf("invalid TOTP code")
	}

	// Check for replay attacks
	if err := m.checkReplayAttack(factor); err != nil {
		return false, err
	}

	return true, nil
}

// verifySMS verifies an SMS code
func (m *MFAManager) verifySMS(factor MFAFactor) (bool, error) {
	// Implement SMS verification
	// This should include:
	// - Rate limiting
	// - Code expiration
	// - Anti-tampering measures
	return false, nil
}

// verifyEmail verifies an email code
func (m *MFAManager) verifyEmail(factor MFAFactor) (bool, error) {
	// Implement email verification
	// This should include:
	// - Rate limiting
	// - Code expiration
	// - Anti-tampering measures
	return false, nil
}

// verifyHardwareToken verifies a hardware token
func (m *MFAManager) verifyHardwareToken(factor MFAFactor) (bool, error) {
	// Implement hardware token verification
	// This should include:
	// - Token validation
	// - Anti-tampering measures
	// - Physical security checks
	return false, nil
}

// verifyBiometric verifies biometric data
func (m *MFAManager) verifyBiometric(factor MFAFactor) (bool, error) {
	// Implement biometric verification
	// This should include:
	// - Biometric data validation
	// - Anti-spoofing measures
	// - Liveness detection
	return false, nil
}

// checkReplayAttack checks for potential replay attacks
func (m *MFAManager) checkReplayAttack(factor MFAFactor) error {
	// Get last used timestamp
	lastUsed, ok := factor.Metadata["last_used"].(time.Time)
	if !ok {
		return fmt.Errorf("missing last used timestamp")
	}

	// Check if code was used too recently
	if time.Since(lastUsed) < time.Second {
		return fmt.Errorf("potential replay attack detected")
	}

	// Update last used timestamp
	factor.Metadata["last_used"] = time.Now()
	return nil
}

// isAllowedFactorType checks if a factor type is allowed
func (m *MFAManager) isAllowedFactorType(factorType string) bool {
	for _, allowed := range m.config.FactorTypes {
		if allowed == factorType {
			return true
		}
	}
	return false
}

// updateRiskLevel updates the risk level based on factor type
func (m *MFAManager) updateRiskLevel(currentLevel, factorType string) string {
	switch factorType {
	case "hardware":
		return "low"
	case "biometric":
		if currentLevel == "high" {
			return "high"
		}
		return "medium"
	case "totp":
		if currentLevel == "high" {
			return "high"
		}
		return "medium"
	case "sms", "email":
		return "high"
	default:
		return currentLevel
	}
}

// GenerateTOTPSecret generates a new TOTP secret
func (m *MFAManager) GenerateTOTPSecret() (string, error) {
	// Generate random secret
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("failed to generate secret: %v", err)
	}

	// Encode as base32
	return base32.StdEncoding.EncodeToString(secret), nil
}

// GenerateTOTPCode generates a TOTP code for a secret
func (m *MFAManager) GenerateTOTPCode(secret string) (string, error) {
	// Generate TOTP code
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %v", err)
	}

	return code, nil
}

// ValidateTOTPCode validates a TOTP code
func (m *MFAManager) ValidateTOTPCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateSecureCode generates a secure random code
func (m *MFAManager) GenerateSecureCode(length int) (string, error) {
	// Generate random bytes
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate code: %v", err)
	}

	// Convert to numeric code
	code := ""
	for _, b := range bytes {
		code += fmt.Sprintf("%d", b%10)
	}

	return code[:length], nil
}

// HashCode securely hashes a code
func (m *MFAManager) HashCode(code string) string {
	h := hmac.New(sha256.New, []byte(code))
	h.Write([]byte(time.Now().String()))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// CheckRateLimit checks if the rate limit has been exceeded
func (m *MFAManager) CheckRateLimit(factor MFAFactor) error {
	// Get last attempt timestamp
	lastAttempt, ok := factor.Metadata["last_attempt"].(time.Time)
	if !ok {
		return nil
	}

	// Check rate limit
	if time.Since(lastAttempt) < time.Minute/time.Duration(m.config.RateLimitPerMinute) {
		return fmt.Errorf("rate limit exceeded")
	}

	// Update last attempt timestamp
	factor.Metadata["last_attempt"] = time.Now()
	return nil
}
