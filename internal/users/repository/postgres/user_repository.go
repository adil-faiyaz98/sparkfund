package postgres

import (
	"fmt"
	"time"

	"github.com/adil-faiyaz98/structgen/internal/users"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) users.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *users.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*users.User, error) {
	var user users.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*users.User, error) {
	var user users.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(user *users.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&users.User{}, "id = ?", id).Error
}

func (r *userRepository) UpdatePassword(id uuid.UUID, passwordHash string) error {
	return r.db.Model(&users.User{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
}

func (r *userRepository) UpdateStatus(id uuid.UUID, status users.UserStatus) error {
	return r.db.Model(&users.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *userRepository) UpdateLastLogin(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&users.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}
