package types

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

// Equals tells whether a and b contain the same data
func (a AccountRow) Equal(b AccountRow) bool {
	return a.Address == b.Address
}

// ________________________________________________

// AccountBalanceRow represents a single row inside the account_balance table
type AccountBalanceRow struct {
	Address string   `db:"address"`
	Coins   *DbCoins `db:"coins"`
}

// NewAccountBalanceRow allows to build a new AccountBalanceRow instance
func NewAccountBalanceRow(address string, coins DbCoins) AccountBalanceRow {
	return AccountBalanceRow{
		Address: address,
		Coins:   &coins,
	}
}

// Equal tells whether a and b contain the same data
func (a AccountBalanceRow) Equal(b AccountBalanceRow) bool {
	return a.Address == b.Address &&
		a.Coins.Equal(b.Coins)
}

// AccountBalanceHistoryRow represents a single row inside the account_balance_history table
type AccountBalanceHistoryRow struct {
	Address string   `db:"address"`
	Coins   *DbCoins `db:"coins"`
	Height  int64    `db:"height"`
}

// NewAccountBalanceHistoryRow allows to build a new AccountBalanceHistoryRow instance
func NewAccountBalanceHistoryRow(address string, coins DbCoins, height int64) AccountBalanceHistoryRow {
	return AccountBalanceHistoryRow{
		Address: address,
		Coins:   &coins,
		Height:  height,
	}
}

// Equal tells whether a and b contain the same data
func (a AccountBalanceHistoryRow) Equal(b AccountBalanceHistoryRow) bool {
	return a.Address == b.Address &&
		a.Height == b.Height &&
		a.Coins.Equal(b.Coins)
}
