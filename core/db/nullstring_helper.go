package db

import (
	"database/sql"
	"database/sql/driver"
)

// ScanNullString function will scan NullString value.
func ScanNullString(nullString sql.NullString) driver.Value {
	if !nullString.Valid {
		return nil
	}
	return nullString.String
}
