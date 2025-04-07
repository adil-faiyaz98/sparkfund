# User Service API Documentation

This document describes the API endpoints provided by the User Service.

## Base URL

```
https://api.sparkfund.com/api/v1/users
```

## Authentication

All API requests require authentication using a JWT token. The token should be included in the `Authorization` header:

```
Authorization: Bearer <token>
```

## API Endpoints

### User Management

#### Create User

```
POST /api/v1/users
```

Creates a new user.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securePassword123!",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+1234567890",
  "date_of_birth": "1990-01-01",
  "address_line1": "123 Main St",
  "address_line2": "Apt 4B",
  "city": "New York",
  "state": "NY",
  "postal_code": "10001",
  "country": "US"
}
```

**Response (201 Created):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "status": "pending",
  "created_at": "2023-01-01T12:00:00Z"
}
```

#### Get User

```
GET /api/v1/users/{id}
```

Retrieves a user by ID.

**Response (200 OK):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+1234567890",
  "date_of_birth": "1990-01-01",
  "address_line1": "123 Main St",
  "address_line2": "Apt 4B",
  "city": "New York",
  "state": "NY",
  "postal_code": "10001",
  "country": "US",
  "status": "active",
  "email_verified": true,
  "phone_verified": true,
  "mfa_enabled": false,
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

#### Update User

```
PUT /api/v1/users/{id}
```

Updates a user.

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Smith",
  "phone_number": "+1987654321",
  "address_line1": "456 Park Ave",
  "address_line2": "Suite 789",
  "city": "New York",
  "state": "NY",
  "postal_code": "10022",
  "country": "US"
}
```

**Response (200 OK):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Smith",
  "phone_number": "+1987654321",
  "date_of_birth": "1990-01-01",
  "address_line1": "456 Park Ave",
  "address_line2": "Suite 789",
  "city": "New York",
  "state": "NY",
  "postal_code": "10022",
  "country": "US",
  "status": "active",
  "email_verified": true,
  "phone_verified": true,
  "mfa_enabled": false,
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-02T12:00:00Z"
}
```

#### Delete User

```
DELETE /api/v1/users/{id}
```

Deletes a user.

**Response (204 No Content)**

#### List Users

```
GET /api/v1/users
```

Lists all users with pagination.

**Query Parameters:**

- `page`: Page number (default: 1)
- `limit`: Number of items per page (default: 10)
- `status`: Filter by status (optional)
- `search`: Search term for email, first name, or last name (optional)

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user1@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "status": "active",
      "created_at": "2023-01-01T12:00:00Z"
    },
    {
      "id": "223e4567-e89b-12d3-a456-426614174000",
      "email": "user2@example.com",
      "first_name": "Jane",
      "last_name": "Smith",
      "status": "active",
      "created_at": "2023-01-02T12:00:00Z"
    }
  ],
  "pagination": {
    "total": 100,
    "page": 1,
    "limit": 10,
    "pages": 10
  }
}
```

### Authentication

#### Login

```
POST /api/v1/auth/login
```

Authenticates a user and returns a JWT token.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securePassword123!"
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2023-01-02T12:00:00Z",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "status": "active"
  }
}
```

#### Refresh Token

```
POST /api/v1/auth/refresh
```

Refreshes a JWT token.

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2023-01-03T12:00:00Z"
}
```

#### Logout

```
POST /api/v1/auth/logout
```

Logs out a user by invalidating their token.

**Request Body:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (204 No Content)**

### User Profile

#### Get Profile

```
GET /api/v1/users/{id}/profile
```

Retrieves a user's profile.

**Response (200 OK):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "profile_picture_url": "https://example.com/profile.jpg",
  "bio": "Software engineer with 10 years of experience",
  "occupation": "Software Engineer",
  "company": "Acme Inc.",
  "website": "https://johndoe.com",
  "social_links": {
    "linkedin": "https://linkedin.com/in/johndoe",
    "twitter": "https://twitter.com/johndoe",
    "github": "https://github.com/johndoe"
  },
  "preferences": {
    "theme": "dark",
    "notifications": {
      "email": true,
      "sms": false
    }
  },
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

#### Update Profile

```
PUT /api/v1/users/{id}/profile
```

Updates a user's profile.

**Request Body:**

```json
{
  "profile_picture_url": "https://example.com/new-profile.jpg",
  "bio": "Senior software engineer with 10+ years of experience",
  "occupation": "Senior Software Engineer",
  "company": "New Company Inc.",
  "website": "https://johndoe.dev",
  "social_links": {
    "linkedin": "https://linkedin.com/in/johndoe",
    "twitter": "https://twitter.com/johndoe",
    "github": "https://github.com/johndoe"
  },
  "preferences": {
    "theme": "light",
    "notifications": {
      "email": true,
      "sms": true
    }
  }
}
```

**Response (200 OK):**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "profile_picture_url": "https://example.com/new-profile.jpg",
  "bio": "Senior software engineer with 10+ years of experience",
  "occupation": "Senior Software Engineer",
  "company": "New Company Inc.",
  "website": "https://johndoe.dev",
  "social_links": {
    "linkedin": "https://linkedin.com/in/johndoe",
    "twitter": "https://twitter.com/johndoe",
    "github": "https://github.com/johndoe"
  },
  "preferences": {
    "theme": "light",
    "notifications": {
      "email": true,
      "sms": true
    }
  },
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-02T12:00:00Z"
}
```

### Password Management

#### Change Password

```
POST /api/v1/users/{id}/password
```

Changes a user's password.

**Request Body:**

```json
{
  "current_password": "securePassword123!",
  "new_password": "evenMoreSecurePassword456!"
}
```

**Response (204 No Content)**

#### Request Password Reset

```
POST /api/v1/auth/password-reset
```

Requests a password reset.

**Request Body:**

```json
{
  "email": "user@example.com"
}
```

**Response (204 No Content)**

#### Reset Password

```
POST /api/v1/auth/password-reset/{token}
```

Resets a password using a reset token.

**Request Body:**

```json
{
  "password": "newSecurePassword789!"
}
```

**Response (204 No Content)**

### Email Verification

#### Request Email Verification

```
POST /api/v1/users/{id}/verify-email
```

Requests an email verification.

**Response (204 No Content)**

#### Verify Email

```
GET /api/v1/auth/verify-email/{token}
```

Verifies an email using a verification token.

**Response (302 Found)**

Redirects to the frontend with a success message.

### MFA Management

#### Enable MFA

```
POST /api/v1/users/{id}/mfa
```

Enables MFA for a user.

**Response (200 OK):**

```json
{
  "secret": "JBSWY3DPEHPK3PXP",
  "qr_code_url": "data:image/png;base64,..."
}
```

#### Verify MFA

```
POST /api/v1/users/{id}/mfa/verify
```

Verifies an MFA code.

**Request Body:**

```json
{
  "code": "123456"
}
```

**Response (204 No Content)**

#### Disable MFA

```
DELETE /api/v1/users/{id}/mfa
```

Disables MFA for a user.

**Request Body:**

```json
{
  "code": "123456"
}
```

**Response (204 No Content)**

### Sessions

#### List Sessions

```
GET /api/v1/users/{id}/sessions
```

Lists all active sessions for a user.

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
      "device_info": "Windows 10, Chrome 91",
      "location": "New York, US",
      "created_at": "2023-01-01T12:00:00Z",
      "last_activity": "2023-01-01T13:00:00Z",
      "expires_at": "2023-01-02T12:00:00Z",
      "current": true
    },
    {
      "id": "223e4567-e89b-12d3-a456-426614174000",
      "ip_address": "192.168.1.2",
      "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
      "device_info": "iPhone, iOS 14.6, Safari",
      "location": "New York, US",
      "created_at": "2023-01-01T10:00:00Z",
      "last_activity": "2023-01-01T11:00:00Z",
      "expires_at": "2023-01-02T10:00:00Z",
      "current": false
    }
  ]
}
```

#### Revoke Session

```
DELETE /api/v1/users/{id}/sessions/{session_id}
```

Revokes a specific session.

**Response (204 No Content)**

#### Revoke All Sessions

```
DELETE /api/v1/users/{id}/sessions
```

Revokes all sessions except the current one.

**Response (204 No Content)**

### Audit Logs

#### List Audit Logs

```
GET /api/v1/users/{id}/audit-logs
```

Lists audit logs for a user.

**Query Parameters:**

- `page`: Page number (default: 1)
- `limit`: Number of items per page (default: 10)
- `action`: Filter by action (optional)
- `from`: Filter by start date (optional)
- `to`: Filter by end date (optional)

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "action": "LOGIN",
      "entity_type": "USER",
      "entity_id": "123e4567-e89b-12d3-a456-426614174000",
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
      "details": {
        "device": "Windows 10, Chrome 91",
        "location": "New York, US"
      },
      "status": "SUCCESS",
      "created_at": "2023-01-01T12:00:00Z"
    },
    {
      "id": "223e4567-e89b-12d3-a456-426614174000",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "action": "PASSWORD_CHANGE",
      "entity_type": "USER",
      "entity_id": "123e4567-e89b-12d3-a456-426614174000",
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
      "details": {
        "device": "Windows 10, Chrome 91",
        "location": "New York, US"
      },
      "status": "SUCCESS",
      "created_at": "2023-01-01T13:00:00Z"
    }
  ],
  "pagination": {
    "total": 50,
    "page": 1,
    "limit": 10,
    "pages": 5
  }
}
```

## Error Responses

### 400 Bad Request

```json
{
  "error": "Bad Request",
  "message": "Invalid request parameters",
  "details": {
    "email": "Email is required",
    "password": "Password must be at least 8 characters long"
  }
}
```

### 401 Unauthorized

```json
{
  "error": "Unauthorized",
  "message": "Invalid credentials"
}
```

### 403 Forbidden

```json
{
  "error": "Forbidden",
  "message": "You do not have permission to access this resource"
}
```

### 404 Not Found

```json
{
  "error": "Not Found",
  "message": "User not found"
}
```

### 409 Conflict

```json
{
  "error": "Conflict",
  "message": "Email already exists"
}
```

### 429 Too Many Requests

```json
{
  "error": "Too Many Requests",
  "message": "Rate limit exceeded",
  "retry_after": 60
}
```

### 500 Internal Server Error

```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred"
}
```
