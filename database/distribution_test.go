package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveCommunityPool() {
	coins := sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(100)))
	err := suite.database.SaveCommunityPool(coins, 10)
	suite.Require().NoError(err)

	expected := dbtypes.NewCommunityPoolRow(dbtypes.NewDbDecCoins(coins), 10)
	var rows []dbtypes.CommunityPoolRow
	err = suite.database.Sqlx.Select(&rows, `SELECT coins,height FROM community_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "community_pool table should contain only one row")

	suite.Require().True(expected.Equals(rows[0]))

}
