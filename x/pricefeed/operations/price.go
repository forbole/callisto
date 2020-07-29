package pricefeed

import (
	"fmt"
	"time"

	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	api "github.com/forbole/bdjuno/x/pricefeed/apiTypes"
	"github.com/forbole/bdjuno/x/utils"
	"github.com/rs/zerolog/log"
)

// UpdatePrice fetch total amount of coins in the system from RPC and store it into database
func UpdatePrice(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "pricefeed").
		Str("operation", "pricefeed").
		Msg("getting token price and market cap")

	//get token names
	var coins api.Coins
	if err := queryCoinGecko("/coins/list", &coins); err != nil {
		return err
	}
	names, err := db.GetTokenNames()
	if err != nil {
		return err
	}
	//requiredCoin point out index that storing the coin we want in coins
	var ids string
	hitcount := 0
	//find the id of the coins
	for _, coin := range coins {
		for _, name := range names {
			if coin.Name == name {
				ids += ids + coin.ID + "&"
				hitcount++
				break //not nesserary to do other check name
			}
		}
		if hitcount == len(names) {
			break
		}
	}
	if hitcount == 0 {
		return fmt.Errorf("cannot find given ids from the API:%s", names)
	}
	ids = ids[:len(ids)-1] //get rid of tail "&""
	timestamp, err := time.Parse(time.RFC3339, time.Now().String())
	if err != nil {
		return err
	}
	//query
	var pricefeeds api.Pricefeeds
	query := fmt.Sprintf("/coins/markets?vs_currency=usd&ids=%s", ids)
	if err = queryCoinGecko(query, &pricefeeds); err != nil {
		return err
	}

	return db.SaveTokensPrice(pricefeeds, timestamp)
}

// QueryCoinGecko can query endpoint from pricefeed
func queryCoinGecko(endpoint string, ptr interface{}) error {
	url := "https://api.coingecko.com/api/v3" + endpoint
	if err := utils.QueryFromAPI(url, ptr); err != nil {
		return err
	}
	return nil
}
