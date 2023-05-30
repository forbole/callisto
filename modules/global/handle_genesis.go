package global

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/forbole/bdjuno/v4/types"

	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "global").Msg("parsing genesis")

	// Read the genesis state
	var genState globaltypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[globaltypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading global genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveGlobalParams(types.NewGlobalParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis global params: %s", err)
	}

	return nil
}
