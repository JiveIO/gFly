package models

import (
	"database/sql"
	mb "github.com/jiveio/fluentmodel"
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
	ID        uuid.UUID    `db:"id" model:"name:id; type:uuid,primary" json:"id"`
	Name      string       `db:"name" model:"name:name" json:"name"`
	Slug      string       `db:"slug" model:"name:slug" json:"slug"`
	CreatedAt time.Time    `db:"created_at" model:"name:created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" model:"name:updated_at" json:"updated_at"`
}
