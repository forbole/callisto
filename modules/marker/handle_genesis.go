package marker

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v3/types"

	markertypes "github.com/MonCatCat/provenance/x/marker/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "marker").Msg("parsing genesis")

	// Read the genesis state
	var genState markertypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[markertypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading marker genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveMarkerParams(types.NewMarkerParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis marker params: %s", err)
	}

	return nil
}
