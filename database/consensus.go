package database

import (
	"time"

	dbtypes "github.com/forbole/bdjuno/database/types"
	constypes "github.com/forbole/bdjuno/x/consensus/types"
	"github.com/rs/zerolog/log"
)

// SaveConsensus allows to properly store the given consensus event into the database.
// Note that only one consensus event is allowed inside the database at any time.
func (db BigDipperDb) SaveConsensus(event constypes.ConsensusEvent) error {
	log.Debug().
		Str("module", "consensus").
		Int64("height", event.Height).
		Int("round", event.Round).
		Str("step", event.Step).
		Msg("saving consensus")

	// Delete all the existing events
	stmt := `DELETE FROM consensus WHERE true`
	if _, err := db.Sql.Exec(stmt); err != nil {
		return err
	}

	stmt = `INSERT INTO consensus (height, round, step) VALUES ($1, $2, $3)`
	_, err := db.Sql.Exec(stmt, event.Height, event.Round, event.Step)
	return err
}

func (db BigDipperDb) getBlockHeightTime(pastTime time.Time) (dbtypes.BlockRow, error) {
	stmt := `SELECT block.timestamp, block.height
	FROM block
	WHERE block.timestamp <= $1
	ORDER BY block.timestamp DESC
	LIMIT 1;`
	var val []dbtypes.BlockRow
	err := db.Sqlx.Select(&val, stmt, pastTime)
	if err != nil {
		return dbtypes.BlockRow{}, err
	}
	return val[0], nil
}

// GetBlockHeightTimeMinuteAgo return block height and time that a block proposals
// about a minute ago from imput date
func (db BigDipperDb) GetBlockHeightTimeMinuteAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Minute * -1)
	return db.getBlockHeightTime(pastTime)
}

// GetBlockHeightTimeHourAgo return block height and time that a block proposals
// about a hour ago from imput date
func (db BigDipperDb) GetBlockHeightTimeHourAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Hour * -1)
	return db.getBlockHeightTime(pastTime)
}

// GetBlockHeightTimeDayAgo return block height and time that a block proposals
// about a day (24hour) ago from imput date
func (db BigDipperDb) GetBlockHeightTimeDayAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Hour * -24)
	return db.getBlockHeightTime(pastTime)
}

// SaveAverageBlockTimePerMin save the average block time in average_block_time_per_minute table
func (db BigDipperDb) SaveAverageBlockTimePerMin(averageTime float64, timestamp time.Time, height int64) error {
	stmt := `INSERT INTO average_block_time_per_minute(average_time,timestamp,height) values ($1,$2,$3)`
	_, err := db.Sqlx.Exec(stmt, averageTime, timestamp, height)
	return err
}

// SaveAverageBlockTimePerHour save the average block time in average_block_time_per_hour table
func (db BigDipperDb) SaveAverageBlockTimePerHour(averageTime float64, timestamp time.Time, height int64) error {
	stmt := `INSERT INTO average_block_time_per_hour(average_time,timestamp,height) values ($1,$2,$3)`
	_, err := db.Sqlx.Exec(stmt, averageTime, timestamp, height)
	return err
}

// SaveAverageBlockTimePerDay save the average block time in average_block_time_per_day table
func (db BigDipperDb) SaveAverageBlockTimePerDay(averageTime float64, timestamp time.Time, height int64) error {
	stmt := `INSERT INTO average_block_time_per_day(average_time,timestamp,height) values ($1,$2,$3)`
	_, err := db.Sqlx.Exec(stmt, averageTime, timestamp, height)
	return err
}

// SaveGenesisHeight save the genesis height
func (db BigDipperDb) SaveGenesisTime(genesisTime time.Time) error {
	stmt := `DELETE FROM genesis`
	_, err := db.Sqlx.Exec(stmt)
	if err != nil {
		return err
	}
	stmt = `INSERT INTO genesis(time) values ($1)`
	_, err = db.Sqlx.Exec(stmt, genesisTime)
	return err
}

// GetGenesisTime get genesis time of chain (only work if x/consensus enabled)
func (db BigDipperDb) GetGenesisTime() (time.Time, error) {
	stmt := `SELECT * from genesis;`
	var val []time.Time
	err := db.Sqlx.Select(&val, stmt)
	if err != nil || len(val) == 0 {
		return time.Time{}, err
	}
	return val[0], nil
}

// SaveAverageBlockTimeGenesis save the average block time in average_block_time_from_genesis table
func (db BigDipperDb) SaveAverageBlockTimeGenesis(averageTime float64, timestamp time.Time, height int64) error {
	stmt := `INSERT INTO average_block_time_from_genesis(average_time,timestamp,height) values ($1,$2,$3)`
	_, err := db.Sqlx.Exec(stmt, averageTime, timestamp, height)
	return err
}
