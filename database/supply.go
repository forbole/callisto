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

//GetTokenNames get token name from  latest height
func (db BigDipperDb) GetTokenNames() ([]string, error) {
	var names []string
	query := `SELECT (coin).denom FROM (
        SELECT unnest(coins) AS coin FROM supply WHERE height = (
            SELECT max(height) FROM supply
            ) 
		) AS unnested`
	if err := db.Sqlx.Select(&names, query); err != nil {
		return nil, err
	}
	return names, nil
}

//return the latest height(has 30 second latency because depend on supply module)
func (db BigDipperDb) GetLatestHeight() (int64, error) {
	var height []int64
	query := `select max(height) FROM supply`
	if err := db.Sqlx.Select(&height, query); err != nil {
		return -1, err
	}
	return height[0], nil
}
