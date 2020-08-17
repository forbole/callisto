package types

import (
	"time"
)

// ModulesRow represents a single row inside the modules table
type ModulesRow struct {
	Staking      bool      `db:"staking"`
	Auth         bool      `db:"auth"`
	Supply       bool      `db:"supply"`
	Distribution bool      `db:"distribution"`
	Pricefeed    bool      `db:"pricefeed"`
	Bank         bool      `db:"bank"`
	Consensus    bool      `db:"consensus"`
	Mint         bool      `db:"mint"`
	Timestamp    time.Time `db:"timestamp"`
}

//NewModulesRow return a new instance of ModulesRow
func NewModulesRow(
	staking bool,
	auth bool,
	supply bool,
	distribution bool,
	pricefeed bool,
	bank bool,
	consensus bool,
	mint bool,
	timestamp time.Time,
) ModulesRow {
	return ModulesRow{
		Staking:      staking,
		Auth:         auth,
		Supply:       supply,
		Distribution: distribution,
		Pricefeed:    pricefeed,
		Bank:         bank,
		Consensus:    consensus,
		Mint:         mint,
		Timestamp:    timestamp,
	}
}

// Equals returns true if two ModulesRow are the same
func (w ModulesRow) Equals(v ModulesRow) bool {
	return w.Staking == v.Staking &&
		w.Auth == v.Auth &&
		w.Supply == v.Supply &&
		w.Distribution == v.Distribution &&
		w.Pricefeed == v.Pricefeed &&
		w.Bank == v.Bank &&
		w.Consensus == v.Consensus &&
		w.Mint == v.Mint &&
		w.Timestamp.Equal(v.Timestamp)
}
