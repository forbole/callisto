package database

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
)

func (db *Db) SaveTopAccountsBalance(column string, bals []types.NativeTokenBalance) error {
	if len(bals) == 0 {
		return nil
	}

	stmt := fmt.Sprintf("INSERT INTO top_accounts (address, %s) VALUES ", column)

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
	return err

}

func (db *Db) GetAccountBalanceSum(address string) (string, error) {
	stmt := `SELECT 
COALESCE(available,0) + COALESCE(delegation,0) + COALESCE(redelegation,0) + COALESCE(unbonding,0) + COALESCE(reward,0) 
as sum FROM top_accounts WHERE address = $1 
`
	var rows []string
	err := db.Sqlx.Select(&rows, stmt, address)
	if err != nil || len(rows) == 0 {
		return "0", err
	}

	return rows[0], nil
}

func (db *Db) UpdateTopAccountsSum(address, sum string) error {
	stmt := `INSERT INTO top_accounts (address, sum) VALUES ($1, $2) 
ON CONFLICT (address) DO UPDATE SET sum = excluded.sum`

	_, err := db.Sql.Exec(stmt, address, sum)
	return err

}
