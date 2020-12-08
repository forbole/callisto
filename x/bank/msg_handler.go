package bank

import (
	"time"

	"github.com/desmos-labs/juno/client"
	"github.com/forbole/bdjuno/x/auth"
	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
)

func Handler(tx types.Tx, index int, msg sdk.Msg, cp *client.Proxy, db *database.BigDipperDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	log.Info().
		Str("module", "bank").
		Str("tx_hash", tx.TxHash).
		Int("msg_index", index).
		Str("msg_type", msg.Type()).
		Msg("found message")

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	switch bankMSg := msg.(type) {
	case bank.MsgSend:
		accounts := []sdk.AccAddress{bankMSg.FromAddress, bankMSg.ToAddress}
		return auth.RefreshAccounts(accounts, tx.Height, timestamp, cp, db)

	case bank.MsgMultiSend:
		var accounts []sdk.AccAddress
		for _, input := range bankMSg.Inputs {
			accounts = append(accounts, input.Address)
		}
		for _, output := range bankMSg.Outputs {
			accounts = append(accounts, output.Address)
		}

		return auth.RefreshAccounts(accounts, tx.Height, timestamp, cp, db)
	}

	return nil
}
