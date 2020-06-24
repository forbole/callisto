package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/handlers"
	"github.com/rs/zerolog/log"
)

// MsgHandler allows to handle the different messages related to the staking module
func MsgHandler(tx types.Tx, index int, msg sdk.Msg, w worker.Worker) error {
	if len(tx.Logs) == 0 {
		log.Info().
			Str("module", "staking").
			Str("tx_hash", tx.TxHash).Int("msg_index", index).
			Msg("skipping message as it was not successful")
		return nil
	}

	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("invalid BigDipper database provided")
	}

	switch cosmosMsg := msg.(type) {
	case staking.MsgCreateValidator:
		return handlers.HandleMsgCreateValidator(cosmosMsg, bigDipperDb)

	case staking.MsgDelegate:
		return handlers.HandleMsgDelegate(tx, cosmosMsg, bigDipperDb)

	case staking.MsgBeginRedelegate:
		return handlers.HandleMsgBeginRedelegate(tx, index, cosmosMsg, bigDipperDb)

	case staking.MsgUndelegate:
		return handlers.HandleMsgUndelegate(tx, index, cosmosMsg, bigDipperDb)
	case staking.MsgEditValidator:
		return handlers.HandleEditValidator(cosmosMsg,tx,bigDipperDb)
	}

	return nil
}
