package database

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
)

func (db *Db) SaveTopAccountsBalance(col string, bals []types.NativeTokenBalance) error {
	if len(bals) == 0 {
		return nil
	}

	stmt := fmt.Sprintf("INSERT INTO top_accounts (address, %s) VALUES ", col)

	var params []interface{}

	for i, bal := range bals {
		bi := i * 2
		stmt += fmt.Sprintf("($%d, $%d),", bi+1, bi+2)

		params = append(params, bal.Address, bal.Balance.String())
	}

	stmt = stmt[:len(stmt)-1]
	stmt += `
ON CONFLICT (address) DO UPDATE 
	SET available = excluded.available`

	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return err
	}
	return nil

}
