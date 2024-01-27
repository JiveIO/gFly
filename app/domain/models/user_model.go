package models

import (
	"database/sql"
	mb "github.com/jiveio/fluentmodel"
	"time"

	"github.com/google/uuid"
)

// TableUser Table name
const TableUser = "users"

// User struct to describe a user object.
type User struct {
	// Table meta data
	MetaData mb.MetaData `db:"-" model:"table:users"`

	// Table fields
	ID           uuid.UUID      `db:"id" model:"name:id; type:uuid,primary" json:"id"`
	Email        string         `db:"email" model:"name:email" json:"email"`
	PasswordHash string         `db:"password_hash" model:"name:password_hash" json:"password_hash"`
	Fullname     string         `db:"fullname" model:"name:fullname" json:"fullname"`
	Phone        string         `db:"phone" model:"name:phone" json:"phone"`
	Token        sql.NullString `db:"token" model:"name:token" json:"token"`
	UserStatus   int            `db:"user_status" model:"name:user_status" json:"user_status"`
	CreatedAt    time.Time      `db:"created_at" model:"name:created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" model:"name:updated_at" json:"updated_at"`
	VerifiedAt   sql.NullTime   `db:"verified_at" model:"name:verified_at" json:"verified_at"`
	BlockedAt    sql.NullTime   `db:"blocked_at" model:"name:blocked_at" json:"blocked_at"`
	DeletedAt    sql.NullTime   `db:"deleted_at" model:"name:deleted_at" json:"deleted_at"`
	LastAccessAt sql.NullTime   `db:"last_access_at" model:"name:last_access_at" json:"last_access_at"`
}
