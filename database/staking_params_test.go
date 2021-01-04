package database_test

import (
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/staking/types"
)

func (suite *DbTestSuite) TestSaveStakingParams() {
	params := types.NewStakingParams("uatom")

	err := suite.database.SaveStakingParams(params)
	suite.Require().NoError(err)

	var rows []dbtypes.StakingParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_params`)
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewStakingParamsRow("uatom")))
}

func (suite *DbTestSuite) TestGetStakingParams() {
	_, err := suite.database.Sql.Exec(`INSERT INTO staking_params (bond_denom) VALUES ($1)`, "uatom")
	suite.Require().NoError(err)

	params, err := suite.database.GetStakingParams()
	suite.Require().NoError(err)

	suite.Require().Equal("uatom", params.BondName)
}
