# SparkFund Shared Packages

This directory contains shared packages that can be used across all microservices in the SparkFund application.

## Package Structure

- **`cache`**: Caching utilities for in-memory and distributed caching
- **`config`**: Configuration loading and management utilities
- **`database`**: Database connection and management utilities
- **`errors`**: Error handling and standardized error responses
- **`handlers`**: Common HTTP handlers (health checks, etc.)
- **`logger`**: Logging utilities
- **`metrics`**: Metrics collection and reporting
- **`middleware`**: HTTP middleware (authentication, rate limiting, etc.)

## Usage

These packages should be imported by the microservices using the following import path:

```go
import "github.com/sparkfund/pkg/[package]"
```

For example:

```go
import (
    "github.com/sparkfund/pkg/config"
    "github.com/sparkfund/pkg/handlers/health"
)
```

## Design Principles

1. **Minimal Dependencies**: Each package should have minimal dependencies on other packages
2. **Clear Interfaces**: Each package should provide clear interfaces for its functionality
3. **Configurability**: Packages should be configurable to meet the needs of different services
4. **Testability**: Packages should be designed for easy testing

## Adding New Packages

When adding a new package to the `pkg` directory, follow these guidelines:

1. Create a new directory with a descriptive name
2. Include a README.md file explaining the package's purpose and usage
3. Write comprehensive tests for the package
4. Document all exported functions, types, and variables
5. Avoid dependencies on service-specific code
