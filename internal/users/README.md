# Money Pulse Users Service

This is the users service for the Money Pulse application, responsible for managing user accounts, authentication, and authorization.

## Features

- User account management (create, read, update, delete)
- Email and phone number verification
- Password management
- User status management (active, inactive, suspended, deleted)
- Role-based access control (admin, user, manager, support)
- Last login tracking
- Soft delete support

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
All endpoints require authentication using JWT tokens in the Authorization header:
```
Authorization: Bearer <token>
```

### Endpoints

#### Create User
```http
POST /users
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "securepassword",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+1234567890"
}
```

#### Get User by ID
```http
GET /users/:id
```

#### Get User by Email
```http
GET /users/email/:email
```

#### Update User
```http
PUT /users/:id
Content-Type: application/json

{
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+1234567890"
}
```

#### Delete User
```http
DELETE /users/:id
```

#### Update Password
```http
PUT /users/:id/password
Content-Type: application/json

{
    "password": "newpassword"
}
```

#### Update Status
```http
PUT /users/:id/status
Content-Type: application/json

{
    "status": "ACTIVE"
}
```

#### Verify Email
```http
PUT /users/:id/verify-email
```

#### Verify Phone
```http
PUT /users/:id/verify-phone
```

## Development

### Prerequisites

- Go 1.16 or later
- PostgreSQL 12 or later
- Make (optional, for using Makefile commands)

### Setup

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up environment variables:
   ```bash
   export DATABASE_URL="postgres://username:password@localhost:5432/moneypulse?sslmode=disable"
   ```
4. Run migrations:
   ```bash
   go run cmd/migrate/main.go up
   ```

### Running the Service

```bash
go run cmd/main.go
```

### Testing

```bash
go test ./...
```

## Deployment

The service can be deployed using Helm:

```bash
helm install users ./helm/users
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 