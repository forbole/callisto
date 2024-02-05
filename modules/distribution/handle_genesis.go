package distribution

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/callisto/v4/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "distribution").Msg("parsing genesis")

	// Read the genesis state
	var genState distrtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[distrtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading distribution genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveDistributionParams(types.NewDistributionParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis distribution params: %s", err)
	}

	return nil
}
