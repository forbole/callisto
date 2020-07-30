package database_test

import (
	"time"

	dbtypes "github.com/forbole/bdjuno/database/types"
	api "github.com/forbole/bdjuno/x/pricefeed/apiTypes"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveTokenPrice() {
	timestamp, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)

	pricefeed := []api.MarketTicker{
		api.NewMarketTicker(
			"udaric",
			100.01,
			10,
		),
		api.NewMarketTicker(
			"utopi",
			200.01,
			20,
		),
	}
	err = suite.database.SaveTokensPrice(pricefeed, timestamp)
	suite.Require().NoError(err)

	expected := []dbtypes.TokenPriceRow{
		dbtypes.NewTokenPriceRow("udaric",
			100.01,
			10,
			timestamp),
		dbtypes.NewTokenPriceRow("utopi",
			200.01,
			20,
			timestamp,
		),
	}
	var rows []dbtypes.TokenPriceRow
	err = suite.database.Sqlx.Select(&rows, `SELECT denom,price,market_cap,timestamp FROM token_price`)
	suite.Require().NoError(err)
	for i, row := range rows {
		suite.Require().True(expected[i].Equals(row))
	}
}
