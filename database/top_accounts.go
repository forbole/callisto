package database

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
)

func (db *Db) SaveTopAccountsBalance(column string, bals []types.NativeTokenAmount) error {
	if len(bals) == 0 {
		return nil
	}

	stmt := fmt.Sprintf("INSERT INTO top_accounts (address, %s, height) VALUES ", column)

	var params []interface{}

	for i, bal := range bals {
		bi := i * 3
		stmt += fmt.Sprintf("($%d, $%d, $%d),", bi+1, bi+2, bi+3)

		params = append(params, bal.Address, bal.Balance.String(), bal.Height)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += fmt.Sprintf("ON CONFLICT (address) DO UPDATE SET %s = excluded.%s, height = excluded.height WHERE top_accounts.height <= excluded.height", column, column)

	_, err := db.SQL.Exec(stmt, params...)
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

func (db *Db) UpdateTopAccountsSum(address, sum string, height int64) error {
	stmt := `INSERT INTO top_accounts (address, sum, height) VALUES ($1, $2, $3) 
ON CONFLICT (address) DO UPDATE SET 
	sum = excluded.sum, 
	height = excluded.height  
WHERE top_accounts.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, address, sum, height)
	return err

}

// SaveTotalAccounts allows to store total accounts params inside the database
func (db *Db) SaveTotalAccounts(totalAccounts int64, height int64) error {
	stmt := `
INSERT INTO top_accounts_params (total_accounts, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET total_accounts = excluded.total_accounts,
        height = excluded.height
WHERE top_accounts_params.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, totalAccounts, height)
	if err != nil {
		return fmt.Errorf("error while storing top accounts params: %s", err)
	}

	return nil
}
