package db

import (
	"app/core/utils"
	"github.com/jmoiron/sqlx"
)

// ===========================================================================================================
// 											Structure & Interface
// ===========================================================================================================

type IDB interface {
	connect() (*sqlx.DB, error)
}

type Driver string

const (
	postgresql = Driver("postgresql")
	mysql      = Driver("mysql")
)

// ===========================================================================================================
// 												Database
// ===========================================================================================================

// DB the database
type DB struct {
	*sqlx.DB // Embed sqlx DB.
}

// Connect func for opening database connection.
func (db *DB) connect() error {
	var err error
	var dbDriver IDB

	// Get DB_TYPE value from .env file.
	driverType := Driver(utils.Getenv("DB_DRIVER", "mysql"))

	// Define a new Database connection with a right DB type.
	switch driverType {
	case postgresql:
		dbDriver = NewPostgreSQL()
	case mysql:
		dbDriver = NewMySQL()
	}

	if db.DB, err = dbDriver.connect(); err != nil {
		return err
	}

	return nil
}

// defaultDB a singleton database instance
var defaultDB = &DB{}

// Instance returns db instance to handle CRUD at somewhere.
func Instance() *DB {
	return defaultDB
}

// Connect sets the db client of a database, The method should be call at initial.
func Connect() error {
	return defaultDB.connect()
}
