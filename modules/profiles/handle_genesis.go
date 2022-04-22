package profiles

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v3/types"

	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "profiles").Msg("parsing genesis")

	// Read the genesis state
	var genState profilestypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[profilestypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshalling profiles genesis params: %s", err)
	}

	// Save the params
	err = m.db.SaveProfilesParams(types.NewProfilesParams(
		genState.Params,
		doc.InitialHeight,
	))
	if err != nil {
		return fmt.Errorf("error while storing genesis profiles params: %s", err)
	}

	return nil
}
