package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewConnection(config *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}
