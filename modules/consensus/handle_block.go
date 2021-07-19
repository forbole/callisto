package consensus

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func HandleBlock(block *tmctypes.ResultBlock, db *database.Db) error {
	err := updateBlockTimeFromGenesis(block, db)
	if err != nil {
		log.Error().Str("module", "consensus").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating block time from genesis")
	}

	return nil
}

// updateBlockTimeFromGenesis insert average block time from genesis
func updateBlockTimeFromGenesis(block *tmctypes.ResultBlock, db *database.Db) error {
	log.Trace().Str("module", "consensus").Int64("height", block.Block.Height).
		Msg("updating block time from genesis")

	genesis, err := db.GetGenesis()
	if genesis == nil {
		return errors.New("Genesis table are empty")
	}
	if err != nil {
		return err
	}

	newBlockTime := block.Block.Time.Sub(genesis.Time).Seconds() / float64(block.Block.Height-genesis.InitialHeight)
	return db.SaveAverageBlockTimeGenesis(newBlockTime, block.Block.Height)
}
