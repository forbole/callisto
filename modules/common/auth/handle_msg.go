package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules/messages"
	"github.com/rs/zerolog/log"

	utils2 "github.com/forbole/bdjuno/modules/common/utils"
)

// HandleMsg handles any message updating the involved accounts
func HandleMsg(msg sdk.Msg, getAddresses messages.MessageAddressesParser, cdc codec.Marshaler, db DB) error {
	addresses, err := getAddresses(cdc, msg)
	if err != nil {
		log.Error().Str("module", "auth").Str("operation", "refresh account").
			Err(err).Msgf("error while refreshing accounts after message of type %s", msg.Type())
	}

	return UpdateAccounts(utils2.FilterNonAccountAddresses(addresses), db)
}
