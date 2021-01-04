package database

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

// SaveSupplyTokenPool allows to save for the given height the given total amount of coins
func (db *BigDipperDb) SaveSupplyToken(coins sdk.Coins, height int64) error {
	query := `INSERT INTO supply(coins,height) VALUES ($1,$2)`

	_, err := db.Sql.Exec(query, pq.Array(dbtypes.NewDbCoins(coins)), height)
	if err != nil {
		return err
	}
	return nil
}

//GetTokenNames returns the list of token names stored inside the supply table
func (db *BigDipperDb) GetTokenNames() ([]string, error) {
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
