package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User struct to describe a user object.
type User struct {
	ID           uuid.UUID    `db:"id" json:"id" validate:"required,uuid"`
	Email        string       `db:"email" json:"email" validate:"required,email,lte=255"`
	PasswordHash string       `db:"password_hash" json:"password_hash" validate:"required,gte=6"`
	Fullname     string       `db:"fullname" json:"fullname" validate:"lte=255"`
	Phone        string       `db:"phone" json:"phone" validate:"lte=20"`
	Token        string       `db:"token" json:"token" validate:"lte=100"`
	UserStatus   int          `db:"user_status" json:"user_status" validate:"required,len=1"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at" validate:"required"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at" validate:"required"`
	VerifiedAt   sql.NullTime `db:"verified_at" json:"verified_at"`
	BlockedAt    sql.NullTime `db:"blocked_at" json:"blocked_at"`
	DeletedAt    sql.NullTime `db:"deleted_at" json:"deleted_at"`
	LastAccessAt sql.NullTime `db:"last_access_at" json:"last_access_at"`
}
