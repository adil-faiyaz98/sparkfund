package repository

import "errors"

var (
	// User errors
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")

	// Profile errors
	ErrProfileNotFound = errors.New("profile not found")

	// Password reset errors
	ErrResetTokenNotFound = errors.New("reset token not found")
	ErrResetTokenExpired  = errors.New("reset token expired")
	ErrResetTokenUsed     = errors.New("reset token already used")

	// Session errors
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)
