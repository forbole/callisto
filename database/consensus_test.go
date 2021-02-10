package database_test

import (
	time "time"

	tmtypes "github.com/tendermint/tendermint/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
	constypes "github.com/forbole/bdjuno/x/consensus/types"
)

func (suite *DbTestSuite) TestSaveConsensus() {
	events := []*constypes.ConsensusEvent{
		constypes.NewConsensusEvent(1, 0, "First"),
		constypes.NewConsensusEvent(2, 0, "Second - Round 0 "),
		constypes.NewConsensusEvent(2, 1, "Second - Round 1"),
	}

	for _, event := range events {
		err := suite.database.SaveConsensus(event)
		suite.Require().NoError(err)

		var rows []dbtypes.ConsensusRow
		err = suite.database.Sqlx.Select(&rows, "SELECT * FROM consensus")
		suite.Require().NoError(err)

		// Make sure the consensus table only holds 1 value at a time with the correct data inside
		suite.Require().Len(rows, 1)
		suite.Require().True(rows[0].Equal(dbtypes.ConsensusRow{
			Height: event.Height,
			Round:  event.Round,
			Step:   event.Step,
		}))
	}
}

func (suite *DbTestSuite) TestSaveConsensus_GetBlockHeightTimeMinuteAgo() {
	timeAgo, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	var height int64 = 1000

	_, err = suite.database.Sql.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
	VALUES ('desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO block(height, hash, num_txs, total_gas, proposer_address, timestamp)
	VALUES ($1, '5EF85F2251F656BA0FBFED9AEFCBC44A9CCBCFD8B96897E74426E07229D2ADE0', '0', '0', 'desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', $2)`, height, timeAgo)
	suite.Require().NoError(err)

	timeNow := timeAgo.Add(time.Minute)
	result, err := suite.database.GetBlockHeightTimeMinuteAgo(timeNow)
	suite.Require().NoError(err)

	suite.Require().True(result.Timestamp.Equal(timeAgo))
	suite.Require().Equal(height, result.Height)
}

func (suite *DbTestSuite) TestSaveConsensus_GetBlockHeightTimeHourAgo() {
	timeAgo, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	var height int64 = 1000
	_, err = suite.database.Sql.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
	VALUES ('desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO block(height, hash, num_txs, total_gas, proposer_address, timestamp)
	VALUES ($1, '5EF85F2251F656BA0FBFED9AEFCBC44A9CCBCFD8B96897E74426E07229D2ADE0', '0', '0', 'desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', $2)`, height, timeAgo)
	suite.Require().NoError(err)

	timeNow := timeAgo.Add(time.Hour)
	result, err := suite.database.GetBlockHeightTimeHourAgo(timeNow)
	suite.Require().NoError(err)

	suite.Require().True(result.Timestamp.Equal(timeAgo))
	suite.Require().Equal(height, result.Height)
}

func (suite *DbTestSuite) TestSaveConsensus_GetBlockHeightTimeDayAgo() {
	timeAgo, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	var height int64 = 1000

	_, err = suite.database.Sql.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
	VALUES ('desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO block(height, hash, num_txs, total_gas, proposer_address, timestamp)
	VALUES ($1, '5EF85F2251F656BA0FBFED9AEFCBC44A9CCBCFD8B96897E74426E07229D2ADE0', '0', '0', 'desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', $2)`, height, timeAgo)
	suite.Require().NoError(err)

	timeNow := timeAgo.Add(time.Hour * 24)
	result, err := suite.database.GetBlockHeightTimeDayAgo(timeNow)
	suite.Require().NoError(err)

	suite.Require().True(result.Timestamp.Equal(timeAgo))
	suite.Require().Equal(height, result.Height)
}

func (suite *DbTestSuite) TestSaveConsensus_SaveAverageBlockTimePerMin() {
	var height int64 = 1000
	var averageTime = 5.05

	err := suite.database.SaveAverageBlockTimePerMin(averageTime, height)
	suite.Require().NoError(err)

	expected := dbtypes.NewBlockTimeRow(averageTime, height)

	var rows []dbtypes.BlockTimeRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_minute")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected))
}

func (suite *DbTestSuite) TestSaveConsensus_SaveAverageBlockTimePerHour() {
	var height int64 = 1000
	var averageTime = 5.05

	err := suite.database.SaveAverageBlockTimePerHour(averageTime, height)
	suite.Require().NoError(err)

	expected := dbtypes.NewBlockTimeRow(averageTime, height)

	var rows []dbtypes.BlockTimeRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_hour")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected))
}

func (suite *DbTestSuite) TestSaveConsensus_SaveAverageBlockTimePerDay() {
	var height int64 = 1000
	var averageTime = 5.05

	err := suite.database.SaveAverageBlockTimePerDay(averageTime, height)
	suite.Require().NoError(err)

	expected := dbtypes.NewBlockTimeRow(averageTime, height)

	var rows []dbtypes.BlockTimeRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_day")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected))
}

func (suite *DbTestSuite) TestSaveConsensus_SaveGenesisTime() {
	err := suite.database.SaveGenesisData(&tmtypes.GenesisDoc{
		ChainID:     "testnet-1",
		GenesisTime: time.Date(2020, 1, 02, 15, 00, 00, 000, time.UTC),
	})
	suite.Require().NoError(err)

	// Should have only one row
	err = suite.database.SaveGenesisData(&tmtypes.GenesisDoc{
		ChainID:     "testnet-2",
		GenesisTime: time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC),
	})

	var rows []*dbtypes.GenesisRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM genesis")
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewGenesisRow(
		"testnet-2",
		time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC),
	)))
}

func (suite *DbTestSuite) TestSaveConsensus_GetGenesisTime() {
	_, err := suite.database.Sqlx.Exec(
		`INSERT INTO genesis(chain_id, time) VALUES ($1, $2)`,
		"testnet-1",
		time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC),
	)

	genTime, err := suite.database.GetGenesisTime()
	suite.Require().NoError(err)
	suite.Require().True(genTime.Equal(
		time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC),
	))
}

func (suite *DbTestSuite) TestSaveConsensus_SaveAverageBlockTimeGenesis() {
	var height int64 = 1000
	var averageTime = 5.05

	err := suite.database.SaveAverageBlockTimeGenesis(averageTime, height)
	suite.Require().NoError(err)

	expected := dbtypes.NewBlockTimeRow(averageTime, height)

	var rows []dbtypes.BlockTimeRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_from_genesis")
	suite.Require().Len(rows, 1)
	suite.Require().True(expected.Equal(rows[0]))
}
