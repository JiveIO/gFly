package db

import (
	"database/sql"
	"database/sql/driver"
)

// ScanNullInt64 function will scan NullInt64 value.
func ScanNullInt64(nullInt sql.NullInt64) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Int64
}

// NullInt64 function will create a NullInt64 object.
func NullInt64() sql.NullInt64 {
	return sql.NullInt64{
		Int64: 0,
		Valid: true,
	}
}

// ScanNullInt32 function will scan NullInt32 value.
func ScanNullInt32(nullInt sql.NullInt32) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Int32
}

// NullInt32 function will create a NullInt32 object.
func NullInt32() sql.NullInt32 {
	return sql.NullInt32{
		Int32: 0,
		Valid: true,
	}
}

// ScanNullInt16 function will scan NullInt16 value.
func ScanNullInt16(nullInt sql.NullInt16) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Int16
}

// NullInt16 function will create a NullInt16 object.
func NullInt16() sql.NullInt16 {
	return sql.NullInt16{
		Int16: 0,
		Valid: true,
	}
}
