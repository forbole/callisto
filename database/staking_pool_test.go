package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveStakingPool() {
	height := int64(100)

	timestamp, err := time.Parse(time.RFC3339, "2020-02-02T15:00:00Z")
	suite.Require().NoError(err)

	pool := stakingtypes.NewPool(sdk.NewInt(100), sdk.NewInt(50))

	// Save the data
	err = suite.database.SaveStakingPool(pool, height, timestamp)
	suite.Require().NoError(err)

	var count int
	err = suite.database.Sqlx.QueryRow(`SELECT COUNT(*) FROM staking_pool_history`).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(1, count, "inserting a single staking pool row should return 1")

	// Perform a double insertion
	err = suite.database.SaveStakingPool(pool, height, timestamp)
	suite.Require().NoError(err)

	err = suite.database.Sqlx.QueryRow(`SELECT COUNT(*) FROM staking_pool_history`).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(1, count, "double inserting the same staking pool should return 1 row")

	// Verify the data
	var rows []dbtypes.StakingPoolRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_pool_history`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewStakingPoolRow(
		50,
		100,
		height,
		timestamp,
	)))
}
