package deployment

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/bdjuno/v4/types"

	deploymenttypes "github.com/akash-network/akash-api/go/node/deployment/v1beta3"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "deployment").Msg("parsing genesis")

	// Read the genesis state
	var genState deploymenttypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[deploymenttypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading deployment genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveDeploymentParams(types.NewDeploymentParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis deployment params: %s", err)
	}

	return nil
}
