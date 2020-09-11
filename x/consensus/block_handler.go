package consensus

import (
	"fmt"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/consensus/operations"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/juno/parse/worker"
	juno "github.com/desmos-labs/juno/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func BlockHandler(block *tmctypes.ResultBlock, txs []juno.Tx, _ *tmctypes.ResultValidators, w worker.Worker) error {
	log.Debug().
		Str("module", "gov").
		Int64("block", block.Block.Height).
		Msg("handling block")
	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("provided database is not a BigDipper database")
	}
	if err := operations.UpdateBlockTimeFromGenesis(block, bigDipperDb); err != nil {
		return err
	}

	return nil
}
