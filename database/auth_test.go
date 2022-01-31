package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/forbole/bdjuno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
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
