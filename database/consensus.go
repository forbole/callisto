package database

import (
	"fmt"
	"time"

	"github.com/forbole/bdjuno/v3/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

// GetLastBlock returns the last block stored inside the database based on the heights
func (db *Db) GetLastBlock() (*dbtypes.BlockRow, error) {
	stmt := `SELECT * FROM block ORDER BY height DESC LIMIT 1`

	var blocks []dbtypes.BlockRow
	if err := db.Sqlx.Select(&blocks, stmt); err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("cannot get block, no blocks saved")
	}

	return &blocks[0], nil
}

// GetLastBlockHeight returns the last block height stored inside the database
func (db *Db) GetLastBlockHeight() (int64, error) {
	block, err := db.GetLastBlock()
	if err != nil {
		return 0, err
	}
	if block == nil {
		return 0, fmt.Errorf("block table is empty")
	}
	return block.Height, nil
}

// -------------------------------------------------------------------------------------------------------------------

// getBlockHeightTime retrieves the block at the specific time
func (db *Db) getBlockHeightTime(pastTime time.Time) (dbtypes.BlockRow, error) {
	stmt := `SELECT * FROM block WHERE block.timestamp <= $1 ORDER BY block.timestamp DESC LIMIT 1;`

	var val []dbtypes.BlockRow
	if err := db.Sqlx.Select(&val, stmt, pastTime); err != nil {
		return dbtypes.BlockRow{}, err
	}

	if len(val) == 0 {
		return dbtypes.BlockRow{}, fmt.Errorf("cannot get block time, no blocks saved")
	}

	return val[0], nil
}

// GetBlockHeightTimeMinuteAgo return block height and time that a block proposals
// about a minute ago from input date
func (db *Db) GetBlockHeightTimeMinuteAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Minute * -1)
	return db.getBlockHeightTime(pastTime)
}

// GetBlockHeightTimeHourAgo return block height and time that a block proposals
// about a hour ago from input date
func (db *Db) GetBlockHeightTimeHourAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Hour * -1)
	return db.getBlockHeightTime(pastTime)
}

// GetBlockHeightTimeDayAgo return block height and time that a block proposals
// about a day (24hour) ago from input date
func (db *Db) GetBlockHeightTimeDayAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Hour * -24)
	return db.getBlockHeightTime(pastTime)
}

// -------------------------------------------------------------------------------------------------------------------

// SaveAverageBlockTimePerMin save the average block time in average_block_time_per_minute table
func (db *Db) SaveAverageBlockTimePerMin(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_per_minute(average_time, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time, 
        height = excluded.height
WHERE average_block_time_per_minute.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time per minute: %s", err)
	}

	return nil
}

// SaveAverageBlockTimePerHour save the average block time in average_block_time_per_hour table
func (db *Db) SaveAverageBlockTimePerHour(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_per_hour(average_time, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time,
        height = excluded.height
WHERE average_block_time_per_hour.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time per hour: %s", err)
	}

	return nil
}

// SaveAverageBlockTimePerDay save the average block time in average_block_time_per_day table
func (db *Db) SaveAverageBlockTimePerDay(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_per_day(average_time, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time,
        height = excluded.height
WHERE average_block_time_per_day.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time per day: %s", err)
	}

	return nil
}

// SaveAverageBlockTimeGenesis save the average block time in average_block_time_from_genesis table
func (db *Db) SaveAverageBlockTimeGenesis(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_from_genesis(average_time ,height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time, 
        height = excluded.height
WHERE average_block_time_from_genesis.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time since genesis: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveGenesis save the given genesis data
func (db *Db) SaveGenesis(genesis *types.Genesis) error {
	stmt := `
INSERT INTO genesis(time, chain_id, initial_height) 
VALUES ($1, $2, $3) ON CONFLICT (one_row_id) DO UPDATE 
    SET time = excluded.time,
        initial_height = excluded.initial_height,
        chain_id = excluded.chain_id`

	_, err := db.Sqlx.Exec(stmt, genesis.Time, genesis.ChainID, genesis.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis: %s", err)
	}

	return nil
}

// GetGenesis returns the genesis information stored inside the database
func (db *Db) GetGenesis() (*types.Genesis, error) {
	var rows []*dbtypes.GenesisRow
	err := db.Sqlx.Select(&rows, `SELECT * FROM genesis;`)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows inside the genesis table")
	}

	row := rows[0]
	return types.NewGenesis(row.ChainID, row.Time, row.InitialHeight), nil
}
