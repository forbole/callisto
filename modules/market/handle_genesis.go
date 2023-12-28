package market

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/bdjuno/v4/types"

	markettypes "github.com/akash-network/akash-api/go/node/market/v1beta4"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "market").Msg("parsing genesis")

	// Read the genesis state
	var genState markettypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[markettypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading market genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveMarketParams(types.NewMarketParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis market params: %s", err)
	}

	return nil
}
