package inflation

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v3/types"

	inflationtypes "github.com/evmos/evmos/v6/x/inflation/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "inflation").Msg("parsing genesis")

	// Read the genesis state
	var genState inflationtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[inflationtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading mint genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveEvmosInflationParams(types.NewEvmosInflationParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis evmos inflation params: %s", err)
	}

	return nil
}
