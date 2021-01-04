package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveInflation() {
	inflation := sdk.NewDecWithPrec(10050, 2)

	err := suite.database.SaveInflation(inflation, 100)
	suite.Require().NoError(err)

	err = suite.database.SaveInflation(inflation, 100)
	suite.Require().NoError(err, "double inflation insertion should return no error")

	var value []dbtypes.InflationRow
	err = suite.database.Sqlx.Select(&value, `SELECT * FROM inflation_history`)
	suite.Require().NoError(err)
	suite.Require().Len(value, 1, "no duplicated inflation rows should be inserted")

	expected := dbtypes.NewInflationRow(100.50, 100)
	suite.Require().True(expected.Equal(value[0]))
}
