package database

import (
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
