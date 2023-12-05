package core

import (
	"encoding/json"

	core "git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", core.ModuleName).Msg("parsing genesis")

	// Unmarshal the core state
	var coreState core.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[core.ModuleName], &coreState); err != nil {
		return err
	}

	return nil // TODO: add table stats and methods
}
