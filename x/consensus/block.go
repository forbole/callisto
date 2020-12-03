package consensus

import (
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func HandleBlock(block *tmctypes.ResultBlock, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "gov").
		Int64("block", block.Block.Height).
		Msg("handling block")

	err := updateBlockTimeFromGenesis(block, db)
	if err != nil {
		return err
	}

	return nil
}

// updateBlockTimeFromGenesis insert average block time from genesis
func updateBlockTimeFromGenesis(block *tmctypes.ResultBlock, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", " tokens").
		Msg("getting total token supply")

	genesis, err := db.GetGenesisTime()
	if err != nil {
		return err
	}

	newBlockTime := block.Block.Time.Sub(genesis).Seconds() / float64(block.Block.Height)
	return db.SaveAverageBlockTimeGenesis(newBlockTime, block.Block.Time, block.Block.Height)
}
