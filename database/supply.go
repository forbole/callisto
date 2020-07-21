package database

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/lib/pq"
)

// SaveSupplyTokenPool allows to save for the given height the given total amount of coins
func (db BigDipperDb) SaveSupplyToken(coins sdk.Coins, height int64) error {
	query := `INSERT INTO supply(coins,height) VALUES ($1,$2)`

	_, err := db.Sql.Exec(query, pq.Array(dbtypes.NewDbCoins(coins)), height)
	if err != nil {
		return err
	}
	return nil
}
