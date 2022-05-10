package slashing

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v3/types"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "slashing").Msg("parsing genesis")

	// Read the genesis state
	var genState slashingtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[slashingtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading mint genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveSlashingParams(types.NewSlashingParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis slashing params: %s", err)
	}

	return nil
}
