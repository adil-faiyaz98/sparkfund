# Money Pulse Common Packages

This directory contains shared packages that are used across multiple services in the Money Pulse application.

## Package Overview

- **analytics**: Integration with BigQuery for data analytics and anomaly detection.
- **auth**: JWT-based authentication and authorization.
- **cache**: Redis-based caching for improved performance.
- **config**: Configuration management using environment variables.
- **health**: Health checking and monitoring for services.
- **logger**: Structured logging with Zap.
- **middleware**: HTTP middleware for common concerns.
- **ml**: Machine learning integration with AWS SageMaker.
- **monitoring**: Prometheus metrics collection.
- **tracing**: Distributed tracing with OpenTelemetry.
- **utils**: Common utility functions.
- **validation**: Request validation using struct tags.

## Usage Examples

### Authentication

```go
import "money-pulse/pkg/auth"

// Create a new token manager
tokenManager, err := auth.NewTokenManager("your-secret-key", 24*time.Hour)
if err != nil {
    log.Fatal(err)
}

// Generate a token
token, err := tokenManager.GenerateToken("user123", []string{"user"}, map[string]interface{}{
    "email": "user@example.com",
})

// Validate a token
claims, err := tokenManager.ValidateToken(token)
```

### Configuration

```go
import "money-pulse/pkg/config"

// Load configuration from environment variables
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatal(err)
}

// Access configuration values
dbDSN := cfg.Database.GetDSN()
serverPort := cfg.Server.Port
```

### Logging

```go
import "money-pulse/pkg/logger"

// Create a new logger
log, err := logger.NewLogger(logger.Config{
    Level:      "info",
    Format:     "json",
    OutputPath: "stdout",
})
if err != nil {
    log.Fatal(err)
}

// Log messages with context and structured data
log.Info("Server starting", "port", 8080)
log.WithFields(map[string]interface{}{
    "user_id": "123",
    "action":  "login",
}).Info("User logged in")

// Log with error
if err != nil {
    log.WithError(err).Error("Failed to process request")
}
```

### Monitoring

```go
import "money-pulse/pkg/monitoring"

// Record HTTP request metrics
monitoring.RecordHTTPRequest("/api/users", "GET", "200", 0.45)

// Record database operation metrics
monitoring.RecordDatabaseOperation("query", "success", 0.1)

// Record business metrics
monitoring.RecordBusinessMetric("active_users", 1250)
```

### Validation

```go
import "money-pulse/pkg/validation"

type User struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Phone    string `json:"phone" validate:"required,phone"`
}

user := User{...}
if err := validation.Validate(user); err != nil {
    // Handle validation errors
    validationErrors := err.(validation.ValidationErrors)
    for _, e := range validationErrors {
        fmt.Printf("Field: %s, Error: %s\n", e.Field, e.Message)
    }
}
```

### Middleware

```go
import (
    "money-pulse/pkg/auth"
    "money-pulse/pkg/logger"
    "money-pulse/pkg/middleware"
)

// Create middleware instance
log, _ := logger.NewLogger(logger.Config{...})
tokenManager, _ := auth.NewTokenManager("secret", 24*time.Hour)
mw := middleware.NewMiddleware(log, tokenManager)

// Use middleware in HTTP server
router := http.NewServeMux()
router.Handle("/api/", mw.AuthMiddleware(apiHandler))
router.Handle("/metrics", mw.LoggingMiddleware(metricsHandler))

// Chain middleware
handler := middleware.Chain(
    mw.LoggingMiddleware,
    mw.RateLimitMiddleware,
    mw.AuthMiddleware,
)(finalHandler)
```

## Best Practices

1. **Consistency**: Use these common packages across all services for consistency.
2. **Error Handling**: Always check and handle errors returned by these packages.
3. **Context**: Pass context through all function calls for proper cancellation and tracing.
4. **Configuration**: Use environment variables for configuration rather than hard-coding values.
5. **Logging**: Include relevant context in logs using structured fields.
6. **Metrics**: Instrument code with metrics for better observability.
7. **Validation**: Validate all user input before processing it. 