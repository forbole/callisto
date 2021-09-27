package database_test

import (
	"encoding/json"
	"time"

	inflationtypes "github.com/e-money/em-ledger/x/inflation/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveEMoneyInflation() {

	// Save the data
	timeNow := time.Now()
	inflationState := inflationtypes.NewInflationState(
		timeNow,
		"echf", "0.005",
		"edkk", "0.005",
		"ungm", "0.1",
	)
	err := suite.database.SaveEMoneyInflation(types.NewEMoneyInfaltion(inflationState, 1))
	suite.Require().NoError(err)

	// Get stored data from DB
	var rows []dbtypes.EMoneyInflationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM emoney_inflation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	// Verify data
	var storedInflationAssets inflationtypes.InflationAssets
	err = json.Unmarshal([]byte(rows[0].Inflation), &storedInflationAssets)
	suite.Require().NoError(err)
	suite.Require().Equal(inflationState.InflationAssets, storedInflationAssets)
	suite.Require().Equal(int64(0), rows[0].LastAppliedHeight)
	suite.Require().Equal(int64(1), rows[0].Height)

}
