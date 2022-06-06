package shield

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"

	tmtypes "github.com/tendermint/tendermint/types"

	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "shield").Msg("parsing genesis")

	// Read the genesis state
	var genState shieldtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[shieldtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshaling shield state: %s", err)
	}

	// Save the params
	err = m.db.SaveShieldPoolParams(types.NewShieldPoolParams(genState.PoolParams, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis shield params: %s", err)
	}

	return nil
}
