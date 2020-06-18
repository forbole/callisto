package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"

	"fmt"
)

func MsgHandler(tx types.Tx, index int, msg sdk.Msg, w worker.Worker) error {
	log.Info().Str("tx_hash", tx.TxHash).Int("msg_index", index).Str("msg_type", msg.Type()).Msg("found message")

	if len(tx.Logs) == 0 {
		log.Info().Msg(fmt.Sprintf("Skipping message at index %d of tx hash %s as it was not successull",
			index, tx.TxHash))
		return nil
	}
	db, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("invalid database")
	}
	switch stakingMsg := msg.(type) {
	case staking.MsgEditValidator:
		// TODO: Handle message here
		//store commission rate
		StoreModifiedCommission(stakingMsg, w.ClientProxy, db)
	}

	return nil
}

func StoreModifiedCommission(msg staking.MsgEditValidator, cp client.ClientProxy, db database.BigDipperDb) error {
	//should I take from REST or store the message?
	//store the message
	address := msg.ValidatorAddress
	commission := msg.CommissionRate
	if found, _ := db.HasValidator(address.String()); !found {
		return nil
	}

	var validator staking.Validator
	endpoint := fmt.Sprintf("/staking/validators/$s", address.String())
	height, ok := cp.QueryLCDWithHeight(endpoint, &validator)
	if ok != nil {
		return ok
	}

	db.SaveVaildatorComission(validator, height)
	return nil
}
