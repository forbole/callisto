package database_test

import (
	dbtypes "github.com/forbole/bdjuno/database/types"
	api "github.com/forbole/bdjuno/x/supply/coinGeckoTypes"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveTokenPrice() {
	market := []api.Market{
		api.NewMarket(
			"udaric",
			100.01,
			10,
		),
		api.NewMarket(
			"utopi",
			200.01,
			20,
		),
	}
	err := suite.database.SaveTokensPrice(market, 30)
	suite.Require().NoError(err)

	expected := []dbtypes.TokenPriceRow{
		dbtypes.NewTokenPriceRow("udaric",
			100.01,
			10,
			30),
		dbtypes.NewTokenPriceRow("utopi",
			200.01,
			20,
			30,
		),
	}
	var rows []dbtypes.TokenPriceRow
	err = suite.database.Sqlx.Select(&rows, `SELECT denom,current_price,market_cap,height FROM token_values`)
	suite.Require().NoError(err)
	for i, row := range rows {
		suite.Require().True(expected[i].Equals(row))
	}
}
