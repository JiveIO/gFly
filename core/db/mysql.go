package db

import (
	"app/core/utils"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql" // load driver for Mysql
)

// ===========================================================================================================
// 											MySQL Connection
// ===========================================================================================================

func NewMySQL() *MySQL {
	return &MySQL{}
}

type MySQL struct{}

// Connect func for connection to Mysql database.
func (db *MySQL) connect() (*sqlx.DB, error) {
	// Define database connection settings.
	maxConn := utils.Getenv("DB_MAX_CONNECTION", 0)                  // the default is 0 (unlimited)
	maxIdleConn := utils.Getenv("DB_MAX_IDLE_CONNECTION", 2)         // default is 2
	maxLifetimeConn := utils.Getenv("DB_MAX_LIFETIME_CONNECTION", 0) // default is 0, connections are reused forever

	// Build Mysql connection URL.
	mysqlConnURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Define database connection for Mysql.
	dbConnection, err := sqlx.Connect("mysql", mysqlConnURL)
	if err != nil {
		return nil, err
	}

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
