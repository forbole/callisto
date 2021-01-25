package database

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
	bbanktypes "github.com/forbole/bdjuno/x/bank/types"
)

// SaveAccountBalances allows to store the balance for the given account associating it to the given height
func (db *BigDipperDb) SaveAccountBalances(balances []bbanktypes.AccountBalance) error {
	paramsNumber := 3
	slices := db.splitBalances(balances, paramsNumber)

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
		stmt += " ON CONFLICT DO NOTHING"
		_, err := db.Sql.Exec(stmt, params...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db BigDipperDb) splitBalances(
	balances []bbanktypes.AccountBalance, paramsNumber int,
) [][]bbanktypes.AccountBalance {
	maxPostgreSQLParams := 65535
	maxBalancesPerSlice := maxPostgreSQLParams / paramsNumber

	slices := make([][]bbanktypes.AccountBalance, len(balances)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, balance := range balances {
		slices[sliceIndex] = append(slices[sliceIndex], balance)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}

// SaveSupplyTokenPool allows to save for the given height the given total amount of coins
func (db *BigDipperDb) SaveSupplyToken(coins sdk.Coins, height int64) error {
	query := `INSERT INTO supply(coins, height) VALUES ($1,$2)`

	_, err := db.Sql.Exec(query, pq.Array(dbtypes.NewDbCoins(coins)), height)
	if err != nil {
		return err
	}
	return nil
}

//GetTokenNames returns the list of token names stored inside the supply table
func (db *BigDipperDb) GetTokenNames() ([]string, error) {
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
