package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// TableRole Table name
const TableRole = "roles"

// Role struct to describe a role object.
type Role struct {
	ID        uuid.UUID    `db:"id" json:"id" validate:"required,uuid"`
	Name      string       `db:"name" json:"name" validate:"lte=100"`
	Slug      string       `db:"slug" json:"slug" validate:"lte=100"`
	CreatedAt time.Time    `db:"created_at" json:"created_at" validate:"required"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
}
