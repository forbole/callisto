package consensus

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v2/types"

	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, _ map[string]json.RawMessage) error {
	// Save the genesis time
	err := m.db.SaveGenesis(types.NewGenesis(doc.ChainID, doc.GenesisTime, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis time: %s", err)
	}

	return nil
}
