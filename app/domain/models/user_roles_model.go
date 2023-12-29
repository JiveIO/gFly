package models

import (
	"time"

	"github.com/google/uuid"
)

// TableUserRole Table name
const TableUserRole = "user_roles"

// UserRole struct to describe a user role object.
type UserRole struct {
	ID        uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	RoleID    uuid.UUID `db:"role_id" json:"role_id" validate:"required,uuid"`
	UserID    uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
