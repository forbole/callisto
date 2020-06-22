package staking

import (
	"time"

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
	log.Info().Str("tx_hash", tx.TxHash).Int("msg_index", index).
		Str("msg_type", msg.Type()).Msg("found message")

	if len(tx.Logs) == 0 {
		log.Info().Str("tx_hash", tx.TxHash).Int("msg_index", index).
			Msg("skipping message as it was not successful")
		return nil
	}

	db, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("given database instance is not a BigDipperDb")
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}
	if len(tx.Logs) == 0 {
		log.Info().Msg(fmt.Sprintf("Skipping message at index %d of tx hash %s as it was not successull",
			index, tx.TxHash))
		return nil
	}
	switch stakingMsg := msg.(type) {
	case staking.MsgEditValidator:
		// TODO: Handle message here
		//store commission rate
		StoreEditValidator(stakingMsg, w.ClientProxy, timestamp, tx.Height, db)
	}

	return nil
}

func StoreEditValidator(msg staking.MsgEditValidator, cp client.ClientProxy, time time.Time, height int64, db database.BigDipperDb) error {
	//should I take from REST or store the message?
	//store the message
	address := msg.ValidatorAddress
	if found, _ := db.HasValidator(address.String()); !found {
		return nil
	}
	db.SaveEditValidator(msg.ValidatorAddress, msg.CommissionRate.Int64(), msg.MinSelfDelegation.Int64(), msg.Description, time, height)
	return nil
}
