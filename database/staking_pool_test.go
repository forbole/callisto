package database_test

import (
	"cosmossdk.io/math"

	dbtypes "github.com/forbole/callisto/v4/database/types"
	"github.com/forbole/callisto/v4/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveStakingPool() {
	// Save the data
	original := types.NewPool(math.NewInt(50), math.NewInt(100), math.NewInt(5), math.NewInt(1), 10)
	err := suite.database.SaveStakingPool(original)
	suite.Require().NoError(err)

	// Verify the data
	expected := dbtypes.NewStakingPoolRow(50, 100, 5, 1, 10)

	var rows []dbtypes.StakingPoolRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating using a lower height
	pool := types.NewPool(math.NewInt(1), math.NewInt(1), math.NewInt(1), math.NewInt(1), 8)
	err = suite.database.SaveStakingPool(pool)
	suite.Require().NoError(err)

	// Verify the data
	rows = []dbtypes.StakingPoolRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with a lower height should not modify the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with the same height
	pool = types.NewPool(math.NewInt(1), math.NewInt(1), math.NewInt(1), math.NewInt(1), 10)
	err = suite.database.SaveStakingPool(pool)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewStakingPoolRow(1, 1, 1, 1, 10)

	rows = []dbtypes.StakingPoolRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with a lower height should not modify the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a higher height
	pool = types.NewPool(math.NewInt(1000000), math.NewInt(1000000), math.NewInt(20), math.NewInt(15), 20)
	err = suite.database.SaveStakingPool(pool)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewStakingPoolRow(1000000, 1000000, 20, 15, 20)

	rows = []dbtypes.StakingPoolRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with a lower height should not modify the data")
}
