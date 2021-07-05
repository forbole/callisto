package database_test

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"time"
)

func (suite *DbTestSuite) TestSaveStakingParams() {
	params := stakingtypes.NewParams(
		time.Duration(259200000000000),
		200,
		7,
		10000,
		"uatom",
	)

	err := suite.database.SaveStakingParams(params)
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
	)))
}

func (suite *DbTestSuite) TestGetStakingParams() {
	_, err := suite.database.Sql.Exec(`
INSERT INTO staking_params (bond_denom, unbonding_time, max_entries, historical_entries, max_validators) 
VALUES ($1, $2, $3, $4, $5)`, "uatom", 259200000000000, 7, 10000, 200)
	suite.Require().NoError(err)

	params, err := suite.database.GetStakingParams()
	suite.Require().NoError(err)

	suite.Require().Equal(&stakingtypes.Params{
		UnbondingTime:     259200000000000,
		MaxValidators:     200,
		MaxEntries:        7,
		HistoricalEntries: 10000,
		BondDenom:         "uatom",
	}, params)
}
