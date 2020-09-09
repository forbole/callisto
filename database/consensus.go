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
/* 
func (db BigDipperDb) getBlockHeightTime(pastTime time.Time) (dbtypes.BlockRow, error) {
	stmt := `SELECT block.timestamp, block.height
	FROM block
	WHERE block.timestamp < $1
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
	pastTime := now.Add(time.Hour * 24 * -1)
	return db.getBlockHeightTime(pastTime)
}
 */

func (db BigDipperDb) SaveAverageBlockTimePerMin (average_time float64,timestamp time.Time,height)error{
	stmt:=`INSERT INTO average_block_time_per_minute(average_time,timestamp,height) values ($1,$2)`
	_,err := db.Sqlx.Exec(stmt,average_time,timestamp,height)
}


// SaveGenesisHeight save the genesis height
func (db BigDipperDb) SaveGenesisTime(genesisTime time.Time)error{
	stmt := `INSERT INTO genesis(time) values ($1)`
	_,err:= db.Sqlx.Exec(stmt, genesisTime)
	if err != nil {
		return err
	}
	return  nil
}

func (db BigDipperDb) GetGenesisTime() (time.Time, error) {
	stmt := `SELECT time from genesis;`
	var val []time.Time
	err := db.Sqlx.Select(&val, stmt)
	if err != nil {
		return time.Time{}, err
	}
	return val[0], nil
}

