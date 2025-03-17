package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates a new PostgreSQL account repository
func NewAccountRepository(connectionURL string) (*AccountRepository, error) {
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize the database schema if needed
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &AccountRepository{db: db}, nil
}

// initSchema creates the necessary tables if they don't exist
func initSchema(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS accounts (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        account_number VARCHAR(50) NOT NULL UNIQUE,
        account_type VARCHAR(20) NOT NULL,
        balance NUMERIC(19,4) NOT NULL DEFAULT 0,
        currency VARCHAR(3) NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
    
    CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
    `

	_, err := db.Exec(query)
	return err
}

// Create adds a new account to the database
func (r *AccountRepository) Create(ctx context.Context, account *model.Account) (*model.Account, error) {
	// Generate a new UUID if not provided
	if account.ID == "" {
		account.ID = uuid.New().String()
	}

	// Generate a unique account number if not provided
	if account.AccountNumber == "" {
		account.AccountNumber = generateAccountNumber()
	}

	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	query := `
    INSERT INTO accounts (id, user_id, account_number, account_type, balance, currency, is_active, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id, user_id, account_number, account_type, balance, currency, is_active, created_at, updated_at
    `

	row := r.db.QueryRowContext(
		ctx,
		query,
		account.ID,
		account.UserID,
		account.AccountNumber,
		account.AccountType,
		account.Balance,
		account.Currency,
		account.IsActive,
		account.CreatedAt,
		account.UpdatedAt,
	)

	var result model.Account
	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.AccountNumber,
		&result.AccountType,
		&result.Balance,
		&result.Currency,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return &result, nil
}

// GetByID retrieves an account by its ID
func (r *AccountRepository) GetByID(ctx context.Context, id string) (*model.Account, error) {
	query := `
    SELECT id, user_id, account_number, account_type, balance, currency, is_active, created_at, updated_at
    FROM accounts
    WHERE id = $1
    `

	row := r.db.QueryRowContext(ctx, query, id)

	var account model.Account
	err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.AccountNumber,
		&account.AccountType,
		&account.Balance,
		&account.Currency,
		&account.IsActive,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

// Update updates an existing account
func (r *AccountRepository) Update(ctx context.Context, account *model.Account) (*model.Account, error) {
	account.UpdatedAt = time.Now()

	query := `
    UPDATE accounts
    SET balance = $1, is_active = $2, updated_at = $3
    WHERE id = $4
    RETURNING id, user_id, account_number, account_type, balance, currency, is_active, created_at, updated_at
    `

	row := r.db.QueryRowContext(
		ctx,
		query,
		account.Balance,
		account.IsActive,
		account.UpdatedAt,
		account.ID,
	)

	var result model.Account
	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.AccountNumber,
		&result.AccountType,
		&result.Balance,
		&result.Currency,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return &result, nil
}

// Delete removes an account from the database
func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// ListByUserID retrieves all accounts for a given user
func (r *AccountRepository) ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.Account, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Count total accounts for pagination
	countQuery := `SELECT COUNT(*) FROM accounts WHERE user_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count accounts: %w", err)
	}

	// Get paginated accounts
	query := `
    SELECT id, user_id, account_number, account_type, balance, currency, is_active, created_at, updated_at
    FROM accounts
    WHERE user_id = $1
    ORDER BY created_at DESC
    LIMIT $2 OFFSET $3
    `

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	accounts := make([]*model.Account, 0)
	for rows.Next() {
		var account model.Account
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.AccountNumber,
			&account.AccountType,
			&account.Balance,
			&account.Currency,
			&account.IsActive,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan account row: %w", err)
		}
		accounts = append(accounts, &account)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating account rows: %w", err)
	}

	return accounts, total, nil
}

// Close closes the database connection
func (r *AccountRepository) Close() error {
	return r.db.Close()
}

// generateAccountNumber creates a unique account number
func generateAccountNumber() string {
	return fmt.Sprintf("ACC-%s", uuid.New().String()[:8])
}
