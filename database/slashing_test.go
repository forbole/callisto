package database_test

import (
	"time"

	slashingtypes "github.com/forbole/bdjuno/x/slashing/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_ValidatorSigningInfo() {
	timestamp, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)
	validator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	input := slashingtypes.NewValidatorSigningInfo(
		validator1.GetConsAddr(),
		10,
		10,
		timestamp,
		true,
		10,
		10,
		timestamp,
	)
	err = suite.database.SaveValidatorSigningInfo(input)
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorSigningInfoRow{
		dbtypes.NewValidatorSigningInfoRow(
			validator1.GetConsAddr().String(),
			10,
			10,
			timestamp,
			true,
			10,
			10,
			timestamp),
	}

	var rows []dbtypes.ValidatorSigningInfoRow
	err = suite.database.Sqlx.Select(&rows, `SELECT validator_address,start_height,index_offset,jailed_until,tombstoned,missed_blocks_counter,height,timestamp FROM validator_signing_info`)

	suite.Require().NoError(err)

	for i, row := range rows {

		suite.Require().True(expected[i].Equal(row))
	}
}
