package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	bbanktypes "github.com/forbole/bdjuno/x/bank/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestSaveAccountBalance() {
	address1 := suite.getAccount("cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7")
	address2 := suite.getAccount("cosmos1tcpsdy9alvucwj0h23n56tey6zmrvkm7sndh9j")

	// Save the data
	err := suite.database.SaveAccountBalances([]bbanktypes.AccountBalance{
		bbanktypes.NewAccountBalance(
			address1.String(),
			sdk.NewCoins(
				sdk.NewCoin("desmos", sdk.NewInt(10)),
				sdk.NewCoin("uatom", sdk.NewInt(20)),
			),
			10,
		),
		bbanktypes.NewAccountBalance(
			address2.String(),
			sdk.NewCoins(
				sdk.NewCoin("uatom", sdk.NewInt(100)),
			),
			10,
		),
	})
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.AccountBalanceRow{
		dbtypes.NewAccountBalanceRow(
			address1.String(),
			dbtypes.NewDbCoins(sdk.NewCoins(
				sdk.NewCoin("desmos", sdk.NewInt(10)),
				sdk.NewCoin("uatom", sdk.NewInt(20)),
			)),
			10,
		),
		dbtypes.NewAccountBalanceRow(
			address2.String(),
			dbtypes.NewDbCoins(sdk.NewCoins(
				sdk.NewCoin("uatom", sdk.NewInt(100)),
			)),
			10,
		),
	}

	var rows []dbtypes.AccountBalanceRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM account_balance ORDER BY address`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]))
	}

	// Update the data
	err = suite.database.SaveAccountBalances([]bbanktypes.AccountBalance{
		bbanktypes.NewAccountBalance(
			address1.String(),
			sdk.NewCoins(
				sdk.NewCoin("desmos", sdk.NewInt(10)),
			),
			9,
		),
		bbanktypes.NewAccountBalance(
			address2.String(),
			sdk.NewCoins(
				sdk.NewCoin("uatom", sdk.NewInt(100)),
				sdk.NewCoin("desmos", sdk.NewInt(200)),
			),
			11,
		),
	})
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.AccountBalanceRow{
		dbtypes.NewAccountBalanceRow(
			address1.String(),
			dbtypes.NewDbCoins(sdk.NewCoins(
				sdk.NewCoin("desmos", sdk.NewInt(10)),
				sdk.NewCoin("uatom", sdk.NewInt(20)),
			)),
			10,
		),
		dbtypes.NewAccountBalanceRow(
			address2.String(),
			dbtypes.NewDbCoins(sdk.NewCoins(
				sdk.NewCoin("uatom", sdk.NewInt(100)),
				sdk.NewCoin("desmos", sdk.NewInt(200)),
			)),
			11,
		),
	}

	rows = []dbtypes.AccountBalanceRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM account_balance ORDER BY address`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveSupply() {
	// Save the data
	original := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
		sdk.NewCoin("uatom", sdk.NewInt(15)),
	)
	err := suite.database.SaveSupply(original, 10)
	suite.Require().NoError(err)

	// Verify the data
	expected := dbtypes.NewSupplyRow(dbtypes.NewDbCoins(original), 10)

	var rows []dbtypes.SupplyRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM supply`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "supply table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a lower height
	coins := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
		sdk.NewCoin("uatom", sdk.NewInt(15)),
	)
	err = suite.database.SaveSupply(coins, 9)
	suite.Require().NoError(err)

	// Verify the data
	rows = []dbtypes.SupplyRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM supply`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "supply table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with same height
	coins = sdk.NewCoins(sdk.NewCoin("uakash", sdk.NewInt(10)))
	err = suite.database.SaveSupply(coins, 10)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewSupplyRow(dbtypes.NewDbCoins(coins), 10)

	rows = []dbtypes.SupplyRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM supply`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "supply table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with higher height
	coins = sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10)))
	err = suite.database.SaveSupply(coins, 20)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewSupplyRow(dbtypes.NewDbCoins(coins), 20)

	rows = []dbtypes.SupplyRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM supply`)
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
