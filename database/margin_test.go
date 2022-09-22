package database_test

import (
	"encoding/json"

	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveMarginParams() {
	var pools []string
	marginParams := dbtypes.NewMarginParams(
		sdk.NewDecWithPrec(4, 1),
		sdk.NewDecWithPrec(8, 1),
		sdk.NewDecWithPrec(4, 1),
		sdk.NewDecWithPrec(8, 1),
		sdk.NewDecWithPrec(8, 1),
		sdk.NewDecWithPrec(8, 1),
		10,
		pools,
		sdk.NewDecWithPrec(8, 1),
		uint64(10),
		sdk.NewDecWithPrec(8, 1),
		sdk.NewDecWithPrec(8, 1),
		"",
		sdk.NewDecWithPrec(8, 1),
		"",
		sdk.NewDecWithPrec(8, 1),
		sdk.NewDecWithPrec(8, 1),
		pools,
		true,
		true,
	)
	err := suite.database.SaveMarginParams(types.NewMarginParams(marginParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.MarginParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM margin_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var storedParams margintypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &storedParams)
	suite.Require().NoError(err)
	suite.Require().Equal(marginParams, storedParams)
	suite.Require().Equal(int64(10), rows[0].Height)
}
