package coingecko

import (
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	"github.com/forbole/bdjuno/coingecko"
	"github.com/forbole/bdjuno/coingecko/apiTypes"

)

// UpdatePrice fetch total amount of coins in the system from RPC and store it into database
func UpdatePrice(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "coingecko").
		Str("operation", "coingecko").
		Msg("getting token price and market cap")

	//get token name	
	type coins []apiTypes.Coin
	if err :=coingecko.GetCoinGeckoReqBody("/coin/list",&coins);err!=nil{
		return err
	}

	var s sdk.Coins
	height, err := cp.QueryLCDWithHeight("/supply/total", &s)
	if err != nil {
		return err
	}
	// Store the signing infos into the database
	err = db.SaveSupplyToken(s,
		height,
	)
	if err != nil {
		return err
	}
	return nil
}
