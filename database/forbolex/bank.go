package forbolex

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"

	"github.com/lib/pq"

	dbutils "github.com/forbole/bdjuno/database/utils"
	"github.com/forbole/bdjuno/modules/common/bank"
	"github.com/forbole/bdjuno/types"
)

var (
	_ bank.DB = &Db{}
)

// SaveAccountBalances implements bank.DB
func (db *Db) SaveAccountBalances(balances []types.AccountBalance) error {
	paramsNumber := 3
	slices := dbutils.SplitBalances(balances, paramsNumber)

	for _, balances := range slices {
		if len(balances) == 0 {
			continue
		}

		stmt := `INSERT INTO account_balance_history (address, coins, height) VALUES `
		var params []interface{}

		for i, bal := range balances {
			bi := i * paramsNumber
			stmt += fmt.Sprintf("($%d, $%d, $%d),", bi+1, bi+2, bi+3)

			coins := pq.Array(dbtypes.NewDbCoins(bal.Balance))
			params = append(params, bal.Address, coins, bal.Height)
		}

		stmt = stmt[:len(stmt)-1]
		stmt += `ON CONFLICT ON CONSTRAINT unique_balance_for_height DO UPDATE SET coins = excluded.coins`

		_, err := db.Sql.Exec(stmt, params...)
		if err != nil {
			return err
		}
	}

	return nil
}
