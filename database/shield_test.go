package database_test

import (
	"github.com/forbole/bdjuno/v3/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_ShieldPool() {

	// Save the data
	shield := types.NewShieldPool(
		1,
		sdk.NewInt(1000000000),
		sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(1000000000))),
		sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(2000000000))),
		"Sponsor",
		"certik1ktt6kej30pycfnm5aq42x0edcm807kqcpw273p",
		"Shield1 Description",
		sdk.NewInt(1000000000000),
		false,
		123654)

	shield2 := types.NewShieldPool(
		2,
		sdk.NewInt(1000000000),
		sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(3000000000))),
		sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(4000000000))),
		"Sponsor",
		"certik1ktt6kej30pycfnm5aq42x0edcm807kqcpw273p",
		"Shield2 Description",
		sdk.NewInt(1000000000000),
		true,
		123654)

	err := suite.database.SaveShieldPool(shield)
	suite.Require().NoError(err)

	err = suite.database.SaveShieldPool(shield2)
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.ShieldPoolRow{
		dbtypes.NewShieldPoolRow(
			1,
			"1000000000",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(1000000000)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(2000000000)))),
			"Sponsor",
			"certik1ktt6kej30pycfnm5aq42x0edcm807kqcpw273p",
			"Shield1 Description",
			"1000000000000",
			false,
			123654,
		),
		dbtypes.NewShieldPoolRow(
			2,
			"1000000000",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(3000000000)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(4000000000)))),
			"Sponsor",
			"certik1ktt6kej30pycfnm5aq42x0edcm807kqcpw273p",
			"Shield2 Description",
			"1000000000000",
			true,
			123654,
		),
	}

	var rows []dbtypes.ShieldPoolRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM shield_pool`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equal(row))
	}

	// ----------------------------------------------------------------------------------------------------------------

	// Update the data
	updateShield1 := types.NewShieldPool(
		1,
		sdk.NewInt(5500000000),
		sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(1000000055))),
		sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(7500000000))),
		"Sponsor",
		"certik1ktt6kej30pycfnm5aq42x0edcm807kqcpw273p",
		"Shield Description",
		sdk.NewInt(1000000000002),
		false,
		123700,
	)

	updateShield2 := types.NewShieldPool(
		2,
		sdk.NewInt(6600000000),
		sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(3000000055))),
		sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(8500000000))),
		"Sponsor",
		"certik1ktt6kej30pycfnm5aq42x0edcm807kqcpw273p",
		"Shield Description",
		sdk.NewInt(1000000000001),
		true,
		123700,
	)

	err = suite.database.SaveShieldPool(updateShield1)
	suite.Require().NoError(err)

	err = suite.database.SaveShieldPool(updateShield2)
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.ShieldPoolRow{
		dbtypes.NewShieldPoolRow(
			1,
			"5500000000",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(1000000055)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(7500000000)))),
			"Sponsor",
			"certik1ktt6kej30pycfnm5aq42x0edcm807kqcpw273p",
			"Shield Description",
			"1000000000002",
			false,
			123700,
		),
		dbtypes.NewShieldPoolRow(
			2,
			"6600000000",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(3000000055)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("ucrtk", sdk.NewInt(8500000000)))),
			"Sponsor",
			"certik1ktt6kej30pycfnm5aq42x0edcm807kqcpw273p",
			"Shield Description",
			"1000000000001",
			true,
			123700,
		),
	}

	rows = []dbtypes.ShieldPoolRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM shield_pool`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equal(row))
	}
}
