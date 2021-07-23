package consensus

import (
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/utils"
)

// Register registers the utils that should be run periodically
func Register(scheduler *gocron.Scheduler, db *database.Db) error {
	log.Debug().Str("module", "consensus").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(1).Minute().Do(func() {
		utils.WatchMethod(func() error { return updateBlockTimeInMinute(db) })
	}); err != nil {
		return err
	}

	if _, err := scheduler.Every(1).Hour().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateBlockTimeInHour(db) })
	}); err != nil {
		return err
	}

	if _, err := scheduler.Every(1).Day().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateBlockTimeInDay(db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateBlockTimeInMinute insert average block time in the latest minute
func updateBlockTimeInMinute(db *database.Db) error {
	log.Trace().Str("module", "consensus").Str("operation", "block time").
		Msg("updating block time in minutes")

	block, err := db.GetLastBlock()
	if err != nil {
		return err
	}

	genesis, err := db.GetGenesis()
	if err != nil {
		return err
	}

	// Check if the chain has been created at least a minute ago
	if block.Timestamp.Sub(genesis.Time).Minutes() < 0 {
		return nil
	}

	minute, err := db.GetBlockHeightTimeMinuteAgo(block.Timestamp)
	if err != nil {
		return err
	}
	newBlockTime := block.Timestamp.Sub(minute.Timestamp).Seconds() / float64(block.Height-minute.Height)

	return db.SaveAverageBlockTimePerMin(newBlockTime, block.Height)
}

// updateBlockTimeInHour insert average block time in the latest hour
func updateBlockTimeInHour(db *database.Db) error {
	log.Trace().Str("module", "consensus").Str("operation", "block time").
		Msg("updating block time in hours")

	block, err := db.GetLastBlock()
	if err != nil {
		return err
	}

	genesis, err := db.GetGenesis()
	if err != nil {
		return err
	}

	// Check if the chain has been created at least an hour ago
	if block.Timestamp.Sub(genesis.Time).Hours() < 0 {
		return nil
	}

	hour, err := db.GetBlockHeightTimeHourAgo(block.Timestamp)
	if err != nil {
		return err
	}
	newBlockTime := block.Timestamp.Sub(hour.Timestamp).Seconds() / float64(block.Height-hour.Height)

	return db.SaveAverageBlockTimePerHour(newBlockTime, block.Height)
}

// updateBlockTimeInDay insert average block time in the latest minute
func updateBlockTimeInDay(db *database.Db) error {
	log.Trace().Str("module", "consensus").Str("operation", "block time").
		Msg("updating block time in days")

	block, err := db.GetLastBlock()
	if err != nil {
		return err
	}

	genesis, err := db.GetGenesis()
	if err != nil {
		return err
	}

	// Check if the chain has been created at least a days ago
	if block.Timestamp.Sub(genesis.Time).Hours() < 24 {
		return nil
	}

	day, err := db.GetBlockHeightTimeDayAgo(block.Timestamp)
	if err != nil {
		return err
	}
	newBlockTime := block.Timestamp.Sub(day.Timestamp).Seconds() / float64(block.Height-day.Height)

	return db.SaveAverageBlockTimePerDay(newBlockTime, block.Height)
}
