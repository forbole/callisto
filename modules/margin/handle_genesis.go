package margin

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v3/types"

	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "margin").Msg("parsing genesis")

	// Read the genesis state
	var genState margintypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[margintypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading margin genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveMarginParams(types.NewMarginParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis margin params: %s", err)
	}

	return nil
}
