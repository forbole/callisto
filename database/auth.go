package database

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/lib/pq"
)

// SaveAccount saves the given account information for the given block height and timestamp
func (db BigDipperDb) SaveAccount(account exported.Account, height int64, timestamp time.Time) error {
	accStmt := `INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(accStmt, account.GetAddress().String())
	if err != nil {
		return err
	}

	balStmt := `INSERT INTO balance (address, coins, height, timestamp) 
				VALUES ($1, $2::coin[], $3, $4) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(balStmt,
		account.GetAddress().String(), pq.Array(dbtypes.NewDbCoins(account.GetCoins())), height, timestamp)
	return err
}

// SaveAccount saves the given accounts information for the given block height and timestamp
func (db BigDipperDb) SaveAccounts(accounts []exported.Account, height int64, timestamp time.Time) error {
	// Do nothing with empty accounts
	if len(accounts) == 0 {
		return nil
	}

	accountsStmt := "INSERT INTO account (address) VALUES "
	var accountParams []interface{}

	balancesStmt := "INSERT INTO balance (address, coins, height, timestamp) VALUES "
	var balanceParams []interface{}

	for i, account := range accounts {
		a1 := i * 1 // Starting position for
		b1 := i * 4 // Starting position for balances insertion

		accountsStmt += fmt.Sprintf("($%d),", a1+1)
		accountParams = append(accountParams, account.GetAddress().String())

		balancesStmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", b1+1, b1+2, b1+3, b1+4)
		balanceParams = append(balanceParams,
			account.GetAddress().String(), pq.Array(dbtypes.NewDbCoins(account.GetCoins())), height, timestamp)
	}

	// Store the accounts
	accountsStmt = accountsStmt[:len(accountsStmt)-1] // Remove trailing ","
	accountsStmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(accountsStmt, accountParams...)
	if err != nil {
		return err
	}

	// Store the balances
	balancesStmt = balancesStmt[:len(balancesStmt)-1] // Remove trailing ","
	balancesStmt += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(balancesStmt, balanceParams...)
	return err
}

// GetAccounts returns all the accounts that are currently stored inside the database.
func (db BigDipperDb) GetAccounts() ([]sdk.AccAddress, error) {
	sqlStmt := `SELECT * FROM account`

	var rows []dbtypes.AccountRow
	err := db.Sqlx.Select(&rows, sqlStmt)
	if err != nil {
		return nil, err
	}

	addresses := make([]sdk.AccAddress, len(rows))
	for index, row := range rows {
		address, err := sdk.AccAddressFromBech32(row.Address)
		if err != nil {
			return nil, err
		}

		addresses[index] = address
	}

	return addresses, nil
}
