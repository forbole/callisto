package slashing

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v2/types"
	tmtypes "github.com/tendermint/tendermint/types"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "distribution").Msg("parsing genesis")

	// Read the genesis state
	var genState slashingtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[slashingtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading mint genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveSlashingParams(types.NewSlashingParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis mint params: %s", err)
	}

	return nil
}
