package coingecko

import (
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	api "github.com/forbole/bdjuno/x/coingecko/apiTypes"
	"github.com/forbole/bdjuno/x/utils"
	"github.com/rs/zerolog/log"
)

// UpdatePrice fetch total amount of coins in the system from RPC and store it into database
func UpdatePrice(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "coingecko").
		Str("operation", "coingecko").
		Msg("getting token price and market cap")

	//get token names
	var coins api.Coins
	if err := utils.QueryCoinGecko("/coins/list", &coins); err != nil {
		return err
	}

	
	/*
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
		} */
	return nil
}
