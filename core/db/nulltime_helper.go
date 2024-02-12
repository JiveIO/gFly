package db

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// ScanNullTime function will scan NullTime value.
func ScanNullTime(nullTime sql.NullTime) driver.Value {
	if !nullTime.Valid {
		return nil
	}
	return nullTime.Time
}

// NowNullTime function will create a NullTime object.
func NowNullTime() sql.NullTime {
	return NullTime(time.Now())
}

// NullTime function will create a NullTime object.
func NullTime(val time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  val,
		Valid: true,
	}
}
