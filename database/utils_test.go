package database_test

import (
	"time"

	"github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_InsertEnableModules() {
	modules := make(map[string]bool)
	modules["staking"] = true
	modules["auth"] = true
	modules["supply"] = true
	modules["distribution"] = true
	modules["pricefeed"] = true
	modules["bank"] = true
	modules["consensus"] = true
	modules["mint"] = false

	timestamp, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)

	err = suite.database.InsertEnableModules(modules, timestamp)
	suite.Require().NoError(err)

	var results []types.ModulesRow
	err = suite.database.Sqlx.Select(&results, "SELECT * FROM modules")
	suite.Require().Len(results, 1, "modules table should contain only one row")
	suite.Require().NoError(err)

	expected := types.NewModulesRow(true, true, true, true, true, true, true, false, timestamp)
	for _, result := range results {
		suite.Require().True(result.Equals(expected))
	}
}
