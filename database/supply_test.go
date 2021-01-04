package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveTotalTokens() {
	coins := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
		sdk.NewCoin("uatom", sdk.NewInt(15)),
	)
	err := suite.database.SaveSupplyToken(coins, 10)
	suite.Require().NoError(err)

	expected := dbtypes.NewTotalSupplyRow(
		dbtypes.NewDbCoins(coins),
		10,
	)
	var rows []dbtypes.TotalSupplyRow
	err = suite.database.Sqlx.Select(&rows, `SELECT coins,height FROM supply`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "supply table should contain only one row")

	suite.Require().True(expected.Equals(rows[0]))

}

func (suite *DbTestSuite) TestBigDipperDb_GetTokenNames() {
	coins := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
		sdk.NewCoin("uatom", sdk.NewInt(15)),
	)
	_, err := suite.database.Sql.Exec("INSERT INTO supply(coins,height) VALUES ($1,$2) ", pq.Array(dbtypes.NewDbCoins(coins)), 10)
	suite.Require().NoError(err)
	expected := [2]string{"desmos", "uatom"}
	result, err := suite.database.GetTokenNames()

	suite.Require().NoError(err)
	for i, row := range expected {
		suite.Require().True(row == (result[i]))
	}
}
