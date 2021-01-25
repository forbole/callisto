package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/forbole/bdjuno/x/auth"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
)

func HandleMsg(
	tx *types.Tx, msg sdk.Msg,
	authClient authtypes.QueryClient, bankClient banktypes.QueryClient, cdc codec.Marshaler,
	db *database.BigDipperDb,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch bankMSg := msg.(type) {
	case *banktypes.MsgSend:
		accounts := []string{bankMSg.FromAddress, bankMSg.ToAddress}
		return auth.RefreshAccounts(accounts, tx.Height, authClient, bankClient, cdc, db)

	case *banktypes.MsgMultiSend:
		var accounts []string
		for _, input := range bankMSg.Inputs {
			accounts = append(accounts, input.Address)
		}
		for _, output := range bankMSg.Outputs {
			accounts = append(accounts, output.Address)
		}

		return auth.RefreshAccounts(accounts, tx.Height, authClient, bankClient, cdc, db)
	}

	return nil
}
