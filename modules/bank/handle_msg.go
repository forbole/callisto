package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules/messages"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	bankutils "github.com/forbole/bdjuno/modules/bank/utils"
	"github.com/forbole/bdjuno/modules/utils"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg handles any message updating the involved addresses balances
func HandleMsg(
	tx *juno.Tx, msg sdk.Msg, getAddresses messages.MessageAddressesParser, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	addresses, err := getAddresses(cdc, msg)
	if err != nil {
		log.Error().Str("module", "bank").Str("operation", "refresh balances").
			Err(err).Msgf("error while refreshing balances after message of type %s", msg.Type())
	}

	return bankutils.UpdateBalances(utils.FilterNonAccountAddresses(addresses), tx.Height, bankClient, db)
}
