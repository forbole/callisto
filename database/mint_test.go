package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveInflation() {

	// Save the data
	err := suite.database.SaveInflation(sdk.NewDecWithPrec(10050, 2), 100)
	suite.Require().NoError(err)

	// Verify the data
	var rows []dbtypes.InflationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM inflation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "no duplicated inflation rows should be inserted")

	expected := dbtypes.NewInflationRow(100.50, 100)
	suite.Require().True(expected.Equal(rows[0]))

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with lower height
	err = suite.database.SaveInflation(sdk.NewDecWithPrec(20000, 2), 90)
	suite.Require().NoError(err, "double inflation insertion should return no error")

	// Verify the data
	rows = []dbtypes.InflationRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM inflation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "no duplicated inflation rows should be inserted")

	expected = dbtypes.NewInflationRow(100.50, 100)
	suite.Require().True(expected.Equal(rows[0]), "data should not change with lower height")

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with same height
	err = suite.database.SaveInflation(sdk.NewDecWithPrec(30000, 2), 100)
	suite.Require().NoError(err, "double inflation insertion should return no error")

	// Verify the data
	rows = []dbtypes.InflationRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM inflation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "no duplicated inflation rows should be inserted")

	expected = dbtypes.NewInflationRow(300.00, 100)
	suite.Require().True(expected.Equal(rows[0]), "data should change with same height")

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with higher height
	err = suite.database.SaveInflation(sdk.NewDecWithPrec(40000, 2), 110)
	suite.Require().NoError(err, "double inflation insertion should return no error")

	// Verify the data
	rows = []dbtypes.InflationRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM inflation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "no duplicated inflation rows should be inserted")

	expected = dbtypes.NewInflationRow(400.00, 110)
	suite.Require().True(expected.Equal(rows[0]), "data should change with higher height")
}

func (suite *DbTestSuite) TestBigDipperDb_SaveMintParams() {
	err := suite.database.SaveMintParams(types.NewMintParams(
		"udaric",
		sdk.NewDecWithPrec(4, 1),
		sdk.NewDecWithPrec(8, 1),
		sdk.NewDecWithPrec(4, 1),
		sdk.NewDecWithPrec(8, 1),
		5006000,
		10,
	))
	suite.Require().NoError(err)

	var rows []dbtypes.MintParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM mint_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(dbtypes.MintParamsRow{
		OneRowID:            true,
		MintDenom:           "udaric",
		InflationRateChange: "0.400000000000000000",
		InflationMax:        "0.800000000000000000",
		InflationMin:        "0.400000000000000000",
		GoalBonded:          "0.800000000000000000",
		BlocksPerYear:       5006000,
		Height:              10,
	}, rows[0])
}
