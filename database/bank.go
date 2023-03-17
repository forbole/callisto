package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v4/database/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"
)

// SaveSupply allows to save for the given height the given total amount of coins
func (db *Db) SaveSupply(coins sdk.Coins, height int64) error {
	query := `
INSERT INTO supply (coins, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET coins = excluded.coins,
    	height = excluded.height
WHERE supply.height <= excluded.height`

	_, err := db.SQL.Exec(query, pq.Array(dbtypes.NewDbCoins(coins)), height)
	if err != nil {
		return fmt.Errorf("error while storing supply: %s", err)
	}

	return nil
}

// GetSupply allows to get the current total supply of certain denom
func (db *Db) GetSupply(denom string) (*dbtypes.DbCoin, error) {
	var rows []dbtypes.DbCoins
	err := db.Sqlx.Select(&rows, `SELECT coins FROM supply`)
	if err != nil {
		return nil, fmt.Errorf("error while getting supply: %s", err)
	}

	for _, row := range rows {
		for _, coin := range row {
			if coin.Denom == denom {
				return coin, nil
			}
		}
	}

	return nil, fmt.Errorf("denom not found: %s", err)
}
