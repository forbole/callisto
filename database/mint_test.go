package database_test

import (
	"encoding/json"

	minttypes "github.com/Stride-Labs/stride/v12/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v4/types"

	dbtypes "github.com/forbole/bdjuno/v4/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveMintParams() {
	mintParams := minttypes.NewParams(
		"udaric",
		sdk.NewDecWithPrec(4, 1),
		"day",
		sdk.NewDecWithPrec(8, 1),
		4,
		minttypes.DistributionProportions{
			Staking:                     sdk.NewDec(0),
			CommunityPoolGrowth:         sdk.NewDec(0),
			CommunityPoolSecurityBudget: sdk.NewDec(0),
			StrategicReserve:            sdk.NewDec(0),
		},
		5006000,
	)
	err := suite.database.SaveMintParams(types.NewMintParams(mintParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.MintParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM mint_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var storedParams minttypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &storedParams)
	suite.Require().NoError(err)
	suite.Require().Equal(mintParams, storedParams)
	suite.Require().Equal(int64(10), rows[0].Height)
}
