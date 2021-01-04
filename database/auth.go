package database

import (
	"github.com/cosmos/cosmos-sdk/x/auth/exported"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

// SaveAccount saves the given account information for the given block height and timestamp
func (db *BigDipperDb) SaveAccount(account exported.Account) error {
	stmt := `INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(stmt, account.GetAddress().String())
	return err
}

// GetAccounts returns all the accounts that are currently stored inside the database.
func (db *BigDipperDb) GetAccounts() ([]string, error) {
	sqlStmt := `SELECT * FROM account`

	var rows []dbtypes.AccountRow
	err := db.Sqlx.Select(&rows, sqlStmt)
	if err != nil {
		return nil, err
	}

	addresses := make([]string, len(rows))
	for index, row := range rows {
		addresses[index] = row.Address
	}

	return addresses, nil
}
