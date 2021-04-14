package coingecko

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	pricefeedtypes "github.com/forbole/bdjuno/x/pricefeed/types"
)

// GetCoinsList allows to fetch from the remote APIs the list of all the supported tokens
func GetCoinsList() (coins Tokens, err error) {
	err = queryCoinGecko("/coins/list", &coins)
	return coins, err
}

// GetTokensPrices queries the remote APIs to get the token prices of all the tokens having the given ids
func GetTokensPrices(ids []string) ([]pricefeedtypes.TokenPrice, error) {
	var prices []MarketTicker
	query := fmt.Sprintf("/coins/markets?vs_currency=usd&ids=%s", strings.Join(ids, "&"))
	err := queryCoinGecko(query, &prices)
	if err != nil {
		return nil, err
	}

	return convertCoingeckoPrices(prices), nil
}

func convertCoingeckoPrices(prices []MarketTicker) []pricefeedtypes.TokenPrice {
	tokenPrices := make([]pricefeedtypes.TokenPrice, len(prices))
	for i, price := range prices {
		tokenPrices[i] = pricefeedtypes.NewTokenPrice(
			price.Symbol,
			price.CurrentPrice,
			price.MarketCap,
			price.LastUpdated,
		)
	}
	return tokenPrices
}

// QueryCoinGecko can query endpoint from pricefeed
func queryCoinGecko(endpoint string, ptr interface{}) error {
	resp, err := http.Get(fmt.Sprintf("%s/%s", "https://api.coingecko.com/api/v3", endpoint))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bz, &ptr); err != nil {
		return err
	}
	return nil
}
