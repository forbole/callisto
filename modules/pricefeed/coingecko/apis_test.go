package coingecko_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/forbole/bdjuno/v4/modules/pricefeed/coingecko"
)

func TestConvertCoingeckoPrices(t *testing.T) {
	result := `
[
	{
	   "id":"cosmos",
	   "symbol":"atom",
	   "name":"Cosmos",
	   "image":"https://assets.coingecko.com/coins/images/1481/large/cosmos_hub.png?1555657960",
	   "current_price":31.16,
	   "market_cap":8809250407,
	   "market_cap_rank":21,
	   "fully_diluted_valuation":null,
	   "total_volume":2121320957,
	   "high_24h":36.58,
	   "low_24h":28.95,
	   "price_change_24h":2.0,
	   "price_change_percentage_24h":6.87504,
	   "market_cap_change_24h":702667080,
	   "market_cap_change_percentage_24h":8.66786,
	   "circulating_supply":279080458.802741,
	   "total_supply":null,
	   "max_supply":null,
	   "ath":36.58,
	   "ath_change_percentage":-13.6991,
	   "ath_date":"2021-09-13T00:15:37.723Z",
	   "atl":1.16,
	   "atl_change_percentage":2621.02411,
	   "atl_date":"2020-03-13T02:27:44.591Z",
	   "roi":{
		  "times":310.6319088856543,
		  "currency":"usd",
		  "percentage":31063.19088856543
	   },
	   "last_updated":"2021-09-13T08:48:15.930Z"
	},
	{
	   "id":"bitcanna",
	   "symbol":"bcna",
	   "name":"BitCanna",
	   "image":"https://assets.coingecko.com/coins/images/4716/large/bcna.png?1547040016",
	   "current_price":0.03998645,
	   "market_cap":0.0,
	   "market_cap_rank":null,
	   "fully_diluted_valuation":null,
	   "total_volume":186.04,
	   "high_24h":0.082037,
	   "low_24h":0.03733862,
	   "price_change_24h":-0.042050606106,
	   "price_change_percentage_24h":-51.25806,
	   "market_cap_change_24h":0.0,
	   "market_cap_change_percentage_24h":0.0,
	   "circulating_supply":0.0,
	   "total_supply":397427654.089788,
	   "max_supply":null,
	   "ath":0.922537,
	   "ath_change_percentage":-95.6656,
	   "ath_date":"2020-02-06T22:44:34.222Z",
	   "atl":0.00216099,
	   "atl_change_percentage":1750.37612,
	   "atl_date":"2020-01-01T11:24:08.916Z",
	   "roi":{
		  "times":-0.6667795841096327,
		  "currency":"usd",
		  "percentage":-66.67795841096327
	   },
	   "last_updated":"2021-09-13T07:08:54.093Z"
	},
	{
	   "id":"bitcoin",
	   "symbol":"btc",
	   "name":"Bitcoin",
	   "image":"https://assets.coingecko.com/coins/images/1/large/bitcoin.png?1547033579",
	   "current_price":44397,
	   "market_cap":836648999243,
	   "market_cap_rank":1,
	   "fully_diluted_valuation":933832354264,
	   "total_volume":30593030219,
	   "high_24h":46494,
	   "low_24h":44313,
	   "price_change_24h":-1735.169587249417,
	   "price_change_percentage_24h":-3.76128,
	   "market_cap_change_24h":-30612716244.69995,
	   "market_cap_change_percentage_24h":-3.52981,
	   "circulating_supply":18814543,
	   "total_supply":21000000,
	   "max_supply":21000000,
	   "ath":64805,
	   "ath_change_percentage":-30.97016,
	   "ath_date":"2021-04-14T11:54:46.763Z",
	   "atl":67.81,
	   "atl_change_percentage":65871.47051,
	   "atl_date":"2013-07-06T00:00:00.000Z",
	   "roi":null,
	   "last_updated":"2021-09-13T08:56:15.583Z"
	}
]
`

	var apisPrices []coingecko.MarketTicker
	err := json.Unmarshal([]byte(result), &apisPrices)
	require.NoError(t, err)

	prices := coingecko.ConvertCoingeckoPrices(apisPrices)
	require.Equal(t, int64(8809250407), prices[0].MarketCap)
	require.Equal(t, int64(0), prices[1].MarketCap)
	require.Equal(t, int64(836648999243), prices[2].MarketCap)
}
