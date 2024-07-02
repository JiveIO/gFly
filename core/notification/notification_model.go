package notification

import (
	mb "app/core/fluentmodel"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Notification struct to describe a notification object.
type Notification struct {
	// Table meta data
	MetaData mb.MetaData `db:"-" model:"table:notifications"`

	// Table fields
	ID              uuid.UUID    `db:"id" model:"name:id; type:uuid,primary" validate:"required,uuid"`
	Type            string       `db:"type" model:"name:type" validate:"required"`
	NotifiableGroup string       `db:"notifiable_group" model:"name:notifiable_group" validate:"required"`
	NotifiableId    string       `db:"notifiable_id" model:"name:notifiable_id" validate:"required,uuid"`
	Data            string       `db:"data" model:"name:data" validate:"required"`
	ReadAt          sql.NullTime `db:"read_at" model:"name:read_at"`
	CreatedAt       time.Time    `db:"created_at" model:"name:created_at" validate:"required"`
	UpdatedAt       time.Time    `db:"updated_at" model:"name:updated_at" validate:"required"`
}
