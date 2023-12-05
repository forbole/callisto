package referral

import (
	"encoding/json"

	tmtypes "github.com/cometbft/cometbft/types"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return nil // TODO: add tables stats, user and CRUD methods
}
