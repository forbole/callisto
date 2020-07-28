package coingecko

import (
	"fmt"

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

	//get ids name referring in coingecko that watches in supply modules
	var coins api.Coins
	if err := utils.QueryCoinGecko("/coins/list", &coins); err != nil {
		return err
	}
	names, err := db.GetTokenNames()
	if err != nil {
		return err
	}
	var ids string
	hitcount := 0 //to stop unnesery checking
	//find the id of the coins
	for _, coin := range coins {
		for _, name := range names {
			if coin.Name == name {
				ids += ids + coin.Id + "&"
				hitcount++
				break //not nesserary to do other check name
			}
		}
		if hitcount == len(names) {
			break
		}
	}
	ids = ids[:len(ids)-1] //get rid of tail "&"

	height, err := db.GetLatestHeight()
	if err != nil {
		return err
	}

	//query
	var markets api.Markets
	query := fmt.Sprintf("/coins/markets?vs_currency=usd&ids=%s", ids)
	if err = utils.QueryCoinGecko(query, &markets); err != nil {
		return err
	}

	if err = db.SaveTokensPrice(markets, height); err != nil {
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
