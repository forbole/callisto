package database_test

import (
	"time"

	"github.com/forbole/bdjuno/v3/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_ShieldPool() {
	sponsorAddress1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	sponsorAddress2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	// Save the data
	shield := types.NewShieldPool(
		1,
		sdk.NewInt(1000000000),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000000000))),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000000))),
		"Sponsor",
		sponsorAddress1.String(),
		"Shield1 Description",
		sdk.NewInt(1000000000000),
		false,
		123654)

	shield2 := types.NewShieldPool(
		2,
		sdk.NewInt(1000000000),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(3000000000))),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000000))),
		"Sponsor",
		sponsorAddress2.String(),
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
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000000000)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000000)))),
			"Sponsor",
			sponsorAddress1.String(),
			"Shield1 Description",
			"1000000000000",
			false,
			123654,
		),
		dbtypes.NewShieldPoolRow(
			2,
			"1000000000",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(3000000000)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000000)))),
			"Sponsor",
			sponsorAddress2.String(),
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
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000000055))),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(7500000000))),
		"Sponsor",
		sponsorAddress1.String(),
		"Shield Description",
		sdk.NewInt(1000000000002),
		false,
		123700,
	)

	updateShield2 := types.NewShieldPool(
		2,
		sdk.NewInt(6600000000),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(3000000055))),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(8500000000))),
		"Sponsor",
		sponsorAddress2.String(),
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
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000000055)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(7500000000)))),
			"Sponsor",
			sponsorAddress1.String(),
			"Shield Description",
			"1000000000002",
			false,
			123700,
		),
		dbtypes.NewShieldPoolRow(
			2,
			"6600000000",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(3000000055)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(8500000000)))),
			"Sponsor",
			sponsorAddress2.String(),
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

func (suite *DbTestSuite) TestBigDipperDb_ShieldProvider() {
	providerAddress1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	providerAddress2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	// Save the data
	shieldProvider1 := types.NewShieldProvider(
		providerAddress1.String(),
		2000000000,
		1000000,
		sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(30000))),
		sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(10000))),
		1000000,
		0,
		18265)

	shieldProvider2 := types.NewShieldProvider(
		providerAddress2.String(),
		5000000000,
		6000000,
		sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(90000))),
		sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(40000))),
		3000000,
		100000,
		18265)

	err := suite.database.SaveShieldProvider(shieldProvider1)
	suite.Require().NoError(err)

	err = suite.database.SaveShieldProvider(shieldProvider2)
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.ShieldProviderRow{
		dbtypes.NewShieldProviderRow(
			providerAddress1.String(),
			2000000000,
			1000000,
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(30000)))),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(10000)))),
			1000000,
			0,
			18265,
		),
		dbtypes.NewShieldProviderRow(
			providerAddress2.String(),
			5000000000,
			6000000,
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(90000)))),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(40000)))),
			3000000,
			100000,
			18265,
		),
	}

	var rows []dbtypes.ShieldProviderRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM shield_provider`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equal(row))
	}

	// ----------------------------------------------------------------------------------------------------------------

	// Update the data
	updatedshieldProvider1 := types.NewShieldProvider(
		providerAddress1.String(),
		4000000000,
		2000000,
		sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(540000))),
		sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(12000))),
		3829111,
		20100,
		18762)

	updatedshieldProvider2 := types.NewShieldProvider(
		providerAddress2.String(),
		7000000000,
		9000000,
		sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(124000))),
		sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(17002))),
		21300000,
		2344444,
		18762)

	err = suite.database.SaveShieldProvider(updatedshieldProvider1)
	suite.Require().NoError(err)

	err = suite.database.SaveShieldProvider(updatedshieldProvider2)
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.ShieldProviderRow{
		dbtypes.NewShieldProviderRow(
			providerAddress1.String(),
			4000000000,
			2000000,
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(540000)))),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(12000)))),
			3829111,
			20100,
			18762),
		dbtypes.NewShieldProviderRow(
			providerAddress2.String(),
			7000000000,
			9000000,
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(124000)))),
			dbtypes.NewDbDecCoins(sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt(17002)))),
			21300000,
			2344444,
			18762,
		),
	}

	rows = []dbtypes.ShieldProviderRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM shield_provider`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equal(row))
	}

}

func (suite *DbTestSuite) TestBigDipperDb_ShieldPurchase() {
	purchaserAddress1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	purchaserAddress2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	// Save pool details
	shield := types.NewShieldPool(
		1,
		sdk.NewInt(1000000000),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000000000))),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000000))),
		"Sponsor",
		purchaserAddress1.String(),
		"Shield Description",
		sdk.NewInt(1000000000000),
		false,
		123654)

	shield2 := types.NewShieldPool(
		2,
		sdk.NewInt(1000000000),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(3000000000))),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000000))),
		"Sponsor",
		purchaserAddress2.String(),
		"Shield Description",
		sdk.NewInt(1000000000000),
		false,
		123654)

	err := suite.database.SaveShieldPool(shield)
	suite.Require().NoError(err)

	err = suite.database.SaveShieldPool(shield2)
	suite.Require().NoError(err)

	// Verify the data
	expectedPools := []dbtypes.ShieldPoolRow{
		dbtypes.NewShieldPoolRow(
			1,
			"1000000000",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000000000)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(2000000000)))),
			"Sponsor",
			purchaserAddress1.String(),
			"Shield Description",
			"1000000000000",
			false,
			123654,
		),
		dbtypes.NewShieldPoolRow(
			2,
			"1000000000",
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(3000000000)))),
			dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000000)))),
			"Sponsor",
			purchaserAddress2.String(),
			"Shield Description",
			"1000000000000",
			false,
			123654,
		),
	}

	var poolRows []dbtypes.ShieldPoolRow
	err = suite.database.Sqlx.Select(&poolRows, `SELECT * FROM shield_pool`)
	suite.Require().NoError(err)

	for i, row := range poolRows {
		suite.Require().True(expectedPools[i].Equal(row))
	}

	// -----------------------------------------------------------------------------------
	// Save the data
	shieldPurchase1 := types.NewShieldPurchase(
		1,
		purchaserAddress1.String(),
		sdk.NewInt(1000000000),
		"Shield Description",
		652211)

	shieldPurchase2 := types.NewShieldPurchase(
		2,
		purchaserAddress2.String(),
		sdk.NewInt(3000000000),
		"Shield Description",
		652211)

	err = suite.database.SaveShieldPurchase(shieldPurchase1)
	suite.Require().NoError(err)

	err = suite.database.SaveShieldPurchase(shieldPurchase2)
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.ShieldPurchaseRow{
		dbtypes.NewShieldPurchaseRow(
			1,
			purchaserAddress1.String(),
			"1000000000",
			"Shield Description",
			652211,
		),
		dbtypes.NewShieldPurchaseRow(
			2,
			purchaserAddress2.String(),
			"3000000000",
			"Shield Description",
			652211,
		),
	}

	var rows []dbtypes.ShieldPurchaseRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM shield_purchase`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equal(row))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_ShieldWithdraws() {
	withdrawAddress1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	withdrawAddress2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	// Save the data
	shieldWithdraw1 := types.NewShieldWithdraw(
		withdrawAddress1.String(),
		10000000,
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		123311)

	shieldWithdraw2 := types.NewShieldWithdraw(
		withdrawAddress2.String(),
		40000000,
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		123311)

	err := suite.database.SaveShieldWithdraw(shieldWithdraw1)
	suite.Require().NoError(err)

	err = suite.database.SaveShieldWithdraw(shieldWithdraw2)
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.ShieldWithdrawRow{
		dbtypes.NewShieldWithdrawRow(
			withdrawAddress1.String(),
			10000000,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			123311,
		),
		dbtypes.NewShieldWithdrawRow(
			withdrawAddress2.String(),
			40000000,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			123311,
		),
	}

	var rows []dbtypes.ShieldWithdrawRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM shield_withdraws`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equal(row))
	}

}
