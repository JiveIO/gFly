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

// NullString function will create a NullString object.
func NullString() sql.NullString {
	return sql.NullString{
		String: "",
		Valid:  true,
	}
}
