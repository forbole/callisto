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

//GetTokenNames get names of the tokens that exist at the latest height
func (db BigDipperDb) GetTokenNames()([]string, error){
	var names []string
	query :=`select (coin).denom from (
        select unnest(coins) as coin from supply where height = (
            select max(height) from supply
            ) 
		) as unnested`
	if err := db.Sqlx.Select(&names,query);err!=nil{
		return nil,err
	}
	return names,nil
}
