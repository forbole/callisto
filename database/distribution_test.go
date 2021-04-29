package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/x/distribution/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveCommunityPool() {
	err := suite.database.SaveCommunityPool(sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(10000))))
	suite.Require().NoError(err)

	coins := sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(100)))
	err = suite.database.SaveCommunityPool(coins)
	suite.Require().NoError(err, "updating community pool should return no error")

	expected := dbtypes.NewCommunityPoolRow(dbtypes.NewDbDecCoins(coins))

	var rows []dbtypes.CommunityPoolRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM community_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "community_pool table should contain only one row")

	suite.Require().True(expected.Equals(rows[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorCommissionAmount() {
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	// Store the value
	err := suite.database.SaveValidatorCommissionAmount(types.NewValidatorCommissionAmount(
		validator.GetConsAddr(),
		sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(10000))),
	))
	suite.Require().NoError(err)

	// Update the value
	amount := sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(20000)))
	err = suite.database.SaveValidatorCommissionAmount(types.NewValidatorCommissionAmount(
		validator.GetConsAddr(),
		amount,
	))
	suite.Require().NoError(err, "updating validator commission amount should return no error")

	var rows []dbtypes.ValidatorCommissionAmountRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission_amount`)
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorCommissionAmountRow{
		dbtypes.NewValidatorCommissionAmountRow(
			validator.GetConsAddr(),
			dbtypes.NewDbDecCoins(amount),
		),
	}
	suite.Require().Len(rows, len(expected))

	for index, row := range rows {
		suite.Require().True(row.Equals(expected[index]))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveDelegatorsRewardsAmounts() {
	delegator := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	// Insert value
	err := suite.database.SaveDelegatorsRewardsAmounts([]types.DelegatorReward{
		types.NewDelegatorRewardAmount(
			validator.GetConsAddr(),
			delegator.String(),
			delegator.String(),
			sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(1000))),
		),
	})
	suite.Require().NoError(err)

	// Update value
	amount := sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(21000)))
	err = suite.database.SaveDelegatorsRewardsAmounts([]types.DelegatorReward{
		types.NewDelegatorRewardAmount(
			validator.GetConsAddr(),
			delegator.String(),
			delegator.String(),
			amount,
		),
	})
	suite.Require().NoError(err)

	var rows []dbtypes.DelegatorRewardRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM delegation_reward`)
	suite.Require().NoError(err)

	expected := []dbtypes.DelegatorRewardRow{
		dbtypes.NewDelegatorRewardRow(
			validator.GetConsAddr(),
			delegator.String(),
			delegator.String(),
			dbtypes.NewDbDecCoins(amount),
		),
	}
	suite.Require().Len(rows, len(expected))

	for index, row := range rows {
		suite.Require().True(row.Equals(expected[index]))
	}
}
