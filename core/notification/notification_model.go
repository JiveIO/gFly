package notification

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Notification struct to describe a notification object.
type Notification struct {
	ID              uuid.UUID    `db:"id" validate:"required,uuid"`
	Type            string       `db:"type" validate:"required"`
	NotifiableGroup string       `db:"notifiable_group" validate:"required"`
	NotifiableId    string       `db:"notifiable_id" validate:"required,uuid"`
	Data            string       `db:"data" validate:"required"`
	ReadAt          sql.NullTime `db:"read_at"`
	CreatedAt       time.Time    `db:"created_at" validate:"required"`
	UpdatedAt       time.Time    `db:"updated_at" validate:"required"`
}
