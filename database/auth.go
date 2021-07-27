package database

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
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
		err := db.saveAccounts(accounts)
		if err != nil {
			return fmt.Errorf("error while storing accounts: %s", err)
		}
	}

	return nil
}

func (db *Db) saveAccounts(accounts []types.Account) error {
	if len(accounts) == 0 {
		return nil
	}
	stmt := `INSERT INTO account (address,details) VALUES `
	var params []interface{}

	for i, account := range accounts {
		ai := i
		stmt += fmt.Sprintf("($%d,$%d),", ai+1, ai+2)
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

		contentBzstring := string(contentBz)

		params = append(params, account.Address, contentBzstring)
		stmt = stmt[:len(stmt)-1]
		stmt += " ON CONFLICT (address) DO UPDATE SET details = excluded.details"
		_, err = db.Sql.Exec(stmt, params...)
		if err != nil {
			return err
		}

	}
	return nil
}

// GetAccounts returns all the accounts that are currently stored inside the database.
func (db *Db) GetAccounts() ([]types.Account, error) {
	var rows []dbtypes.AccountRow
	err := db.Sqlx.Select(&rows, `SELECT address,details FROM account`)
	if err != nil {
		return nil, err
	}

	returnRows := make([]types.Account, len(rows))
	for i, row := range rows {
		b := []byte(row.Details)

		if len(b) == 0 {
			returnRows[i] = types.NewAccount(row.Address, nil)
		} else {
			//var inter interface{}
			var account authtypes.AccountI
			err = db.EncodingConfig.Marshaler.UnmarshalInterfaceJSON(b, &account)
			if err != nil {
				return nil, err
			}
			returnRows[i] = types.NewAccount(row.Address, account)
		}
	}
	return returnRows, nil
}
