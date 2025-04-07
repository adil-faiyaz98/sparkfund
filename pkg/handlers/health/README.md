# Health Check Handlers

This package provides standardized health check handlers for microservices.

## Features

- **Health Check**: Comprehensive health check endpoint (`/health`)
- **Liveness Check**: Simple liveness check endpoint (`/live`)
- **Readiness Check**: Readiness check endpoint (`/ready`)
- **Customizable Checks**: Add custom health checks for databases, caches, and external services

## Usage

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/sparkfund/pkg/handlers/health"
)

func main() {
    router := gin.Default()

    // Create a health handler
    healthHandler := health.NewHealthHandler("1.0.0", "abc123")

    // Add custom health checks
    healthHandler.AddCheck("database", &health.DatabaseChecker{
        Check: func() (bool, string) {
            // Check database connection
            return true, "Database connection is healthy"
        },
    })

    healthHandler.AddCheck("cache", &health.CacheChecker{
        Check: func() (bool, string) {
            // Check cache connection
            return true, "Cache connection is healthy"
        },
    })

    // Register health check endpoints
    router.GET("/health", healthHandler.HealthCheck())
    router.GET("/live", healthHandler.LivenessCheck())
    router.GET("/ready", healthHandler.ReadinessCheck())

    router.Run(":8080")
}
```

## Health Check Response

The health check endpoint returns a JSON response with the following structure:

```json
{
    "status": "healthy",
    "version": "1.0.0",
    "commitSha": "abc123",
    "timestamp": "2023-06-01T12:00:00Z",
    "uptime": "1h30m45s",
    "checks": {
        "database": {
            "status": "healthy",
            "message": "Database connection is healthy"
        },
        "cache": {
            "status": "healthy",
            "message": "Cache connection is healthy"
        }
    }
}
```

## Status Codes

- **200 OK**: All checks are healthy
- **503 Service Unavailable**: One or more checks are unhealthy
