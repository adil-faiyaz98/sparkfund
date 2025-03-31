package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/services/user-service/internal/models"
	"github.com/sparkfund/services/user-service/internal/repository"
)

// UserRepository implements repository.UserRepository interface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create implements repository.UserRepository.Create
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			id, email, password, first_name, last_name, phone_number,
			country, status, is_locked, failed_attempts, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Country,
		user.Status,
		user.IsLocked,
		user.FailedAttempts,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

// Get implements repository.UserRepository.Get
func (r *UserRepository) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone_number,
			country, status, is_locked, failed_attempts, last_login_at,
			created_at, updated_at
		FROM users WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Country,
		&user.Status,
		&user.IsLocked,
		&user.FailedAttempts,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrUserNotFound
	}
	return user, err
}

// GetByEmail implements repository.UserRepository.GetByEmail
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone_number,
			country, status, is_locked, failed_attempts, last_login_at,
			created_at, updated_at
		FROM users WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Country,
		&user.Status,
		&user.IsLocked,
		&user.FailedAttempts,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrUserNotFound
	}
	return user, err
}

// Update implements repository.UserRepository.Update
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET
			email = $1,
			first_name = $2,
			last_name = $3,
			phone_number = $4,
			country = $5,
			status = $6,
			is_locked = $7,
			failed_attempts = $8,
			last_login_at = $9,
			updated_at = $10
		WHERE id = $11
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Country,
		user.Status,
		user.IsLocked,
		user.FailedAttempts,
		user.LastLoginAt,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}

// UpdateStatus implements repository.UserRepository.UpdateStatus
func (r *UserRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE users SET status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}

// UpdatePassword implements repository.UserRepository.UpdatePassword
func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, hashedPassword, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}

// GetProfile implements repository.UserRepository.GetProfile
func (r *UserRepository) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	query := `
		SELECT user_id, address, city, state, postal_code,
			date_of_birth, occupation, income, updated_at
		FROM user_profiles WHERE user_id = $1
	`

	profile := &models.UserProfile{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.Address,
		&profile.City,
		&profile.State,
		&profile.PostalCode,
		&profile.DateOfBirth,
		&profile.Occupation,
		&profile.Income,
		&profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrProfileNotFound
	}
	return profile, err
}

// UpdateProfile implements repository.UserRepository.UpdateProfile
func (r *UserRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, profile *models.UserProfile) error {
	query := `
		INSERT INTO user_profiles (
			user_id, address, city, state, postal_code,
			date_of_birth, occupation, income, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
			address = EXCLUDED.address,
			city = EXCLUDED.city,
			state = EXCLUDED.state,
			postal_code = EXCLUDED.postal_code,
			date_of_birth = EXCLUDED.date_of_birth,
			occupation = EXCLUDED.occupation,
			income = EXCLUDED.income,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		userID,
		profile.Address,
		profile.City,
		profile.State,
		profile.PostalCode,
		profile.DateOfBirth,
		profile.Occupation,
		profile.Income,
		profile.UpdatedAt,
	)
	return err
}

// StoreResetToken implements repository.UserRepository.StoreResetToken
func (r *UserRepository) StoreResetToken(ctx context.Context, email string, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at, created_at)
		SELECT id, $1, $2, $3 FROM users WHERE email = $4
	`

	result, err := r.db.ExecContext(ctx, query, token, expiresAt, time.Now(), email)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}

// GetResetToken implements repository.UserRepository.GetResetToken
func (r *UserRepository) GetResetToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	query := `
		SELECT user_id, token, expires_at, created_at, used
		FROM password_resets WHERE token = $1
	`

	reset := &models.PasswordReset{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&reset.UserID,
		&reset.Token,
		&reset.ExpiresAt,
		&reset.CreatedAt,
		&reset.Used,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrResetTokenNotFound
	}
	return reset, err
}

// MarkResetTokenUsed implements repository.UserRepository.MarkResetTokenUsed
func (r *UserRepository) MarkResetTokenUsed(ctx context.Context, token string) error {
	query := `UPDATE password_resets SET used = true WHERE token = $1`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrResetTokenNotFound
	}
	return nil
}

// CreateSession implements repository.UserRepository.CreateSession
func (r *UserRepository) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (
			id, user_id, token, ip, user_agent,
			created_at, expires_at, last_used
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.Token,
		session.IP,
		session.UserAgent,
		session.CreatedAt,
		session.ExpiresAt,
		session.LastUsed,
	)
	return err
}

// GetSession implements repository.UserRepository.GetSession
func (r *UserRepository) GetSession(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	query := `
		SELECT id, user_id, token, ip, user_agent,
			created_at, expires_at, last_used
		FROM sessions WHERE id = $1
	`

	session := &models.Session{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.IP,
		&session.UserAgent,
		&session.CreatedAt,
		&session.ExpiresAt,
		&session.LastUsed,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrSessionNotFound
	}
	return session, err
}

// UpdateSession implements repository.UserRepository.UpdateSession
func (r *UserRepository) UpdateSession(ctx context.Context, session *models.Session) error {
	query := `
		UPDATE sessions SET
			last_used = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), session.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrSessionNotFound
	}
	return nil
}

// DeleteSession implements repository.UserRepository.DeleteSession
func (r *UserRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrSessionNotFound
	}
	return nil
}

// DeleteExpiredSessions implements repository.UserRepository.DeleteExpiredSessions
func (r *UserRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < $1`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	return err
}

// IncrementFailedAttempts implements repository.UserRepository.IncrementFailedAttempts
func (r *UserRepository) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users 
		SET failed_attempts = failed_attempts + 1,
			is_locked = CASE 
				WHEN failed_attempts + 1 >= $1 THEN true 
				ELSE is_locked 
			END,
			updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, 5, time.Now(), userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}

// ResetFailedAttempts implements repository.UserRepository.ResetFailedAttempts
func (r *UserRepository) ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users 
		SET failed_attempts = 0,
			is_locked = false,
			updated_at = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}

// GetFailedAttempts implements repository.UserRepository.GetFailedAttempts
func (r *UserRepository) GetFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT failed_attempts FROM users WHERE id = $1`

	var attempts int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&attempts)
	if err == sql.ErrNoRows {
		return 0, repository.ErrUserNotFound
	}
	return attempts, err
}

// Delete implements repository.UserRepository.Delete
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrUserNotFound
	}
	return nil
}
