package database_test

import (
	"encoding/json"

	stakeibctypes "github.com/Stride-Labs/stride/v5/x/stakeibc/types"

	"github.com/forbole/bdjuno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveStakeIBCParams() {
	icsParams := stakeibctypes.NewParams(
		100,
		10,
		10,
		10,
		11,
		12,
		15,
		10,
		12,
		11,
		10,
		11,
		9,
		10,
		10,
		10,
	)
	err := suite.database.SaveStakeIBCParams(types.NewStakeIBCParams(icsParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.StakeIBCParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM stakeibc_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var storedParams stakeibctypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &storedParams)
	suite.Require().NoError(err)
	suite.Require().Equal(icsParams, storedParams)
	suite.Require().Equal(int64(10), rows[0].Height)
}
