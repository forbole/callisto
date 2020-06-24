package types

import (
	"time"
)

// AccountRow represents a single row inside the account table
type AccountRow struct {
	Address string `db:"address"`
}

// NewAccountRow allows to easily build a new AccountRow
func NewAccountRow(address string) AccountRow {
	return AccountRow{
		Address: address,
	}
}

// Equals tells whether a and b represent the same row
func (a AccountRow) Equal(b AccountRow) bool {
	return a.Address == b.Address
}

// ________________________________________________

// BalanceRow represents a single database row inside the account table
type BalanceRow struct {
	Address   string    `db:"address"`
	Coins     *DbCoins  `db:"coins"`
	Height    int64     `db:"height"`
	Timestamp time.Time `db:"timestamp"`
}

// NewBalanceRow allows to build a new BalanceRow instance
func NewBalanceRow(address string, coins DbCoins, height int64, timestamp time.Time) BalanceRow {
	return BalanceRow{
		Address:   address,
		Coins:     &coins,
		Height:    height,
		Timestamp: timestamp,
	}
}

// Equal tells whether a and b represent the same database rows
func (a BalanceRow) Equal(b BalanceRow) bool {
	return a.Address == b.Address &&
		a.Height == b.Height &&
		a.Coins.Equal(b.Coins) &&
		a.Timestamp.Equal(b.Timestamp)
}
