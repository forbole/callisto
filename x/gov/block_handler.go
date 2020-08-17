package gov

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/parse/worker"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/handlers"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// OneShotOperation represents a method that is called only once when setting up bdjuno
func BlockHandler(block *tmctypes.ResultBlock, txs []juno.Tx, _ *tmctypes.ResultValidators, w worker.Worker) error {
	log.Debug().
		Str("module", "gov").
		Str("operation", "Block Handler").
		Str("block", string(block.Block.Height)).
		Msg("Block Handler")
	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("provided database is not a BigDipper database")
	}

	for _, tx := range txs {
		for _, msg := range tx.Messages {
			switch cosmosMsg := msg.(type) {
			case gov.MsgSubmitProposal:
				//when the proposal first submitted
				return handlers.HandleMsgSubmitProposal(tx, cosmosMsg, bigDipperDb, w.ClientProxy)
			case gov.MsgDeposit:
				return handlers.HandleMsgDeposit(tx, cosmosMsg, bigDipperDb, w.ClientProxy)
			case gov.MsgVote:
				return handlers.HandleMsgVote(tx, cosmosMsg, bigDipperDb, w.ClientProxy)
			}
		}
	}

	return nil
}
