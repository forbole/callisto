package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
)

func (suite *DbTestSuite) TestSaveTopAccountsBalance() {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	// Test saving balances
	amount := types.NewNativeTokenAmount(
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		sdk.NewInt(100),
		10,
	)

	err := suite.database.SaveTopAccountsBalance("available", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("delegation", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("redelegation", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("unbonding", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("reward", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.UpdateTopAccountsSum("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", "500")
	suite.Require().NoError(err)

	// Verify data
	expected := dbtypes.NewTopAccountsRow("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", 100, 100, 100, 100, 100, 500)

	var rows []dbtypes.TopAccountsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM top_accounts`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(expected.Equals(rows[0]))

	// Test saving higher values
	newAmount := types.NewNativeTokenAmount(
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		sdk.NewInt(200),
		10,
	)

	err = suite.database.SaveTopAccountsBalance("available", []types.NativeTokenAmount{newAmount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("delegation", []types.NativeTokenAmount{newAmount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("redelegation", []types.NativeTokenAmount{newAmount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("unbonding", []types.NativeTokenAmount{newAmount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("reward", []types.NativeTokenAmount{newAmount})
	suite.Require().NoError(err)

	err = suite.database.UpdateTopAccountsSum("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", "1000")
	suite.Require().NoError(err)

	// Verify data
	expected = dbtypes.NewTopAccountsRow("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", 200, 200, 200, 200, 200, 1000)
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM top_accounts`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(expected.Equals(rows[0]))

}

func (suite *DbTestSuite) TestGetAccountBalanceSum() {

}

func (suite *DbTestSuite) TestUpdateTopAccountsSum() {

}
