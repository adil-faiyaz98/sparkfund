# Database Package

This package provides database connectivity and migration management for the Money Pulse application.

## Features

- PostgreSQL connection management with connection pooling
- Environment-based configuration
- Migration management system
- Support for multiple services
- Transaction support for migrations
- Rollback capability

## Usage

### Database Connection

```go
import "github.com/adil-faiyaz98/structgen/database"

// Using connection URL
db, err := database.ConnectFromURL("postgres://user:password@localhost:5432/dbname?sslmode=disable")
if err != nil {
    log.Fatal(err)
}
defer database.Close(db)

// Using configuration
config := database.NewConfig()
db, err := database.Connect(config)
if err != nil {
    log.Fatal(err)
}
defer database.Close(db)
```

### Running Migrations

The package includes a command-line tool for managing migrations:

```bash
# Run migrations for a specific service
go run database/cmd/migrate/main.go -service users -action up

# Rollback migrations for a specific service
go run database/cmd/migrate/main.go -service users -action down
```

## Environment Variables

The following environment variables can be used to configure the database connection:

- `DATABASE_URL`: Complete database connection URL
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: moneypulse)
- `DB_PASSWORD`: Database password (default: moneypulse)
- `DB_NAME`: Database name (default: moneypulse)
- `DB_SSL_MODE`: SSL mode (default: disable)
- `ENV`: Environment (development/production)

## Migration Files

Migration files should be placed in the `internal/<service>/migrations` directory with the following naming convention:

- Up migrations: `YYYYMMDDHHMMSS_description.up.sql`
- Down migrations: `YYYYMMDDHHMMSS_description.down.sql`

Example:
```
internal/users/migrations/
├── 20240315000000_create_users_table.up.sql
└── 20240315000000_create_users_table.down.sql
```

## Migration Table

The package automatically creates a `migrations` table to track applied migrations:

```sql
CREATE TABLE migrations (
    id VARCHAR(255) PRIMARY KEY,
    service VARCHAR(255) NOT NULL,
    timestamp BIGINT NOT NULL,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Best Practices

1. Always include both up and down migrations
2. Use transactions in migrations when possible
3. Keep migrations idempotent
4. Use meaningful timestamps in migration filenames
5. Test migrations in both directions (up and down)
6. Back up the database before running migrations in production

## Error Handling

The package provides detailed error messages for common issues:

- Connection failures
- Migration errors
- Invalid migration files
- Missing environment variables

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 