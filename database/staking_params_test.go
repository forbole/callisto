package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
)

func (suite *DbTestSuite) TestSaveStakingParams() {
	err := suite.database.SaveStakingParams(types.NewStakingParams(
		stakingtypes.NewParams(
			time.Duration(259200000000000),
			200,
			7,
			10000,
			"uatom",
			sdk.NewDec(10),
		),
		10,
	))
	suite.Require().NoError(err)

	var rows []dbtypes.StakingParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_params`)
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewStakingParamsRow(
		"uatom",
		259200000000000,
		7,
		10000,
		200,
		"10.000000000000000000",
		10,
	)))
}

func (suite *DbTestSuite) TestGetStakingParams() {
	_, err := suite.database.Sql.Exec(`
INSERT INTO staking_params (bond_denom, unbonding_time, max_entries, historical_entries, max_validators, min_commission_rate, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7)`, "uatom", 259200000000000, 7, 10000, 200, "10", 10)
	suite.Require().NoError(err)

	params, err := suite.database.GetStakingParams()
	suite.Require().NoError(err)

	suite.Require().Equal(&types.StakingParams{
		Params: stakingtypes.NewParams(
			259200000000000,
			200,
			7,
			10000,
			"uatom",
			sdk.NewDec(10),
		),
		Height: 10,
	}, params)
}
