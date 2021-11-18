package profiles

import (
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v2/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "profiles").Msg("parsing genesis")

	// Read the genesis state
	var genState profilestypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[profilestypes.ModuleName], &genState)
	if err != nil {
		return err
	}

	return m.db.SaveProfilesParams(types.NewProfilesParams(genState.Params, genesisDoc.InitialHeight))
}
