package database

import (
	"fmt"

	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
)

// SaveAccountBalanceHistory saves the given entry inside the account_balance_history table
func (db *Db) SaveAccountBalanceHistory(entry types.AccountBalanceHistory) error {
	stmt := `
INSERT INTO account_balance_history (address, balance, delegated, unbonding, redelegating, commission, reward, timestamp) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT ON CONSTRAINT unique_balance_for_height DO UPDATE 
    SET balance = excluded.balance,
        delegated = excluded.delegated,
        unbonding = excluded.unbonding,
        redelegating = excluded.redelegating,
        commission = excluded.commission, 
        reward = excluded.reward`

	_, err := db.Sql.Exec(stmt,
		entry.Account,
		pq.Array(dbtypes.NewDbCoins(entry.Balance)),
		pq.Array(dbtypes.NewDbCoins(entry.Delegations)),
		pq.Array(dbtypes.NewDbCoins(entry.Unbonding)),
		pq.Array(dbtypes.NewDbCoins(entry.Redelegations)),
		pq.Array(dbtypes.NewDbDecCoins(entry.Commission)),
		pq.Array(dbtypes.NewDbDecCoins(entry.Reward)),
		entry.Timestamp,
	)
	return err
}

// SaveTokenPricesHistory stores the given prices as historic ones
func (db *Db) SaveTokenPricesHistory(prices []types.TokenPrice) error {
	if len(prices) == 0 {
		return nil
	}

	query := `INSERT INTO token_price_history (unit_name, price, market_cap, timestamp) VALUES`
	var param []interface{}

	for i, ticker := range prices {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, ticker.UnitName, ticker.Price, ticker.MarketCap, ticker.Timestamp)
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT ON CONSTRAINT unique_price_for_timestamp DO UPDATE 
	SET price = excluded.price,
	    market_cap = excluded.market_cap`

	_, err := db.Sql.Exec(query, param...)
	return err
}
