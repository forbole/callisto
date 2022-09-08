package database_test

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	"github.com/forbole/bdjuno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveInterchainStakingParams() {
	icsParams := icstypes.NewParams(
		100,
		10,
		10,
		sdk.NewDecWithPrec(4, 1),
	)
	err := suite.database.SaveInterchainStakingParams(types.NewInterchainStakingParams(icsParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.InterchainStakingParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM interchain_staking_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var storedParams icstypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &storedParams)
	suite.Require().NoError(err)
	suite.Require().Equal(icsParams, storedParams)
	suite.Require().Equal(int64(10), rows[0].Height)
}
