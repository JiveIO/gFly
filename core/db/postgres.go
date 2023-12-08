package db

import (
	"app/core/utils"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // load `pgx` driver for PostgreSQL
	"github.com/jmoiron/sqlx"
)

// ===========================================================================================================
// 											MySQL Connection
// ===========================================================================================================

func NewPostgreSQL() *PostgreSQL {
	return &PostgreSQL{}
}

type PostgreSQL struct{}

// Connect func for connection to PostgreSQL database.
func (db *PostgreSQL) connect() (*sqlx.DB, error) {
	// Define database connection settings.
	maxConn := utils.Getenv("DB_MAX_CONNECTION", 0)                  // the default is 0 (unlimited)
	maxIdleConn := utils.Getenv("DB_MAX_IDLE_CONNECTION", 2)         // default is 2
	maxLifetimeConn := utils.Getenv("DB_MAX_LIFETIME_CONNECTION", 0) // default is 0, connections are reused forever

	// Build PostgreSQL connection URL.
	postgresConnURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	// Define database connection for PostgreSQL.
	dbConnection, err := sqlx.Connect("pgx", postgresConnURL)
	if err != nil {
		return nil, err
	}

	// Set database connection settings:
	dbConnection.SetMaxOpenConns(maxConn)
	dbConnection.SetMaxIdleConns(maxIdleConn)
	dbConnection.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	// Try to ping database.
	if err := dbConnection.Ping(); err != nil {
		defer func(db *sqlx.DB) {
			_ = db.Close()
		}(dbConnection)
		return nil, err
	}

	return dbConnection, nil
}
