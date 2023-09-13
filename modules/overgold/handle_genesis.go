package overgold

import (
	"encoding/json"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/forbole/juno/v5/modules"
)

// HandleGenesis implements modules.GenesisModule
func (m *module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	for _, module := range m.overgoldModules {
		if genesisModule, ok := module.(modules.GenesisModule); ok {
			if err := genesisModule.HandleGenesis(doc, appState); err != nil {
				m.logger.GenesisError(module, err)
			}
		}
	}

	return nil
}
