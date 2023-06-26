package database_test

import (
	"encoding/json"
	"time"

	"github.com/forbole/bdjuno/v4/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	dbtypes "github.com/forbole/bdjuno/v4/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_ValidatorSigningInfo() {
	validator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	validator2 := suite.getValidator(
		"cosmosvalcons1rtst6se0nfgjy362v33jt5d05crgdyhfvvvvay",
		"cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr",
		"cosmosvalconspub1zcjduepq5e8w7t7k9pwfewgrwy8vn6cghk0x49chx64vt0054yl4wwsmjgrqfackxm",
	)

	// Save the data
	infos := []types.ValidatorSigningInfo{
		types.NewValidatorSigningInfo(
			validator1.GetConsAddr(),
			10,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
			true,
			10,
			10,
		),
		types.NewValidatorSigningInfo(
			validator2.GetConsAddr(),
			10,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
			true,
			10,
			10,
		),
	}
	err := suite.database.SaveValidatorsSigningInfos(infos)
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.ValidatorSigningInfoRow{
		dbtypes.NewValidatorSigningInfoRow(
			validator1.GetConsAddr(),
			10,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
			true,
			10,
			10,
		),
		dbtypes.NewValidatorSigningInfoRow(
			validator2.GetConsAddr(),
			10,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
			true,
			10,
			10,
		),
	}

	var rows []dbtypes.ValidatorSigningInfoRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_signing_info`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equal(row))
	}

	// ----------------------------------------------------------------------------------------------------------------

	// Update the data
	infos = []types.ValidatorSigningInfo{
		types.NewValidatorSigningInfo(
			validator1.GetConsAddr(),
			100,
			10000,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
			true,
			70,
			9,
		),
		types.NewValidatorSigningInfo(
			validator2.GetConsAddr(),
			10,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
			false,
			11,
			11,
		),
	}
	err = suite.database.SaveValidatorsSigningInfos(infos)
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.ValidatorSigningInfoRow{
		dbtypes.NewValidatorSigningInfoRow(
			validator1.GetConsAddr(),
			10,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
			true,
			10,
			10,
		),
		dbtypes.NewValidatorSigningInfoRow(
			validator2.GetConsAddr(),
			10,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
			false,
			11,
			11,
		),
	}

	rows = []dbtypes.ValidatorSigningInfoRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_signing_info`)
	suite.Require().NoError(err)

	for i, row := range rows {
		suite.Require().True(expected[i].Equal(row))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveSlashingParams() {
	// Save data
	slashingParams := slashingtypes.Params{
		SignedBlocksWindow:      10,
		MinSignedPerWindow:      sdk.NewDecWithPrec(100, 2),
		DowntimeJailDuration:    10000,
		SlashFractionDoubleSign: sdk.NewDecWithPrec(100, 2),
		SlashFractionDowntime:   sdk.NewDecWithPrec(100, 4),
	}
	params := types.NewSlashingParams(slashingParams, 10)
	err := suite.database.SaveSlashingParams(params)
	suite.Require().NoError(err)

	// Verify data
	var rows []dbtypes.SlashingParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM slashing_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var stored slashingtypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &stored)
	suite.Require().NoError(err)
	suite.Require().Equal(slashingParams, stored)

	// --- Try updating with a lower height ---
	err = suite.database.SaveSlashingParams(types.NewSlashingParams(
		slashingtypes.Params{
			SignedBlocksWindow:      5,
			MinSignedPerWindow:      sdk.NewDecWithPrec(50, 2),
			DowntimeJailDuration:    10000,
			SlashFractionDoubleSign: sdk.NewDecWithPrec(50, 2),
			SlashFractionDowntime:   sdk.NewDecWithPrec(50, 4),
		},
		9,
	))
	suite.Require().NoError(err)

	rows = []dbtypes.SlashingParamsRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM slashing_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	err = json.Unmarshal([]byte(rows[0].Params), &stored)
	suite.Require().NoError(err)
	suite.Require().Equal(slashingParams, stored)

	// Try updating with same height
	slashingParams = slashingtypes.Params{
		SignedBlocksWindow:      5,
		MinSignedPerWindow:      sdk.NewDecWithPrec(50, 2),
		DowntimeJailDuration:    10000,
		SlashFractionDoubleSign: sdk.NewDecWithPrec(50, 2),
		SlashFractionDowntime:   sdk.NewDecWithPrec(50, 4),
	}
	err = suite.database.SaveSlashingParams(types.NewSlashingParams(slashingParams, 10))
	suite.Require().NoError(err)

	rows = []dbtypes.SlashingParamsRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM slashing_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	err = json.Unmarshal([]byte(rows[0].Params), &stored)
	suite.Require().NoError(err)
	suite.Require().Equal(slashingParams, stored)

	// Try updating with higher height
	slashingParams = slashingtypes.Params{
		SignedBlocksWindow:      6,
		MinSignedPerWindow:      sdk.NewDecWithPrec(60, 2),
		DowntimeJailDuration:    10000,
		SlashFractionDoubleSign: sdk.NewDecWithPrec(60, 2),
		SlashFractionDowntime:   sdk.NewDecWithPrec(60, 4),
	}
	err = suite.database.SaveSlashingParams(types.NewSlashingParams(slashingParams, 11))
	suite.Require().NoError(err)

	rows = []dbtypes.SlashingParamsRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM slashing_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	err = json.Unmarshal([]byte(rows[0].Params), &stored)
	suite.Require().NoError(err)
	suite.Require().Equal(slashingParams, stored)
}
