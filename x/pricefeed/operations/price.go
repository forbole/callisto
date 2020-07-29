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
func UpdatePrice(cp client.ClientProxy, db database.BigDipperDb ) error {
	log.Debug().
		Str("module", "pricefeed").
		Str("operation", "pricefeed").
		Msg("getting token price and market cap")

	//get token names
	var coins api.Coins
	if err := utils.QueryCoinGecko("/coins/list", &coins); err != nil {
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

	timestamp := time.Now()
	//query
	var pricefeeds api.Pricefeeds
	query := fmt.Sprintf("/coins/markets?vs_currency=usd&ids=%s", ids)
	if err = utils.QueryCoinGecko(query, &pricefeeds); err != nil {
		return err
	}

	return db.SaveTokensPrice(markets, timestamp); err != nil {
	
}


// QueryCoinGecko can query endpoint from pricefeed
func QueryCoinGecko(endpoint string,ptr interface{})error{
	url:="https://api.pricefeed.com/api/v3"+endpoint
	if err:= queryFromApi(url,ptr);err!=nil{
		return err
	}
	return nil
}