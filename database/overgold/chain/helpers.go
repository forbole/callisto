package chain

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// ToNullString - helper for creating null string from string.
func ToNullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  value != "",
	}
}

// ToNullInt64 - helper for creating null int64 from int64.
func ToNullInt64(value int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: value,
		Valid: value != 0,
	}
}

// ToNullInt32 - helper for creating null int32 from int32.
func ToNullInt32(value int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: value,
		Valid: value != 0,
	}
}

// ToNullTime - helper for creating null time from time.
func ToNullTime(value time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  value,
		Valid: !value.IsZero(),
	}
}

// IsAlreadyExists - helper for checking if record in the db is already exists.
func IsAlreadyExists(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" { // unique_violation
			return true
		}
	}
	return err != nil && err.Error() == "pq: duplicate key value violates unique constraint"
}
