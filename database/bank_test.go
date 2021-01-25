package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	bbanktypes "github.com/forbole/bdjuno/x/bank/types"

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
	err := suite.database.SaveAccountBalances([]bbanktypes.AccountBalance{
		bbanktypes.NewAccountBalance(address.String(), coins, height),
	})
	suite.Require().NoError(err)

	// Current balances
	var balRows []dbtypes.AccountBalanceRow
	err = suite.database.Sqlx.Select(&balRows, `SELECT * FROM account_balance ORDER BY address`)
	suite.Require().NoError(err)
	suite.Require().Len(balRows, 1)
	suite.Require().True(balRows[0].Equal(dbtypes.NewAccountBalanceRow(
		"cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7",
		dbtypes.NewDbCoins(coins),
		height,
	)))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveTotalTokens() {
	coins := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
		sdk.NewCoin("uatom", sdk.NewInt(15)),
	)
	err := suite.database.SaveSupplyToken(coins, 10)
	suite.Require().NoError(err)

	expected := dbtypes.NewTotalSupplyRow(
		dbtypes.NewDbCoins(coins),
		10,
	)
	var rows []dbtypes.TotalSupplyRow
	err = suite.database.Sqlx.Select(&rows, `SELECT coins,height FROM supply`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "supply table should contain only one row")

	suite.Require().True(expected.Equals(rows[0]))

}

func (suite *DbTestSuite) TestBigDipperDb_GetTokenNames() {
	coins := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
		sdk.NewCoin("uatom", sdk.NewInt(15)),
	)
	_, err := suite.database.Sql.Exec("INSERT INTO supply(coins,height) VALUES ($1,$2) ", pq.Array(dbtypes.NewDbCoins(coins)), 10)
	suite.Require().NoError(err)
	expected := [2]string{"desmos", "uatom"}
	result, err := suite.database.GetTokenNames()

	suite.Require().NoError(err)
	for i, row := range expected {
		suite.Require().True(row == (result[i]))
	}
}
