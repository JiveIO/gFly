package models

import (
	mb "app/core/fluentmodel"
	"time"

	"github.com/google/uuid"
)

// TableUserRole Table name
const TableUserRole = "user_roles"

// UserRole struct to describe a user role object.
type UserRole struct {
	// Table meta data
	MetaData mb.MetaData `db:"-" model:"table:user_roles"`

	// Table fields
	ID        uuid.UUID `db:"id" model:"name:id; type:uuid,primary"`
	RoleID    uuid.UUID `db:"role_id" model:"name:role_id; type:uuid"`
	UserID    uuid.UUID `db:"user_id" model:"name:user_id; type:uuid"`
	CreatedAt time.Time `db:"created_at" model:"name:created_at"`
}
