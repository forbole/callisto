package database

import (
	"fmt"

	"github.com/forbole/bdjuno/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"


)

// SaveAccounts saves the given accounts inside the database
func (db *Db) SaveAccounts(accounts []types.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	stmt := `INSERT INTO account (address,details) VALUES `
	var params []interface{}

	for i, account := range accounts {
		ai := i * 2
		stmt += fmt.Sprintf("($%d,$%d),", ai+1,ai+2)
		protoContent, ok := account.Details.(authtypes.AccountI)
		if !ok {
			return fmt.Errorf("invalid proposal content types: %T", account.Details)
		}

		anyContent, err := codectypes.NewAnyWithValue(protoContent)
		if err != nil {
			return err
		}

		contentBz, err := db.EncodingConfig.Marshaler.MarshalJSON(anyContent)
		if err != nil {
			return err
		}

		params = append(params, account.Address,string(contentBz))

	}

	stmt = stmt[:len(stmt)-1]
	stmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(stmt, params...)
	return err
}

// GetAccounts returns all the accounts that are currently stored inside the database.
func (db *Db) GetAccounts() ([]string, error) {
	var rows []string
	err := db.Sqlx.Select(&rows, `SELECT address FROM account`)
	return rows, err
}
