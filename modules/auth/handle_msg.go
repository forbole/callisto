package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/modules/utils"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, _ *juno.Tx) error {
	addresses, err := m.messagesParser(m.cdc, msg)
	if err != nil {
		log.Error().Str("module", "auth").Err(err).
			Str("operation", "refresh account").
			Msgf("error while refreshing accounts after message of type %s", msg.Type())
	}

	return RefreshAccounts(utils.FilterNonAccountAddresses(addresses), m.db)
}
