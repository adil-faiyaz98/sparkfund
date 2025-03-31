# User Service

The User Service is a microservice responsible for managing user accounts, authentication, and authorization in the SparkFund platform.

## Features

- User registration and authentication
- Profile management
- Password reset functionality
- Session management
- Multi-factor authentication support
- Account security features (login attempts, account locking)

## API Endpoints

### User Management

- `POST /api/v1/users` - Register a new user
- `POST /api/v1/users/login` - Authenticate user
- `GET /api/v1/users/{id}` - Get user details
- `PUT /api/v1/users/{id}` - Update user details
- `DELETE /api/v1/users/{id}` - Delete user account

### Profile Management

- `GET /api/v1/users/{id}/profile` - Get user profile
- `PUT /api/v1/users/{id}/profile` - Update user profile

### Password Management

- `PUT /api/v1/users/{id}/password` - Change password
- `POST /api/v1/users/password/reset` - Request password reset
- `POST /api/v1/users/password/reset/confirm` - Confirm password reset

## Configuration

The service can be configured using environment variables:

```env
# Server settings
PORT=8080
READ_TIMEOUT=5s
WRITE_TIMEOUT=10s
IDLE_TIMEOUT=120s

# Database settings
DATABASE_URL=postgres://postgres:postgres@localhost:5432/sparkfund?sslmode=disable

# Security settings
JWT_SECRET=your-secret-key
PASSWORD_HASH_COST=10
SESSION_EXPIRY=24h
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=30m
REQUIRE_MFA=false
ALLOWED_COUNTRIES=US,GB,CA
MIN_PASSWORD_LENGTH=8

# Email settings
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-username
SMTP_PASSWORD=your-password
FROM_EMAIL=noreply@sparkfund.com
```

## Development

### Prerequisites

- Go 1.21 or later
- PostgreSQL 13 or later
- Make (optional, for using Makefile commands)

### Setup

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up environment variables
4. Run the service:
   ```bash
   go run cmd/api/main.go
   ```

### Database Migrations

The service uses SQL migrations for database schema management. To run migrations:

```bash
go run cmd/migrate/main.go
```

## Testing

Run tests:

```bash
go test ./...
```

## Security

The service implements several security features:

- Password hashing using bcrypt
- JWT-based authentication
- Rate limiting for login attempts
- Account locking after failed attempts
- Session management
- Multi-factor authentication support
- Input validation and sanitization
- Secure password reset flow

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 