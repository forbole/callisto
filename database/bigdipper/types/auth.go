package types

import "github.com/forbole/bdjuno/database/types"

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

// Equal tells whether a and b contain the same data
func (a AccountRow) Equal(b AccountRow) bool {
	return a.Address == b.Address
}

// ________________________________________________

// AccountBalanceRow represents a single row inside the account_balance table
type AccountBalanceRow struct {
	Address string         `db:"address"`
	Coins   *types.DbCoins `db:"coins"`
	Height  int64          `db:"height"`
}

// NewAccountBalanceRow allows to build a new AccountBalanceRow instance
func NewAccountBalanceRow(address string, coins types.DbCoins, height int64) AccountBalanceRow {
	return AccountBalanceRow{
		Address: address,
		Coins:   &coins,
		Height:  height,
	}
}

// Equal tells whether a and b contain the same data
func (a AccountBalanceRow) Equal(b AccountBalanceRow) bool {
	return a.Address == b.Address &&
		a.Coins.Equal(b.Coins) &&
		a.Height == b.Height
}
