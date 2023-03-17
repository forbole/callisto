package database_test

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	creminttypes "github.com/crescent-network/crescent/v4/x/mint/types"

	"github.com/forbole/bdjuno/v4/types"

	dbtypes "github.com/forbole/bdjuno/v4/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveInflation() {

	// Save the data
	err := suite.database.SaveInflation(sdk.NewDecWithPrec(10050, 2).String(), 100)
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
	err = suite.database.SaveInflation(sdk.NewDecWithPrec(20000, 2).String(), 90)
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
	err = suite.database.SaveInflation(sdk.NewDecWithPrec(30000, 2).String(), 100)
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
	err = suite.database.SaveInflation(sdk.NewDecWithPrec(40000, 2).String(), 110)
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
	mintParams := creminttypes.Params{
		MintDenom:          "udaric",
		MintPoolAddress:    "cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		BlockTimeThreshold: 10 * time.Second,
		InflationSchedules: []creminttypes.InflationSchedule{
			{
				StartTime: time.Date(2022, time.Month(4), 13, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, time.Month(4), 13, 0, 0, 0, 0, time.UTC),
				Amount:    sdk.NewInt(108700000000000),
			},
			{
				StartTime: time.Date(2023, time.Month(4), 13, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2024, time.Month(4), 13, 0, 0, 0, 0, time.UTC),
				Amount:    sdk.NewInt(216100000000000),
			},
		},
	}

	err := suite.database.SaveMintParams(types.NewMintParams(mintParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.MintParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM mint_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var storedParams creminttypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &storedParams)
	suite.Require().NoError(err)
	suite.Require().Equal(mintParams, storedParams)
	suite.Require().Equal(int64(10), rows[0].Height)
}
