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
	suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	validator := dbtypes.NewValidatorData(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
	)

	// First inserting
	err := suite.database.SaveSingleValidatorData(validator)

	// Test double inserting
	err = suite.database.SaveSingleValidatorData(validator)
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
	)))

}

func (suite *DbTestSuite) TestBigDipperDb_GetValidatorData() {
	suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	// Insert test data
	_, err := suite.database.Sql.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address) 
VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl','cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl','cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a')`)
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

	suite.Require().Equal("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a", validator.GetSelfDelegateAddress().String())

}

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorsData() {
	suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	validators := []types.Validator{
		dbtypes.NewValidatorData(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		),
		dbtypes.NewValidatorData(
			"cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk",
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
		),
	}

	expectedValidatorInfo := []dbtypes.ValidatorInfoRow{
		dbtypes.NewValidatorInfoRow("cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		),
		dbtypes.NewValidatorInfoRow("cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y",
			"cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn",
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
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
		suite.Require().True(wInfo == expectedValidatorInfo[index])
	}
}

func (suite *DbTestSuite) TestBigDipperDb_GetValidatorsData() {
	suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	// Inser the test data
	queries := []string{
		`INSERT INTO validator (consensus_address, consensus_pubkey) VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`,
		`INSERT INTO validator (consensus_address, consensus_pubkey) VALUES ('cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y', 'cosmosvalconspub1zcjduepqe93asg05nlnj30ej2pe3r8rkeryyuflhtfw3clqjphxn4j3u27msrr63nk')`,
		`INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address) VALUES ('cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl', 'cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl','cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs')`,
		`INSERT INTO validator_info (consensus_address, operator_address,self_delegate_address) VALUES ('cosmosvalcons1qq92t2l4jz5pt67tmts8ptl4p0jhr6utx5xa8y', 'cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn','cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a')`,
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
			"cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a",
		),
		dbtypes.NewValidatorData(
			"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
			"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
			"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
			"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		),
	}

	suite.Require().Len(data, len(expected))
	for index, validator := range data {
		suite.Require().Equal(expected[index], validator)
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

	validator := types.NewValidator(constAddrObj, valAddrObj, pubKey, selfDelegation)
	err = suite.database.SaveSingleValidatorData(validator)
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

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorCommission() {
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	var height int64 = 1000

	timestamp, err := time.Parse(time.RFC3339, "2020-01-01T10:00:00Z")
	suite.Require().NoError(err)

	commission := types.NewValidatorCommission(
		validator.GetOperator(),
		newDecPts(11, 3),
		newIntPtr(12),
		height,
		timestamp,
	)

	err = suite.database.SaveEditCommission(commission)
	suite.Require().NoError(err)

	var rows []dbtypes.ValidatorCommission
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewValidatorCommission(
		validator.GetOperator().String(),
		"0.011000000000000000",
		"12",
		height,
		timestamp,
	)))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorCommissions() {
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

	timestamp, err := time.Parse(time.RFC3339, "2020-01-01T10:00:00Z")
	suite.Require().NoError(err)

	commissions := []types.ValidatorCommission{
		types.NewValidatorCommission(validator1.GetOperator(), newDecPts(1, 2), newIntPtr(30), 0, timestamp),
		types.NewValidatorCommission(validator2.GetOperator(), newDecPts(2, 2), newIntPtr(40), 0, timestamp),
	}

	err = suite.database.SaveValidatorCommissions(commissions)
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorCommission{
		dbtypes.NewValidatorCommission(
			validator1.GetOperator().String(), "0.010000000000000000", "30", 0, timestamp,
		),
		dbtypes.NewValidatorCommission(
			validator2.GetOperator().String(), "0.020000000000000000", "40", 0, timestamp,
		),
	}

	var rows []dbtypes.ValidatorCommission
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM validator_commission`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 2)
	for index, row := range rows {
		suite.Require().True(row.Equal(expected[index]))
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

func (suite *DbTestSuite) TestBigDipperDb_SaveSelfDelegation() {
	delegator1 := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	delegator2 := suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	validator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	timestamp1, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	// Save data
	delegations := []types.DelegationShare{
		types.NewDelegationShare(
			//self delegation
			validator1.GetOperator(),
			delegator1,
			"1000.00001",
			1000,
			timestamp1,
		),
		types.NewDelegationShare(
			validator1.GetOperator(),
			delegator2,
			"1000.0002",
			1000,
			timestamp1,
		),
	}

	err = suite.database.SaveDelegationsShares(delegations)
	suite.Require().NoError(err)

	//expected
	delegationsExpected := []dbtypes.ValidatorDelegationSharesRow{
		dbtypes.NewValidatorDelegationSharesRow(
			//self delegation
			validator1.GetOperator().String(),
			delegator1.String(),
			1000.00001,
			timestamp1,
			1000,
		),
		dbtypes.NewValidatorDelegationSharesRow(
			validator1.GetOperator().String(),
			delegator2.String(),
			1000.0002,
			timestamp1,
			1000,
		),
	}

	//read data
	var delegationrows []dbtypes.ValidatorDelegationSharesRow
	err = suite.database.Sqlx.Select(&delegationrows, `SELECT * FROM validator_delegation_shares`)
	suite.Require().NoError(err)

	for index, row := range delegationrows {
		suite.Require().True(row.Equal(delegationsExpected[index]))
	}

}
func (suite *DbTestSuite) TestBigDipperDb_SaveVotingPower() {
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
		),
		types.NewValidatorVotingPower(
			validator2.GetConsAddr(),
			2000,
			100,
		),
	}
	err := suite.database.SaveVotingPowers(votingPowers)
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorVotingPowerRow{
		dbtypes.NewValidatorVotingPowerRow(
			validator1.GetConsAddr().String(),
			1000,
			100,
		),
		dbtypes.NewValidatorVotingPowerRow(
			validator2.GetConsAddr().String(),
			2000,
			100,
		),
	}

	var result []dbtypes.ValidatorVotingPowerRow
	err = suite.database.Sqlx.Select(&result, "SELECT * FROM validator_voting_power")
	suite.Require().NoError(err)

	for index, row := range result {
		suite.Require().True(row.Equal(expected[index]))
	}

}

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorDescription() {
	validator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	timestamp, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	var height int64 = 1
	description := types.NewValidatorDescription(validator1.GetOperator(), stakingtypes.NewDescription(
		"moniker",
		"identity",
		"", //test null value
		"securityContact",
		"details",
	), timestamp, height)
	err = suite.database.SaveValidatorDescription(description)
	suite.Require().NoError(err)

	expected := dbtypes.NewValidatorDescriptionRow(validator1.GetOperator().String(),
		"moniker",
		"identity",
		"", //test null value
		"securityContact",
		"details", height, timestamp)

	var rows []dbtypes.ValidatorDescriptionRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM validator_description")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equals(expected))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveValidatorsDescription() {
	validator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	validator2 := suite.getValidator(
		"cosmosvalcons1px0zkz2cxvc6lh34uhafveea9jnaagckmrlsye",
		"cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn",
		"cosmosvalconspub1zcjduepq0dc9apn3pz2x2qyujcnl2heqq4aceput2uaucuvhrjts75q0rv5smjjn7v",
	)
	timestamp, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	var height int64 = 1
	descriptions := []types.ValidatorDescription{
		types.NewValidatorDescription(validator1.GetOperator(), stakingtypes.NewDescription(
			"moniker1",
			"identity1",
			"", //test null value
			"securityContact1",
			"details1",
		), timestamp, height),
		types.NewValidatorDescription(validator2.GetOperator(), stakingtypes.NewDescription(
			"moniker2",
			"identity2",
			"", //test null value
			"securityContact2",
			"details2",
		), timestamp, height),
	}

	err = suite.database.SaveValidatorsDescription(descriptions)
	suite.Require().NoError(err)

	expected := []dbtypes.ValidatorDescriptionRow{
		dbtypes.NewValidatorDescriptionRow(validator1.GetOperator().String(),
			"moniker1",
			"identity1",
			"", //test null value
			"securityContact1",
			"details1", height, timestamp),
		dbtypes.NewValidatorDescriptionRow(validator2.GetOperator().String(),
			"moniker2",
			"identity2",
			"", //test null value
			"securityContact2",
			"details2", height, timestamp),
	}

	var rows []dbtypes.ValidatorDescriptionRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM validator_description order by moniker ASC")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 2)
	for index, row := range rows {
		suite.Require().True(row.Equals(expected[index]))
	}
}
