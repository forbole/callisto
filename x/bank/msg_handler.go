package bank

import (
	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/x/auth"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
)

func Handler(tx *types.Tx, index int, msg sdk.Msg, cp *client.Proxy, db *database.BigDipperDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch bankMSg := msg.(type) {
	case bank.MsgSend:
		accounts := []string{bankMSg.FromAddress.String(), bankMSg.ToAddress.String()}
		return auth.RefreshAccounts(accounts, tx.Height, cp, db)

	case bank.MsgMultiSend:
		var accounts []string
		for _, input := range bankMSg.Inputs {
			accounts = append(accounts, input.Address.String())
		}
		for _, output := range bankMSg.Outputs {
			accounts = append(accounts, output.Address.String())
		}

		return auth.RefreshAccounts(accounts, tx.Height, cp, db)
	}

	return nil
}
