package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	ClientID      string `mapstructure:"client_id"`
	ClientSecret  string `mapstructure:"client_secret"`
	IssuerURL     string `mapstructure:"issuer_url"`
	RedirectURL   string `mapstructure:"redirect_url"`
	Scopes        string `mapstructure:"scopes"`
	CookieName    string `mapstructure:"cookie_name"`
	CookieSecure  bool   `mapstructure:"cookie_secure"`
	CookieMaxAge  int    `mapstructure:"cookie_max_age"`
	CookieDomain  string `mapstructure:"cookie_domain"`
	CookiePath    string `mapstructure:"cookie_path"`
	CookieHTTPOnly bool   `mapstructure:"cookie_http_only"`
}

// OAuthManager manages OAuth authentication
type OAuthManager struct {
	config       OAuthConfig
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
	logger       *logrus.Logger
	jwksURL      string
	mu           sync.RWMutex
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager(config OAuthConfig, logger *logrus.Logger) (*OAuthManager, error) {
	if !config.Enabled {
		return &OAuthManager{
			config: config,
			logger: logger,
		}, nil
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, config.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// Get JWKS URL for token verification
	var providerJSON map[string]interface{}
	if err := provider.Claims(&providerJSON); err != nil {
		return nil, fmt.Errorf("failed to parse provider claims: %w", err)
	}

	jwksURL, ok := providerJSON["jwks_uri"].(string)
	if !ok {
		return nil, errors.New("failed to get JWKS URL from provider")
	}

	// Configure OAuth2
	scopes := strings.Split(config.Scopes, ",")
	if len(scopes) == 0 || (len(scopes) == 1 && scopes[0] == "") {
		scopes = []string{oidc.ScopeOpenID, "profile", "email"}
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	return &OAuthManager{
		config:       config,
		provider:     provider,
		verifier:     verifier,
		oauth2Config: oauth2Config,
		logger:       logger,
		jwksURL:      jwksURL,
	}, nil
}

// IsEnabled returns whether OAuth is enabled
func (m *OAuthManager) IsEnabled() bool {
	return m.config.Enabled
}

// GetAuthURL returns the OAuth authorization URL
func (m *OAuthManager) GetAuthURL(state string) string {
	if !m.config.Enabled {
		return ""
	}
	return m.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange exchanges an authorization code for tokens
func (m *OAuthManager) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	if !m.config.Enabled {
		return nil, errors.New("OAuth is not enabled")
	}
	return m.oauth2Config.Exchange(ctx, code)
}

// VerifyIDToken verifies an ID token
func (m *OAuthManager) VerifyIDToken(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	if !m.config.Enabled {
		return nil, errors.New("OAuth is not enabled")
	}
	return m.verifier.Verify(ctx, rawIDToken)
}

// Middleware returns a middleware that handles OAuth authentication
func (m *OAuthManager) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.config.Enabled {
			c.Next()
			return
		}

		// Skip authentication for certain paths
		path := c.Request.URL.Path
		if path == "/auth/login" || path == "/auth/callback" || path == "/health" || strings.HasPrefix(path, "/metrics") {
			c.Next()
			return
		}

		// Check for token in cookie
		cookie, err := c.Cookie(m.config.CookieName)
		if err != nil {
			c.Redirect(http.StatusFound, "/auth/login?redirect="+c.Request.URL.String())
			c.Abort()
			return
		}

		// Verify token
		ctx := c.Request.Context()
		idToken, err := m.verifier.Verify(ctx, cookie)
		if err != nil {
			m.logger.Warnf("Invalid ID token: %v", err)
			c.SetCookie(m.config.CookieName, "", -1, m.config.CookiePath, m.config.CookieDomain, m.config.CookieSecure, m.config.CookieHTTPOnly)
			c.Redirect(http.StatusFound, "/auth/login?redirect="+c.Request.URL.String())
			c.Abort()
			return
		}

		// Extract claims
		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			m.logger.Errorf("Failed to parse claims: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token claims"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user", claims)
		c.Set("token", cookie)

		c.Next()
	}
}

// SetupRoutes sets up OAuth routes
func (m *OAuthManager) SetupRoutes(router *gin.Engine) {
	if !m.config.Enabled {
		return
	}

	auth := router.Group("/auth")
	{
		auth.GET("/login", m.handleLogin)
		auth.GET("/callback", m.handleCallback)
		auth.GET("/logout", m.handleLogout)
	}
}

// handleLogin handles the login request
func (m *OAuthManager) handleLogin(c *gin.Context) {
	state := generateRandomState()
	c.SetCookie("oauth_state", state, 600, "/", m.config.CookieDomain, m.config.CookieSecure, true)
	
	redirect := c.Query("redirect")
	if redirect != "" {
		c.SetCookie("oauth_redirect", redirect, 600, "/", m.config.CookieDomain, m.config.CookieSecure, true)
	}

	authURL := m.GetAuthURL(state)
	c.Redirect(http.StatusFound, authURL)
}

// handleCallback handles the OAuth callback
func (m *OAuthManager) handleCallback(c *gin.Context) {
	// Verify state
	state, err := c.Cookie("oauth_state")
	if err != nil || state != c.Query("state") {
		m.logger.Warn("Invalid OAuth state")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	// Clear state cookie
	c.SetCookie("oauth_state", "", -1, "/", m.config.CookieDomain, m.config.CookieSecure, true)

	// Exchange code for token
	code := c.Query("code")
	token, err := m.Exchange(c.Request.Context(), code)
	if err != nil {
		m.logger.Errorf("Failed to exchange code for token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	// Get ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		m.logger.Error("No ID token in OAuth response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No ID token in OAuth response"})
		return
	}

	// Verify ID token
	idToken, err := m.VerifyIDToken(c.Request.Context(), rawIDToken)
	if err != nil {
		m.logger.Errorf("Failed to verify ID token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID token"})
		return
	}

	// Set ID token in cookie
	c.SetCookie(m.config.CookieName, rawIDToken, m.config.CookieMaxAge, m.config.CookiePath, m.config.CookieDomain, m.config.CookieSecure, m.config.CookieHTTPOnly)

	// Get redirect URL
	redirect, _ := c.Cookie("oauth_redirect")
	c.SetCookie("oauth_redirect", "", -1, "/", m.config.CookieDomain, m.config.CookieSecure, true)

	if redirect == "" {
		redirect = "/"
	}

	// Extract user info for logging
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err == nil {
		if email, ok := claims["email"].(string); ok {
			m.logger.Infof("User %s logged in", email)
		}
	}

	c.Redirect(http.StatusFound, redirect)
}

// handleLogout handles the logout request
func (m *OAuthManager) handleLogout(c *gin.Context) {
	// Clear ID token cookie
	c.SetCookie(m.config.CookieName, "", -1, m.config.CookiePath, m.config.CookieDomain, m.config.CookieSecure, m.config.CookieHTTPOnly)

	// Redirect to home
	c.Redirect(http.StatusFound, "/")
}

// generateRandomState generates a random state for OAuth
func generateRandomState() string {
	// In a real implementation, use a secure random generator
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
