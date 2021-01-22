package database_test

import (
	"fmt"

	consensustypes "github.com/tendermint/tendermint/types"

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

// _________________________________________________________

func (suite *DbTestSuite) TestSaveValidator() {
	expectedMaxRate := sdk.NewDec(int64(1))
	expectedMaxChangeRate := sdk.NewDec(int64(2))

	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
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
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
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
		suite.Require().Equal(v.ConsAddress, w.GetConsAddr())
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
	suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	// Insert test data
	_, err := suite.database.Sql.Exec(`
INSERT INTO validator (consensus_address, consensus_pubkey) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 
        'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`
INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address,max_change_rate,max_rate) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl',
        'cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl',
        'cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a','2','1')`)
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
		sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, validator.GetConsPubKey()),
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

	var height int64 = 1
	description := types.NewValidatorDescription(
		validator.GetOperator(),
		stakingtypes.NewDescription(
			"moniker",
			"identity",
			"",
			"securityContact",
			"details",
		),
		height,
	)
	err := suite.database.SaveValidatorDescription(description)
	suite.Require().NoError(err)

	var rows []dbtypes.ValidatorDescriptionRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM validator_description")
	suite.Require().NoError(err)

	expectedRows := []dbtypes.ValidatorDescriptionRow{
		dbtypes.NewValidatorDescriptionRow(
			validator.GetConsAddr(),
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
	var height int64 = 1000
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	err := suite.database.SaveValidatorCommission(types.NewValidatorCommission(
		validator.GetOperator(),
		newDecPts(11, 3),
		newIntPtr(12),
		height,
	))
	suite.Require().NoError(err)

	var rows []dbtypes.ValidatorCommissionRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission`)
	suite.Require().NoError(err)

	expectedRows := []dbtypes.ValidatorCommissionRow{
		dbtypes.NewValidatorCommissionRow(
			validator.GetConsAddr(),
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

	_, err = suite.database.Sql.Exec(`
INSERT INTO validator (consensus_address, consensus_pubkey) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 
        'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	// Save the data
	uptime := types.NewValidatorUptime(valAddr.String(), 10, 100, 500)

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
	block := suite.getBlock(100)

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
		types.NewValidatorVotingPower(validator1.GetConsAddr(), 1000, block.Block.Height),
		types.NewValidatorVotingPower(validator2.GetConsAddr(), 2000, block.Block.Height),
	}

	for _, power := range votingPowers {
		err := suite.database.SaveValidatorVotingPower(power)
		suite.Require().NoError(err)
	}

	expected := []dbtypes.ValidatorVotingPowerRow{
		dbtypes.NewValidatorVotingPowerRow(
			validator1.GetConsAddr(),
			1000,
		),
		dbtypes.NewValidatorVotingPowerRow(
			validator2.GetConsAddr(),
			2000,
		),
	}

	var result []dbtypes.ValidatorVotingPowerRow
	err := suite.database.Sqlx.Select(&result, "SELECT * FROM validator_voting_power")
	suite.Require().NoError(err)

	for index, row := range result {
		suite.Require().True(row.Equal(expected[index]))
	}

}

//-----------------------------------------------------------

func (suite *DbTestSuite) TestSaveValidatorStatus() {
	block1 := suite.getBlock(10)
	block2 := suite.getBlock(20)

	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	err := suite.database.SaveValidatorStatus(types.NewValidatorStatus(
		validator.GetConsAddr(),
		1,
		false,
		block1.Block.Height,
	))
	suite.Require().NoError(err)

	var result []dbtypes.ValidatorStatusRow
	err = suite.database.Sqlx.Select(&result, "SELECT * FROM validator_status")
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(result[0].Equal(dbtypes.NewValidatorStatusRow(1, false, validator.GetConsAddr())))

	firstRow := dbtypes.NewValidatorStatusHistoryRow(
		1,
		false,
		block1.Block.Height,
		validator.GetConsAddr(),
	)

	var result2 []dbtypes.ValidatorStatusHistoryRow
	err = suite.database.Sqlx.Select(&result2, "SELECT * FROM validator_status_history")
	suite.Require().NoError(err)
	suite.Require().Len(result2, 1)
	suite.Require().True(result2[0].Equal(firstRow))

	// Second insert
	err = suite.database.SaveValidatorStatus(types.NewValidatorStatus(
		validator.GetConsAddr(),
		2,
		true,
		block2.Block.Height,
	))
	suite.Require().NoError(err)

	var result3 []dbtypes.ValidatorStatusRow
	err = suite.database.Sqlx.Select(&result3, "SELECT * FROM validator_status")
	suite.Require().NoError(err)
	suite.Require().Len(result3, 1)
	suite.Require().True(result3[0].Equal(dbtypes.NewValidatorStatusRow(2, true, validator.GetConsAddr())))

	expected := []dbtypes.ValidatorStatusHistoryRow{
		firstRow,
		dbtypes.NewValidatorStatusHistoryRow(
			2,
			true,
			block2.Block.Height,
			validator.GetConsAddr(),
		),
	}

	var result4 []dbtypes.ValidatorStatusHistoryRow
	err = suite.database.Sqlx.Select(&result4, "SELECT * FROM validator_status_history")
	suite.Require().NoError(err)
	suite.Require().Len(result4, 2)
	for index, row := range result4 {
		suite.Require().True(row.Equal(expected[index]))
	}

}

//--------------------------------------------
func (suite *DbTestSuite) TestSaveDoubleVoteEvidence() {
	// Insert the validator
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	// Insert data
	evidence := types.NewDoubleSignEvidence(
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
		types.NewDoubleSignVote(
			int(consensustypes.PrevoteType),
			10,
			1,
			"A42C9492F5DE01BFA6117137102C3EF909F1A46C2F56915F542D12AC2D0A5BCA",
			validator.GetConsAddr(),
			1,
			"1qwPQjPrc7DH7+f6YAE3fOkq6phDAJ60dEyhmcZ7dx2ZgGvi9DbVLsn4leYqRNA/63ZeeH5kVly8zI1jCh4iBg==",
		),
		types.NewDoubleSignVote(
			int(consensustypes.PrevoteType),
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
	suite.Require().Equal(dbtypes.NewDoubleSignEvidenceRow(
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
		1,
		2,
	), evidenceRows[0])

	expectVotes := []dbtypes.DoubleSignVoteRow{
		dbtypes.NewDoubleSignVoteRow(
			1,
			int(consensustypes.PrevoteType),
			10,
			1,
			"A42C9492F5DE01BFA6117137102C3EF909F1A46C2F56915F542D12AC2D0A5BCA",
			validator.GetConsAddr(),
			1,
			"1qwPQjPrc7DH7+f6YAE3fOkq6phDAJ60dEyhmcZ7dx2ZgGvi9DbVLsn4leYqRNA/63ZeeH5kVly8zI1jCh4iBg==",
		),
		dbtypes.NewDoubleSignVoteRow(
			2,
			int(consensustypes.PrevoteType),
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
