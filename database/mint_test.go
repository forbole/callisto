package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveInflation() {
	timezone, err := time.LoadLocation("UTC")
	suite.Require().NoError(err)

	inflation := sdk.NewDecWithPrec(10050, 2)
	timestamp := time.Date(2020, 02, 02, 15, 00, 00, 00, timezone)

	err = suite.database.SaveInflation(inflation, 100, timestamp)
	suite.Require().NoError(err)

	err = suite.database.SaveInflation(inflation, 100, timestamp)
	suite.Require().NoError(err, "double inflation insertion should return no error")

	var value []dbtypes.InflationRow
	err = suite.database.Sqlx.Select(&value, `SELECT * FROM inflation_history`)
	suite.Require().NoError(err)
	suite.Require().Len(value, 1, "no duplicated inflation rows should be inserted")

	expected := dbtypes.NewInflationRow(100.50, 100, timestamp)
	suite.Require().True(expected.Equal(value[0]))
}
