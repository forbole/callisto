package interchainstaking

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v3/types"

	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "interchainstaking").Msg("parsing genesis")

	// Read the genesis state
	var genState icstypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[icstypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading interchainstaking genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveInterchainStakingParams(types.NewInterchainStakingParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis interchainstaking params: %s", err)
	}

	return nil
}
