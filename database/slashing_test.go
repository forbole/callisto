package database_test

import (
	"time"

	stakingtypes "github.com/forbole/bdjuno/x/staking/types"

	slashingtypes "github.com/forbole/bdjuno/x/slashing/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
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

	usecases := []struct {
		name     string
		input    []slashingtypes.ValidatorSigningInfo
		expected []dbtypes.ValidatorSigningInfoRow
	}{
		{
			name: "different validators are added properly",
			input: []slashingtypes.ValidatorSigningInfo{
				slashingtypes.NewValidatorSigningInfo(
					validator1.GetConsAddr(),
					10,
					10,
					time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
					true,
					10,
					10,
				),
				slashingtypes.NewValidatorSigningInfo(
					validator2.GetConsAddr(),
					10,
					10,
					time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
					true,
					10,
					10,
				),
			},
			expected: []dbtypes.ValidatorSigningInfoRow{
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
			},
		},
		{
			name: "same validator does not error and is not replaced",
			input: []slashingtypes.ValidatorSigningInfo{
				slashingtypes.NewValidatorSigningInfo(
					validator1.GetConsAddr(),
					10,
					10,
					time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
					true,
					10,
					10,
				),
				slashingtypes.NewValidatorSigningInfo(
					validator1.GetConsAddr(),
					11,
					15,
					time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
					false,
					500,
					10,
				),
			},
			expected: []dbtypes.ValidatorSigningInfoRow{
				dbtypes.NewValidatorSigningInfoRow(
					validator1.GetConsAddr(),
					10,
					10,
					time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
					true,
					10,
					10,
				),
			},
		},
	}

	for _, uc := range usecases {
		uc := uc
		suite.Run(uc.name, func() {
			suite.SetupTest()

			for _, input := range uc.input {
				suite.getBlock(input.Height)
			}

			err := suite.database.SaveValidators([]stakingtypes.Validator{validator1, validator2})
			suite.Require().NoError(err)

			err = suite.database.SaveValidatorsSigningInfos(uc.input)
			suite.Require().NoError(err)

			var rows []dbtypes.ValidatorSigningInfoRow
			err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_signing_info`)
			suite.Require().NoError(err)

			for i, row := range rows {
				suite.Require().True(uc.expected[i].Equal(row))
			}
		})
	}

}
