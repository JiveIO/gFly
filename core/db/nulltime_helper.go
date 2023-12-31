package db

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// ScanNullTime function will scan NullTime value.
// Reference https://stackoverflow.com/questions/24564619/nullable-time-time
func ScanNullTime(nullTime sql.NullTime) driver.Value {
	if !nullTime.Valid {
		return nil
	}
	return nullTime.Time
}

// NowNullTime function will create a NullTime object.
func NowNullTime() sql.NullTime {
	return sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
}
