package models

import (
	mb "app/core/fluentmodel"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// TableRole Table name
const TableRole = "roles"

// Role struct to describe a role object.
type Role struct {
	// Table meta data
	MetaData mb.MetaData `db:"-" model:"table:roles"`

	// Table fields
	ID        uuid.UUID    `db:"id" model:"name:id; type:uuid,primary"`
	Name      string       `db:"name" model:"name:name"`
	Slug      string       `db:"slug" model:"name:slug"`
	CreatedAt time.Time    `db:"created_at" model:"name:created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" model:"name:updated_at"`
}
