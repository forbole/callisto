package database_test

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/staking/types"
)

func newDecPts(value int64, prec int64) *sdk.Dec {
	dec := sdk.NewDecWithPrec(value, prec)
	return &dec
}

func newIntPtr(value int64) *sdk.Int {
	val := sdk.NewInt(value)
	return &val
}

func (suite *DbTestSuite) getValidator(consAddr, valAddr, pubkey string) types.Validator {
	selfDelegation := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	valAddrObj, err := sdk.ValAddressFromBech32(valAddr)
	suite.Require().NoError(err)

	constAddrObj, err := sdk.ConsAddressFromBech32(consAddr)
	suite.Require().NoError(err)

	pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pubkey)
	suite.Require().NoError(err)

	maxRate := sdk.NewDec(10)
	maxChangeRate := sdk.NewDec(20)

	validator := types.NewValidator(constAddrObj, valAddrObj, pubKey, selfDelegation, &maxChangeRate, &maxRate)
	err = suite.database.SaveValidatorData(validator)
	suite.Require().NoError(err)

	return validator
}

func (suite *DbTestSuite) getDelegator(addr string) sdk.AccAddress {
	delegator, err := sdk.AccAddressFromBech32(addr)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`, delegator.String())
	suite.Require().NoError(err)

	return delegator
}

// _________________________________________________________

func (suite *DbTestSuite) TestSaveValidator() {
	expectedMaxRate := sdk.NewDec(int64(1))
	expectedMaxChangeRate := sdk.NewDec(int64(2))

	suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	validator := dbtypes.NewValidatorData(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"1", "2",
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
	fmt.Print(valInfoRows[0])
	suite.Require().True(valInfoRows[0].Equal(dbtypes.NewValidatorInfoRow(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		expectedMaxChangeRate.String(), expectedMaxRate.String(),
	)))

}

func (suite *DbTestSuite) TestSaveValidators() {
	suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	expectedMaxRate := sdk.NewDec(int64(1))
	expectedMaxChangeRate := sdk.NewDec(int64(2))

	validators := []types.Validator{
		dbtypes.NewValidatorData(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
			"1", "2",
		),
		dbtypes.NewValidatorData(
			"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
			"1", "2",
		),
	}

	expectedValidatorInfo := []dbtypes.ValidatorInfoRow{
		dbtypes.NewValidatorInfoRow("cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
			expectedMaxChangeRate.String(), expectedMaxRate.String(),
		),
		dbtypes.NewValidatorInfoRow("cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
			expectedMaxChangeRate.String(), expectedMaxRate.String(),
		),
	}

	// Insert the data
	err := suite.database.SaveValidators(validators)

	suite.Require().NoError(err)

	// Verify the inserted data
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
		suite.Require().Equal(v.ConsAddress, w.GetConsAddr().String())
		suite.Require().Equal(v.ConsPubKey, sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, w.GetConsPubKey()))

		wInfo := validatorInfoRows[index]
		suite.Require().True(wInfo.Equal(expectedValidatorInfo[index]))
	}
}

func (suite *DbTestSuite) TestGetValidator() {
	var i int64 = 1
	var ii int64 = 2
	maxRate := sdk.NewDec(i)
	maxChangeRate := sdk.NewDec(ii)
	suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	// Insert test data
	_, err := suite.database.Sql.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address,max_change_rate,max_rate) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl','cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl','cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a','2','1')`)
	suite.Require().NoError(err)

	// Get the data
	valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl")
	validator, err := suite.database.GetValidator(valAddr)
	suite.Require().NoError(err)
	suite.Require().Equal(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		validator.GetConsAddr().String(),
	)
	suite.Require().Equal(
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		validator.GetOperator().String(),
	)
	suite.Require().Equal(
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
		sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, validator.GetConsPubKey()),
	)

	suite.Require().Equal("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a", validator.GetSelfDelegateAddress().String())
	suite.Require().True(validator.GetMaxChangeRate().Equal(maxChangeRate))
	suite.Require().True(validator.GetMaxRate().Equal(maxRate))

}

func (suite *DbTestSuite) TestGetValidators() {
	suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	// Inser the test data
	queries := []string{
		`INSERT INTO validator (consensus_address, consensus_pubkey) VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`,
		`INSERT INTO validator (consensus_address, consensus_pubkey) VALUES ('cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y', 'cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk')`,
		`INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address,max_rate,max_change_rate) VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl','cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs','1','2')`,
		`INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address,max_rate,max_change_rate) VALUES ('cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y', 'cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn','cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a','1','2')`,
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
			"1", "2",
		),
		dbtypes.NewValidatorData(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
			"1", "2",
		),
	}

	suite.Require().Len(data, len(expected))
	for index, validator := range data {
		suite.Require().Equal(expected[index], validator)
	}
}

// _________________________________________________________

func (suite *DbTestSuite) TestSaveValidatorDescription() {
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	timestamp, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	var height int64 = 1
	description := types.NewValidatorDescription(validator.GetOperator(), stakingtypes.NewDescription(
		"moniker",
		"identity",
		"",
		"securityContact",
		"details",
	), height, timestamp)
	err = suite.database.SaveValidatorDescription(description)
	suite.Require().NoError(err)

	var rows []dbtypes.ValidatorDescriptionRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM validator_description")
	suite.Require().NoError(err)

	expectedRows := []dbtypes.ValidatorDescriptionRow{
		dbtypes.NewValidatorDescriptionRow(
			validator.GetConsAddr().String(),
			"moniker",
			"identity",
			"",
			"securityContact",
			"details",
		),
	}
	suite.Require().Len(rows, len(expectedRows))
	for index, expected := range expectedRows {
		suite.Require().True(expected.Equals(rows[index]))
	}
}

// _________________________________________________________

func (suite *DbTestSuite) TestSaveValidatorCommission() {
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	var height int64 = 1000

	timestamp, err := time.Parse(time.RFC3339, "2020-01-01T10:00:00Z")
	suite.Require().NoError(err)

	err = suite.database.SaveValidatorCommission(types.NewValidatorCommission(
		validator.GetOperator(),
		newDecPts(11, 3),
		newIntPtr(12),
		height,
		timestamp,
	))
	suite.Require().NoError(err)

	var rows []dbtypes.ValidatorCommissionRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission`)
	suite.Require().NoError(err)

	expectedRows := []dbtypes.ValidatorCommissionRow{
		dbtypes.NewValidatorCommissionRow(
			validator.GetConsAddr().String(),
			"0.011000000000000000",
			"12",
		),
	}
	suite.Require().Len(rows, len(expectedRows))
	for index, expected := range expectedRows {
		suite.Require().True(expected.Equal(rows[index]))
	}
}

// _________________________________________________________
func (suite *DbTestSuite) TestSaveValidatorUptime() {
	valAddr, err := sdk.ConsAddressFromBech32("cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl")
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	// Save the data
	timestamp := time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC)
	uptime := types.NewValidatorUptime(valAddr, 10, 100, 500, timestamp)

	err = suite.database.SaveValidatorUptime(uptime)
	suite.Require().NoError(err, "validator uptime should not error while inserting")

	err = suite.database.SaveValidatorUptime(uptime)
	suite.Require().NoError(err, "double validator uptime insertion should not error")

	// Verify the data
	var validatorData []dbtypes.ValidatorUptimeRow
	err = suite.database.Sqlx.Select(&validatorData, `SELECT * FROM validator_uptime`)
	suite.Require().NoError(err)
	suite.Require().Len(validatorData, 1)
	suite.Require().Equal(validatorData[0], dbtypes.NewValidatorUptimeRow(
		valAddr.String(),
		10,
		100,
	))
}

// _________________________________________________________

func (suite *DbTestSuite) TestSaveValidatorsVotingPowers() {
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

	votingPowers := []types.ValidatorVotingPower{
		types.NewValidatorVotingPower(
			validator1.GetConsAddr(),
			1000,
			100,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		),
		types.NewValidatorVotingPower(
			validator2.GetConsAddr(),
			2000,
			100,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		),
	}
	err := suite.database.SaveValidatorsVotingPowers(votingPowers)
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorVotingPowerRow{
		dbtypes.NewValidatorVotingPowerRow(
			validator1.GetConsAddr().String(),
			1000,
		),
		dbtypes.NewValidatorVotingPowerRow(
			validator2.GetConsAddr().String(),
			2000,
		),
	}

	var result []dbtypes.ValidatorVotingPowerRow
	err = suite.database.Sqlx.Select(&result, "SELECT * FROM validator_voting_power")
	suite.Require().NoError(err)

	for index, row := range result {
		suite.Require().True(row.Equal(expected[index]))
	}

}
