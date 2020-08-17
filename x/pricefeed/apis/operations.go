package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/forbole/bdjuno/x/pricefeed/types"
)

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

// GetCoinsList allows to fetch from the remote APIs the list of all the supported tokens
func GetCoinsList() (coins types.Tokens, err error) {
	err = queryCoinGecko("/coins/list", &coins)
	return coins, err
}

// GetTokensPrices queries the remote APIs to get the token prices of all the tokens having the given ids
func GetTokensPrices(ids []string) (prices types.MarketTickers, err error) {
	query := fmt.Sprintf("/coins/markets?vs_currency=usd&ids=%s", strings.Join(ids, "&"))
	err = queryCoinGecko(query, &prices)
	return prices, err
}
