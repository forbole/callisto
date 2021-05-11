package bigdipper_test

import (
	types2 "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
)

func (suite *DbTestSuite) TestSaveStakingParams() {
	params := types.NewStakingParams("uatom")

	err := suite.database.SaveStakingParams(params)
	suite.Require().NoError(err)

	var rows []types2.StakingParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_params`)
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(types2.NewStakingParamsRow("uatom")))
}

func (suite *DbTestSuite) TestGetStakingParams() {
	_, err := suite.database.Sql.Exec(`INSERT INTO staking_params (bond_denom) VALUES ($1)`, "uatom")
	suite.Require().NoError(err)

	params, err := suite.database.GetStakingParams()
	suite.Require().NoError(err)

	suite.Require().Equal("uatom", params.BondName)
}
