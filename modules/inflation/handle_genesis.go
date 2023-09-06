package inflation

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/bdjuno/v4/types"

	inflationtypes "github.com/evmos/evmos/v14/x/inflation/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "inflation").Msg("parsing genesis")

	// Read the genesis state
	var genState inflationtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[inflationtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading evmos inflation genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveEvmosInflationParams(types.NewEvmosInflationParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis evmos inflation params: %s", err)
	}

	return nil
}
