package bigdipper

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"

	dbutils "github.com/forbole/bdjuno/database/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	"github.com/forbole/bdjuno/types"
)

// SaveAccountBalances allows to store the balance for the given account associating it to the given height
func (db *Db) SaveAccountBalances(balances []types.AccountBalance) error {
	paramsNumber := 3
	slices := dbutils.SplitBalances(balances, paramsNumber)

	for _, balances := range slices {
		if len(balances) == 0 {
			continue
		}

		stmt := `INSERT INTO account_balance (address, coins, height) VALUES `
		var params []interface{}

		for i, bal := range balances {
			bi := i * paramsNumber
			stmt += fmt.Sprintf("($%d, $%d, $%d),", bi+1, bi+2, bi+3)

			coins := pq.Array(dbtypes.NewDbCoins(bal.Balance))
			params = append(params, bal.Address, coins, bal.Height)
		}

		stmt = stmt[:len(stmt)-1]
		stmt += `
ON CONFLICT (address) DO UPDATE 
	SET coins = excluded.coins, 
	    height = excluded.height 
WHERE account_balance.height <= excluded.height`

		_, err := db.Sql.Exec(stmt, params...)
		if err != nil {
			return err
		}
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveSupply allows to save for the given height the given total amount of coins
func (db *Db) SaveSupply(coins sdk.Coins, height int64) error {
	query := `
INSERT INTO supply (coins, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET coins = excluded.coins,
    	height = excluded.height
WHERE supply.height <= excluded.height`

	_, err := db.Sql.Exec(query, pq.Array(dbtypes.NewDbCoins(coins)), height)
	if err != nil {
		return err
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetTokenNames returns the list of token names stored inside the supply table
func (db *Db) GetTokenNames() ([]string, error) {
	var names []string
	query := `
SELECT (coin).denom FROM (
    SELECT unnest(coins) AS coin FROM supply WHERE height = (
        SELECT max(height) FROM supply
	) 
) AS unnested`
	if err := db.Sqlx.Select(&names, query); err != nil {
		return nil, err
	}
	return names, nil
}
