package database

import (
	"fmt"
	"time"

	"github.com/sparkfund/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config represents PostgreSQL configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Client represents a PostgreSQL client
type Client struct {
	db *gorm.DB
}

// NewClient creates a new PostgreSQL client
func NewClient(cfg *Config) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.ErrInternalServer(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.ErrInternalServer(err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Client{db: db}, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return errors.ErrInternalServer(err)
	}
	return sqlDB.Close()
}

// DB returns the underlying GORM DB instance
func (c *Client) DB() *gorm.DB {
	return c.db
}

// Create creates a new record
func (c *Client) Create(value interface{}) error {
	if err := c.db.Create(value).Error; err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// Find finds records
func (c *Client) Find(dest interface{}, conds ...interface{}) error {
	if err := c.db.Find(dest, conds...).Error; err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// First finds the first record
func (c *Client) First(dest interface{}, conds ...interface{}) error {
	if err := c.db.First(dest, conds...).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound(err)
		}
		return errors.ErrInternalServer(err)
	}
	return nil
}

// Update updates a record
func (c *Client) Update(value interface{}) error {
	if err := c.db.Save(value).Error; err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// Delete deletes a record
func (c *Client) Delete(value interface{}) error {
	if err := c.db.Delete(value).Error; err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// Transaction executes a function within a transaction
func (c *Client) Transaction(fc func(tx *gorm.DB) error) error {
	return c.db.Transaction(fc)
}

// Migrate runs database migrations
func (c *Client) Migrate(dst ...interface{}) error {
	if err := c.db.AutoMigrate(dst...); err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// Where adds a WHERE clause
func (c *Client) Where(query interface{}, args ...interface{}) *gorm.DB {
	return c.db.Where(query, args...)
}

// Order adds an ORDER BY clause
func (c *Client) Order(value interface{}) *gorm.DB {
	return c.db.Order(value)
}

// Limit adds a LIMIT clause
func (c *Client) Limit(limit int) *gorm.DB {
	return c.db.Limit(limit)
}

// Offset adds an OFFSET clause
func (c *Client) Offset(offset int) *gorm.DB {
	return c.db.Offset(offset)
}

// Joins adds JOIN clauses
func (c *Client) Joins(query string, args ...interface{}) *gorm.DB {
	return c.db.Joins(query, args...)
}

// Group adds a GROUP BY clause
func (c *Client) Group(name string) *gorm.DB {
	return c.db.Group(name)
}

// Having adds a HAVING clause
func (c *Client) Having(query interface{}, args ...interface{}) *gorm.DB {
	return c.db.Having(query, args...)
}

// Distinct adds a DISTINCT clause
func (c *Client) Distinct(args ...interface{}) *gorm.DB {
	return c.db.Distinct(args...)
}

// Count returns the count of records
func (c *Client) Count(count *int64) error {
	if err := c.db.Count(count).Error; err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// Pluck returns a slice of values for a column
func (c *Client) Pluck(column string, dest interface{}) error {
	if err := c.db.Pluck(column, dest).Error; err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// Raw executes a raw SQL query
func (c *Client) Raw(sql string, values ...interface{}) *gorm.DB {
	return c.db.Raw(sql, values...)
}

// Exec executes a raw SQL query without returning rows
func (c *Client) Exec(sql string, values ...interface{}) error {
	if err := c.db.Exec(sql, values...).Error; err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
} 