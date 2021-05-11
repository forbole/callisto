package coingecko

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/forbole/bdjuno/types"
)

// GetCoinsList allows to fetch from the remote APIs the list of all the supported tokens
func GetCoinsList() (coins Tokens, err error) {
	err = queryCoinGecko("/coins/list", &coins)
	return coins, err
}

// GetTokensPrices queries the remote APIs to get the token prices of all the tokens having the given ids
func GetTokensPrices(ids []string) ([]types.TokenPrice, error) {
	var prices []MarketTicker
	query := fmt.Sprintf("/coins/markets?vs_currency=usd&ids=%s", strings.Join(ids, "&"))
	err := queryCoinGecko(query, &prices)
	if err != nil {
		return nil, err
	}

	return convertCoingeckoPrices(prices), nil
}

func convertCoingeckoPrices(prices []MarketTicker) []types.TokenPrice {
	tokenPrices := make([]types.TokenPrice, len(prices))
	for i, price := range prices {
		tokenPrices[i] = types.NewTokenPrice(
			price.Symbol,
			price.CurrentPrice,
			price.MarketCap,
			price.LastUpdated,
		)
	}
	return tokenPrices
}

// queryCoinGecko queries the CoinGecko APIs for the given endpoint
func queryCoinGecko(endpoint string, ptr interface{}) error {
	resp, err := http.Get("https://api.coingecko.com/api/v3" + endpoint)
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
