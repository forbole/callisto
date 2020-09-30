package operations

import (
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UpdateBlockTimeInMinute insert average block time in the latest minute
func UpdateBlockTimeInMinute(cp client.ClientProxy, db database.BigDipperDb) error {
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
	newBlockTime := block.Block.Time.Sub(minute.Timestamp).Seconds() / float64((block.Block.Height - minute.Height))

	return db.SaveAverageBlockTimePerMin(newBlockTime, block.Block.Time, block.Block.Height)
}

// UpdateBlockTimeInHour insert average block time in the latest hour
func UpdateBlockTimeInHour(cp client.ClientProxy, db database.BigDipperDb) error {
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

// UpdateBlockTimeInDay insert average block time in the latest minute
func UpdateBlockTimeInDay(cp client.ClientProxy, db database.BigDipperDb) error {
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

// UpdateBlockTimeFromGenesis insert average block time from genesis
func UpdateBlockTimeFromGenesis(block *tmctypes.ResultBlock, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", " tokens").
		Msg("getting total token supply")

	genesis, err := db.GetGenesisTime()
	if err != nil {
		return err
	}
	newBlockTime := block.Block.Time.Sub(genesis).Seconds() / float64((block.Block.Height))

	return db.SaveAverageBlockTimeGenesis(newBlockTime, block.Block.Time, block.Block.Height)
}
