package database_test

import (
	"encoding/json"

	"cosmossdk.io/math"
	"github.com/forbole/bdjuno/v4/types"

	markertypes "github.com/MonCatCat/provenance/x/marker/types"
	dbtypes "github.com/forbole/bdjuno/v4/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveMarkerParams() {
	markerParams := markertypes.NewParams(
		1000000, true, "abc", math.NewInt(1000),
	)
	err := suite.database.SaveMarkerParams(types.NewMarkerParams(markerParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.MarkerParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM marker_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var storedParams markertypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &storedParams)
	suite.Require().NoError(err)
	suite.Require().Equal(markerParams, storedParams)
	suite.Require().Equal(int64(10), rows[0].Height)
}
