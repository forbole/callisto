package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------------------------- EMoneyInflationRow ----------------------------

// EMoneyInflationRow represents a single row of the emoney_inflation table
type EMoneyInflationRow struct {
	OneRowID          bool      `db:"one_row_id"`
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
		OneRowID:          true,
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

// ---------------------------- EMoneyGasPrices ----------------------------

// EMoneyGasPricesRow represents a single row of the emoney_gas_prices table
type EMoneyGasPricesRow struct {
	AuthorityKey string       `db:"authority_key"`
	GasPrices    sdk.DecCoins `db:"gas_prices"`
	Height       int64        `db:"height"`
}

// EMoneyGasPricesRow allows to build a new EmoneyGasPricesRow
func NewEMoneyGasPricesRow(
	authorityKey string,
	gasPrices sdk.DecCoins,
	height int64,
) EMoneyGasPricesRow {
	return EMoneyGasPricesRow{
		AuthorityKey: authorityKey,
		GasPrices:    gasPrices,
		Height:       height,
	}
}

// Equal tells whether v and w represent the same rows
func (v EMoneyGasPricesRow) Equal(w EMoneyGasPricesRow) bool {
	return v.AuthorityKey == w.AuthorityKey &&
		v.Height == w.Height
}
