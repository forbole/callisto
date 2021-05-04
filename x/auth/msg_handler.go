package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/desmos-labs/juno/modules/messages"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/auth/common"
	"github.com/forbole/bdjuno/x/utils"
)

// HandleMsg handles any message updating the involved accounts
func HandleMsg(
	tx *juno.Tx, msg sdk.Msg, getAddresses messages.MessageAddressesParser, authClient authtypes.QueryClient,
	cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	addresses, err := getAddresses(cdc, msg)
	if err != nil {
		log.Error().Str("module", "auth").Str("operation", "refresh account").
			Err(err).Msgf("error while refreshing accounts after message of type %s", msg.Type())
	}

	return common.UpdateAccounts(utils.FilterNonAccountAddresses(addresses), tx.Height, authClient, cdc, db)
}
