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
// This marsha json return trash not string
		contentBz, err := db.EncodingConfig.Marshaler.MarshalJSON(anyContent)
		if err != nil {
			return err
		}

		contentBzstring:=string(contentBz)

		params = append(params, account.Address,string(contentBzstring))

	}

	stmt = stmt[:len(stmt)-1]
	stmt += " ON CONFLICT (address) DO UPDATE SET details = excluded.details"
	_, err := db.Sql.Exec(stmt, params...)
	return err
}

// GetAccounts returns all the accounts that are currently stored inside the database.
func (db *Db) GetAccounts() ([]types.Account, error) {
	var rows []dbtypes.AccountRow
	var returnRows []types.Account
	err := db.Sqlx.Select(&rows, `SELECT address,details FROM account`)
	if err!=nil{
		return nil,err
	}
	for i,row:=range rows {
		var accountI authtypes.AccountI
		err=db.EncodingConfig.Marshaler.UnmarshalJSON([]byte(row.Details),accountI)
		if err!=nil{
			return nil,err
		}
		returnRows[i]=types.NewAccount(row.Address,accountI)
	}
	return returnRows, err
}
