package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestSaveAccount() {
	address, err := sdk.AccAddressFromBech32("cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7")
	suite.Require().NoError(err)

	height := int64(100)

	timestamp, err := time.Parse(time.RFC3339, "2020-02-02T15:00:00Z")
	suite.Require().NoError(err)

	coins := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
		sdk.NewCoin("uatom", sdk.NewInt(15)),
	)

	account := auth.NewBaseAccountWithAddress(address)
	suite.Require().NoError(account.SetCoins(coins))

	// ------------------------------
	// --- Save the data
	// ------------------------------

	err = suite.database.SaveAccount(&account, height, timestamp)
	suite.Require().NoError(err)

	err = suite.database.SaveAccount(&account, height, timestamp)
	suite.Require().NoError(err, "double account insertion should not insert and returns no error")

	// ------------------------------
	// --- Verify the data
	// ------------------------------

	// Accounts row
	var accountRows []dbtypes.AccountRow
	err = suite.database.Sqlx.Select(&accountRows, `SELECT * FROM account`)
	suite.Require().NoError(err)
	suite.Require().Len(accountRows, 1, "account table should contain only one row")

	expectedAccountRow := dbtypes.NewAccountRow("cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7")
	suite.Require().True(expectedAccountRow.Equal(accountRows[0]))

	// Current balances
	var balRows []dbtypes.AccountBalanceRow
	err = suite.database.Sqlx.Select(&balRows, `SELECT * FROM account_balance ORDER BY address`)
	suite.Require().NoError(err)

	expectedBalRows := []dbtypes.AccountBalanceRow{
		dbtypes.NewAccountBalanceRow(
			"cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7",
			dbtypes.NewDbCoins(coins),
		),
	}
	suite.Require().Len(balRows, len(expectedBalRows))
	for index, expected := range expectedBalRows {
		suite.Require().True(expected.Equal(balRows[index]))
	}

	// Balance histories
	var balHisRows []dbtypes.AccountBalanceHistoryRow
	err = suite.database.Sqlx.Select(&balHisRows, `SELECT * FROM account_balance_history ORDER BY address`)
	suite.Require().NoError(err)

	expectedBalHisRows := []dbtypes.AccountBalanceHistoryRow{
		dbtypes.NewAccountBalanceHistoryRow(
			"cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7",
			dbtypes.NewDbCoins(coins),
			height,
			timestamp,
		),
	}
	suite.Require().Len(balHisRows, len(expectedBalHisRows))
	for index, expected := range expectedBalHisRows {
		suite.Require().True(expected.Equal(balHisRows[index]))
	}
}

func (suite *DbTestSuite) TestSaveAccounts() {
	firstAddr, err := sdk.AccAddressFromBech32("cosmos150zkt7g7kf3ymnzl28dksqvhjuxs9newc9uaq4")
	suite.Require().NoError(err)

	secondAddr, err := sdk.AccAddressFromBech32("cosmos1ngpsastyerhhpj72lvy38kn56cmuspfdwu7lg2")
	suite.Require().NoError(err)

	firstAcc := auth.NewBaseAccountWithAddress(firstAddr)
	suite.Require().NoError(firstAcc.SetCoins(sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))))

	secondAcc := auth.NewBaseAccountWithAddress(secondAddr)
	suite.Require().NoError(secondAcc.SetCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(50000)))))

	accounts := []exported.Account{&firstAcc, &secondAcc}

	height := int64(100)

	timestamp, err := time.Parse(time.RFC3339, "2020-02-02T15:00:00Z")
	suite.Require().NoError(err)

	// ------------------------------
	// --- Save the data
	// ------------------------------

	err = suite.database.SaveAccounts(accounts, height, timestamp)
	suite.Require().NoError(err, "storing accounts should return no error")

	// ------------------------------
	// --- Verify the data
	// ------------------------------

	// Accounts data
	var accountsRows []dbtypes.AccountRow
	err = suite.database.Sqlx.Select(&accountsRows, `SELECT * FROM account ORDER BY address`)
	suite.Require().NoError(err)

	expAccRows := []dbtypes.AccountRow{
		dbtypes.NewAccountRow("cosmos150zkt7g7kf3ymnzl28dksqvhjuxs9newc9uaq4"),
		dbtypes.NewAccountRow("cosmos1ngpsastyerhhpj72lvy38kn56cmuspfdwu7lg2"),
	}
	suite.Require().Len(accountsRows, len(expAccRows))
	for index, expected := range expAccRows {
		suite.Require().True(expected.Equal(accountsRows[index]))
	}

	// Current balances
	var balRows []dbtypes.AccountBalanceRow
	err = suite.database.Sqlx.Select(&balRows, `SELECT * FROM account_balance ORDER BY address`)
	suite.Require().NoError(err)

	expectedBalRows := []dbtypes.AccountBalanceRow{
		dbtypes.NewAccountBalanceRow(
			"cosmos150zkt7g7kf3ymnzl28dksqvhjuxs9newc9uaq4",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))),
		),
		dbtypes.NewAccountBalanceRow(
			"cosmos1ngpsastyerhhpj72lvy38kn56cmuspfdwu7lg2",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(50000)))),
		),
	}
	suite.Require().Len(balRows, len(expectedBalRows))
	for index, expected := range expectedBalRows {
		suite.Require().True(expected.Equal(balRows[index]))
	}

	// Balance histories
	var balHisRows []dbtypes.AccountBalanceHistoryRow
	err = suite.database.Sqlx.Select(&balHisRows, `SELECT * FROM account_balance_history ORDER BY address`)
	suite.Require().NoError(err)

	expectedBalHisRows := []dbtypes.AccountBalanceHistoryRow{
		dbtypes.NewAccountBalanceHistoryRow(
			"cosmos150zkt7g7kf3ymnzl28dksqvhjuxs9newc9uaq4",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))),
			height,
			timestamp,
		),
		dbtypes.NewAccountBalanceHistoryRow(
			"cosmos1ngpsastyerhhpj72lvy38kn56cmuspfdwu7lg2",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(50000)))),
			height,
			timestamp,
		),
	}
	suite.Require().Len(balHisRows, len(expectedBalHisRows))
	for index, expected := range expectedBalHisRows {
		suite.Require().True(expected.Equal(balHisRows[index]))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_GetAccounts() {
	// Insert the data
	queries := []string{
		`INSERT INTO account (address) VALUES ('cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt')`,
		`INSERT INTO account (address) VALUES ('cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn')`,
		`INSERT INTO account (address) VALUES ('cosmos1eg47ue0l85lzkfgc4leske6hcah8cz3qajpjy2')`,
		`INSERT INTO account (address) VALUES ('cosmos1495ghynrns8sxfnw8mj887pgh0c9z6c4lqkzme')`,
		`INSERT INTO account (address) VALUES ('cosmos18fzr6adp3gjw43xu62vfhg248lepfwpf0pj2dm')`,
	}

	for _, query := range queries {
		_, err := suite.database.Sql.Exec(query)
		suite.Require().NoError(err)
	}

	// Get the data
	accounts, err := suite.database.GetAccounts()
	suite.Require().NoError(err)

	// Verify the get
	expectedAccs := []string{
		"cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt",
		"cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn",
		"cosmos1eg47ue0l85lzkfgc4leske6hcah8cz3qajpjy2",
		"cosmos1495ghynrns8sxfnw8mj887pgh0c9z6c4lqkzme",
		"cosmos18fzr6adp3gjw43xu62vfhg248lepfwpf0pj2dm",
	}

	for index, acc := range expectedAccs {
		suite.Require().Equal(acc, accounts[index].String())
	}
}
