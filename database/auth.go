package database

import (
	"fmt"

	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

// SaveAccount saves the given account information for the given block height and timestamp
func (db *BigDipperDb) SaveAccounts(accounts []authttypes.AccountI) error {
	stmt := `INSERT INTO account (address) VALUES `
	var params []interface{}

	for i, account := range accounts {
		ai := i * 1
		stmt += fmt.Sprintf("($%d),", ai+1)
		params = append(params, account.GetAddress().String())
	}

	stmt = stmt[:len(stmt)-1]
	stmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(stmt, params...)
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
