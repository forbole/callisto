package bank

import (
	"encoding/json"

	tmtypes "github.com/cometbft/cometbft/types"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// TODO: implement method if need
	return nil
}
