package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/handlers"
	"github.com/rs/zerolog/log"
)

// MsgHandler allows to handle the different messages related to the staking module
func MsgHandler(tx types.Tx, index int, msg sdk.Msg, w worker.Worker) error {
	if len(tx.Logs) == 0 {
		log.Info().
			Str("module", "gov").
			Str("tx_hash", tx.TxHash).Int("msg_index", index).
			Msg("skipping message as it was not successful")
		return nil
	}
	log.Info().
		Str("module", "gov").
		Str("tx_hash", tx.TxHash).Int("msg_index", index).
		Msg(msg.Type())

	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("invalid BigDipper database provided")
	}

	switch cosmosMsg := msg.(type) {
	case gov.MsgSubmitProposal:
		return handlers.HandleMsgSubmitProposal(tx, cosmosMsg, bigDipperDb, w.ClientProxy)

	case gov.MsgDeposit:
		return handlers.HandleMsgDeposit(tx, cosmosMsg, bigDipperDb, w.ClientProxy)

	case gov.MsgVote:
		return handlers.HandleMsgVote(tx, cosmosMsg, bigDipperDb, w.ClientProxy)
	}

	return nil
}
