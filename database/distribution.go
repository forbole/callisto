package database

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

// SaveCommunityPool allows to save for the given height the given total amount of coins
func (db *BigDipperDb) SaveCommunityPool(coin sdk.DecCoins, height int64) error {
	query := `INSERT INTO community_pool(coins, height) VALUES ($1, $2)`
	_, err := db.Sql.Exec(query, pq.Array(dbtypes.NewDbDecCoins(coin)), height)
	return err
}
