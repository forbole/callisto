package bank

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/auth"
	"github.com/rs/zerolog/log"
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

	timestamp, err := time.Parse("2006-01-02T15:04:05Z", tx.Timestamp)
	if err != nil {
		return err
	}

	switch bankMSg := msg.(type) {
	case bank.MsgSend:
		accounts := []sdk.AccAddress{bankMSg.FromAddress, bankMSg.FromAddress}
		return auth.RefreshAccounts(accounts, tx.Height, timestamp, w.ClientProxy, db)
	case bank.MsgMultiSend:
		var accounts []sdk.AccAddress
		for _, input := range bankMSg.Inputs {
			accounts = append(accounts, input.Address)
		}
		for _, output := range bankMSg.Outputs {
			accounts = append(accounts, output.Address)
		}

		return auth.RefreshAccounts(accounts, tx.Height, timestamp, w.ClientProxy, db)
	}

	return nil
}
