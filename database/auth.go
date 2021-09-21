package database

import (
	"fmt"

	dbutils "github.com/forbole/bdjuno/database/utils"

	"github.com/forbole/bdjuno/types"
)

// SaveAccounts saves the given accounts inside the database
func (db *Db) SaveAccounts(accounts []types.Account) error {
	paramsNumber := 1
	slices := dbutils.SplitAccounts(accounts, paramsNumber)

	for _, accounts := range slices {
		if len(accounts) == 0 {
			continue
		}

		// Store up-to-date data
		err := db.saveAccounts(paramsNumber, accounts)
		if err != nil {
			return fmt.Errorf("error while storing accounts: %s", err)
		}
	}

	return nil
}

func (db *Db) saveAccounts(paramsNumber int, accounts []types.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	stmt := `INSERT INTO account (address) VALUES `
	var params []interface{}

	for i, account := range accounts {
		ai := i * paramsNumber
		stmt += fmt.Sprintf("($%d),", ai+1)
		params = append(params, account.Address)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing accounts: %s", err)
	}

	return nil
}

// GetAccounts returns all the accounts that are currently stored inside the database.
func (db *Db) GetAccounts() ([]string, error) {
	var rows []string
	err := db.Sqlx.Select(&rows, `SELECT address FROM account`)
	return rows, err
}
