package bigdipper_test

import (
	"fmt"
	"time"

	dbtypes "github.com/forbole/bdjuno/database/types"

	"github.com/forbole/bdjuno/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	bddbtypes "github.com/forbole/bdjuno/database/bigdipper/types"
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
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			100,
		),
		types.NewDelegation(
			delegator1.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			50,
		),
		types.NewDelegation(
			delegator2.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(200)),
			101,
		),
	}
	err := suite.database.SaveDelegations(delegations)
	suite.Require().NoError(err, "inserting delegations should return no error")

	// ------------------------------
	// --- Verify the data
	// ------------------------------

	// Verify delegation rows
	var delRows []bddbtypes.DelegationRow
	err = suite.database.Sqlx.Select(&delRows, `SELECT * FROM delegation`)
	suite.Require().NoError(err)

	expectedDelRows := []bddbtypes.DelegationRow{
		bddbtypes.NewDelegationRow(
			delegator1.String(),
			validator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			100,
		),
		bddbtypes.NewDelegationRow(
			delegator1.String(),
			validator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			50,
		),
		bddbtypes.NewDelegationRow(
			delegator2.String(),
			validator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			101,
		),
	}

	suite.Require().Len(delRows, len(expectedDelRows))
	for index, delegation := range expectedDelRows {
		suite.Require().True(delegation.Equal(delRows[index]))
	}

	// ------------------------------
	// --- Update the data
	// ------------------------------

	delegations = []types.Delegation{
		types.NewDelegation(
			delegator1.String(),
			validator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(150)),
			80,
		),
		types.NewDelegation(
			delegator1.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(120)),
			50,
		),
		types.NewDelegation(
			delegator2.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(180)),
			102,
		),
	}
	err = suite.database.SaveDelegations(delegations)
	suite.Require().NoError(err, "updating delegations should return no error")

	// ------------------------------
	// --- Verify the data
	// ------------------------------

	delRows = []bddbtypes.DelegationRow{}
	err = suite.database.Sqlx.Select(&delRows, `SELECT * FROM delegation`)
	suite.Require().NoError(err)

	expectedDelRows = []bddbtypes.DelegationRow{
		bddbtypes.NewDelegationRow(
			delegator1.String(),
			validator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			100,
		),
		bddbtypes.NewDelegationRow(
			delegator1.String(),
			validator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(120))),
			50,
		),
		bddbtypes.NewDelegationRow(
			delegator2.String(),
			validator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(180))),
			102,
		),
	}

	suite.Require().Len(delRows, len(expectedDelRows))
	for index, delegation := range expectedDelRows {
		suite.Require().True(delegation.Equal(delRows[index]))
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestSaveRedelegations() {
	delegator1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	delegator2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	srcValidator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
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

	// Save the data

	reDelegations := []types.Redelegation{
		types.NewRedelegation(
			delegator1.String(),
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			time.Date(2021, 1, 1, 12, 00, 01, 000, time.UTC),
			10,
		),
		types.NewRedelegation(
			delegator1.String(),
			srcValidator1.GetOperator(),
			dstValidator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(120)),
			time.Date(2021, 1, 1, 12, 00, 02, 000, time.UTC),
			10,
		),
		types.NewRedelegation(
			delegator1.String(),
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			time.Date(2021, 2, 1, 12, 00, 03, 000, time.UTC),
			10,
		),
		types.NewRedelegation(
			delegator2.String(),
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(200)),
			time.Date(2021, 2, 1, 12, 00, 04, 000, time.UTC),
			10,
		),
	}
	err := suite.database.SaveRedelegations(reDelegations)
	suite.Require().NoError(err)

	// Verify the data

	var rows []bddbtypes.RedelegationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM redelegation ORDER BY completion_time`)
	suite.Require().NoError(err)

	expected := []bddbtypes.RedelegationRow{
		bddbtypes.NewRedelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr(),
			dstValidator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			time.Date(2021, 1, 1, 12, 00, 01, 000, time.UTC),
			10,
		),
		bddbtypes.NewRedelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr(),
			dstValidator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(120))),
			time.Date(2021, 1, 1, 12, 00, 02, 000, time.UTC),
			10,
		),
		bddbtypes.NewRedelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr(),
			dstValidator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			time.Date(2021, 2, 1, 12, 00, 03, 000, time.UTC),
			10,
		),
		bddbtypes.NewRedelegationRow(
			delegator2.String(),
			srcValidator1.GetConsAddr(),
			dstValidator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			time.Date(2021, 2, 1, 12, 00, 04, 000, time.UTC),
			10,
		),
	}

	suite.Require().Len(rows, len(expected))
	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]))
	}

	// Update the data

	reDelegations = []types.Redelegation{
		types.NewRedelegation(
			delegator1.String(),
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(120)),
			time.Date(2021, 1, 1, 12, 00, 01, 000, time.UTC),
			10,
		),
		types.NewRedelegation(
			delegator1.String(),
			srcValidator1.GetOperator(),
			dstValidator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(80)),
			time.Date(2021, 1, 1, 12, 00, 02, 000, time.UTC),
			9,
		),
		types.NewRedelegation(
			delegator1.String(),
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(90)),
			time.Date(2021, 2, 1, 12, 00, 03, 000, time.UTC),
			11,
		),
		types.NewRedelegation(
			delegator2.String(),
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(260)),
			time.Date(2021, 2, 1, 12, 00, 04, 000, time.UTC),
			15,
		),
	}
	err = suite.database.SaveRedelegations(reDelegations)
	suite.Require().NoError(err)

	// ------------------------------
	// --- Verify the data
	// ------------------------------
	rows = []bddbtypes.RedelegationRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM redelegation ORDER BY completion_time`)
	suite.Require().NoError(err, "updating rows should not return an error")

	expected = []bddbtypes.RedelegationRow{
		bddbtypes.NewRedelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr(),
			dstValidator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(120))),
			time.Date(2021, 1, 1, 12, 00, 01, 000, time.UTC),
			10,
		),
		bddbtypes.NewRedelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr(),
			dstValidator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(120))),
			time.Date(2021, 1, 1, 12, 00, 02, 000, time.UTC),
			10,
		),
		bddbtypes.NewRedelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr(),
			dstValidator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(90))),
			time.Date(2021, 2, 1, 12, 00, 03, 000, time.UTC),
			11,
		),
		bddbtypes.NewRedelegationRow(
			delegator2.String(),
			srcValidator1.GetConsAddr(),
			dstValidator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(260))),
			time.Date(2021, 2, 1, 12, 00, 04, 000, time.UTC),
			15,
		),
	}

	suite.Require().Len(rows, len(expected))
	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]))
	}
}

// --------------------------------------------------------------------------------------------------------------------

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

	// Save the data

	unbondingDelegations := []types.UnbondingDelegation{
		types.NewUnbondingDelegation(
			delegator1.String(),
			validator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			time.Date(2021, 1, 1, 12, 00, 01, 000, time.UTC),
			10,
		),
		types.NewUnbondingDelegation(
			delegator1.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(120)),
			time.Date(2021, 1, 1, 12, 00, 02, 000, time.UTC),
			10,
		),
		types.NewUnbondingDelegation(
			delegator2.String(),
			validator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			time.Date(2021, 2, 1, 12, 00, 03, 000, time.UTC),
			10,
		),
		types.NewUnbondingDelegation(
			delegator2.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(200)),
			time.Date(2021, 2, 1, 12, 00, 04, 000, time.UTC),
			10,
		),
	}
	err := suite.database.SaveUnbondingDelegations(unbondingDelegations)
	suite.Require().NoError(err)

	// Verify the data

	var rows []bddbtypes.UnbondingDelegationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM unbonding_delegation ORDER BY completion_timestamp`)
	suite.Require().NoError(err)

	expected := []bddbtypes.UnbondingDelegationRow{
		bddbtypes.NewUnbondingDelegationRow(
			delegator1.String(),
			validator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			time.Date(2021, 1, 1, 12, 00, 01, 000, time.UTC),
			10,
		),
		bddbtypes.NewUnbondingDelegationRow(
			delegator1.String(),
			validator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(120))),
			time.Date(2021, 1, 1, 12, 00, 02, 000, time.UTC),
			10,
		),
		bddbtypes.NewUnbondingDelegationRow(
			delegator2.String(),
			validator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			time.Date(2021, 2, 1, 12, 00, 03, 000, time.UTC),
			10,
		),
		bddbtypes.NewUnbondingDelegationRow(
			delegator2.String(),
			validator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			time.Date(2021, 2, 1, 12, 00, 04, 000, time.UTC),
			10,
		),
	}

	suite.Require().Len(rows, len(expected))
	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]), fmt.Sprintf("%d", index))
	}

	// Update the data

	unbondingDelegations = []types.UnbondingDelegation{
		types.NewUnbondingDelegation(
			delegator1.String(),
			validator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			time.Date(2021, 1, 1, 12, 00, 01, 000, time.UTC),
			10,
		),
		types.NewUnbondingDelegation(
			delegator1.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(130)),
			time.Date(2021, 1, 1, 12, 00, 02, 000, time.UTC),
			10,
		),
		types.NewUnbondingDelegation(
			delegator2.String(),
			validator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(90)),
			time.Date(2021, 2, 1, 12, 00, 03, 000, time.UTC),
			11,
		),
		types.NewUnbondingDelegation(
			delegator2.String(),
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(250)),
			time.Date(2021, 2, 1, 12, 00, 04, 000, time.UTC),
			9,
		),
	}
	err = suite.database.SaveUnbondingDelegations(unbondingDelegations)
	suite.Require().NoError(err)

	// Verify the data

	rows = []bddbtypes.UnbondingDelegationRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM unbonding_delegation ORDER BY completion_timestamp`)
	suite.Require().NoError(err, "updating rows should not return an error")

	expected = []bddbtypes.UnbondingDelegationRow{
		bddbtypes.NewUnbondingDelegationRow(
			delegator1.String(),
			validator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			time.Date(2021, 1, 1, 12, 00, 01, 000, time.UTC),
			10,
		),
		bddbtypes.NewUnbondingDelegationRow(
			delegator1.String(),
			validator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(130))),
			time.Date(2021, 1, 1, 12, 00, 02, 000, time.UTC),
			10,
		),
		bddbtypes.NewUnbondingDelegationRow(
			delegator2.String(),
			validator1.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(90))),
			time.Date(2021, 2, 1, 12, 00, 03, 000, time.UTC),
			11,
		),
		bddbtypes.NewUnbondingDelegationRow(
			delegator2.String(),
			validator2.GetConsAddr(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			time.Date(2021, 2, 1, 12, 00, 04, 000, time.UTC),
			10,
		),
	}

	suite.Require().Len(rows, len(expected))
	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]))
	}
}
