package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveStakingPool() {
	pool := stakingtypes.NewPool(sdk.NewInt(100), sdk.NewInt(50))
	err := suite.database.SaveStakingPool(1000, pool)
	suite.Require().NoError(err)
}
