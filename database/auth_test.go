package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/forbole/bdjuno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

func (suite *DbTestSuite) TestSaveAccount() {
	address, err := sdk.AccAddressFromBech32("cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7")
	suite.Require().NoError(err)

	account := authttypes.NewBaseAccountWithAddress(address)

	// ------------------------------
	// --- Save the data
	// ------------------------------

	err = suite.database.SaveAccounts([]types.Account{types.NewAccount(account.Address)})
	suite.Require().NoError(err)

	err = suite.database.SaveAccounts([]types.Account{types.NewAccount(account.Address)})
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
		suite.Require().Equal(acc, accounts[index])
	}
}
