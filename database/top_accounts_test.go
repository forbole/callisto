package database_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/v4/database/types"
	"github.com/forbole/bdjuno/v4/types"
)

func (suite *DbTestSuite) TestSaveTopAccountsBalance() {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	// Test saving balances
	amount := types.NewNativeTokenAmount(
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		sdk.NewInt(100),
		100,
	)
	
	account := types.NewAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	err := suite.database.SaveAccounts([]types.Account{account})
	suite.Require().NoError(err)

	topAccount := types.NewTopAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", "/cosmos.auth.v1beta1.BaseAccount")
	err = suite.database.SaveTopAccounts([]types.TopAccount{topAccount}, 100)
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("available", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("delegation", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("redelegation", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("unbonding", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.SaveTopAccountsBalance("reward", []types.NativeTokenAmount{amount})
	suite.Require().NoError(err)

	err = suite.database.UpdateTopAccountsSum("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", "500", 100)
	suite.Require().NoError(err)

	// Verify data
	expected := dbtypes.NewTopAccountsRow("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", "cosmos.auth.v1beta1.BaseAccount", 100, 100, 100, 100, 100, 500, 100)

	var rows []dbtypes.TopAccountsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM top_accounts`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(expected.Equals(rows[0]))

	// Test saving higher values
	newAmount := types.NewNativeTokenAmount(
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		sdk.NewInt(200),
		300,
	)

	err = suite.database.SaveTopAccounts([]types.TopAccount{topAccount}, 200)
	suite.Require().NoError(err)

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

	err = suite.database.UpdateTopAccountsSum("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", "1000", 300)
	suite.Require().NoError(err)

	// Verify data
	expected = dbtypes.NewTopAccountsRow("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", "cosmos.auth.v1beta1.BaseAccount", 200, 200, 200, 200, 200, 1000, 300)
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM top_accounts`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(expected.Equals(rows[0]))

}

func (suite *DbTestSuite) TestGetAccountBalanceSum() {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	// Store balances
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

	// Verify Data
	expectedSum := "500"
	sum, err := suite.database.GetAccountBalanceSum("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	suite.Require().NoError(err)
	suite.Require().Equal(expectedSum, sum)

	// Verify getting 0 amount
	expectedSum = "0"
	sum, err = suite.database.GetAccountBalanceSum("")
	suite.Require().NoError(err)
	suite.Require().Equal(expectedSum, sum)
}

func (suite *DbTestSuite) TestUpdateTopAccountsSum() {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	// Store top accounts sum
	amount := "100"
	err := suite.database.UpdateTopAccountsSum("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", amount, 500)
	suite.Require().NoError(err)

	// Verify data
	var rows []string
	err = suite.database.Sqlx.Select(&rows, `SELECT sum FROM top_accounts`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(amount, rows[0])

	// Store different amount
	amount = "200"
	err = suite.database.UpdateTopAccountsSum("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs", amount, 500)
	suite.Require().NoError(err)

	// Verify data
	err = suite.database.Sqlx.Select(&rows, `SELECT sum FROM top_accounts`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(amount, rows[0])
}
