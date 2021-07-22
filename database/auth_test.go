package database_test

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/forbole/bdjuno/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) getAccountString(add string) string{
	address, err := sdk.AccAddressFromBech32(add)
	suite.Require().NoError(err)

	coin:=sdk.Coin{
		Denom: "daric",
		Amount: sdk.NewInt(10),
	}
	account := authtypes.NewBaseAccountWithAddress(address)
	baseVestingAccount := authvestingtypes.BaseVestingAccount{
		BaseAccount: account,
		OriginalVesting: sdk.NewCoins(coin),
		DelegatedFree: sdk.NewCoins(coin),
		DelegatedVesting: sdk.NewCoins(coin),
	}
	continuousVestingAccount :=authvestingtypes.ContinuousVestingAccount{
		BaseVestingAccount: &baseVestingAccount,
		StartTime: 10,
	}
	saveAccount:=[]types.Account{
		types.NewAccount(
			continuousVestingAccount.Address,
			&continuousVestingAccount)}

}

func (suite *DbTestSuite) TestSaveAccount() {
	address, err := sdk.AccAddressFromBech32("cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7")
	suite.Require().NoError(err)

	coin:=sdk.Coin{
		Denom: "daric",
		Amount: sdk.NewInt(10),
	}
	account := authtypes.NewBaseAccountWithAddress(address)
	baseVestingAccount := authvestingtypes.BaseVestingAccount{
		BaseAccount: account,
		OriginalVesting: sdk.NewCoins(coin),
		DelegatedFree: sdk.NewCoins(coin),
		DelegatedVesting: sdk.NewCoins(coin),
	}
	continuousVestingAccount :=authvestingtypes.ContinuousVestingAccount{
		BaseVestingAccount: &baseVestingAccount,
		StartTime: 10,
	}
	saveAccount:=[]types.Account{
		types.NewAccount(
			continuousVestingAccount.Address,
			&continuousVestingAccount)}

	
	// ------------------------------
	// --- Save the data
	// ------------------------------

	err = suite.database.SaveAccounts(saveAccount)
	suite.Require().NoError(err)

	err = suite.database.SaveAccounts(saveAccount)
	suite.Require().NoError(err, "double account insertion should not insert and returns no error")

	// ------------------------------
	// --- Verify the data
	// ------------------------------

	// Accounts row
	var accountRows []dbtypes.AccountRow
	err = suite.database.Sqlx.Select(&accountRows, `SELECT * FROM account`)
	suite.Require().NoError(err)
	suite.Require().Len(accountRows, 1, "account table should contain only one row")

	protoContent, ok := saveAccount[0].Details.(authtypes.AccountI)
	suite.Require().True(ok)
	anyContent, err := codectypes.NewAnyWithValue(protoContent)
	suite.Require().NoError(err)
	js,err:=suite.database.EncodingConfig.Marshaler.MarshalJSON(anyContent)
	suite.Require().NoError(err)
	expectedAccountRow := dbtypes.NewAccountRow("cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7",string(js))
	
	suite.Require().True(expectedAccountRow.Equal(accountRows[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_GetAccounts() {
	account:=`{
		"@type": "/cosmos.vesting.v1beta1.ContinuousVestingAccount",
		"base_vesting_account": {
		  "base_account": {
			"address": "cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7",
			"pub_key": {
			  "@type": "/cosmos.crypto.secp256k1.PubKey",
			  "key": "AsmSLYJM5CoIwuYQF+XvOFMSK1HeijFHF9XehSZGfET9"
			},
			"account_number": "143",
			"sequence": "3"
		  },
		  "original_vesting": [
			{
			  "denom": "udaric",
			  "amount": "30000000000"
			}
		  ],
		  "delegated_free": [],
		  "delegated_vesting": [],
		  "end_time": "1619766000"
		},
		"start_time": "1619506800",
		"vesting_periods": [
		  {
			"length": "86400",
			"amount": [
			  {
				"denom": "udaric",
				"amount": "10000000000"
			  }
			]
		  },
		  {
			"length": "86400",
			"amount": [
			  {
				"denom": "udaric",
				"amount": "10000000000"
			  }
			]
		  },
		  {
			"length": "86400",
			"amount": [
			  {
				"denom": "udaric",
				"amount": "10000000000"
			  }
			]
		  }
		]
	  }`
	

	
	// Insert the data
	queries := []string{
		fmt.Sprintf("INSERT INTO account (address) VALUES ('cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt',%s)",account),
		fmt.Sprintf("INSERT INTO account (address) VALUES ('cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn',%s)",account),
		fmt.Sprintf("INSERT INTO account (address) VALUES ('cosmos1eg47ue0l85lzkfgc4leske6hcah8cz3qajpjy2',%s)",account),
		fmt.Sprintf("INSERT INTO account (address) VALUES ('cosmos1495ghynrns8sxfnw8mj887pgh0c9z6c4lqkzme',%s)",account),
		fmt.Sprintf("INSERT INTO account (address) VALUES ('cosmos18fzr6adp3gjw43xu62vfhg248lepfwpf0pj2dm',%s)",account),
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
