package types

import (
	"database/sql"
	"time"
)

// ModuleRow represents a single row inside the modules table
type ModuleRow struct {
	Module string `db:"module_name"`
}

// Equal return true if two moduleRow is equal
func (v ModuleRow) Equal(w ModuleRow) bool {
	return v.Module == w.Module
}

// ModuleRows represent an array of ModulerRow
type ModuleRows []*ModuleRow

// NewModuleRows return a new instance of ModuleRows
func NewModuleRows(names []string) ModuleRows {
	rows := make([]*ModuleRow, 0)
	for _, name := range names {
		rows = append(rows, &ModuleRow{Module: name})
	}
	return rows
}

// Equal return true if two ModulesRow is equal
func (v ModuleRows) Equal(w *ModuleRows) bool {
	if w == nil {
		return false
	}

	if len(v) != len(*w) {
		return false
	}

	for index, val := range v {
		if !val.Equal(*(*w)[index]) {
			return false
		}
	}
	return true
}

func TimeToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

func NullTimeToTime(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func AreNullTimesEqual(first sql.NullTime, second sql.NullTime) bool {
	return first.Valid == second.Valid && first.Time.Equal(second.Time)
}
