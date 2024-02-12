package db

import (
	"database/sql"
	"database/sql/driver"
)

// ScanNullByte function will scan NullByte value.
func ScanNullByte(nullBool sql.NullByte) driver.Value {
	if !nullBool.Valid {
		return nil
	}
	return nullBool.Byte
}

// NullByte function will create a NullBool object.
func NullByte(val byte) sql.NullByte {
	return sql.NullByte{
		Byte:  val,
		Valid: true,
	}
}
