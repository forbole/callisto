package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveAccount() {
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

	err = suite.database.SaveAccount(&account, height, timestamp)
	suite.Require().NoError(err)

	err = suite.database.SaveAccount(&account, height, timestamp)
	suite.Require().NoError(err, "double account insertion should not insert and returns no error")

	// Verify the insertion
	var accountRows []dbtypes.AccountRow
	err = suite.database.Sqlx.Select(&accountRows, `SELECT * FROM account`)
	suite.Require().NoError(err)
	suite.Require().Len(accountRows, 1, "account table should contain only one row")

	expectedAccountRow := dbtypes.NewAccountRow("cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7")
	suite.Require().True(expectedAccountRow.Equal(accountRows[0]))

	var balancesRows []dbtypes.BalanceRow
	err = suite.database.Sqlx.Select(&balancesRows, `SELECT address,coins,height,timestamp FROM balance`)
	suite.Require().NoError(err)
	suite.Require().Len(balancesRows, 1, "balance table should contain only one row")

	expectedBalanceRow := dbtypes.NewBalanceRow(
		"cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7",
		dbtypes.NewDbCoins(coins),
		height,
		timestamp,
	)
	suite.Require().True(expectedBalanceRow.Equal(balancesRows[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveAccounts() {
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

	// Save the accounts
	err = suite.database.SaveAccounts(accounts, height, timestamp)
	suite.Require().NoError(err, "storing accounts should return no error")

	// Verify the data
	var accountsRows []dbtypes.AccountRow
	err = suite.database.Sqlx.Select(&accountsRows, `SELECT * FROM account ORDER BY address`)
	suite.Require().NoError(err)
	suite.Require().Len(accountsRows, 2)

	suite.Require().Equal("cosmos150zkt7g7kf3ymnzl28dksqvhjuxs9newc9uaq4", accountsRows[0].Address)
	suite.Require().Equal("cosmos1ngpsastyerhhpj72lvy38kn56cmuspfdwu7lg2", accountsRows[1].Address)

	var balancesRows []dbtypes.BalanceRow
	err = suite.database.Sqlx.Select(&balancesRows, `SELECT * FROM balance ORDER BY address`)
	suite.Require().NoError(err)
	suite.Require().Len(balancesRows, 2)

	suite.Require().True(balancesRows[0].Equal(dbtypes.NewBalanceRow(
		"cosmos150zkt7g7kf3ymnzl28dksqvhjuxs9newc9uaq4",
		dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))),
		height,
		timestamp,
	)))

	suite.Require().True(balancesRows[1].Equal(dbtypes.NewBalanceRow(
		"cosmos1ngpsastyerhhpj72lvy38kn56cmuspfdwu7lg2",
		dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(50000)))),
		height,
		timestamp,
	)))
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
