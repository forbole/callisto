package database

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

// SaveAccount saves the given account information for the given block height and timestamp
func (db *BigDipperDb) SaveAccount(account exported.Account, height int64, timestamp time.Time) error {
	stmt := `INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(stmt, account.GetAddress().String())
	if err != nil {
		return err
	}

	coins := pq.Array(dbtypes.NewDbCoins(account.GetCoins()))

	stmt = `
INSERT INTO account_balance (address, coins) 
VALUES ($1, $2) 
ON CONFLICT (address) DO UPDATE SET coins = excluded.coins`
	_, err = db.Sql.Exec(stmt, account.GetAddress().String(), coins)
	if err != nil {
		return err
	}

	stmt = `
INSERT INTO account_balance_history (address, coins, height, timestamp) 
VALUES ($1, $2::coin[], $3, $4) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt, account.GetAddress().String(), coins, height, timestamp)
	return err
}

// SaveAccount saves the given accounts information for the given block height and timestamp
func (db *BigDipperDb) SaveAccounts(accounts []exported.Account, height int64, timestamp time.Time) error {
	// Do nothing with empty accounts
	if len(accounts) == 0 {
		return nil
	}

	accQry := `INSERT INTO account (address) VALUES `
	var accParams []interface{}

	balQry := `INSERT INTO account_balance (address, coins) VALUES `
	var balParams []interface{}

	balHisQry := `INSERT INTO account_balance_history (address, coins, height, timestamp) VALUES `
	var bHisParams []interface{}

	for i, account := range accounts {
		ai := i * 1
		accQry += fmt.Sprintf("($%d),", ai+1)
		accParams = append(accParams, account.GetAddress().String())

		coins := pq.Array(dbtypes.NewDbCoins(account.GetCoins()))

		bi := i * 2
		balQry += fmt.Sprintf("($%d, $%d),", bi+1, bi+2)
		balParams = append(balParams, account.GetAddress().String(), coins)

		bHi := i * 4
		balHisQry += fmt.Sprintf("($%d,$%d,$%d,$%d),", bHi+1, bHi+2, bHi+3, bHi+4)
		bHisParams = append(bHisParams, account.GetAddress().String(), coins, height, timestamp)
	}

	// Store the accounts
	accQry = accQry[:len(accQry)-1] // Remove trailing ","
	accQry += " ON CONFLICT (address) DO NOTHING"
	_, err := db.Sql.Exec(accQry, accParams...)
	if err != nil {
		return err
	}

	// Remove all the existing balances
	_, err = db.Sql.Exec(`DELETE FROM account_balance WHERE TRUE`)
	if err != nil {
		return err
	}

	// Insert the current balances
	balQry = balQry[:len(balQry)-1] // Remove trailing ","
	balQry += " ON CONFLICT (address) DO NOTHING"
	_, err = db.Sql.Exec(balQry, balParams...)
	if err != nil {
		return err
	}

	// Store the balances histories
	balHisQry = balHisQry[:len(balHisQry)-1] // Remove trailing ","
	balHisQry += " ON CONFLICT (address, height) DO NOTHING"
	_, err = db.Sql.Exec(balHisQry, bHisParams...)
	return err
}

// GetAccounts returns all the accounts that are currently stored inside the database.
func (db *BigDipperDb) GetAccounts() ([]sdk.AccAddress, error) {
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
