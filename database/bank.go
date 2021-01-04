package database

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

// SaveAccountBalance allows to store the balance for the given account associating it to the given height
func (db *BigDipperDb) SaveAccountBalance(address string, balance sdk.Coins, height int64) error {
	coins := pq.Array(dbtypes.NewDbCoins(balance))

	stmt := `
INSERT INTO account_balance (address, coins)
VALUES ($1, $2) ON CONFLICT (address) DO UPDATE 
    SET coins = excluded.coins`
	_, err := db.Sql.Exec(stmt, address, coins)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO account_balance_history (address, coins, height) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt, address, coins, height)
	return err
}
