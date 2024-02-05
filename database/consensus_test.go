package database_test

import (
	"time"

	dbtypes "github.com/forbole/callisto/v4/database/types"
	"github.com/forbole/callisto/v4/types"
)

func (suite *DbTestSuite) TestSaveConsensus_GetBlockHeightTimeMinuteAgo() {
	timeAgo, err := time.Parse(time.RFC3339, "2020-01-01T15:00:00Z")
	suite.Require().NoError(err)

	var height int64 = 1000

	_, err = suite.database.SQL.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
	VALUES ('desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.SQL.Exec(`INSERT INTO block(height, hash, num_txs, total_gas, proposer_address, timestamp)
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
	_, err = suite.database.SQL.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
	VALUES ('desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.SQL.Exec(`INSERT INTO block(height, hash, num_txs, total_gas, proposer_address, timestamp)
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

	_, err = suite.database.SQL.Exec(`INSERT INTO validator (consensus_address, consensus_pubkey) 
	VALUES ('desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', 'cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8')`)
	suite.Require().NoError(err)

	_, err = suite.database.SQL.Exec(`INSERT INTO block(height, hash, num_txs, total_gas, proposer_address, timestamp)
	VALUES ($1, '5EF85F2251F656BA0FBFED9AEFCBC44A9CCBCFD8B96897E74426E07229D2ADE0', '0', '0', 'desmosvalcons1mxrd5cyjgpx5vfgltrdufq9wq4ynwc799ndrg8', $2)`, height, timeAgo)
	suite.Require().NoError(err)

	timeNow := timeAgo.Add(time.Hour * 24)
	result, err := suite.database.GetBlockHeightTimeDayAgo(timeNow)
	suite.Require().NoError(err)

	suite.Require().True(result.Timestamp.Equal(timeAgo))
	suite.Require().Equal(height, result.Height)
}

func (suite *DbTestSuite) TestSaveConsensus_SaveAverageBlockTimePerMin() {
	// Save the data
	err := suite.database.SaveAverageBlockTimePerMin(5.05, 10)
	suite.Require().NoError(err)

	original := dbtypes.NewAverageTimeRow(5.05, 10)

	// Verify the data
	var rows []dbtypes.AverageTimeRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_minute")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(original))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a lower height
	err = suite.database.SaveAverageBlockTimePerMin(6, 9)
	suite.Require().NoError(err)

	// Verify the data
	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_minute")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(original), "updating with a lower height should not change the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with the same height
	err = suite.database.SaveAverageBlockTimePerMin(10, 10)
	suite.Require().NoError(err)

	// Verify the data
	expected := dbtypes.NewAverageTimeRow(10, 10)

	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_minute")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with same height should change the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a higher height
	err = suite.database.SaveAverageBlockTimePerMin(20, 15)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewAverageTimeRow(20, 15)

	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_minute")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with higher height should change the data")
}

func (suite *DbTestSuite) TestSaveConsensus_SaveAverageBlockTimePerHour() {
	// Save the data
	err := suite.database.SaveAverageBlockTimePerHour(5.05, 10)
	suite.Require().NoError(err)

	original := dbtypes.NewAverageTimeRow(5.05, 10)

	// Verify the data
	var rows []dbtypes.AverageTimeRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_hour")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(original))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a lower height
	err = suite.database.SaveAverageBlockTimePerHour(6, 9)
	suite.Require().NoError(err)

	// Verify the data
	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_hour")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(original), "updating with a lower height should not change the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with the same height
	err = suite.database.SaveAverageBlockTimePerHour(10, 10)
	suite.Require().NoError(err)

	// Verify the data
	expected := dbtypes.NewAverageTimeRow(10, 10)

	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_hour")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with same height should change the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a higher height
	err = suite.database.SaveAverageBlockTimePerHour(20, 15)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewAverageTimeRow(20, 15)

	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_hour")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with higher height should change the data")
}

func (suite *DbTestSuite) TestSaveConsensus_SaveAverageBlockTimePerDay() {
	// Save the data
	err := suite.database.SaveAverageBlockTimePerDay(5.05, 10)
	suite.Require().NoError(err)

	original := dbtypes.NewAverageTimeRow(5.05, 10)

	// Verify the data
	var rows []dbtypes.AverageTimeRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_day")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(original))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a lower height
	err = suite.database.SaveAverageBlockTimePerDay(6, 9)
	suite.Require().NoError(err)

	// Verify the data
	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_day")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(original), "updating with a lower height should not change the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with the same height
	err = suite.database.SaveAverageBlockTimePerDay(10, 10)
	suite.Require().NoError(err)

	// Verify the data
	expected := dbtypes.NewAverageTimeRow(10, 10)

	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_day")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with same height should change the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a higher height
	err = suite.database.SaveAverageBlockTimePerDay(20, 15)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewAverageTimeRow(20, 15)

	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_per_day")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with higher height should change the data")
}

func (suite *DbTestSuite) TestSaveConsensus_SaveAverageBlockTimeGenesis() {
	// Save the data
	err := suite.database.SaveAverageBlockTimeGenesis(5.05, 10)
	suite.Require().NoError(err)

	original := dbtypes.NewAverageTimeRow(5.05, 10)

	// Verify the data
	var rows []dbtypes.AverageTimeRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_from_genesis")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(original))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a lower height
	err = suite.database.SaveAverageBlockTimeGenesis(6, 9)
	suite.Require().NoError(err)

	// Verify the data
	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_from_genesis")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(original), "updating with a lower height should not change the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with the same height
	err = suite.database.SaveAverageBlockTimeGenesis(10, 10)
	suite.Require().NoError(err)

	// Verify the data
	expected := dbtypes.NewAverageTimeRow(10, 10)

	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_from_genesis")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with same height should change the data")

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a higher height
	err = suite.database.SaveAverageBlockTimeGenesis(20, 15)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewAverageTimeRow(20, 15)

	rows = []dbtypes.AverageTimeRow{}
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM average_block_time_from_genesis")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(expected), "updating with higher height should change the data")
}

func (suite *DbTestSuite) TestSaveConsensus_SaveGenesisData() {
	err := suite.database.SaveGenesis(types.NewGenesis(
		"testnet-1",
		time.Date(2020, 1, 02, 15, 00, 00, 000, time.UTC),
		0,
	))
	suite.Require().NoError(err)

	// Should have only one row
	err = suite.database.SaveGenesis(types.NewGenesis(
		"testnet-2",
		time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC),
		0,
	))

	var rows []*dbtypes.GenesisRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM genesis")
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewGenesisRow(
		"testnet-2",
		time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC),
		0,
	)))
}

func (suite *DbTestSuite) TestSaveConsensus_GetGenesis() {
	_, err := suite.database.Sqlx.Exec(
		`INSERT INTO genesis(chain_id, time, initial_height) VALUES ($1, $2, $3)`,
		"testnet-1",
		time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC),
		0,
	)

	genesis, err := suite.database.GetGenesis()
	suite.Require().NoError(err)
	suite.Require().True(genesis.Equal(types.NewGenesis(
		"testnet-1",
		time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC),
		0,
	)))
}
