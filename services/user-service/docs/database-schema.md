# User Service Database Schema

This document describes the database schema for the User Service.

## Tables

### users

Stores user account information.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| email | VARCHAR(255) | User's email address (unique) |
| password_hash | VARCHAR(255) | Hashed password |
| first_name | VARCHAR(100) | User's first name |
| last_name | VARCHAR(100) | User's last name |
| phone_number | VARCHAR(20) | User's phone number |
| date_of_birth | DATE | User's date of birth |
| address_line1 | VARCHAR(255) | Address line 1 |
| address_line2 | VARCHAR(255) | Address line 2 |
| city | VARCHAR(100) | City |
| state | VARCHAR(100) | State or province |
| postal_code | VARCHAR(20) | Postal code |
| country | VARCHAR(2) | Country code (ISO 3166-1 alpha-2) |
| status | VARCHAR(20) | Account status (pending, active, suspended, locked) |
| role | VARCHAR(20) | User role (user, admin) |
| email_verified | BOOLEAN | Whether the email is verified |
| phone_verified | BOOLEAN | Whether the phone is verified |
| mfa_enabled | BOOLEAN | Whether MFA is enabled |
| mfa_secret | VARCHAR(255) | MFA secret key |
| login_attempts | INT | Number of failed login attempts |
| locked_until | TIMESTAMP WITH TIME ZONE | When the account will be unlocked |
| last_login | TIMESTAMP WITH TIME ZONE | Last login timestamp |
| password_changed_at | TIMESTAMP WITH TIME ZONE | When the password was last changed |
| created_at | TIMESTAMP WITH TIME ZONE | When the user was created |
| updated_at | TIMESTAMP WITH TIME ZONE | When the user was last updated |
| deleted_at | TIMESTAMP WITH TIME ZONE | When the user was deleted (soft delete) |

**Indexes:**
- `idx_users_email`: Index on `email`
- `idx_users_status`: Index on `status`
- `idx_users_role`: Index on `role`
- `idx_users_created_at`: Index on `created_at`
- `idx_users_deleted_at`: Index on `deleted_at`

### user_sessions

Stores user session information.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users.id |
| token | VARCHAR(255) | Session token (unique) |
| refresh_token | VARCHAR(255) | Refresh token (unique) |
| ip_address | VARCHAR(45) | IP address |
| user_agent | TEXT | User agent string |
| device_info | TEXT | Device information |
| location | TEXT | Location information |
| expires_at | TIMESTAMP WITH TIME ZONE | When the session expires |
| last_activity | TIMESTAMP WITH TIME ZONE | Last activity timestamp |
| created_at | TIMESTAMP WITH TIME ZONE | When the session was created |
| revoked | BOOLEAN | Whether the session is revoked |
| revoked_at | TIMESTAMP WITH TIME ZONE | When the session was revoked |

**Indexes:**
- `idx_user_sessions_user_id`: Index on `user_id`
- `idx_user_sessions_token`: Index on `token`
- `idx_user_sessions_refresh_token`: Index on `refresh_token`
- `idx_user_sessions_expires_at`: Index on `expires_at`
- `idx_user_sessions_revoked`: Index on `revoked`

### password_reset_tokens

Stores password reset tokens.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users.id |
| token | VARCHAR(255) | Reset token (unique) |
| expires_at | TIMESTAMP WITH TIME ZONE | When the token expires |
| created_at | TIMESTAMP WITH TIME ZONE | When the token was created |
| used | BOOLEAN | Whether the token has been used |
| used_at | TIMESTAMP WITH TIME ZONE | When the token was used |

**Indexes:**
- `idx_password_reset_tokens_user_id`: Index on `user_id`
- `idx_password_reset_tokens_token`: Index on `token`
- `idx_password_reset_tokens_expires_at`: Index on `expires_at`

### verification_tokens

Stores verification tokens for email and phone verification.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users.id |
| token | VARCHAR(255) | Verification token (unique) |
| type | VARCHAR(20) | Token type (email, phone) |
| expires_at | TIMESTAMP WITH TIME ZONE | When the token expires |
| created_at | TIMESTAMP WITH TIME ZONE | When the token was created |
| used | BOOLEAN | Whether the token has been used |
| used_at | TIMESTAMP WITH TIME ZONE | When the token was used |

**Indexes:**
- `idx_verification_tokens_user_id`: Index on `user_id`
- `idx_verification_tokens_token`: Index on `token`
- `idx_verification_tokens_type`: Index on `type`
- `idx_verification_tokens_expires_at`: Index on `expires_at`

### user_profiles

Stores user profile information.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users.id |
| profile_picture_url | TEXT | URL to profile picture |
| bio | TEXT | User's bio |
| occupation | VARCHAR(255) | User's occupation |
| company | VARCHAR(255) | User's company |
| website | VARCHAR(255) | User's website |
| social_links | JSONB | Social media links |
| preferences | JSONB | User preferences |
| settings | JSONB | User settings |
| created_at | TIMESTAMP WITH TIME ZONE | When the profile was created |
| updated_at | TIMESTAMP WITH TIME ZONE | When the profile was last updated |

**Indexes:**
- `idx_user_profiles_user_id`: Index on `user_id`

### audit_logs

Stores audit logs for user actions.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Foreign key to users.id |
| action | VARCHAR(50) | Action performed |
| entity_type | VARCHAR(50) | Type of entity affected |
| entity_id | UUID | ID of entity affected |
| ip_address | VARCHAR(45) | IP address |
| user_agent | TEXT | User agent string |
| details | JSONB | Additional details |
| status | VARCHAR(20) | Status of the action |
| created_at | TIMESTAMP WITH TIME ZONE | When the action was performed |

**Indexes:**
- `idx_audit_logs_user_id`: Index on `user_id`
- `idx_audit_logs_action`: Index on `action`
- `idx_audit_logs_entity_type`: Index on `entity_type`
- `idx_audit_logs_entity_id`: Index on `entity_id`
- `idx_audit_logs_created_at`: Index on `created_at`

## Relationships

- `user_sessions.user_id` -> `users.id` (CASCADE DELETE)
- `password_reset_tokens.user_id` -> `users.id` (CASCADE DELETE)
- `verification_tokens.user_id` -> `users.id` (CASCADE DELETE)
- `user_profiles.user_id` -> `users.id` (CASCADE DELETE)
- `audit_logs.user_id` -> `users.id` (SET NULL)

## Triggers

- `update_updated_at_column()`: Updates the `updated_at` column on record update
  - Applied to `users` and `user_profiles` tables

## Extensions

- `uuid-ossp`: For UUID generation
