package types

import "time"

// ---------------------------- EMoneyInflationRow ----------------------------

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
	one_row_id bool,
	inflation string,
	last_applied_time time.Time,
	last_applied_height int64,
	height int64,
) EMoneyInflationRow {
	return EMoneyInflationRow{
		OneRowId:          true,
		Inflation:         inflation,
		LastAppliedTime:   last_applied_time,
		LastAppliedHeight: last_applied_height,
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

// ---------------------------- EMoneyGasPrices ----------------------------

// EMoneyGasPricesRow represents a single row of the emoney_gas_prices table
type EMoneyGasPricesRow struct {
	AuthorityKey string `db:"authority_key"`
	GasPrices    string `db:"gas_prices"`
	Height       int64  `db:"height"`
}

// EMoneyGasPricesRow allows to build a new EmoneyGasPricesRow
func NewEMoneyGasPricesRow(
	authority_key string,
	gas_prices string,
	height int64,
) EMoneyGasPricesRow {
	return EMoneyGasPricesRow{
		AuthorityKey: authority_key,
		GasPrices:    gas_prices,
		Height:       height,
	}
}

// Equal tells whether v and w represent the same rows
func (v EMoneyGasPricesRow) Equal(w EMoneyGasPricesRow) bool {
	return v.AuthorityKey == w.AuthorityKey &&
		v.GasPrices == w.GasPrices &&
		v.Height == w.Height
}
