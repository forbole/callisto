package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/v2/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/modules/utils"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	addresses, err := m.messagesParser(m.cdc, msg)
	if err != nil {
		log.Error().Str("module", "auth").Err(err).
			Str("operation", "refresh account").
			Msgf("error while refreshing accounts after message of type %s", msg.Type())
	}

	return m.RefreshAccounts(tx.Height, utils.FilterNonAccountAddresses(addresses))
}
