package vipcoin

import (
	"encoding/json"

	"github.com/forbole/juno/v2/modules"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	for _, module := range m.vipcoinModules {
		if genesisModule, ok := module.(modules.GenesisModule); ok {
			if err := genesisModule.HandleGenesis(doc, appState); err != nil {
				m.logger.GenesisError(module, err)
			}
		}
	}

	return nil
}
