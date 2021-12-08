package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	"github.com/forbole/bdjuno/v2/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveFeeGrantAllowance() {

	// save the data
	var allowance feegranttypes.FeeAllowanceI
	feeGrant, err := feegranttypes.NewGrant(sdk.AccAddress("cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt"), sdk.AccAddress("cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn"), allowance)
	feeGrantAllowance := types.NewFeeGrant(feeGrant, 121622)

	err = suite.database.SaveFeeGrantAllowance(feeGrantAllowance)
	suite.Require().NoError(err)

	// test dobule insertion
	err = suite.database.SaveFeeGrantAllowance(types.NewFeeGrant(feeGrant, 121622))
	suite.Require().NoError(err, "double feegrant allowance insertion should not insert the values again and returns no error")

	// verify the data
	expected := []dbtypes.FeeAllowanceRow{dbtypes.NewFeeAllowanceRow("cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn", "cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt", "", 122222)}

	var result []dbtypes.FeeAllowanceRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM fee_grant_allowance`)
	suite.Require().NoError(err)
	suite.Require().Len(result, len(expected))
	for index, row := range result {
		suite.Require().True(row.Equals(expected[index]))
	}

}

func (suite *DbTestSuite) TestBigDipperDb_RemoveFeeGrantAllowance() {

	// save the data
	var allowance feegranttypes.FeeAllowanceI
	feeGrant, err := feegranttypes.NewGrant(sdk.AccAddress("cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt"), sdk.AccAddress("cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn"), allowance)
	feeGrantAllowance := types.NewFeeGrant(feeGrant, 121622)

	err = suite.database.SaveFeeGrantAllowance(feeGrantAllowance)
	suite.Require().NoError(err)

	// delete the data
	allowanceToDelete := types.NewGrantRemoval("cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn", "cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt", 122222)
	err = suite.database.DeleteFeeGrantAllowance(allowanceToDelete)
	suite.Require().NoError(err)

	// verify the data
	var count int
	err = suite.database.Sql.QueryRow(`SELECT COUNT(*) FROM redelegation`).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(0, count)

}
