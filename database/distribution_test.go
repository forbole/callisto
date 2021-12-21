package database_test

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"

	"github.com/forbole/bdjuno/v2/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	bddbtypes "github.com/forbole/bdjuno/v2/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveCommunityPool() {
	// Save data
	original := sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(100)))
	err := suite.database.SaveCommunityPool(original, 10)
	suite.Require().NoError(err)

	// Verify data
	expected := bddbtypes.NewCommunityPoolRow(dbtypes.NewDbDecCoins(original), 10)
	var rows []bddbtypes.CommunityPoolRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM community_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "community_pool table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]))

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with lower height
	coins := sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(50)))
	err = suite.database.SaveCommunityPool(coins, 5)
	suite.Require().NoError(err)

	// Verify data
	expected = bddbtypes.NewCommunityPoolRow(dbtypes.NewDbDecCoins(original), 10)
	rows = []bddbtypes.CommunityPoolRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM community_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "community_pool table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]), "updating with lower height should not modify the data")

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with equal height
	coins = sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(120)))
	err = suite.database.SaveCommunityPool(coins, 10)
	suite.Require().NoError(err)

	// Verify data
	expected = bddbtypes.NewCommunityPoolRow(dbtypes.NewDbDecCoins(coins), 10)
	rows = []bddbtypes.CommunityPoolRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM community_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "community_pool table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]), "updating with same height should modify the data")

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with higher height
	coins = sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(200)))
	err = suite.database.SaveCommunityPool(coins, 11)
	suite.Require().NoError(err)

	// Verify data
	expected = bddbtypes.NewCommunityPoolRow(dbtypes.NewDbDecCoins(coins), 11)
	rows = []bddbtypes.CommunityPoolRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM community_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "community_pool table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]), "updating with higher height should modify the data")
}

func (suite *DbTestSuite) TestBigDipperDb_SaveDistributionParams() {
	distrParams := distrtypes.Params{
		CommunityTax:        sdk.NewDecWithPrec(2, 2),
		BaseProposerReward:  sdk.NewDecWithPrec(1, 2),
		BonusProposerReward: sdk.NewDecWithPrec(4, 2),
		WithdrawAddrEnabled: true,
	}
	err := suite.database.SaveDistributionParams(types.NewDistributionParams(distrParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.DistributionParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM distribution_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var stored distrtypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &stored)
	suite.Require().NoError(err)
	suite.Require().Equal(distrParams, stored)
	suite.Require().Equal(int64(10), rows[0].Height)
}

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorCommissionAmount() {
	_ = suite.getBlock(9)
	_ = suite.getBlock(10)
	_ = suite.getBlock(11)

	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	// Save the data
	original := sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(100)))
	amount := types.NewValidatorCommissionAmount(
		validator.GetOperator(),
		validator.GetSelfDelegateAddress(),
		original,
		10,
	)
	err := suite.database.SaveValidatorCommissionAmount(amount)
	suite.Require().NoError(err)

	// Verify the data
	originalRow := bddbtypes.NewValidatorCommissionAmountRow(validator.GetConsAddr(), dbtypes.NewDbDecCoins(original), 10)
	var rows []bddbtypes.ValidatorCommissionAmountRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission_amount`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(originalRow, rows[0])

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with lower height
	coins := sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(120)))
	amount = types.NewValidatorCommissionAmount(validator.GetOperator(), validator.GetSelfDelegateAddress(), coins, 9)
	err = suite.database.SaveValidatorCommissionAmount(amount)
	suite.Require().NoError(err)

	// Verify the data
	rows = []bddbtypes.ValidatorCommissionAmountRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission_amount`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(originalRow, rows[0], "updating with lower height should not modify the data")

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with same height
	coins = sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(200)))
	amount = types.NewValidatorCommissionAmount(validator.GetOperator(), validator.GetSelfDelegateAddress(), coins, 10)
	err = suite.database.SaveValidatorCommissionAmount(amount)
	suite.Require().NoError(err)

	// Verify the data
	expected := bddbtypes.NewValidatorCommissionAmountRow(validator.GetConsAddr(), dbtypes.NewDbDecCoins(coins), 10)
	rows = []bddbtypes.ValidatorCommissionAmountRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission_amount`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(expected, rows[0], "updating with same height should modify the data")

	// ---------------------------------------------------------------------------------------------------------------

	// Try updating with higher height
	coins = sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(500)))
	amount = types.NewValidatorCommissionAmount(validator.GetOperator(), validator.GetSelfDelegateAddress(), coins, 11)
	err = suite.database.SaveValidatorCommissionAmount(amount)
	suite.Require().NoError(err)

	// Verify the data
	expected = bddbtypes.NewValidatorCommissionAmountRow(validator.GetConsAddr(), dbtypes.NewDbDecCoins(coins), 11)
	rows = []bddbtypes.ValidatorCommissionAmountRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission_amount`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(expected, rows[0], "updating with higher height should modify the data")
}

func (suite *DbTestSuite) TestBigDipperDb_SaveDelegatorsRewardsAmounts() {
	_ = suite.getBlock(9)
	_ = suite.getBlock(10)
	_ = suite.getBlock(11)
	_ = suite.getBlock(12)
	_ = suite.getBlock(13)

	delegator := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	validator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	validator2 := suite.getValidator(
		"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
		"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
		"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
	)

	// Save the data
	rewards := []types.DelegatorRewardAmount{
		types.NewDelegatorRewardAmount(
			validator1.GetOperator(),
			delegator.String(),
			sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(100))),
		),
		types.NewDelegatorRewardAmount(
			validator2.GetOperator(),
			delegator.String(),
			sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(200))),
		),
	}
	err := suite.database.SaveDelegatorsRewardsAmounts(10, delegator.String(), rewards)
	suite.Require().NoError(err)

	// Verify the data
	expected := []bddbtypes.DelegationRewardRow{
		bddbtypes.NewDelegationRewardRow(
			delegator.String(),
			validator1.GetConsAddr(),
			delegator.String(),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(100)))),
			10,
		),
		bddbtypes.NewDelegationRewardRow(
			delegator.String(),
			validator2.GetConsAddr(),
			delegator.String(),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(200)))),
			10,
		),
	}

	var rows []bddbtypes.DelegationRewardRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM delegation_reward ORDER BY height`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, row := range rows {
		suite.Require().True(row.Equals(expected[index]))
	}

	// -------------------------------------------------------------------------------------------------------------------

	// Update the data (older height)
	rewards = []types.DelegatorRewardAmount{
		types.NewDelegatorRewardAmount(
			validator1.GetOperator(),
			delegator.String(),
			sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(120))),
		),
	}
	err = suite.database.SaveDelegatorsRewardsAmounts(9, delegator.String(), rewards)
	suite.Require().NoError(err)

	// Verify the data
	expected = []bddbtypes.DelegationRewardRow{
		bddbtypes.NewDelegationRewardRow(
			delegator.String(),
			validator1.GetConsAddr(),
			delegator.String(),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(100)))),
			10,
		),
		bddbtypes.NewDelegationRewardRow(
			delegator.String(),
			validator2.GetConsAddr(),
			delegator.String(),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(200)))),
			10,
		),
	}

	rows = []bddbtypes.DelegationRewardRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM delegation_reward ORDER BY height`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, row := range rows {
		suite.Require().True(row.Equals(expected[index]))
	}

	// -------------------------------------------------------------------------------------------------------------------

	// Update the data (same height)
	rewards = []types.DelegatorRewardAmount{
		types.NewDelegatorRewardAmount(
			validator1.GetOperator(),
			delegator.String(),
			sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(120))),
		),
	}
	err = suite.database.SaveDelegatorsRewardsAmounts(10, delegator.String(), rewards)
	suite.Require().NoError(err)

	// Verify the data
	expected = []bddbtypes.DelegationRewardRow{
		bddbtypes.NewDelegationRewardRow(
			delegator.String(),
			validator1.GetConsAddr(),
			delegator.String(),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(120)))),
			10,
		),
	}

	rows = []bddbtypes.DelegationRewardRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM delegation_reward ORDER BY height`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, row := range rows {
		suite.Require().True(row.Equals(expected[index]))
	}

	// -------------------------------------------------------------------------------------------------------------------

	// Update the data (new height)
	rewards = []types.DelegatorRewardAmount{
		types.NewDelegatorRewardAmount(
			validator1.GetOperator(),
			delegator.String(),
			sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(500))),
		),
	}
	err = suite.database.SaveDelegatorsRewardsAmounts(11, delegator.String(), rewards)
	suite.Require().NoError(err)

	// Verify the data
	expected = []bddbtypes.DelegationRewardRow{
		bddbtypes.NewDelegationRewardRow(
			delegator.String(),
			validator1.GetConsAddr(),
			delegator.String(),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("cosmos", sdk.NewInt(500)))),
			11,
		),
	}

	rows = []bddbtypes.DelegationRewardRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM delegation_reward ORDER BY height`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, row := range rows {
		suite.Require().True(row.Equals(expected[index]))
	}
}
