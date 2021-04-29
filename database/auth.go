package database

import (
	"fmt"

	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// SaveAccounts saves the given accounts inside the database
func (db *BigDipperDb) SaveAccounts(accounts []authttypes.AccountI) error {
	if len(accounts) == 0 {
		return nil
	}

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
	var rows []string
	err := db.Sqlx.Select(&rows, `SELECT address FROM account`)
	return rows, err
}
