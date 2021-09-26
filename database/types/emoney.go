package types

import "time"

// EmoneyInflationRow represents a single row of the emoney_inflation table
type EmoneyInflationRow struct {
	OneRowId          bool      `db:"one_row_id"`
	Inflation         string    `db:"inflation"`
	LastAppliedTime   time.Time `db:"last_applied_time"`
	LastAppliedHeight int64     `db:"last_applied_height"`
	Height            int64     `db:"height"`
}

// EmoneyInflationRow allows to build a new EmoneyInflationRow
func NewEmoneyInflationRow(
	one_row_id bool,
	inflation string,
	last_applied_time time.Time,
	last_applied_height int64,
	height int64,
) EmoneyInflationRow {
	return EmoneyInflationRow{
		OneRowId:          true,
		Inflation:         inflation,
		LastAppliedTime:   last_applied_time,
		LastAppliedHeight: last_applied_height,
		Height:            height,
	}
}

// Equal tells whether v and w represent the same rows
func (v EmoneyInflationRow) Equal(w EmoneyInflationRow) bool {
	return v.Inflation == w.Inflation &&
		v.LastAppliedTime == w.LastAppliedTime &&
		v.LastAppliedHeight == w.LastAppliedHeight &&
		v.Height == w.Height
}
