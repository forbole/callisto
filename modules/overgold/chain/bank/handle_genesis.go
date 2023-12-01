package bank

import (
	"encoding/json"

	tmtypes "github.com/cometbft/cometbft/types"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return nil // don't need to do anything, SDK module
}
