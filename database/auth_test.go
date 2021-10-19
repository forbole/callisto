package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

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

func (suite *DbTestSuite) Test_SaveVestingAccounts() {
	// --- Save account addresses for foreign key contraint ---
	queries := []string{
		`INSERT INTO account (address) VALUES ('cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt')`,
		`INSERT INTO account (address) VALUES ('cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn')`,
		`INSERT INTO account (address) VALUES ('cosmos1eg47ue0l85lzkfgc4leske6hcah8cz3qajpjy2')`,
	}
	for _, query := range queries {
		_, err := suite.database.Sql.Exec(query)
		suite.Require().NoError(err)
	}
	var accountRows []dbtypes.AccountRow
	suite.database.Sqlx.Select(&accountRows, `SELECT * FROM account`)

	// --- Prepare VestingAccounts for saving into DB ---
	sdkCoins := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10)))
	// ContinuousVestingAccount
	address1, err := sdk.AccAddressFromBech32("cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt")
	suite.Require().NoError(err)
	ContinuousVestingAccount := vestingtypes.NewContinuousVestingAccount(
		authttypes.NewBaseAccountWithAddress(address1),
		sdkCoins,
		time.Date(1990, 9, 9, 00, 00, 00, 000, time.UTC).Unix(), // Start Time
		time.Date(2020, 9, 9, 00, 00, 00, 000, time.UTC).Unix(), // End Time
	)

	// DelayedVestingAccount
	address2, _ := sdk.AccAddressFromBech32("cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn")
	DelayedVestingAccount := vestingtypes.NewDelayedVestingAccount(
		authttypes.NewBaseAccountWithAddress(address2),
		sdkCoins,
		time.Date(2020, 9, 9, 00, 00, 00, 000, time.UTC).Unix(), // End Time
	)

	// PeriodicVestingAccount
	address3, _ := sdk.AccAddressFromBech32("cosmos1eg47ue0l85lzkfgc4leske6hcah8cz3qajpjy2")
	periods := []vestingtypes.Period{
		{
			Length: 2629743,
			Amount: sdkCoins,
		},
		{
			Length: 7889229,
			Amount: sdkCoins,
		},
	}
	PeriodicVestingAccount := vestingtypes.NewPeriodicVestingAccount(
		authttypes.NewBaseAccountWithAddress(address3),
		sdkCoins,
		time.Date(1990, 9, 9, 00, 00, 00, 000, time.UTC).Unix(), // Start Time
		periods,
	)

	// VestingAccounts
	vestingAccounts := []exported.VestingAccount{
		ContinuousVestingAccount,
		DelayedVestingAccount,
		PeriodicVestingAccount,
	}

	// --- Save the data into DB ---
	err = suite.database.SaveVestingAccounts(vestingAccounts)
	suite.Require().NoError(err)

	// --- Verify Continuous Vesting Account ---
	var continuousVestingAccountRow []dbtypes.ContinuousVestingAccountRow
	err = suite.database.Sqlx.Select(&continuousVestingAccountRow, `SELECT * FROM vesting_account WHERE type = 'ContinuousVestingAccount'`)
	suite.Require().NoError(err)
	suite.Require().Len(continuousVestingAccountRow, 1, "ContinuousVestingAccount type should contain only one row")

	expectedContinuousVestingAccountRow := dbtypes.NewContinuousVestingAccountRow(
		1,
		"ContinuousVestingAccount",
		"cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt",
		dbtypes.NewDbCoins(sdkCoins),
		time.Date(2020, 9, 9, 00, 00, 00, 000, time.UTC), // EndTime
		time.Date(1990, 9, 9, 00, 00, 00, 000, time.UTC), // StartTime
	)
	suite.Require().True(expectedContinuousVestingAccountRow.Equal(continuousVestingAccountRow[0]))

	// --- Verify Delayed Vesting Account ---
	var delayedVestingAccountRow []dbtypes.DelayedVestingAccountRow
	err = suite.database.Sqlx.Select(&delayedVestingAccountRow, `SELECT id, type, address, original_vesting, end_time FROM vesting_account WHERE type = 'DelayedVestingAccount'`)
	suite.Require().NoError(err)
	suite.Require().Len(delayedVestingAccountRow, 1, "DelayedVestingAccountRow type should contain only one row")

	expectedDelayedVestingAccountRow := dbtypes.NewDelayedVestingAccountRow(
		2,
		"DelayedVestingAccount",
		"cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn",
		dbtypes.NewDbCoins(sdkCoins),
		time.Date(2020, 9, 9, 00, 00, 00, 000, time.UTC), // EndTime
	)
	suite.Require().True(expectedDelayedVestingAccountRow.Equal(delayedVestingAccountRow[0]))

	// --- Verify Periodic Vesting Account ---
	var periodicVestingAccountRow []dbtypes.PeriodicVestingAccountRow
	err = suite.database.Sqlx.Select(&periodicVestingAccountRow, `SELECT * FROM vesting_account WHERE type = 'PeriodicVestingAccount'`)
	suite.Require().NoError(err)
	suite.Require().Len(delayedVestingAccountRow, 1, "DelayedVestingAccountRow type should contain only one row")

	expectedPeriodicVestingAccountRow := dbtypes.NewPeriodicVestingAccountRow(
		3,
		"PeriodicVestingAccount",
		"cosmos1eg47ue0l85lzkfgc4leske6hcah8cz3qajpjy2",
		dbtypes.NewDbCoins(sdkCoins),
		time.Date(2020, 9, 9, 00, 00, 00, 000, time.UTC), // EndTime
		time.Date(1990, 9, 9, 00, 00, 00, 000, time.UTC), // StartTime
	)
	suite.Require().True(expectedPeriodicVestingAccountRow.Equal(periodicVestingAccountRow[0]))

	// --- Verify vesting periods ---
	var vestingPeriodRows []dbtypes.VestingPeriodRow
	err = suite.database.Sqlx.Select(&vestingPeriodRows, `SELECT * FROM vesting_period`)
	suite.Require().NoError(err)
	suite.Require().Len(vestingPeriodRows, 2, "vestingPeriodRows should contain only 2 rows")

	expectedVestingPeriods := []dbtypes.VestingPeriodRow{
		dbtypes.NewVestingPeriodRow(3, 0, "2629743", dbtypes.NewDbCoins(sdkCoins)),
		dbtypes.NewVestingPeriodRow(3, 1, "7889229", dbtypes.NewDbCoins(sdkCoins)),
	}
	for index, vestingPeriod := range expectedVestingPeriods {
		suite.Require().True(vestingPeriod.Equal(vestingPeriodRows[index]))
	}

}
