package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusInactive  UserStatus = "INACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	FirstName string
	LastName  string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, password, firstName, lastName string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}
	if firstName == "" {
		return nil, fmt.Errorf("first name cannot be empty")
	}
	if lastName == "" {
		return nil, fmt.Errorf("last name cannot be empty")
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password, // Note: Password should be hashed before storage
		FirstName: firstName,
		LastName:  lastName,
		Status:    UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) Update(firstName, lastName string) error {
	if firstName == "" {
		return fmt.Errorf("first name cannot be empty")
	}
	if lastName == "" {
		return fmt.Errorf("last name cannot be empty")
	}

	u.FirstName = firstName
	u.LastName = lastName
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) UpdateStatus(status UserStatus) error {
	switch status {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended:
		u.Status = status
		u.UpdatedAt = time.Now()
		return nil
	default:
		return fmt.Errorf("invalid user status: %s", status)
	}
}
