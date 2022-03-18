package market

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "provider").Msg("parsing genesis")

	// Read the genesis state
	var genState markettypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[markettypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshaling genesis market state: %s", err)
	}

	err = m.db.SaveGenesisLeases(genState.Leases, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while saving genesis leases: %s", err)
	}

	err = m.db.SaveMarketParams(genState.Params, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while saving genesis leases: %s", err)
	}

	return nil
}
