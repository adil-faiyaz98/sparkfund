package security

import (
	"context"
	"time"
)

// AdaptiveAccessControl handles combined RBAC and ABAC for KYC service
type AdaptiveAccessControl struct {
	config *AccessControlConfig
	store  AccessControlStore
}

// AccessControlConfig defines configuration for adaptive access control
type AccessControlConfig struct {
	DefaultRole        string
	RoleHierarchy      map[string][]string
	PermissionMatrix   map[string][]string
	DynamicRules       []DynamicRule
	EvaluationInterval time.Duration
	RiskThreshold      float64
	MaxFailedAttempts  int
	LockoutDuration    time.Duration
}

// DynamicRule represents a dynamic access control rule
type DynamicRule struct {
	ID         string
	Name       string
	Condition  string
	Action     string
	Priority   int
	Attributes map[string]interface{}
	ExpiresAt  time.Time
}

// AccessControlStore defines the interface for access control storage
type AccessControlStore interface {
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	GetRolePermissions(ctx context.Context, role string) ([]string, error)
	GetUserAttributes(ctx context.Context, userID string) (map[string]interface{}, error)
	GetAccessHistory(ctx context.Context, userID string) ([]AccessEvent, error)
	UpdateAccessHistory(ctx context.Context, event AccessEvent) error
}

// AccessEvent represents an access control event
type AccessEvent struct {
	UserID     string
	Timestamp  time.Time
	Action     string
	Resource   string
	Status     string
	Attributes map[string]interface{}
	RiskScore  float64
	Location   Location
	Device     DeviceInfo
}

// NewAdaptiveAccessControl creates a new adaptive access control instance
func NewAdaptiveAccessControl(config AccessControlConfig, store AccessControlStore) *AdaptiveAccessControl {
	return &AdaptiveAccessControl{
		config: &config,
		store:  store,
	}
}

// CheckAccess determines if a user has access to a resource
func (a *AdaptiveAccessControl) CheckAccess(ctx context.Context, userID string, resource string, action string) (*AccessDecision, error) {
	// Get user's roles
	roles, err := a.store.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// If no roles assigned, use default role
	if len(roles) == 0 {
		roles = []string{a.config.DefaultRole}
	}

	// Get base permissions from roles
	basePermissions, err := a.getBasePermissions(ctx, roles)
	if err != nil {
		return nil, err
	}

	// Get user attributes
	attributes, err := a.store.GetUserAttributes(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get access history
	history, err := a.store.GetAccessHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Apply dynamic rules
	permissions := a.applyDynamicRules(basePermissions, attributes, history)

	// Apply risk-based restrictions
	permissions = a.applyRiskRestrictions(permissions, attributes, history)

	// Apply time-based restrictions
	permissions = a.applyTimeRestrictions(permissions, attributes)

	// Apply location-based restrictions
	permissions = a.applyLocationRestrictions(permissions, attributes)

	// Apply device-based restrictions
	permissions = a.applyDeviceRestrictions(permissions, attributes)

	// Determine access decision
	riskLevel := a.determineRiskLevel(attributes, history)
	decision := &AccessDecision{
		Allowed:     len(permissions) > 0,
		Permissions: permissions,
		RiskLevel:   riskLevel,
		ExpiresAt:   time.Now().Add(a.config.EvaluationInterval),
	}

	// Record access event
	event := AccessEvent{
		UserID:     userID,
		Timestamp:  time.Now(),
		Action:     action,
		Resource:   resource,
		Status:     a.determineAccessStatus(decision),
		Attributes: attributes,
		RiskScore:  riskLevel.Score,
	}
	if err := a.store.UpdateAccessHistory(ctx, event); err != nil {
		return nil, err
	}

	return decision, nil
}

// Helper functions

func (a *AdaptiveAccessControl) getBasePermissions(ctx context.Context, roles []string) ([]string, error) {
	var permissions []string
	for _, role := range roles {
		rolePermissions, err := a.store.GetRolePermissions(ctx, role)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, rolePermissions...)
	}
	return permissions, nil
}

func (a *AdaptiveAccessControl) applyDynamicRules(permissions []string, attributes map[string]interface{}, history []AccessEvent) []string {
	var filteredPermissions []string

	for _, permission := range permissions {
		allowed := true
		for _, rule := range a.config.DynamicRules {
			if !a.evaluateRule(rule, attributes, history) {
				allowed = false
				break
			}
		}
		if allowed {
			filteredPermissions = append(filteredPermissions, permission)
		}
	}

	return filteredPermissions
}

func (a *AdaptiveAccessControl) applyRiskRestrictions(permissions []string, attributes map[string]interface{}, history []AccessEvent) []string {
	var filteredPermissions []string

	riskLevel := a.determineRiskLevel(attributes, history)
	for _, permission := range permissions {
		if a.isPermissionAllowedForRisk(permission, riskLevel) {
			filteredPermissions = append(filteredPermissions, permission)
		}
	}

	return filteredPermissions
}

func (a *AdaptiveAccessControl) applyTimeRestrictions(permissions []string, attributes map[string]interface{}) []string {
	var filteredPermissions []string

	hour := time.Now().Hour()
	for _, permission := range permissions {
		if a.isPermissionAllowedForTime(permission, hour) {
			filteredPermissions = append(filteredPermissions, permission)
		}
	}

	return filteredPermissions
}

func (a *AdaptiveAccessControl) applyLocationRestrictions(permissions []string, attributes map[string]interface{}) []string {
	var filteredPermissions []string

	for _, permission := range permissions {
		if a.isPermissionAllowedForLocation(permission, attributes) {
			filteredPermissions = append(filteredPermissions, permission)
		}
	}

	return filteredPermissions
}

func (a *AdaptiveAccessControl) applyDeviceRestrictions(permissions []string, attributes map[string]interface{}) []string {
	var filteredPermissions []string

	for _, permission := range permissions {
		if a.isPermissionAllowedForDevice(permission, attributes) {
			filteredPermissions = append(filteredPermissions, permission)
		}
	}

	return filteredPermissions
}

func (a *AdaptiveAccessControl) evaluateRule(rule DynamicRule, attributes map[string]interface{}, history []AccessEvent) bool {
	// Implement rule evaluation logic
	// This would evaluate the rule's condition against attributes and history
	return true
}

func (a *AdaptiveAccessControl) determineRiskLevel(attributes map[string]interface{}, history []AccessEvent) RiskLevel {
	// Implement risk level determination
	// This would evaluate various risk factors from attributes and history
	return RiskLevel{
		Level:  "medium",
		Score:  0.5,
		Reason: "Standard risk level",
	}
}

func (a *AdaptiveAccessControl) isPermissionAllowedForRisk(permission string, riskLevel RiskLevel) bool {
	// Implement risk-based permission filtering
	return true
}

func (a *AdaptiveAccessControl) isPermissionAllowedForTime(permission string, hour int) bool {
	// Implement time-based permission filtering
	return true
}

func (a *AdaptiveAccessControl) isPermissionAllowedForLocation(permission string, attributes map[string]interface{}) bool {
	// Implement location-based permission filtering
	return true
}

func (a *AdaptiveAccessControl) isPermissionAllowedForDevice(permission string, attributes map[string]interface{}) bool {
	// Implement device-based permission filtering
	return true
}

func (a *AdaptiveAccessControl) determineAccessStatus(decision *AccessDecision) string {
	if !decision.Allowed {
		return "denied"
	}
	switch decision.RiskLevel.Level {
	case "high":
		return "restricted"
	case "medium":
		return "monitored"
	default:
		return "allowed"
	}
}

// AccessDecision represents the result of an access control check
type AccessDecision struct {
	Allowed     bool
	Permissions []string
	RiskLevel   RiskLevel
	ExpiresAt   time.Time
}
