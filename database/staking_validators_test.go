package database_test

import (
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/forbole/bdjuno/v3/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

func newDecPts(value int64, prec int64) *sdk.Dec {
	dec := sdk.NewDecWithPrec(value, prec)
	return &dec
}

func newIntPtr(value int64) *sdk.Int {
	val := sdk.NewInt(value)
	return &val
}

// -----------------------------------------------------------

func (suite *DbTestSuite) TestSaveValidator() {
	expectedMaxRate := sdk.NewDec(int64(1))
	expectedMaxChangeRate := sdk.NewDec(int64(2))

	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	validator := dbtypes.NewValidatorData(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"1",
		"2",
		1,
	)

	// First inserting
	err := suite.database.SaveValidatorData(validator)

	// Test double inserting
	err = suite.database.SaveValidatorData(validator)
	suite.Require().NoError(err, "inserting the same validator info twice should return no error")

	// Verify the data
	var valRows []dbtypes.ValidatorRow
	err = suite.database.Sqlx.Select(&valRows, `SELECT * FROM validator`)
	suite.Require().Len(valRows, 1)
	suite.Require().True(valRows[0].Equal(dbtypes.NewValidatorRow(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)))

	var valInfoRows []dbtypes.ValidatorInfoRow
	err = suite.database.Sqlx.Select(&valInfoRows, `SELECT * FROM validator_info`)
	suite.Require().Len(valInfoRows, 1)
	suite.Require().True(valInfoRows[0].Equal(dbtypes.NewValidatorInfoRow(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		expectedMaxRate.String(),
		expectedMaxChangeRate.String(),
		1,
	)))
}

func (suite *DbTestSuite) TestSaveValidators() {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	// Insert the data
	validators := []types.Validator{
		dbtypes.NewValidatorData(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
			"1",
			"2",
			10,
		),
		dbtypes.NewValidatorData(
			"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
			"1",
			"2",
			10,
		),
	}
	err := suite.database.SaveValidatorsData(validators)
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.ValidatorInfoRow{
		dbtypes.NewValidatorInfoRow(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
			sdk.NewDec(int64(1)).String(),
			sdk.NewDec(int64(2)).String(),
			10,
		),
		dbtypes.NewValidatorInfoRow(
			"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
			sdk.NewDec(int64(1)).String(),
			sdk.NewDec(int64(2)).String(),
			10,
		),
	}

	var validatorRows []dbtypes.ValidatorRow
	err = suite.database.Sqlx.Select(&validatorRows, `SELECT * FROM validator`)
	suite.Require().NoError(err)
	suite.Require().Len(validatorRows, 2)

	var validatorInfoRows []dbtypes.ValidatorInfoRow
	err = suite.database.Sqlx.Select(&validatorInfoRows, `SELECT * FROM validator_info`)
	suite.Require().NoError(err)
	suite.Require().Len(validatorInfoRows, 2)

	for index, v := range validatorRows {
		w := validators[index]
		suite.Require().Equal(v.ConsAddress, w.GetConsAddr())
		suite.Require().Equal(v.ConsPubKey, w.GetConsPubKey())

		wInfo := validatorInfoRows[index]
		suite.Require().True(wInfo.Equal(expected[index]))
	}

	// --------------------------------------------------------------------------------------------------------------

	// Update the data
	validators = []types.Validator{
		dbtypes.NewValidatorData(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
			"100",
			"200",
			9,
		),
		dbtypes.NewValidatorData(
			"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
			"10",
			"5",
			11,
		),
	}
	err = suite.database.SaveValidatorsData(validators)
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.ValidatorInfoRow{
		dbtypes.NewValidatorInfoRow(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
			sdk.NewDec(int64(1)).String(),
			sdk.NewDec(int64(2)).String(),
			10,
		),
		dbtypes.NewValidatorInfoRow(
			"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
			sdk.NewDec(int64(10)).String(),
			sdk.NewDec(int64(5)).String(),
			11,
		),
	}

	validatorRows = []dbtypes.ValidatorRow{}
	err = suite.database.Sqlx.Select(&validatorRows, `SELECT * FROM validator`)
	suite.Require().NoError(err)
	suite.Require().Len(validatorRows, 2)

	validatorInfoRows = []dbtypes.ValidatorInfoRow{}
	err = suite.database.Sqlx.Select(&validatorInfoRows, `SELECT * FROM validator_info`)
	suite.Require().NoError(err)
	suite.Require().Len(validatorInfoRows, 2)

	for index, v := range validatorRows {
		w := validators[index]
		suite.Require().Equal(v.ConsAddress, w.GetConsAddr())
		suite.Require().Equal(v.ConsPubKey, w.GetConsPubKey())

		wInfo := validatorInfoRows[index]
		suite.Require().True(wInfo.Equal(expected[index]))
	}
}

func (suite *DbTestSuite) TestGetValidator() {
	var i int64 = 1
	var ii int64 = 2
	maxRate := sdk.NewDec(i)
	maxChangeRate := sdk.NewDec(ii)
	suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	// Insert test data
	_, err := suite.database.Sql.Exec(`
INSERT INTO validator (consensus_address, consensus_pubkey) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 
        'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`
INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address,max_change_rate,max_rate,height) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl',
        'cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl',
        'cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a',
        '2','1', '1')`)
	suite.Require().NoError(err)

	// Get the data
	validator, err := suite.database.GetValidator("cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl")
	suite.Require().NoError(err)
	suite.Require().Equal(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		validator.GetConsAddr(),
	)
	suite.Require().Equal(
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		validator.GetOperator(),
	)
	suite.Require().Equal(
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
		validator.GetConsPubKey(),
	)

	suite.Require().Equal("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a", validator.GetSelfDelegateAddress())
	suite.Require().True(validator.GetMaxChangeRate().Equal(maxChangeRate))
	suite.Require().True(validator.GetMaxRate().Equal(maxRate))

}

func (suite *DbTestSuite) TestGetValidators() {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	// Inser the test data
	queries := []string{
		`INSERT INTO validator (consensus_address, consensus_pubkey) VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`,
		`INSERT INTO validator (consensus_address, consensus_pubkey) VALUES ('cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y', 'cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk')`,
		`INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address,max_rate,max_change_rate,height) VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl','cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs','1','2',1)`,
		`INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address,max_rate,max_change_rate,height) VALUES ('cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y', 'cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn','cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a','1','2',1)`,
	}

	for _, query := range queries {
		_, err := suite.database.Sql.Exec(query)
		suite.Require().NoError(err)
	}

	// Get the data
	data, err := suite.database.GetValidators()
	suite.Require().NoError(err)

	// Verify
	expected := []dbtypes.ValidatorData{
		dbtypes.NewValidatorData(
			"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
			"1",
			"2",
			1,
		),
		dbtypes.NewValidatorData(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
			"1",
			"2",
			1,
		),
	}

	suite.Require().Len(data, len(expected))
	for index, validator := range data {
		suite.Require().Equal(expected[index], validator)
	}
}

// -----------------------------------------------------------

func (suite *DbTestSuite) TestSaveValidatorDescription() {
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	// Save the data
	description := types.NewValidatorDescription(
		validator.GetOperator(),
		stakingtypes.NewDescription(
			"moniker",
			"identity",
			"",
			"securityContact",
			"details",
		),
		"avatar-url",
		10,
	)
	err := suite.database.SaveValidatorDescription(description)
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.ValidatorDescriptionRow{
		dbtypes.NewValidatorDescriptionRow(
			validator.GetConsAddr(),
			"moniker",
			"identity",
			"avatar-url",
			"",
			"securityContact",
			"details",
			10,
		),
	}

	var rows []dbtypes.ValidatorDescriptionRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM validator_description")
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, expected := range expected {
		suite.Require().True(expected.Equals(rows[index]))
	}

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with lower height
	description = types.NewValidatorDescription(
		validator.GetOperator(),
		stakingtypes.NewDescription("moniker", "", "", "", ""),
		"lower-avatar-url",
		9,
	)
	err = suite.database.SaveValidatorDescription(description)
	suite.Require().NoError(err)

	// Verify the data
	rows = []dbtypes.ValidatorDescriptionRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM validator_description")
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, expected := range expected {
		suite.Require().True(expected.Equals(rows[index]))
	}

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with same height
	description = types.NewValidatorDescription(
		validator.GetOperator(),
		stakingtypes.NewDescription("moniker", "", "", "", ""),
		"new-avatar-url",
		10,
	)
	err = suite.database.SaveValidatorDescription(description)
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.ValidatorDescriptionRow{
		dbtypes.NewValidatorDescriptionRow(
			validator.GetConsAddr(),
			"moniker",
			"",
			"new-avatar-url",
			"",
			"",
			"",
			10,
		),
	}

	rows = []dbtypes.ValidatorDescriptionRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM validator_description")
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, expected := range expected {
		suite.Require().True(expected.Equals(rows[index]))
	}

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with higher height
	description = types.NewValidatorDescription(
		validator.GetOperator(),
		stakingtypes.NewDescription("moniker", "higher-identity", "higher-website", "", ""),
		"higher-avatar-url",
		11,
	)
	err = suite.database.SaveValidatorDescription(description)
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.ValidatorDescriptionRow{
		dbtypes.NewValidatorDescriptionRow(
			validator.GetConsAddr(),
			"moniker",
			"higher-identity",
			"higher-avatar-url",
			"higher-website",
			"",
			"",
			11,
		),
	}

	rows = []dbtypes.ValidatorDescriptionRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM validator_description")
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, expected := range expected {
		suite.Require().True(expected.Equals(rows[index]))
	}
}

// -----------------------------------------------------------

func (suite *DbTestSuite) TestSaveValidatorCommission() {
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	// Save the data
	err := suite.database.SaveValidatorCommission(types.NewValidatorCommission(
		validator.GetOperator(),
		newDecPts(11, 3),
		newIntPtr(12),
		10,
	))
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.ValidatorCommissionRow{
		dbtypes.NewValidatorCommissionRow(
			validator.GetConsAddr(),
			"0.011000000000000000",
			"12",
			10,
		),
	}

	var rows []dbtypes.ValidatorCommissionRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, expected := range expected {
		suite.Require().True(expected.Equal(rows[index]))
	}

	// -------------------------------------------------------------------------------------------------------------

	// Try updating with a lowe height
	err = suite.database.SaveValidatorCommission(types.NewValidatorCommission(
		validator.GetOperator(),
		newDecPts(50, 3),
		newIntPtr(100),
		9,
	))
	suite.Require().NoError(err)

	// Verify the data
	rows = []dbtypes.ValidatorCommissionRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, expected := range expected {
		suite.Require().True(expected.Equal(rows[index]))
	}

	// -------------------------------------------------------------------------------------------------------------

	// Try updating with the same height
	err = suite.database.SaveValidatorCommission(types.NewValidatorCommission(
		validator.GetOperator(),
		newDecPts(50, 3),
		newIntPtr(100),
		10,
	))
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.ValidatorCommissionRow{
		dbtypes.NewValidatorCommissionRow(
			validator.GetConsAddr(),
			"0.050000000000000000",
			"100",
			10,
		),
	}

	rows = []dbtypes.ValidatorCommissionRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, expected := range expected {
		suite.Require().True(expected.Equal(rows[index]))
	}

	// -------------------------------------------------------------------------------------------------------------

	// Try updating with a higher height
	err = suite.database.SaveValidatorCommission(types.NewValidatorCommission(
		validator.GetOperator(),
		newDecPts(70, 2),
		newIntPtr(200),
		11,
	))
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.ValidatorCommissionRow{
		dbtypes.NewValidatorCommissionRow(
			validator.GetConsAddr(),
			"0.700000000000000000",
			"200",
			11,
		),
	}

	rows = []dbtypes.ValidatorCommissionRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, len(expected))

	for index, expected := range expected {
		suite.Require().True(expected.Equal(rows[index]))
	}
}

// -----------------------------------------------------------

func (suite *DbTestSuite) TestSaveValidatorsVotingPowers() {
	_ = suite.getBlock(9)
	_ = suite.getBlock(10)
	_ = suite.getBlock(11)

	validator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	validator2 := suite.getValidator(
		"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
		"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
		"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
	)

	// Save data
	err := suite.database.SaveValidatorsVotingPowers([]types.ValidatorVotingPower{
		types.NewValidatorVotingPower(validator1.GetConsAddr(), 1000, 10),
		types.NewValidatorVotingPower(validator2.GetConsAddr(), 2000, 10),
	})
	suite.Require().NoError(err)

	// Verify the data
	expected := []dbtypes.ValidatorVotingPowerRow{
		dbtypes.NewValidatorVotingPowerRow(validator1.GetConsAddr(), 1000, 10),
		dbtypes.NewValidatorVotingPowerRow(validator2.GetConsAddr(), 2000, 10),
	}

	var result []dbtypes.ValidatorVotingPowerRow
	err = suite.database.Sqlx.Select(&result, "SELECT * FROM validator_voting_power")
	suite.Require().NoError(err)

	for index, row := range result {
		suite.Require().True(row.Equal(expected[index]))
	}

	// Update the data
	err = suite.database.SaveValidatorsVotingPowers([]types.ValidatorVotingPower{
		types.NewValidatorVotingPower(validator1.GetConsAddr(), 5, 9),
		types.NewValidatorVotingPower(validator2.GetConsAddr(), 10, 11),
	})
	suite.Require().NoError(err)

	// Verify the data
	expected = []dbtypes.ValidatorVotingPowerRow{
		dbtypes.NewValidatorVotingPowerRow(validator1.GetConsAddr(), 1000, 10),
		dbtypes.NewValidatorVotingPowerRow(validator2.GetConsAddr(), 10, 11),
	}

	result = []dbtypes.ValidatorVotingPowerRow{}
	err = suite.database.Sqlx.Select(&result, "SELECT * FROM validator_voting_power")
	suite.Require().NoError(err)

	for index, row := range result {
		suite.Require().True(row.Equal(expected[index]))
	}
}

// -----------------------------------------------------------

func (suite *DbTestSuite) TestSaveValidatorStatus() {
	validator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	validator2 := suite.getValidator(
		"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
		"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
		"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
	)

	// Save the data
	err := suite.database.SaveValidatorsStatuses([]types.ValidatorStatus{
		types.NewValidatorStatus(
			validator1.GetConsAddr(),
			validator1.GetConsPubKey(),
			1,
			false,
			false,
			10,
		),
		types.NewValidatorStatus(
			validator2.GetConsAddr(),
			validator2.GetConsPubKey(),
			2,
			true,
			true,
			10,
		),
	})
	suite.Require().NoError(err)

	// Verify the data
	var stored []dbtypes.ValidatorStatusRow
	err = suite.database.Sqlx.Select(&stored, "SELECT * FROM validator_status")
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorStatusRow{
		dbtypes.NewValidatorStatusRow(
			1,
			false,
			false,
			validator1.GetConsAddr(),
			10,
		),
		dbtypes.NewValidatorStatusRow(
			2,
			true,
			true,
			validator2.GetConsAddr(),
			10,
		),
	}
	suite.Require().Len(stored, len(expected))
	for index, stored := range stored {
		suite.Require().True(stored.Equal(expected[index]))
	}

	// Update the data
	err = suite.database.SaveValidatorsStatuses([]types.ValidatorStatus{
		types.NewValidatorStatus(
			validator1.GetConsAddr(),
			validator1.GetConsPubKey(),
			3,
			true,
			true,
			9,
		),
		types.NewValidatorStatus(
			validator2.GetConsAddr(),
			validator2.GetConsPubKey(),
			3,
			true,
			true,
			11,
		),
	})
	suite.Require().NoError(err)

	// Verify the data
	stored = []dbtypes.ValidatorStatusRow{}
	err = suite.database.Sqlx.Select(&stored, "SELECT * FROM validator_status")
	suite.Require().NoError(err)

	expected = []dbtypes.ValidatorStatusRow{
		dbtypes.NewValidatorStatusRow(
			1,
			false,
			false,
			validator1.GetConsAddr(),
			10,
		),
		dbtypes.NewValidatorStatusRow(
			3,
			true,
			true,
			validator2.GetConsAddr(),
			11,
		),
	}
	suite.Require().Len(stored, len(expected))
	for index, stored := range stored {
		suite.Require().True(stored.Equal(expected[index]))
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestSaveDoubleVoteEvidence() {
	// Insert the validator
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	// Insert data
	evidence := types.NewDoubleSignEvidence(
		10,
		types.NewDoubleSignVote(
			int(tmtypes.PrevoteType),
			10,
			1,
			"A42C9492F5DE01BFA6117137102C3EF909F1A46C2F56915F542D12AC2D0A5BCA",
			validator.GetConsAddr(),
			1,
			"1qwPQjPrc7DH7+f6YAE3fOkq6phDAJ60dEyhmcZ7dx2ZgGvi9DbVLsn4leYqRNA/63ZeeH5kVly8zI1jCh4iBg==",
		),
		types.NewDoubleSignVote(
			int(tmtypes.PrevoteType),
			10,
			1,
			"418A20D12F45FC9340BE0CD2EDB0FFA1E4316176B8CE11E123EF6CBED23C8423",
			validator.GetConsAddr(),
			1,
			"A5m7SVuvZ8YNXcUfBKLgkeV+Vy5ea+7rPfzlbkEvHOPPce6B7A2CwOIbCmPSVMKUarUdta+HiyTV+IELaOYyDA==",
		),
	)
	err := suite.database.SaveDoubleSignEvidence(evidence)
	suite.Require().NoError(err)

	// Verify insertion
	var evidenceRows []dbtypes.DoubleSignEvidenceRow
	err = suite.database.Sqlx.Select(&evidenceRows, "SELECT * FROM double_sign_evidence")
	suite.Require().NoError(err)
	suite.Require().Len(evidenceRows, 1)
	suite.Require().Equal(dbtypes.NewDoubleSignEvidenceRow(10, 1, 2), evidenceRows[0])

	expectVotes := []dbtypes.DoubleSignVoteRow{
		dbtypes.NewDoubleSignVoteRow(
			1,
			int(tmtypes.PrevoteType),
			10,
			1,
			"A42C9492F5DE01BFA6117137102C3EF909F1A46C2F56915F542D12AC2D0A5BCA",
			validator.GetConsAddr(),
			1,
			"1qwPQjPrc7DH7+f6YAE3fOkq6phDAJ60dEyhmcZ7dx2ZgGvi9DbVLsn4leYqRNA/63ZeeH5kVly8zI1jCh4iBg==",
		),
		dbtypes.NewDoubleSignVoteRow(
			2,
			int(tmtypes.PrevoteType),
			10,
			1,
			"418A20D12F45FC9340BE0CD2EDB0FFA1E4316176B8CE11E123EF6CBED23C8423",
			validator.GetConsAddr(),
			1,
			"A5m7SVuvZ8YNXcUfBKLgkeV+Vy5ea+7rPfzlbkEvHOPPce6B7A2CwOIbCmPSVMKUarUdta+HiyTV+IELaOYyDA==",
		),
	}

	var votesRows []dbtypes.DoubleSignVoteRow
	err = suite.database.Sqlx.Select(&votesRows, "SELECT * FROM double_sign_vote")
	suite.Require().NoError(err)

	suite.Require().Len(votesRows, len(expectVotes))
	for index, row := range votesRows {
		suite.Require().True(expectVotes[index].Equal(row))
	}
}
