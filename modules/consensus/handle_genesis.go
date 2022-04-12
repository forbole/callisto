package consensus

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"

	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, _ map[string]json.RawMessage) error {
	log.Debug().Str("module", "consensus").Msg("parsing genesis")

	// Save the genesis time
	err := m.db.SaveGenesis(types.NewGenesis(doc.ChainID, doc.GenesisTime, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis time: %s", err)
	}

	return nil
}
