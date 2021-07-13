package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules/messages"
	"github.com/rs/zerolog/log"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/forbole/bdjuno/database"
	authutils "github.com/forbole/bdjuno/modules/auth/utils"
	"github.com/forbole/bdjuno/modules/utils"
)

// HandleMsg handles any message updating the involved accounts
func HandleMsg(msg sdk.Msg, getAddresses messages.MessageAddressesParser, cdc codec.Marshaler, db *database.Db,height int64,authClient authttypes.QueryClient) error {
	addresses, err := getAddresses(cdc, msg)
	if err != nil {
		log.Error().Str("module", "auth").Err(err).
			Str("operation", "refresh account").
			Msgf("error while refreshing accounts after message of type %s", msg.Type())
	}

	return authutils.UpdateAccounts(utils.FilterNonAccountAddresses(addresses), cdc,db,height,authClient)

}
