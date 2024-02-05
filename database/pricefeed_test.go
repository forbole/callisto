package database_test

import (
	"fmt"
	"time"

	"github.com/forbole/callisto/v4/types"

	dbtypes "github.com/forbole/callisto/v4/database/types"
)

func (suite *DbTestSuite) insertToken(name string) {
	query := fmt.Sprintf(
		`INSERT INTO token (name) VALUES ('%s')`, name)
	_, err := suite.database.SQL.Query(query)
	suite.Require().NoError(err)

	query = fmt.Sprintf(
		`INSERT INTO token_unit (token_name, denom, exponent, price_id) VALUES ('%[1]s', 'u%[1]s', 0, 'u%[1]s'), ('%[1]s', 'm%[1]s', 3, 'm%[1]s'), ('%[1]s', '%[1]s', 6, '%[1]s')`,
		name)
	_, err = suite.database.SQL.Query(query)
	suite.Require().NoError(err)
}

func (suite *DbTestSuite) Test_GetTokensPriceID() {
	suite.insertToken("desmos")
	suite.insertToken("daric")

	units, err := suite.database.GetTokensPriceID()
	suite.Require().NoError(err)

	var expected = []string{"udesmos", "mdesmos", "desmos", "udaric", "mdaric", "daric"}

	suite.Require().Len(units, len(expected))
	for _, name := range expected {
		suite.Require().Contains(units, name)
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveTokenPrice() {
	suite.insertToken("desmos")
	suite.insertToken("atom")

	// Save data
	tickers := []types.TokenPrice{
		types.NewTokenPrice(
			"desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		types.NewTokenPrice(
			"atom",
			200.01,
			20,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
	}
	err := suite.database.SaveTokensPrices(tickers)
	suite.Require().NoError(err)

	// Verify data
	expected := []dbtypes.TokenPriceRow{
		dbtypes.NewTokenPriceRow(
			"desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		dbtypes.NewTokenPriceRow(
			"atom",
			200.01,
			20,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
	}

	var rows []dbtypes.TokenPriceRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM token_price`)
	suite.Require().NoError(err)
	for i, row := range rows {
		suite.Require().True(expected[i].Equals(row))
	}

	// Update data
	tickers = []types.TokenPrice{
		types.NewTokenPrice(
			"desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		types.NewTokenPrice(
			"atom",
			1,
			20,
			time.Date(2020, 10, 10, 15, 05, 00, 000, time.UTC),
		),
	}
	err = suite.database.SaveTokensPrices(tickers)
	suite.Require().NoError(err)

	// Verify data
	expected = []dbtypes.TokenPriceRow{
		dbtypes.NewTokenPriceRow(
			"desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		dbtypes.NewTokenPriceRow(
			"atom",
			1,
			20,
			time.Date(2020, 10, 10, 15, 05, 00, 000, time.UTC),
		),
	}

	rows = []dbtypes.TokenPriceRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM token_price ORDER BY timestamp`)
	suite.Require().NoError(err)
	for i, row := range rows {
		suite.Require().True(expected[i].Equals(row))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveTokenPriceHistory() {
	suite.insertToken("desmos")
	suite.insertToken("atom")

	// Save data
	tickers := []types.TokenPrice{
		types.NewTokenPrice(
			"desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		types.NewTokenPrice(
			"desmos",
			200.01,
			20,
			time.Date(2020, 10, 10, 15, 02, 00, 000, time.UTC),
		),
		types.NewTokenPrice(
			"atom",
			1,
			20,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		types.NewTokenPrice(
			"atom",
			1,
			20,
			time.Date(2020, 10, 10, 15, 02, 00, 000, time.UTC),
		),
	}
	err := suite.database.SaveTokenPricesHistory(tickers)
	suite.Require().NoError(err)

	// Verify data
	expected := []dbtypes.TokenPriceRow{
		dbtypes.NewTokenPriceRow(
			"desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		dbtypes.NewTokenPriceRow(
			"desmos",
			200.01,
			20,
			time.Date(2020, 10, 10, 15, 02, 00, 000, time.UTC),
		),
		dbtypes.NewTokenPriceRow(
			"atom",
			1,
			20,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		dbtypes.NewTokenPriceRow(
			"atom",
			1,
			20,
			time.Date(2020, 10, 10, 15, 02, 00, 000, time.UTC),
		),
	}

	var rows []dbtypes.TokenPriceRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM token_price_history`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equals(row))
	}

	// Update data
	tickers = []types.TokenPrice{
		types.NewTokenPrice(
			"desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		types.NewTokenPrice(
			"desmos",
			300.01,
			20,
			time.Date(2020, 10, 10, 15, 02, 00, 000, time.UTC),
		),
		types.NewTokenPrice(
			"atom",
			1,
			20,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		types.NewTokenPrice(
			"atom",
			10,
			20,
			time.Date(2020, 10, 10, 15, 02, 00, 000, time.UTC),
		),
	}
	err = suite.database.SaveTokenPricesHistory(tickers)
	suite.Require().NoError(err)

	// Verify data
	expected = []dbtypes.TokenPriceRow{
		dbtypes.NewTokenPriceRow(
			"desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		dbtypes.NewTokenPriceRow(
			"atom",
			1,
			20,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		dbtypes.NewTokenPriceRow(
			"desmos",
			300.01,
			20,
			time.Date(2020, 10, 10, 15, 02, 00, 000, time.UTC),
		),

		dbtypes.NewTokenPriceRow(
			"atom",
			10,
			20,
			time.Date(2020, 10, 10, 15, 02, 00, 000, time.UTC),
		),
	}

	rows = []dbtypes.TokenPriceRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM token_price_history ORDER BY timestamp`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equals(row))
	}
}
