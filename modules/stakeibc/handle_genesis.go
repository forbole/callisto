package stakeibc

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v4/types"

	stakeibctypes "github.com/Stride-Labs/stride/v6/x/stakeibc/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "stakeibc").Msg("parsing genesis")

	// Read the genesis state
	var genState stakeibctypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[stakeibctypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading stakeibc genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveStakeIBCParams(types.NewStakeIBCParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis stakeibc params: %s", err)
	}

	return nil
}
