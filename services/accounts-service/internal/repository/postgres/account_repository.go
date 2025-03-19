package postgres

import (
	"fmt"
	"time"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AccountEntity struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Name          string    `gorm:"not null"`
	Type          string    `gorm:"not null"`
	Balance       float64   `gorm:"type:decimal(10,2);not null"`
	Currency      string    `gorm:"not null"`
	AccountNumber string    `gorm:"unique;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(dataSourceURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&AccountEntity{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return &Adapter{db: db}, nil
}

// Convert domain model to entity
func toEntity(account *domain.Account) *AccountEntity {
	return &AccountEntity{
		ID:            account.ID,
		UserID:        account.UserID,
		Name:          account.Name,
		Type:          string(account.Type),
		Balance:       account.Balance,
		Currency:      account.Currency,
		AccountNumber: account.AccountNumber,
		CreatedAt:     account.CreatedAt,
		UpdatedAt:     account.UpdatedAt,
	}
}

// Convert entity to domain model
func toDomain(entity *AccountEntity) *domain.Account {
	return &domain.Account{
		ID:            entity.ID,
		UserID:        entity.UserID,
		Name:          entity.Name,
		Type:          domain.AccountType(entity.Type),
		Balance:       entity.Balance,
		Currency:      entity.Currency,
		AccountNumber: entity.AccountNumber,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

// Create saves a new account
func (a *Adapter) Create(account *domain.Account) error {
	entity := toEntity(account)
	result := a.db.Create(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to create account: %v", result.Error)
	}
	return nil
}

// GetByID retrieves an account by ID
func (a *Adapter) GetByID(id uuid.UUID) (*domain.Account, error) {
	var entity AccountEntity
	result := a.db.First(&entity, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get account: %v", result.Error)
	}
	return toDomain(&entity), nil
}

// GetByUserID retrieves all accounts for a user
func (a *Adapter) GetByUserID(userID uuid.UUID) ([]*domain.Account, error) {
	var entities []AccountEntity
	result := a.db.Where("user_id = ?", userID).Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user accounts: %v", result.Error)
	}

	accounts := make([]*domain.Account, len(entities))
	for i, entity := range entities {
		accounts[i] = toDomain(&entity)
	}
	return accounts, nil
}

// Update updates an existing account
func (a *Adapter) Update(account *domain.Account) error {
	entity := toEntity(account)
	result := a.db.Save(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to update account: %v", result.Error)
	}
	return nil
}

// Delete removes an account
func (a *Adapter) Delete(id uuid.UUID) error {
	result := a.db.Delete(&AccountEntity{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete account: %v", result.Error)
	}
	return nil
}

// GetByAccountNumber retrieves an account by account number
func (a *Adapter) GetByAccountNumber(accountNumber string) (*domain.Account, error) {
	var entity AccountEntity
	result := a.db.First(&entity, "account_number = ?", accountNumber)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get account: %v", result.Error)
	}
	return toDomain(&entity), nil
}
