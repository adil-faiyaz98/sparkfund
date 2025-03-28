package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"bufio"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password hashing and verification
type PasswordService struct {
	// Using both Argon2 and bcrypt for maximum security
	argon2Params *argon2Params
	dictionary   map[string]bool
}

type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewPasswordService() *PasswordService {
	ps := &PasswordService{
		argon2Params: &argon2Params{
			memory:      64 * 1024, // 64 MB
			iterations:  3,
			parallelism: 4,
			saltLength:  16,
			keyLength:   32,
		},
		dictionary: make(map[string]bool),
	}
	ps.loadDictionary()
	return ps
}

// loadDictionary loads a dictionary of common passwords
func (s *PasswordService) loadDictionary() {
	file, err := os.Open(filepath.Join("internal", "dictionary", "common_passwords.txt"))
	if err != nil {
		fmt.Printf("Warning: Could not load dictionary: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s.dictionary[strings.ToLower(scanner.Text())] = true
	}
}

// HashPassword creates a secure hash of the password using both Argon2 and bcrypt
func (s *PasswordService) HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, s.argon2Params.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash using Argon2
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		s.argon2Params.iterations,
		s.argon2Params.memory,
		s.argon2Params.parallelism,
		s.argon2Params.keyLength,
	)

	// Encode the salt and hash
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	// Create the Argon2 hash string
	argon2Hash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		s.argon2Params.memory,
		s.argon2Params.iterations,
		s.argon2Params.parallelism,
		encodedSalt,
		encodedHash,
	)

	// Hash using bcrypt as an additional layer
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(argon2Hash), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate bcrypt hash: %w", err)
	}

	return string(bcryptHash), nil
}

// VerifyPassword verifies a password against its hash
func (s *PasswordService) VerifyPassword(password, hash string) (bool, error) {
	// First verify the bcrypt hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, nil
	}

	// Parse the Argon2 hash string
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	// Decode the salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Hash the password with the same parameters
	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		s.argon2Params.iterations,
		s.argon2Params.memory,
		s.argon2Params.parallelism,
		s.argon2Params.keyLength,
	)

	// Compare hashes in constant time
	return subtle.ConstantTimeCompare(expectedHash, actualHash) == 1, nil
}

// GeneratePassword generates a cryptographically secure random password
func (s *PasswordService) GeneratePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?"
	password := make([]byte, length)
	
	for i := range password {
		n, err := rand.Int(rand.Reader, int64(len(charset)))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = charset[n.Int64()]
	}

	return string(password), nil
}

// ValidatePasswordStrength checks if a password meets security requirements
func (s *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var (
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	// Check if the password is in the dictionary
	if s.dictionary[strings.ToLower(password)] {
		return fmt.Errorf("password is too common")
	}

	return nil
}

// RotatePasswordSalt generates a new salt for an existing password hash
func (s *PasswordService) RotatePasswordSalt(password, oldHash string) (string, error) {
	// Verify the old password first
	valid, err := s.VerifyPassword(password, oldHash)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", fmt.Errorf("invalid password")
	}

	// Generate new hash with new salt
	return s.HashPassword(password)
} 