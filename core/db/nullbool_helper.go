package db

import (
	"database/sql"
	"database/sql/driver"
)

// ScanNullBool function will scan NullBool value.
func ScanNullBool(nullBool sql.NullBool) driver.Value {
	if !nullBool.Valid {
		return nil
	}
	return nullBool.Bool
}

// NullBool function will create a NullBool object.
func NullBool() sql.NullBool {
	return sql.NullBool{
		Bool:  true,
		Valid: true,
	}
}
