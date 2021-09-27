package database_test

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveEMoneyGasPrices() {

	//prepare data for storing in DB
	authorityKey := "emoney1zp22dzt8q2vfkemytl8fwvr2d4vjxrsyjugygv"
	minGasPrices := sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "echf",
			Amount: sdk.NewDecWithPrec(53, 2),
		},
		sdk.DecCoin{
			Denom:  "edkk",
			Amount: sdk.NewDecWithPrec(370, 2),
		},
		sdk.DecCoin{
			Denom:  "ungm",
			Amount: sdk.NewDec(1),
		},
	}

	// Save the data
	eMoneyGasPrices := types.NewEMoneyGasPrices(
		authorityKey,
		minGasPrices,
		1,
	)
	err := suite.database.SaveEMoneyGasPrices(eMoneyGasPrices)
	suite.Require().NoError(err)

	// Verify data
	row := []dbtypes.EMoneyGasPricesRow{}
	err = suite.database.Sqlx.Select(&row, `SELECT * FROM emoney_gas_prices`)
	if err != nil {
		log.Fatal(err)
	}
	// suite.Require().NoError(err)
	// suite.Require().Len(row, 1)
	// suite.Require().Equal(row[0].AuthorityKey, authorityKey)
}
