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
	address, err := sdk.AccAddressFromBech32("cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7")
	suite.Require().NoError(err)

	coin:=sdk.Coin{
		Denom: "daric",
		Amount: sdk.NewInt(10),
	}
	baseAccount := authtypes.NewBaseAccountWithAddress(address)
	baseVestingAccount := authvestingtypes.BaseVestingAccount{
		BaseAccount: baseAccount,
		OriginalVesting: sdk.NewCoins(coin),
		DelegatedFree: sdk.NewCoins(coin),
		DelegatedVesting: sdk.NewCoins(coin),
	}
	expectedAccount :=authvestingtypes.ContinuousVestingAccount{
		BaseVestingAccount: &baseVestingAccount,
		StartTime: 10,
	}
	
	account:=`'{"@type":"/cosmos.vesting.v1beta1.ContinuousVestingAccount","base_vesting_account":{"base_account":{"address":"cosmos140xsjjg6pwkjp0xjz8zru7ytha60l5aee9nlf7","pub_key":null,"account_number":"0","sequence":"0"},"original_vesting":[{"denom":"daric","amount":"10"}],"delegated_free":[{"denom":"daric","amount":"10"}],"delegated_vesting":[{"denom":"daric","amount":"10"}],"end_time":"0"},"start_time":"10"}'`
	/* account=strings.ReplaceAll(account,`"`,`""`)
	account=strings.ReplaceAll(account,`/`,`//`) */


	// Insert the data
	queries := []string{
		fmt.Sprintf("INSERT INTO account (address,details) VALUES ('cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt',%s)",account),
		fmt.Sprintf("INSERT INTO account (address,details) VALUES ('cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn',%s)",account),
		fmt.Sprintf("INSERT INTO account (address,details) VALUES ('cosmos1eg47ue0l85lzkfgc4leske6hcah8cz3qajpjy2',%s)",account),
		fmt.Sprintf("INSERT INTO account (address,details) VALUES ('cosmos1495ghynrns8sxfnw8mj887pgh0c9z6c4lqkzme',%s)",account),
		fmt.Sprintf("INSERT INTO account (address,details) VALUES ('cosmos18fzr6adp3gjw43xu62vfhg248lepfwpf0pj2dm',%s)",account),
	}

	for _, query := range queries {
		_, err := suite.database.Sql.Exec(query)
		suite.Require().NoError(err)
	}

	// Get the data
	accounts, err := suite.database.GetAccounts()
	suite.Require().NoError(err)

	// Verify the get
	expectedAccs := []types.Account{
		types.NewAccount("cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt",&expectedAccount),
		types.NewAccount("cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn",&expectedAccount),
		types.NewAccount("cosmos1eg47ue0l85lzkfgc4leske6hcah8cz3qajpjy2",&expectedAccount),
		types.NewAccount("cosmos1495ghynrns8sxfnw8mj887pgh0c9z6c4lqkzme",&expectedAccount),
		types.NewAccount("cosmos18fzr6adp3gjw43xu62vfhg248lepfwpf0pj2dm",&expectedAccount),
	}

	for index, acc := range expectedAccs {
		suite.Require().True(acc.Equal(accounts[index]))
	}
}
