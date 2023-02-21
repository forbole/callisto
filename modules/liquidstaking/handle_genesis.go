package liquidstaking

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v3/types"

	liquidstakingtypes "github.com/crescent-network/crescent/v4/x/liquidstaking/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "liquidstaking").Msg("parsing genesis")

	// Read the genesis state
	var genState liquidstakingtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[liquidstakingtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading liquid staking genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveLiquidStakingParams(types.NewLiquidStakingParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis liquid staking params: %s", err)
	}

	return nil
}
