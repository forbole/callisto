package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/staking/types"
)

func (suite *DbTestSuite) TestDelegations() {
	delegator1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	delegator2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
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

	// ------------------------------
	// --- Save the data
	// ------------------------------

	delegations := []types.Delegation{
		types.NewDelegation(
			delegator1.String(),
			validator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			"1000",
			1000,
		),
		types.NewDelegation(
			delegator1.String(),
			validator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			"1000",
			1000,
		),
		types.NewDelegation(
			delegator1.String(),
			validator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			"1001",
			1001,
		),
		types.NewDelegation(
			delegator2.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(200)),
			"1500",
			1500,
		),
	}
	err := suite.database.SaveDelegations(delegations)
	suite.Require().NoError(err, "inserting delegations should return no error")

	// ------------------------------
	// --- Verify the data
	// ------------------------------

	// Verify delegation rows
	var delRows []dbtypes.DelegationRow
	err = suite.database.Sqlx.Select(&delRows, `SELECT * FROM delegation`)
	suite.Require().NoError(err)

	expectedDelRows := []dbtypes.DelegationRow{
		dbtypes.NewDelegationRow(
			validator1.GetConsAddr(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			1000,
			1000,
		),
		dbtypes.NewDelegationRow(
			validator1.GetConsAddr(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			1000,
			1000,
		),
		dbtypes.NewDelegationRow(
			validator1.GetConsAddr(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			1001,
			1001,
		),
		dbtypes.NewDelegationRow(
			validator2.GetConsAddr(),
			delegator2.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			1500,
			1500,
		),
	}

	suite.Require().Len(delRows, len(expectedDelRows))
	for index, delegation := range expectedDelRows {
		suite.Require().True(delegation.Equal(delRows[index]))
	}
}

// ________________________________________________

func (suite *DbTestSuite) TestSaveUnbondingDelegations() {
	delegator1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	delegator2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
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

	completionTimestamp1, err := time.Parse(time.RFC3339, "2020-08-10T16:00:00Z")
	suite.Require().NoError(err)

	completionTimestamp2, err := time.Parse(time.RFC3339, "2020-08-20T16:00:00Z")
	suite.Require().NoError(err)

	// ------------------------------
	// --- Save the data
	// ------------------------------

	delegations := []types.UnbondingDelegation{
		types.NewUnbondingDelegation(
			delegator1.String(),
			validator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			completionTimestamp1,
			1000,
		),
		types.NewUnbondingDelegation(
			delegator1.String(),
			validator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			completionTimestamp1,
			1000,
		),
		types.NewUnbondingDelegation(
			delegator1.String(),
			validator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			completionTimestamp2,
			1001,
		),
		types.NewUnbondingDelegation(
			delegator2.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(200)),
			completionTimestamp2,
			1500,
		),
	}
	err = suite.database.SaveUnbondingDelegations(delegations)
	suite.Require().NoError(err)

	// ------------------------------
	// --- Verify the data
	// ------------------------------

	var rows []dbtypes.UnbondingDelegationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM unbonding_delegation`)
	suite.Require().NoError(err)

	expected := []dbtypes.UnbondingDelegationRow{
		dbtypes.NewUnbondingDelegationRow(
			validator1.GetConsAddr(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			completionTimestamp1,
			1000,
		),
		dbtypes.NewUnbondingDelegationRow(
			validator1.GetConsAddr(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			completionTimestamp1,
			1000,
		),
		dbtypes.NewUnbondingDelegationRow(
			validator1.GetConsAddr(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			completionTimestamp2,
			1001,
		),
		dbtypes.NewUnbondingDelegationRow(
			validator2.GetConsAddr(),
			delegator2.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			completionTimestamp2,
			1500,
		),
	}

	suite.Require().Len(rows, len(expected))
	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]), "errored index: %d", index)
	}
}

// ________________________________________________

func (suite *DbTestSuite) TestSaveRedelegations() {
	delegator1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	delegator2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	srcValidator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	srcValidator2 := suite.getValidator(
		"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
		"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
		"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
	)
	dstValidator1 := suite.getValidator(
		"cosmosvalcons1px0zkz2cxvc6lh34uhafveea9jnaagckmrlsye",
		"cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn",
		"cosmosvalconspub1zcjduepq0dc9apn3pz2x2qyujcnl2heqq4aceput2uaucuvhrjts75q0rv5smjjn7v",
	)
	dstValidator2 := suite.getValidator(
		"cosmosvalcons1rtst6se0nfgjy362v33jt5d05crgdyhfvvvvay",
		"cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr",
		"cosmosvalconspub1zcjduepq5e8w7t7k9pwfewgrwy8vn6cghk0x49chx64vt0054yl4wwsmjgrqfackxm",
	)

	completionTimestamp1, err := time.Parse(time.RFC3339, "2020-08-10T16:00:00Z")
	suite.Require().NoError(err)

	completionTimestamp2, err := time.Parse(time.RFC3339, "2020-08-20T16:00:00Z")
	suite.Require().NoError(err)

	// ------------------------------
	// --- Save the data
	// ------------------------------

	reDelegations := []types.Redelegation{
		types.NewRedelegation(
			delegator1.String(),
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			completionTimestamp1,
			1000,
		),
		types.NewRedelegation(
			delegator1.String(),
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			completionTimestamp1,
			1000,
		),
		types.NewRedelegation(
			delegator1.String(),
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			completionTimestamp2,
			1001,
		),
		types.NewRedelegation(
			delegator2.String(),
			srcValidator2.GetOperator(),
			dstValidator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(200)),
			completionTimestamp2,
			1500,
		),
	}
	err = suite.database.SaveRedelegations(reDelegations)
	suite.Require().NoError(err)

	// ------------------------------
	// --- Verify the data
	// ------------------------------
	var rows []dbtypes.ReDelegationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM redelegation`)
	suite.Require().NoError(err)

	expected := []dbtypes.ReDelegationRow{
		dbtypes.NewReDelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr(),
			dstValidator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			completionTimestamp1,
			1000,
		),
		dbtypes.NewReDelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr(),
			dstValidator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			completionTimestamp1,
			1000,
		),
		dbtypes.NewReDelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr(),
			dstValidator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			completionTimestamp2,
			1001,
		),
		dbtypes.NewReDelegationRow(
			delegator2.String(),
			srcValidator2.GetConsAddr(),
			dstValidator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			completionTimestamp2,
			1500,
		),
	}

	suite.Require().Len(rows, len(expected))
	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]))
	}
}
