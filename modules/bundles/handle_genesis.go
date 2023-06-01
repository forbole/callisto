package bundles

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/bdjuno/v4/types"

	bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "bundles").Msg("parsing genesis")

	// Read the genesis state
	var genState bundlestypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[bundlestypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading bundles genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveBundlesParams(types.NewBundlesParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis bundles params: %s", err)
	}

	return nil
}
