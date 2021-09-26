package types

import "time"

// EMoneyInflationRow represents a single row of the emoney_inflation table
type EMoneyInflationRow struct {
	OneRowId          bool      `db:"one_row_id"`
	Inflation         string    `db:"inflation"`
	LastAppliedTime   time.Time `db:"last_applied_time"`
	LastAppliedHeight int64     `db:"last_applied_height"`
	Height            int64     `db:"height"`
}

// NewEMoneyInflationRow allows to build a new EMoneyInflationRow
func NewEMoneyInflationRow(
	inflation string,
	lastAppliedTime time.Time,
	lastAppliedHeight int64,
	height int64,
) EMoneyInflationRow {
	return EMoneyInflationRow{
		OneRowId:          true,
		Inflation:         inflation,
		LastAppliedTime:   lastAppliedTime,
		LastAppliedHeight: lastAppliedHeight,
		Height:            height,
	}
}

// Equal tells whether v and w represent the same rows
func (v EMoneyInflationRow) Equal(w EMoneyInflationRow) bool {
	return v.Inflation == w.Inflation &&
		v.LastAppliedTime == w.LastAppliedTime &&
		v.LastAppliedHeight == w.LastAppliedHeight &&
		v.Height == w.Height
}
