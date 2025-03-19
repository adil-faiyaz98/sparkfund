package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/yourusername/money-pulse/pkg/discovery"
	apperrors "github.com/yourusername/money-pulse/pkg/errors"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	accessTokenSecret  = getEnv("ACCESS_TOKEN_SECRET", "your-access-secret-key")
	refreshTokenSecret = getEnv("REFRESH_TOKEN_SECRET", "your-refresh-secret-key")
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Initialize service discovery
	consulClient, err := discovery.NewConsulClient(getEnv("CONSUL_ADDR", "consul:8500"))
	if err != nil {
		log.Fatalf("Failed to create consul client: %v", err)
	}

	// Register service with Consul
	serviceHost := getEnv("SERVICE_HOST", "auth-service")
	servicePort := 8085 // Auth service port
	err = consulClient.Register("auth-service", serviceHost, servicePort, []string{"auth", "api"})
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	defer consulClient.Deregister()

	// Routes
	r.Post("/auth/login", handleLogin)
	r.Post("/auth/refresh", handleRefreshToken)
	r.Get("/auth/validate", handleValidateToken)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":8085",
		Handler: r,
	}

	// Start server
	go func() {
		log.Printf("Auth Service starting on port 8085")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperrors.RespondWithError(w, apperrors.NewBadRequestError("Invalid request body", err))
		return
	}

	// In production, fetch user from database and verify password
	// This is a simplified version
	user := User{
		ID:    uuid.New().String(),
		Email: req.Email,
		Role:  "user",
	}

	// Generate token pair
	tokens, err := generateTokenPair(user)
	if err != nil {
		apperrors.RespondWithError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from request
	var tokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&tokenRequest); err != nil {
		apperrors.RespondWithError(w, apperrors.NewBadRequestError("Invalid request body", err))
		return
	}

	// Parse and validate token
	token, err := jwt.Parse(tokenRequest.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		apperrors.RespondWithError(w, apperrors.NewUnauthorizedError("Invalid refresh token"))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		apperrors.RespondWithError(w, apperrors.NewInternalError(nil))
		return
	}

	// Create user from claims
	user := User{
		ID:    claims["sub"].(string),
		Email: claims["email"].(string),
		Role:  claims["role"].(string),
	}

	// Generate new token pair
	tokens, err := generateTokenPair(user)
	if err != nil {
		apperrors.RespondWithError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func handleValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		apperrors.RespondWithError(w, apperrors.NewUnauthorizedError("Invalid authorization header"))
		return
	}

	tokenString := authHeader[7:]

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessTokenSecret), nil
	})

	if err != nil || !token.Valid {
		apperrors.RespondWithError(w, apperrors.NewUnauthorizedError("Invalid token"))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		apperrors.RespondWithError(w, apperrors.NewInternalError(nil))
		return
	}

	// Return user info from token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": claims["sub"],
		"email":   claims["email"],
		"role":    claims["role"],
	})
}

func generateTokenPair(user User) (TokenPair, error) {
	// Create access token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["sub"] = user.ID
	accessClaims["email"] = user.Email
	accessClaims["role"] = user.Role
	accessClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	accessTokenString, err := accessToken.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return TokenPair{}, err
	}

	// Create refresh token (longer lived)
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["sub"] = user.ID
	refreshClaims["email"] = user.Email
	refreshClaims["role"] = user.Role
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // 7 days

	refreshTokenString, err := refreshToken.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    15 * 60, // 15 minutes in seconds
		TokenType:    "Bearer",
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
