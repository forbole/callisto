package stakers

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/bdjuno/v4/types"

	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "stakers").Msg("parsing genesis")

	// Read the genesis state
	var genState stakerstypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[stakerstypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading stakers genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveStakersParams(types.NewStakersParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis stakers params: %s", err)
	}

	return nil
}
