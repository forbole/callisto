package allowed

import (
	"encoding/json"

	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", allowed.ModuleName).Msg("parsing genesis")

	// Unmarshal the bank state
	var allowedState allowed.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[allowed.ModuleName], &allowedState); err != nil {
		return err
	}

	return m.allowedRepo.InsertToAddresses(allowedState.AddressesList...)
}
