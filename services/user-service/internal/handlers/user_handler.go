package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sparkfund/services/user-service/internal/errors"
	"github.com/sparkfund/services/user-service/internal/logger"
	"github.com/sparkfund/services/user-service/internal/models"
	"github.com/sparkfund/services/user-service/internal/service"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes registers all user routes
func (h *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/users", h.handleRegister).Methods("POST")
	router.HandleFunc("/api/v1/users/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/api/v1/users/{id}", h.handleGetUser).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}", h.handleUpdateUser).Methods("PUT")
	router.HandleFunc("/api/v1/users/{id}/profile", h.handleGetProfile).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}/profile", h.handleUpdateProfile).Methods("PUT")
	router.HandleFunc("/api/v1/users/{id}/password", h.handleChangePassword).Methods("PUT")
	router.HandleFunc("/api/v1/users/password/reset", h.handleResetPassword).Methods("POST")
	router.HandleFunc("/api/v1/users/password/reset/confirm", h.handleConfirmResetPassword).Methods("POST")

	// New security routes
	router.HandleFunc("/api/v1/users/{id}/sessions", h.handleListSessions).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}/sessions/{sessionId}", h.handleRevokeSession).Methods("DELETE")
	router.HandleFunc("/api/v1/users/{id}/sessions", h.handleRevokeAllSessions).Methods("DELETE")
	router.HandleFunc("/api/v1/users/{id}/mfa/enable", h.handleEnableMFA).Methods("POST")
	router.HandleFunc("/api/v1/users/{id}/mfa/disable", h.handleDisableMFA).Methods("POST")
	router.HandleFunc("/api/v1/users/{id}/mfa/verify", h.handleVerifyMFA).Methods("POST")
	router.HandleFunc("/api/v1/users/{id}/security/audit", h.handleGetSecurityAudit).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}/security/activity", h.handleGetSecurityActivity).Methods("GET")
}

// handleRegister handles user registration
func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.RegisterUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
		"id":      user.ID.String(),
	})
}

// handleLogin handles user login
func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.AuthenticateUser(r.Context(), credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Create session
	session := &models.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     uuid.New().String(),
		IP:        r.RemoteAddr,
		UserAgent: r.UserAgent(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		LastUsed:  time.Now(),
	}

	if err := h.userService.CreateSession(r.Context(), session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": session.Token,
		"user":  user,
	})
}

// handleGetUser handles getting user details
func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// handleUpdateUser handles updating user details
func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user.ID = userID
	if err := h.userService.UpdateUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
	})
}

// handleGetProfile handles getting user profile
func (h *UserHandler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := h.userService.GetUserProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// handleUpdateProfile handles updating user profile
func (h *UserHandler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var profile models.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	profile.UserID = userID
	if err := h.userService.UpdateUserProfile(r.Context(), userID, &profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile updated successfully",
	})
}

// handleChangePassword handles changing user password
func (h *UserHandler) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var passwords struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&passwords); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.ChangePassword(r.Context(), userID, passwords.OldPassword, passwords.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password changed successfully",
	})
}

// handleResetPassword handles initiating password reset
func (h *UserHandler) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.ResetPassword(r.Context(), request.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset email sent successfully",
	})
}

// handleConfirmResetPassword handles confirming password reset
func (h *UserHandler) handleConfirmResetPassword(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.ConfirmResetPassword(r.Context(), request.Token, request.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successfully",
	})
}

// handleListSessions lists all active sessions for a user
func (h *UserHandler) handleListSessions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid user ID"))
		return
	}

	sessions, err := h.userService.ListUserSessions(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sessions)
}

// handleRevokeSession revokes a specific session
func (h *UserHandler) handleRevokeSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid user ID"))
		return
	}

	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid session ID"))
		return
	}

	if err := h.userService.RevokeSession(r.Context(), userID, sessionID); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Session revoked successfully",
	})
}

// handleRevokeAllSessions revokes all sessions for a user
func (h *UserHandler) handleRevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid user ID"))
		return
	}

	if err := h.userService.RevokeAllSessions(r.Context(), userID); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "All sessions revoked successfully",
	})
}

// handleEnableMFA enables MFA for a user
func (h *UserHandler) handleEnableMFA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid user ID"))
		return
	}

	secret, err := h.userService.EnableMFA(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"secret": secret,
	})
}

// handleDisableMFA disables MFA for a user
func (h *UserHandler) handleDisableMFA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid user ID"))
		return
	}

	if err := h.userService.DisableMFA(r.Context(), userID); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "MFA disabled successfully",
	})
}

// handleVerifyMFA verifies MFA code
func (h *UserHandler) handleVerifyMFA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid user ID"))
		return
	}

	var request struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid request body"))
		return
	}

	if err := h.userService.VerifyMFA(r.Context(), userID, request.Code); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "MFA code verified successfully",
	})
}

// handleGetSecurityAudit gets security audit logs for a user
func (h *UserHandler) handleGetSecurityAudit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid user ID"))
		return
	}

	auditLogs, err := h.userService.GetSecurityAudit(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(auditLogs)
}

// handleGetSecurityActivity gets recent security activity for a user
func (h *UserHandler) handleGetSecurityActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		h.handleError(w, errors.Wrap(err, "Invalid user ID"))
		return
	}

	activity, err := h.userService.GetSecurityActivity(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activity)
}

// handleError handles error responses
func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var statusCode int
	var message string

	if e, ok := err.(*errors.Error); ok {
		statusCode = e.Code
		message = e.Message
	} else {
		statusCode = http.StatusInternalServerError
		message = "Internal server error"
	}

	// Log the error
	logger.Error(err, "Request failed", map[string]interface{}{
		"status_code": statusCode,
		"message":     message,
	})

	// Send error response
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
