package consensus

import (
	"github.com/desmos-labs/juno/client"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/utils"
)

// PeriodicConcensusOperations returns the AdditionalOperation that periodically runs fetches from
// the LCD to make sure that constantly changing data are synced properly.
func Register(scheduler *gocron.Scheduler, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().Str("module", "consensus").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(1).Minute().Do(func() {
		utils.WatchMethod(func() error { return updateBlockTimeInMinute(cp, db) })
	}); err != nil {
		return err
	}

	if _, err := scheduler.Every(1).Hour().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateBlockTimeInHour(cp, db) })
	}); err != nil {
		return err
	}

	if _, err := scheduler.Every(1).Day().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateBlockTimeInDay(cp, db) })
	}); err != nil {
		return err
	}

	return nil

}

// updateBlockTimeInMinute insert average block time in the latest minute
func updateBlockTimeInMinute(cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", " tokens").
		Msg("getting total token supply")

	var block tmctypes.ResultBlock
	err := cp.QueryLCD("/blocks/latest", &block)
	if err != nil {
		return err
	}

	genesis, err := db.GetGenesisTime()
	if err != nil {
		return err
	}

	//check if chain is not created minutes ago
	if block.Block.Time.Sub(genesis).Minutes() < 0 {
		return nil
	}

	minute, err := db.GetBlockHeightTimeMinuteAgo(block.Block.Time)
	if err != nil {
		return err
	}
	newBlockTime := block.Block.Time.Sub(minute.Timestamp).Seconds() / float64(block.Block.Height-minute.Height)

	return db.SaveAverageBlockTimePerMin(newBlockTime, block.Block.Time, block.Block.Height)
}

// updateBlockTimeInHour insert average block time in the latest hour
func updateBlockTimeInHour(cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", " tokens").
		Msg("getting total token supply")

	var block tmctypes.ResultBlock
	err := cp.QueryLCD("/blocks/latest", &block)
	if err != nil {
		return err
	}

	genesis, err := db.GetGenesisTime()
	if err != nil {
		return err
	}

	//check if chain is not created minutes ago
	if block.Block.Time.Sub(genesis).Hours() < 0 {
		return nil
	}

	hour, err := db.GetBlockHeightTimeHourAgo(block.Block.Time)
	if err != nil {
		return err
	}
	newBlockTime := block.Block.Time.Sub(hour.Timestamp).Seconds() / float64(block.Block.Height-hour.Height)

	return db.SaveAverageBlockTimePerHour(newBlockTime, block.Block.Time, block.Block.Height)
}

// updateBlockTimeInDay insert average block time in the latest minute
func updateBlockTimeInDay(cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", " tokens").
		Msg("getting total token supply")

	var block tmctypes.ResultBlock
	err := cp.QueryLCD("/blocks/latest", &block)
	if err != nil {
		return err
	}

	genesis, err := db.GetGenesisTime()
	if err != nil {
		return err
	}

	//check if chain is not created days ago
	if block.Block.Time.Sub(genesis).Hours() < 24 {
		return nil
	}

	day, err := db.GetBlockHeightTimeDayAgo(block.Block.Time)
	if err != nil {
		return err
	}
	newBlockTime := block.Block.Time.Sub(day.Timestamp).Seconds() / float64((block.Block.Height - day.Height))

	return db.SaveAverageBlockTimePerDay(newBlockTime, block.Block.Time, block.Block.Height)
}
