package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestSaveAccountBalance() {
	address := suite.getAccount("cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7")
	height := int64(100)
	coins := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
		sdk.NewCoin("uatom", sdk.NewInt(15)),
	)

	// Save the balance
	err := suite.database.SaveAccountBalance(address.String(), coins, height)
	suite.Require().NoError(err)

	// Current balances
	var balRows []dbtypes.AccountBalanceRow
	err = suite.database.Sqlx.Select(&balRows, `SELECT * FROM account_balance ORDER BY address`)
	suite.Require().NoError(err)
	suite.Require().Len(balRows, 1)
	suite.Require().True(balRows[0].Equal(dbtypes.NewAccountBalanceRow(
		"cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7",
		dbtypes.NewDbCoins(coins),
	)))

	// Balance histories
	var balHisRows []dbtypes.AccountBalanceHistoryRow
	err = suite.database.Sqlx.Select(&balHisRows, `SELECT * FROM account_balance_history ORDER BY address`)
	suite.Require().NoError(err)
	suite.Require().Len(balHisRows, 1)
	suite.Require().True(balHisRows[0].Equal(dbtypes.NewAccountBalanceHistoryRow(
		"cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7",
		dbtypes.NewDbCoins(coins),
		height,
	)))
}
