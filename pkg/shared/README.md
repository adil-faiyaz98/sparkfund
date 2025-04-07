# Shared Services Code

**DEPRECATED**: This directory is being migrated to the root-level `pkg` directory. Please use packages from `pkg` for new code.

## Migration Plan

See [Migration Plan](../../docs/migration-plan.md) for details on the migration from `services/shared` to `pkg`.

## Current Packages

- **`config`**: Configuration loading and management utilities
- **`handlers/health`**: Health check handlers for microservices
- **`middleware`**: HTTP middleware (authentication, rate limiting, etc.)
- **`utils`**: Utility functions

## Usage (Legacy)

These packages should be imported by the microservices using the following import path:

```go
import "github.com/sparkfund/shared/[package]"
```

For example:

```go
import (
    "github.com/sparkfund/shared/config"
    "github.com/sparkfund/shared/handlers/health"
)
```

## New Code

For new code, please use the packages from the `pkg` directory:

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
