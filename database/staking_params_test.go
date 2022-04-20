package database_test

import (
	"encoding/json"
	"time"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
)

func (suite *DbTestSuite) TestSaveStakingParams() {
	stakingParams := stakingtypes.NewParams(
		time.Duration(259200000000000),
		200,
		7,
		10000,
		"uatom",
	)
	err := suite.database.SaveStakingParams(types.NewStakingParams(stakingParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.StakingParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_params`)
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)

	var stored stakingtypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &stored)
	suite.Require().NoError(err)
	suite.Require().Equal(stakingParams, stored)
}

func (suite *DbTestSuite) TestGetStakingParams() {
	stakingParams := stakingtypes.NewParams(
		259200000000000,
		200,
		7,
		10000,
		"uatom",
	)

	paramsBz, err := json.Marshal(&stakingParams)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(
		`INSERT INTO staking_params (params, height) VALUES ($1, $2)`,
		string(paramsBz), 10,
	)
	suite.Require().NoError(err)

	params, err := suite.database.GetStakingParams()
	suite.Require().NoError(err)

	suite.Require().Equal(&types.StakingParams{
		Params: stakingParams,
		Height: 10,
	}, params)
}
