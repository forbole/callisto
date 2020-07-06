package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/staking/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveStakingPool() {
	height := int64(100)

	timestamp, err := time.Parse(time.RFC3339, "2020-02-02T15:00:00Z")
	suite.Require().NoError(err)

	pool := stakingtypes.NewPool(sdk.NewInt(100), sdk.NewInt(50))

	// Save the data
	err = suite.database.SaveStakingPool(pool, height, timestamp)
	suite.Require().NoError(err)

	var count int
	err = suite.database.Sqlx.QueryRow(`SELECT COUNT(*) FROM staking_pool`).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(1, count, "inserting a single staking pool row should return 1")

	// Perform a double insertion
	err = suite.database.SaveStakingPool(pool, height, timestamp)
	suite.Require().NoError(err)

	err = suite.database.Sqlx.QueryRow(`SELECT COUNT(*) FROM staking_pool`).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(1, count, "double inserting the same staking pool should return 1 row")

	// Verify the data
	var rows []dbtypes.StakingPoolRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM staking_pool`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewStakingPoolRow(
		50,
		100,
		height,
		timestamp,
	)))
}

// _________________________________________________________

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorData() {
	validator := dbtypes.NewValidatorData(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	// First inserting
	err := suite.database.SaveValidatorData(validator)
	suite.Require().NoError(err, "inserting a validator info should return no error")

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
	)))
}

func (suite *DbTestSuite) TestBigDipperDb_GetValidatorData() {
	// Insert test data
	_, err := suite.database.Sql.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO validator_info (consensus_address, operator_address) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl','cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl')`)
	suite.Require().NoError(err)

	// Get the data
	valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl")
	validator, err := suite.database.GetValidatorData(valAddr)
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
}

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorsData() {
	validators := []types.Validator{
		dbtypes.NewValidatorData(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
		),
		dbtypes.NewValidatorData(
			"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
		),
	}

	// Insert the data
	err := suite.database.SaveValidatorsData(validators)
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
		suite.Require().Equal(wInfo.ConsAddress, w.GetConsAddr().String())
		suite.Require().Equal(wInfo.ValAddress, w.GetOperator().String())
	}
}

func (suite *DbTestSuite) TestBigDipperDb_GetValidatorsData() {
	// Inser the test data
	queries := []string{
		`INSERT INTO validator (consensus_address, consensus_pubkey) VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`,
		`INSERT INTO validator (consensus_address, consensus_pubkey) VALUES ('cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y', 'cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk')`,
		`INSERT INTO validator_info (consensus_address, operator_address) VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl')`,
		`INSERT INTO validator_info (consensus_address, operator_address) VALUES ('cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y', 'cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn')`,
	}

	for _, query := range queries {
		_, err := suite.database.Sql.Exec(query)
		suite.Require().NoError(err)
	}

	// Get the data
	data, err := suite.database.GetValidatorsData()
	suite.Require().NoError(err)

	// Verify
	expected := []dbtypes.ValidatorData{
		dbtypes.NewValidatorData(
			"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
		),
		dbtypes.NewValidatorData(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
		),
	}

	suite.Require().Len(data, len(expected))
	for index, validator := range data {
		suite.Require().Equal(validator, expected[index])
	}
}

// _________________________________________________________

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorUptime() {
	valAddr, err := sdk.ConsAddressFromBech32("cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl")
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	// Save the data
	uptime := types.NewValidatorUptime(valAddr, 10, 100, 500)

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
		500,
	))
}

// _________________________________________________________

func (suite *DbTestSuite) getValidator(consAddr, valAddr, pubkey string) types.Validator {
	valAddrObj, err := sdk.ValAddressFromBech32(valAddr)
	suite.Require().NoError(err)

	constAddrObj, err := sdk.ConsAddressFromBech32(consAddr)
	suite.Require().NoError(err)

	pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pubkey)
	suite.Require().NoError(err)

	validator := types.NewValidator(constAddrObj, valAddrObj, pubKey)
	err = suite.database.SaveValidatorData(validator)
	suite.Require().NoError(err)

	return validator
}

func (suite *DbTestSuite) getDelegator(addr string) sdk.AccAddress {
	delegator, err := sdk.AccAddressFromBech32(addr)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO account (address) VALUES ($1)`, delegator.String())
	suite.Require().NoError(err)

	return delegator
}

func (suite *DbTestSuite) TestBigDipperDb_SaveDelegation() {
	// Setup
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	delegator := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	var height int64 = 1000
	amount := sdk.NewCoin("cosmos", sdk.NewInt(10000))

	timestamp, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)

	delegation := types.NewDelegation(delegator, validator.GetOperator(), amount, height, timestamp)

	// Save data
	err = suite.database.SaveDelegation(delegation)
	suite.Require().NoError(err, "saving a delegation should return no error")

	// Get data
	var delegationRows []dbtypes.ValidatorDelegationRow
	err = suite.database.Sqlx.Select(&delegationRows, `SELECT * FROM validator_delegation`)
	suite.Require().NoError(err)

	suite.Require().Len(delegationRows, 1)
	suite.Require().True(delegationRows[0].Equal(dbtypes.NewValidatorDelegationRow(
		validator.GetConsAddr().String(),
		delegator.String(),
		dbtypes.NewDbCoin(amount),
		height,
		timestamp,
	)))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveDelegations() {
	// Setup
	delegator1 := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	delegator2 := suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
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

	time1, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	time2, err := time.Parse(time.RFC3339, "2020-05-05T18:00:00Z")
	suite.Require().NoError(err)

	// Save data
	delegations := []types.Delegation{
		types.NewDelegation(
			delegator1,
			validator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			1000,
			time1,
		),
		types.NewDelegation(
			delegator1,
			validator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			1000,
			time1,
		),
		types.NewDelegation(
			delegator1,
			validator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			1001,
			time1,
		),
		types.NewDelegation(
			delegator2,
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(200)),
			1500,
			time2,
		),
	}

	err = suite.database.SaveDelegations(delegations)
	suite.Require().NoError(err, "inserting delegations should return no error")

	// Verify the data
	var rows []dbtypes.ValidatorDelegationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_delegation`)
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorDelegationRow{
		dbtypes.NewValidatorDelegationRow(
			validator1.GetConsAddr().String(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			1000,
			time1,
		),
		dbtypes.NewValidatorDelegationRow(
			validator1.GetConsAddr().String(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			1000,
			time1,
		),
		dbtypes.NewValidatorDelegationRow(
			validator1.GetConsAddr().String(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			1001,
			time1,
		),
		dbtypes.NewValidatorDelegationRow(
			validator2.GetConsAddr().String(),
			delegator2.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			1500,
			time2,
		),
	}

	suite.Require().Len(rows, len(expected))
	for index, delegation := range expected {
		suite.Require().True(delegation.Equal(rows[index]))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveUnbondingDelegation() {
	// Setup
	delegator := suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	var height int64 = 1000
	amount := sdk.NewCoin("udesmos", sdk.NewInt(1000))

	completionTimestamp, err := time.Parse(time.RFC3339, "2020-08-10T16:00:00Z")
	suite.Require().NoError(err)

	timestamp, err := time.Parse(time.RFC3339, "2020-01-01T10:00:00Z")
	suite.Require().NoError(err)

	// Save data
	unbondingDelegation := types.NewUnbondingDelegation(
		delegator,
		validator.GetOperator(),
		amount,
		completionTimestamp,
		height,
		timestamp,
	)
	err = suite.database.SaveUnbondingDelegation(unbondingDelegation)
	suite.Require().NoError(err)

	// Get inserted data
	var rows []dbtypes.ValidatorUnbondingDelegationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_unbonding_delegation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewValidatorUnbondingDelegationRow(
		validator.GetConsAddr().String(),
		delegator.String(),
		dbtypes.NewDbCoin(amount),
		completionTimestamp, height,
		timestamp,
	)))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveUnbondingDelegations() {
	// Setup
	delegator1 := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	delegator2 := suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
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

	completionTimestamp1, err := time.Parse(time.RFC3339, "2020-08-10T16:00:00Z")
	suite.Require().NoError(err)

	completionTimestamp2, err := time.Parse(time.RFC3339, "2020-08-20T16:00:00Z")
	suite.Require().NoError(err)

	timestamp1, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	timestamp2, err := time.Parse(time.RFC3339, "2020-05-05T18:00:00Z")
	suite.Require().NoError(err)

	// Save data
	delegations := []types.UnbondingDelegation{
		types.NewUnbondingDelegation(
			delegator1,
			validator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			completionTimestamp1,
			1000,
			timestamp1,
		),
		types.NewUnbondingDelegation(
			delegator1,
			validator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			completionTimestamp1,
			1000,
			timestamp1,
		),
		types.NewUnbondingDelegation(
			delegator1,
			validator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			completionTimestamp2,
			1001,
			timestamp1,
		),
		types.NewUnbondingDelegation(
			delegator2,
			validator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(200)),
			completionTimestamp2,
			1500,
			timestamp2,
		),
	}
	err = suite.database.SaveUnbondingDelegations(delegations)
	suite.Require().NoError(err)

	// Read the data
	var rows []dbtypes.ValidatorUnbondingDelegationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_unbonding_delegation`)
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorUnbondingDelegationRow{
		dbtypes.NewValidatorUnbondingDelegationRow(
			validator1.GetConsAddr().String(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			completionTimestamp1,
			1000,
			timestamp1,
		),
		dbtypes.NewValidatorUnbondingDelegationRow(
			validator1.GetConsAddr().String(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			completionTimestamp1,
			1000,
			timestamp1,
		),
		dbtypes.NewValidatorUnbondingDelegationRow(
			validator1.GetConsAddr().String(),
			delegator1.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			completionTimestamp2,
			1001,
			timestamp1,
		),
		dbtypes.NewValidatorUnbondingDelegationRow(
			validator2.GetConsAddr().String(),
			delegator2.String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			completionTimestamp2,
			1500,
			timestamp2,
		),
	}

	suite.Require().Len(rows, len(expected))
	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveRedelegation() {
	// Setup
	delegator := suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	srcValidator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	dstValidator := suite.getValidator(
		"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
		"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
		"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
	)

	var height int64 = 1000
	amount := sdk.NewCoin("udesmos", sdk.NewInt(1000))

	completionTimestamp, err := time.Parse(time.RFC3339, "2020-08-10T16:00:00Z")
	suite.Require().NoError(err)

	// Save data
	reDelegation := types.NewRedelegation(
		delegator,
		srcValidator.GetOperator(),
		dstValidator.GetOperator(),
		amount,
		completionTimestamp,
		height,
	)
	err = suite.database.SaveRedelegation(reDelegation)
	suite.Require().NoError(err)

	// Get inserted data
	var rows []dbtypes.ValidatorReDelegationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_redelegation`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewValidatorReDelegationRow(
		delegator.String(),
		srcValidator.GetConsAddr().String(),
		dstValidator.GetConsAddr().String(),
		dbtypes.NewDbCoin(amount),
		height,
		completionTimestamp,
	)))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveReDelegations() {
	// Setup
	delegator1 := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	delegator2 := suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	srcValidator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	srcValidator2 := suite.getValidator(
		"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
		"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
		"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
	)
	dstValidator1 := suite.getValidator(
		"cosmosvalcons1px0zkz2cxvc6lh34uhafveea9jnaagckmrlsye",
		"cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn",
		"cosmosvalconspub1zcjduepq0dc9apn3pz2x2qyujcnl2heqq4aceput2uaucuvhrjts75q0rv5smjjn7v",
	)
	dstValidator2 := suite.getValidator(
		"cosmosvalcons1rtst6se0nfgjy362v33jt5d05crgdyhfvvvvay",
		"cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr",
		"cosmosvalconspub1zcjduepq5e8w7t7k9pwfewgrwy8vn6cghk0x49chx64vt0054yl4wwsmjgrqfackxm",
	)

	completionTimestamp1, err := time.Parse(time.RFC3339, "2020-08-10T16:00:00Z")
	suite.Require().NoError(err)

	completionTimestamp2, err := time.Parse(time.RFC3339, "2020-08-20T16:00:00Z")
	suite.Require().NoError(err)

	// Save data
	reDelegations := []types.Redelegation{
		types.NewRedelegation(
			delegator1,
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			completionTimestamp1,
			1000,
		),
		types.NewRedelegation(
			delegator1,
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(100)),
			completionTimestamp1,
			1000,
		),
		types.NewRedelegation(
			delegator1,
			srcValidator1.GetOperator(),
			dstValidator1.GetOperator(),
			sdk.NewCoin("desmos", sdk.NewInt(100)),
			completionTimestamp2,
			1001,
		),
		types.NewRedelegation(
			delegator2,
			srcValidator2.GetOperator(),
			dstValidator2.GetOperator(),
			sdk.NewCoin("cosmos", sdk.NewInt(200)),
			completionTimestamp2,
			1500,
		),
	}
	err = suite.database.SaveRedelegations(reDelegations)
	suite.Require().NoError(err)

	// Read the data
	var rows []dbtypes.ValidatorReDelegationRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_redelegation`)
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorReDelegationRow{
		dbtypes.NewValidatorReDelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr().String(),
			dstValidator1.GetConsAddr().String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			1000,
			completionTimestamp1,
		),
		dbtypes.NewValidatorReDelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr().String(),
			dstValidator1.GetConsAddr().String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(100))),
			1000,
			completionTimestamp1,
		),
		dbtypes.NewValidatorReDelegationRow(
			delegator1.String(),
			srcValidator1.GetConsAddr().String(),
			dstValidator1.GetConsAddr().String(),
			dbtypes.NewDbCoin(sdk.NewCoin("desmos", sdk.NewInt(100))),
			1001,
			completionTimestamp2,
		),
		dbtypes.NewValidatorReDelegationRow(
			delegator2.String(),
			srcValidator2.GetConsAddr().String(),
			dstValidator2.GetConsAddr().String(),
			dbtypes.NewDbCoin(sdk.NewCoin("cosmos", sdk.NewInt(200))),
			1500,
			completionTimestamp2,
		),
	}

	suite.Require().Len(rows, len(expected))
	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]))
	}
}
