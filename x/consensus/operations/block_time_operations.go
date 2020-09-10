package operations

import (
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

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
